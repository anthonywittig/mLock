package main

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/climatecontrol"
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

	climateControlRepository := climatecontrol.NewRepository()
	// deviceRepository := device.NewRepository()
	haRepository, err := homeassistant.NewRepository()
	if err != nil {
		return Response{}, fmt.Errorf("error creating climate control repository: %s", err.Error())
	}
	// unitsRepository := unit.NewRepository()

	rawClimateControls, err := haRepository.ListClimateControls(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting climate controls: %s", err.Error())
	}

	existingClimateControls, err := climateControlRepository.List(ctx)
	if err != nil {
		return Response{}, fmt.Errorf("error getting existing climate controls: %s", err.Error())
	}

	idNamespace := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	for _, rawClimateControl := range rawClimateControls {
		fmt.Printf("rawClimateControl: %+v\n", rawClimateControl)

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

	/*
		existingClimateControlsByFriendlyName := climateControlRepository.GroupByFriendlyNamePrefix(existingClimateControls)
		units, err := unitsRepository.List(ctx)
		if err != nil {
			return Response{}, fmt.Errorf("error getting units: %s", err.Error())
		}
		devicesByUnit, err := deviceRepository.ListByUnit(ctx)
		if err != nil {
			return Response{}, fmt.Errorf("error getting devices by unit: %s", err.Error())
		}

		for _, unit := range units {
			unitClimateControls, ok := existingClimateControlsByFriendlyName[unit.Name]
			if !ok {
				continue
			}

			devices, ok := devicesByUnit[unit.ID]
			if !ok {
				continue
			}

			occupiedNow, err := unit.OccupiedStatusForDay(devices, time.Now())
			if err != nil {
				return Response{}, fmt.Errorf("error getting occupied status for day: %s", err.Error())
			}

			fmt.Printf("test --- occupiedNow: %+v\n", occupiedNow)
			fmt.Printf("test --- unitClimateControls: %+v\n", unitClimateControls)

			// If it's not currently occupied and won't be occupied in 6 hours from now, use the vacant settings.

			// If it's not currently occupied but will be in 6 hours from now, use the occupied settings.

		}
	*/

	/*
		for _, existingClimateControl := range existingClimateControls {
			// If it's not currently occupied and won't be occupied in 6 hours from now, use the vacant settings. Otherwise, use the occupied settings.
			existingClimateControl.
		}
	*/

	return Response{
		Message: "ok",
	}, nil
}
