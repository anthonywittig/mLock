package ezlo

import (
	"context"
	"fmt"
)

func GetControllers(ctx context.Context) ([]deviceResponse, error) {
	ad, err := getAuthData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting auth data: %s", err.Error())
	}

	ds, err := getDevices(ctx, ad)
	if err != nil {
		return nil, fmt.Errorf("error getting devices: %s", err.Error())
	}

	drs := []deviceResponse{}
	for _, d := range ds {
		dr, err := getDevice(ctx, ad, d)
		if err != nil {
			return nil, fmt.Errorf("error getting device: %s", err.Error())
		}

		if dr.NMAControllerStatus != 1 {
			// It's not online.
			continue
		}

		drs = append(drs, dr)
	}

	return drs, nil
}
