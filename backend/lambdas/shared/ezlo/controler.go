package ezlo

import (
	"context"
	"fmt"
)

func GetControllers(ctx context.Context) ([]deviceResponse, []deviceResponse, error) {
	ad, err := getAuthData(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting auth data: %s", err.Error())
	}

	ds, err := getDevices(ctx, ad)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting devices: %s", err.Error())
	}

	online := []deviceResponse{}
	offline := []deviceResponse{}
	for _, d := range ds {
		dr, err := getDevice(ctx, ad, d)
		if err != nil {
			return nil, nil, fmt.Errorf("error getting device: %s", err.Error())
		}

		if dr.NMAControllerStatus != 1 {
			offline = append(offline, dr)
		} else {
			online = append(online, dr)
		}
	}

	return online, offline, nil
}
