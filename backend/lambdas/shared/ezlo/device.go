package ezlo

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type DeviceController struct {
	connectionPool *ConnectionPool
}

func NewDeviceController(cp *ConnectionPool) *DeviceController {
	return &DeviceController{
		connectionPool: cp,
	}
}

func (d *DeviceController) AddLockCode(ctx context.Context, device shared.Device, code string) error {
	if device.ControllerID == "" {
		return fmt.Errorf("device doesn't have a controller ID")
	}

	ws, err := d.connectionPool.GetConnection(ctx, device.ControllerID)
	if err != nil {
		return fmt.Errorf("error getting websocket: %s", err.Error())
	}

	lockCodes, item, err := wsGetLockCodesForDevice(ws, device.RawDevice.ID)
	if err != nil {
		return fmt.Errorf("error getting lock codes for device \"%s\": %s", device.RawDevice.Name, err.Error())
	}

	if len(lockCodes) >= item.ElementsMaxNumber {
		return fmt.Errorf("max number of lock codes already set (%d)", item.ElementsMaxNumber)
	}

	for _, lc := range lockCodes {
		if lc.Code == code {
			// Assume this is a retry of some sort.
			return nil
		}
	}

	lc := shared.RawDeviceLockCode{
		Code: code,
		Mode: "enabled",
		Name: code,
	}

	err = wsAddLockCodeForItem(ws, item, lc)
	if err != nil {
		return fmt.Errorf("error adding lock code: %s", err.Error())
	}

	return nil
}

func (d *DeviceController) GetDevices(ctx context.Context, controllerID string) ([]shared.RawDevice, error) {
	if controllerID == "" {
		return nil, nil
	}

	ws, err := d.connectionPool.GetConnection(ctx, controllerID)
	if err != nil {
		return []shared.RawDevice{}, fmt.Errorf("error getting websocket: %s", err.Error())
	}

	devices, err := getRawDevices(ws)
	if err != nil {
		return []shared.RawDevice{}, fmt.Errorf("error getting raw devices: %s", err.Error())
	}

	return devices, nil
}

func (d *DeviceController) RediscoverDevice(ctx context.Context, device shared.Device) error {
	if device.ControllerID == "" {
		return fmt.Errorf("device doesn't have a controller ID")
	}

	ws, err := d.connectionPool.GetConnection(ctx, device.ControllerID)
	if err != nil {
		return fmt.Errorf("error getting websocket: %s", err.Error())
	}

	settings, err := wsGetDeviceSettings(ws, device.RawDevice.ID)
	if err != nil {
		return fmt.Errorf("error getting device settings for \"%s\": %s", device.RawDevice.Name, err.Error())
	}

	rediscoverSettingID := ""
	for _, setting := range settings {
		if setting.Label.Text == "Rediscover device" {
			rediscoverSettingID = setting.ID
			break
		}
	}
	if rediscoverSettingID == "" {
		return fmt.Errorf("couldn't find rediscover setting")
	}

	if err := wsSetDeviceSetting(ws, rediscoverSettingID, ""); err != nil {
		return fmt.Errorf("error rediscovering device: %s", err.Error())
	}

	return nil
}

func (d *DeviceController) RebootController(ctx context.Context, device shared.Device) error {
	if device.ControllerID == "" {
		return fmt.Errorf("device doesn't have a controller ID")
	}

	ws, err := d.connectionPool.GetConnection(ctx, device.ControllerID)
	if err != nil {
		return fmt.Errorf("error getting websocket: %s", err.Error())
	}

	if err := wsRebootHub(ws); err != nil {
		return fmt.Errorf("error rebooting controller: %s", err.Error())
	}

	return nil
}

func (d *DeviceController) RemoveLockCode(ctx context.Context, device shared.Device, code string) error {
	// TODO: if multiple codes for the same device are getting removed within a short period of time, might we end up removing the wrong code?
	if device.ControllerID == "" {
		return fmt.Errorf("device doesn't have a controller ID")
	}

	ws, err := d.connectionPool.GetConnection(ctx, device.ControllerID)
	if err != nil {
		return fmt.Errorf("error getting websocket: %s", err.Error())
	}

	lockCodes, item, err := wsGetLockCodesForDevice(ws, device.RawDevice.ID)
	if err != nil {
		return fmt.Errorf("error getting lock codes for device \"%s\": %s", device.RawDevice.Name, err.Error())
	}

	slot := -1
	for _, lc := range lockCodes {
		if lc.Code == code {
			slot = lc.Slot
			break
		}
	}
	if slot == -1 {
		// If we can't find it, assume this was part of a retry and it's gone now.
		return nil
	}

	err = wsRemoveLockCodeForItem(ws, item, fmt.Sprintf("%d", slot))
	if err != nil {
		return fmt.Errorf("error removing lock code: %s", err.Error())
	}

	return nil
}

