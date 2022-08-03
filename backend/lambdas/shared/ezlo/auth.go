package ezlo

import (
	"context"
	"fmt"
	mshared "mlock/shared"
)

func getAuthData(ctx context.Context) (authData, error) {
	username, err := mshared.GetConfig("EZLO_USERNAME")
	if err != nil {
		return authData{}, fmt.Errorf("error getting username: %s", err.Error())
	}

	password, err := mshared.GetConfig("EZLO_PASSWORD")
	if err != nil {
		return authData{}, fmt.Errorf("error getting password: %s", err.Error())
	}

	ad, err := authenticate(ctx, username, password)
	if err != nil {
		return authData{}, fmt.Errorf("error authenticating: %s", err.Error())
	}

	return ad, nil
}
