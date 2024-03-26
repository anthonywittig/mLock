package main

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"mlock/lambdas/shared/dynamo/climatecontrol"
	"mlock/lambdas/shared/dynamo/device"
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

	climateControlRepository := climatecontrol.NewRepository()
	deviceRepository := device.NewRepository()
	haRepository, err := homeassistant.NewRepository()
	if err != nil {
		return Response{}, fmt.Errorf("error creating climate control repository: %s", err.Error())
	}
	unitsRepository := unit.NewRepository()

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

	return Response{
		Message: "ok",
	}, nil
}
