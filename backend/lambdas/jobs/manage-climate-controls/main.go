package main

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/climatecontrol"
	"mlock/lambdas/shared/dynamo/device"
	"mlock/lambdas/shared/dynamo/miscellaneous"
	"mlock/lambdas/shared/dynamo/unit"
	"mlock/lambdas/shared/homeassistant"
	mshared "mlock/shared"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

type MyEvent struct {
}

type Response struct {
	Message string `json:"message"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest(ctx context.Context, event MyEvent) (Response, error) {
	ctx = shared.CreateContextData(ctx)

	if err := mshared.LoadConfig(); err != nil {
		return Response{}, fmt.Errorf("error loading config: %s", err.Error())
	}

	tzName, err := mshared.GetConfig("TIME_ZONE")
	if err != nil {
		return Response{}, fmt.Errorf("error getting time zone name: %s", err.Error())
	}

	tz, err := time.LoadLocation(tzName)
	if err != nil {
		return Response{}, fmt.Errorf("error getting time zone %s", err.Error())
	}

	climateControlRepository := climatecontrol.NewRepository()
	haRepository, err := homeassistant.NewRepository()
	if err != nil {
		return Response{}, fmt.Errorf("error creating climate control repository: %s", err.Error())
	}

	devices, err := device.NewRepository().List(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting devices: %s", err.Error())
	}

	miscellaneous, ok, err := miscellaneous.NewRepository().Get(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting miscellaneous: %s", err.Error())
	}
	if !ok {
		return Response{}, fmt.Errorf("miscellaneous not found: %s", err.Error())
	}

	units, err := unit.NewRepository().ListByName(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting units: %s", err.Error())
	}

	if err := refreshClimateControls(ctx, climateControlRepository, haRepository); err != nil {
		return Response{}, fmt.Errorf("error refreshing climate controls: %s", err.Error())
	}

	now := time.Now().In(tz)

	fourPM := time.Date(now.Year(), now.Month(), now.Day(), 16, 0, 0, 0, tz)

	// We'll only try to set the desired state if it's between 12pm and 4pm.
	if now.Hour() >= 12 && now.Before(fourPM) {
		existingClimateControls, err := climateControlRepository.List(ctx)
		if err != nil {
			return Response{}, fmt.Errorf("error getting existing climate controls: %s", err.Error())
		}
		for _, ecc := range existingClimateControls {
			u := units[ecc.GetFriendlyNamePrefix()]
			os := u.OccupancyStatusForDay(devices, now)

			if os.At.Occupied {
				// It's currently occupied, let's kill off the desired state.
				if ecc.DesiredState.EndAt.After(now) {
					ecc.DesiredState.EndAt = now
					climateControlRepository.AppendToAuditLog(ctx, ecc, "Killed off the desired state since the unit is occupied.")
					climateControlRepository.Put(ctx, ecc)
				}
				continue
			}

			if ecc.DesiredState.EndAt.After(now) && !ecc.DesiredState.SyncWithSettings {
				// There's a non-syncing setting in place, don't make a change.
				continue
			}

			var newDesiredState *shared.ClimateControlDesiredState = nil

			if !os.FourPM.Occupied {
				// It's not occupied now or at 4pm, let's set it to vacant.
				newDesiredState = &shared.ClimateControlDesiredState{
					EndAt:            fourPM,
					HVACMode:         miscellaneous.ClimateControlVacantSettings.HVACMode,
					Note:             "Setting to vacant",
					SyncWithSettings: true,
					Temperature:      miscellaneous.ClimateControlVacantSettings.Temperature,
				}
			}

			if !os.Noon.Occupied && os.FourPM.Occupied {
				// It'll change from not occupied to occupied.
				newDesiredState = &shared.ClimateControlDesiredState{
					EndAt:            fourPM,
					HVACMode:         miscellaneous.ClimateControlOccupiedSettings.HVACMode,
					Note:             fmt.Sprintf("Setting to occupied for reservation: %s", os.FourPM.ManagedLockCodes[0].Reservation.ID),
					SyncWithSettings: true,
					Temperature:      miscellaneous.ClimateControlOccupiedSettings.Temperature,
				}
			}

			if newDesiredState != nil {
				newIsDifferent := false
				if !ecc.DesiredState.EndAt.Equal(newDesiredState.EndAt) {
					newIsDifferent = true
				}
				if ecc.DesiredState.HVACMode != newDesiredState.HVACMode {
					newIsDifferent = true
				}
				if ecc.DesiredState.Note != newDesiredState.Note {
					newIsDifferent = true
				}
				if ecc.DesiredState.SyncWithSettings != newDesiredState.SyncWithSettings {
					newIsDifferent = true
				}
				if ecc.DesiredState.Temperature != newDesiredState.Temperature {
					newIsDifferent = true
				}
				if newIsDifferent {
					ecc.DesiredState = *newDesiredState
					climateControlRepository.AppendToAuditLog(ctx, ecc, ecc.DesiredState.Note)
					climateControlRepository.Put(ctx, ecc)
				}
			}
		}

		// Pull in the new controls since we just added/updated them.
		existingClimateControls, err = climateControlRepository.List(ctx)
		if err != nil {
			return Response{}, fmt.Errorf("error getting existing climate controls: %s", err.Error())
		}
		attemptedToUpdateAClimateControl := false
		for _, ecc := range existingClimateControls {
			if ecc.RawClimateControl.State == "unavailable" {
				continue
			}
			if ecc.DesiredState.EndAt.Before(now) {
				continue
			}
			if ecc.DesiredState.HVACMode == ecc.RawClimateControl.State && ecc.DesiredState.Temperature == ecc.RawClimateControl.Attributes.Temperature {
				continue
			}
			if ecc.GetFriendlyNamePrefix() != "03A" {
				fmt.Printf("Skipping climate control: %+v\n", ecc.RawClimateControl.Attributes.FriendlyName)
				continue
			}
			fmt.Printf("Updating climate control: %+v\n", ecc.RawClimateControl.Attributes.FriendlyName)
			climateControlRepository.AppendToAuditLog(ctx, ecc, fmt.Sprintf("Attempting to update the climate control's settings; HVAC mode: %s, temperature: %d", ecc.DesiredState.HVACMode, ecc.DesiredState.Temperature))
			if err := haRepository.SetToDesiredState(ctx, ecc); err != nil {
				// Might need to swallow these errors or at least try all the others before returning.
				return Response{}, fmt.Errorf("error setting to desired state: %s", err.Error())
			}
			attemptedToUpdateAClimateControl = true
		}

		if attemptedToUpdateAClimateControl {
			if err := refreshClimateControls(ctx, climateControlRepository, haRepository); err != nil {
				return Response{}, fmt.Errorf("error refreshing climate controls: %s", err.Error())
			}
		}
	}

	return Response{
		Message: "ok",
	}, nil
}

func refreshClimateControls(
	ctx context.Context,
	climateControlRepository *climatecontrol.Repository,
	haRepository *homeassistant.Repository,
) error {
	idNamespace := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	existingClimateControls, err := climateControlRepository.List(ctx)
	if err != nil {
		return fmt.Errorf("error getting existing climate controls: %s", err.Error())
	}

	rawClimateControls, err := haRepository.ListClimateControls(ctx)
	if err != nil {
		return fmt.Errorf("error getting climate controls: %s", err.Error())
	}

	for _, rawClimateControl := range rawClimateControls {
		if rawClimateControl.State == "unavailable" {
			// For now, let's skip these.
			continue
		}

		climateControl := shared.ClimateControl{
			ID:      uuid.NewSHA1(idNamespace, []byte(rawClimateControl.EntityID)),
			History: []shared.ClimateControlHistory{},
		}

		for _, existingClimateControl := range existingClimateControls {
			if existingClimateControl.ID == climateControl.ID {
				// At the time of writing, there's no data that we actually care to preserve, but hopefully this is a good pattern.
				climateControl = existingClimateControl
			}
		}

		climateControl.LastRefreshedAt = time.Now()
		climateControl.RawClimateControl = rawClimateControl

		climateControlRepository.Put(ctx, climateControl)
	}

	return nil
}