func wsAddLockCodeForItem(ws *websocket.Conn, item wsItem, lockCode shared.RawDeviceLockCode) error {
	// https://api.ezlo.com/hub/items_api/#hubitemdictionaryvalueadd
	// https://api.ezlo.com/devices/item_value_types/index.html

	id := fmt.Sprintf("hub.item.dictionary.value.add.%s", uuid.New())
	resp := wsItemsListResponse{}

	type paramsValue struct {
		Code string `json:"code"`
		Mode string `json:"mode"`
		Name string `json:"name"`
	}
	type paramsElement struct {
		Type  string      `json:"type"`
		Value paramsValue `json:"value"`
	}
	type params struct {
		ID      string        `json:"_id"`
		Element paramsElement `json:"element"`
	}

	err := wsSendCommand(
		ws,
		id,
		struct {
			Method string `json:"method"`
			ID     string `json:"id"`
			Params params `json:"params"`
		}{
			Method: "hub.item.dictionary.value.add",
			ID:     id,
			Params: params{
				ID: item.ID,
				Element: paramsElement{
					Type: "userCode",
					Value: paramsValue{
						Code: lockCode.Code,
						Mode: lockCode.Mode,
						Name: lockCode.Name,
					},
				},
			},
		},
		&resp,
	)
	if err != nil {
		return fmt.Errorf("error sending command: %s", err.Error())
	}

	return nil
}

func wsGetLockCodesForDevice(ws *websocket.Conn, deviceID string) ([]shared.RawDeviceLockCode, wsItem, error) {
	id := fmt.Sprintf("hub.items.list.%s", uuid.New())
	resp := wsItemsListResponse{}
	type params struct {
		DeviceIDs []string `json:"deviceIds"`
	}
	err := wsSendCommand(
		ws,
		id,
		struct {
			Method string `json:"method"`
			ID     string `json:"id"`
			Params params `json:"params"`
		}{
			Method: "hub.items.list",
			ID:     id,
			Params: params{
				DeviceIDs: []string{deviceID},
			},
		},
		&resp,
	)
	if err != nil {
		return []shared.RawDeviceLockCode{}, wsItem{}, fmt.Errorf("error sending command: %s", err.Error())
	}

	for _, item := range resp.Result.Items {
		if item.Name == "user_codes" {
			lockCodes, err := item.getLockCodes()
			if err != nil {
				return []shared.RawDeviceLockCode{}, wsItem{}, fmt.Errorf("error getting lock codes: %s", err.Error())
			}

			if item.ElementsMaxNumber == 0 {
				// We didn't get an max number, some locks are as low as 6, but it's probably better to not artificially limit them.
				item.ElementsMaxNumber = 30
			}

			return lockCodes, item, nil
		}
	}

	return []shared.RawDeviceLockCode{}, wsItem{}, fmt.Errorf("couldn't find lock codes for deviceID: %s", deviceID)
}

func wsRebootHub(ws *websocket.Conn) error {
	/*
		Request:

		method: "hub.reboot"
		id: "1688852614989"
		params: Object {}


		Response:

		error: null
		id: "1688852614989"
		method: "hub.reboot"
		result: Object {}
		sender: Object {"conn_id":"84da5a58-1ab3-410b-9fe2-11181f92f06f","type":"ui"}
	*/
	id := fmt.Sprintf("hub.reboot.%s", uuid.New())

	type response struct{}
	resp := response{}

	type params struct{}
	err := wsSendCommand(
		ws,
		id,
		struct {
			Method string `json:"method"`
			ID     string `json:"id"`
			Params params `json:"params"`
		}{
			Method: "hub.reboot",
			ID:     id,
			Params: params{},
		},
		&resp,
	)
	if err != nil {
		return fmt.Errorf("error sending command: %s", err.Error())
	}

	return nil
}

func wsRemoveLockCodeForItem(ws *websocket.Conn, item wsItem, slot string) error {
	id := fmt.Sprintf("hub.item.dictionary.value.remove.%s", uuid.New())
	resp := wsItemsListResponse{}

	type params struct {
		ID  string `json:"_id"`
		Key string `json:"key"`
	}

	err := wsSendCommand(
		ws,
		id,
		struct {
			Method string `json:"method"`
			ID     string `json:"id"`
			Params params `json:"params"`
		}{
			Method: "hub.item.dictionary.value.remove",
			ID:     id,
			Params: params{
				ID:  item.ID,
				Key: slot,
			},
		},
		&resp,
	)
	if err != nil {
		return fmt.Errorf("error sending command: %s", err.Error())
	}

	return nil
}
