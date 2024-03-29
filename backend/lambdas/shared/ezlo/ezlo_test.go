package ezlo

import (
	"context"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Test_GetDevicesNoController(t *testing.T) {
	assert.Nil(t, loadConfig())

	cp := NewConnectionPool()
	defer cp.Close()
	dc := NewDeviceController(cp)

	ds, err := dc.GetDevices(context.Background(), "")
	assert.Nil(t, err)
	assert.Empty(t, ds)
}

func Test_GetDevices(t *testing.T) {
	assert.Nil(t, loadConfig())

	cp := NewConnectionPool()
	defer cp.Close()
	dc := NewDeviceController(cp)

	ds, err := dc.GetDevices(context.Background(), "92001809")
	assert.Nil(t, err)

	// Just check that we have at least one device (fragile).
	assert.Greater(t, len(ds), 0)
	assert.Equal(t, ds[0].Category, "siren")

	// Uncomment to see everything in an error
	//assert.Nil(t, ds)
}

func Test_TemporarySearchDevices(t *testing.T) {
	assert.Nil(t, loadConfig())

	/*
		for _, cID := range []string{
			"84901957",
			"84902125",
			"84902278",
			"84911082",
			"90010778",
			"90010799",
			"92001809",
		} {
			cp := NewConnectionPool()
			defer cp.Close()
			dc := NewDeviceController(cp)

			ds, err := dc.GetDevices(context.Background(), cID)
			assert.Nil(t, err)

			for _, d := range ds {
				if d.DeviceTypeID == "0_0_0" || d.Name == "06B Lock 849749 (1BE469)" {
					fmt.Printf("%s - %s - %+v\n", d.DeviceTypeID, d.Status, d)
				}
			}
		}

		// Uncomment to see everything in an error
		assert.Nil(t, "hi")
	*/
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
