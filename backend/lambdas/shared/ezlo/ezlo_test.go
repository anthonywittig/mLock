package ezlo

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Test_GetDevicesNoController(t *testing.T) {
	assert.Nil(t, loadConfig())

	cp := NewConnectionPool()
	defer cp.Close()
	dc := NewDeviceController(cp)

	ds, err := dc.GetDevices(context.Background(), shared.Property{})
	assert.Nil(t, err)
	assert.Empty(t, ds)
}

func Test_GetDevices(t *testing.T) {
	assert.Nil(t, loadConfig())

	cp := NewConnectionPool()
	defer cp.Close()
	dc := NewDeviceController(cp)

	ds, err := dc.GetDevices(context.Background(), shared.Property{ControllerID: "92001809"})
	assert.Nil(t, err)

	// Just check that we have one device (fragile).
	assert.Greater(t, len(ds), 1)
	assert.Equal(t, ds[0].Category, "siren")

	// Uncomment to see everything in an error
	//assert.Nil(t, ds)
}

func Test_AddLockCode(t *testing.T) {
	assert.Nil(t, loadConfig())

	/*
		// We should really do some audit logging here... :(
		err := AddLockCode(
			context.Background(),
			shared.Property{ControllerID: "92001809"},
			shared.Device{RawDevice: shared.RawDevice{ID: "6159fc8e939a9313ca423558"}}, // "ZC3 (In Box)"
			"3001",
		)
		assert.Nil(t, err)
		assert.NotNil(t, err)
	*/
}

func Test_RemoveLockCode(t *testing.T) {
	assert.Nil(t, loadConfig())

	/*
		// We should really do some audit logging here... :(
		err := RemoveLockCode(
			context.Background(),
			shared.Property{ControllerID: "92001809"},
			shared.Device{RawDevice: shared.RawDevice{ID: "6159fc8e939a9313ca423558"}}, // "ZC3 (In Box)"
			shared.DeviceManagedLockCode{
				Code: "3001",
			},
		)
		assert.Nil(t, err)
		assert.NotNil(t, err)
	*/
}

func loadConfig() error {
	if err := godotenv.Load(".env.test"); err != nil {
		return fmt.Errorf("error loading .env file: %s", err.Error())
	}
	return nil
}
