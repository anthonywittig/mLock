package ezlo

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Test_GetDevices(t *testing.T) {
	assert.Nil(t, loadConfig())

	ds, err := GetDevices(context.Background(), shared.Property{ControllerID: "92001809"})
	assert.Nil(t, err)

	// Just check that we have one device (fragile).
	assert.Greater(t, len(ds), 1)
	assert.Equal(t, ds[0].Category, "siren")

	// Uncomment to see everything in an error
	//assert.Nil(t, ds)
}

func loadConfig() error {
	if err := godotenv.Load(".env.test"); err != nil {
		return fmt.Errorf("error loading .env file: %s", err.Error())
	}
	return nil
}
