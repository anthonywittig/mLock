package ezlo

import (
	"context"
	"fmt"
	"mlock/lambdas/shared"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type DeviceController struct{}

func NewDeviceController() *DeviceController {
	return &DeviceController{}
}

func (d *DeviceController) AddLockCode(ctx context.Context, prop shared.Property, device shared.Device, code string) error {
	if prop.ControllerID == "" {
		return fmt.Errorf("property doesn't have a controller ID")
	}

	ws, err := getConnection(ctx, prop.ControllerID)
	if err != nil {
		return fmt.Errorf("error getting websocket: %s", err.Error())
	}
	defer ws.Close()

	lockCodes, item, err := wsGetLockCodesForDevice(ws, device.RawDevice.ID)
	if err != nil {
		return fmt.Errorf("error getting lock codes: %s", err.Error())
	}

	if len(lockCodes) >= item.ElementsMaxNumber {
		return fmt.Errorf("max number of lock codes already set")
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

func (d *DeviceController) GetDevices(ctx context.Context, prop shared.Property) ([]shared.RawDevice, error) {
	if prop.ControllerID == "" {
		return nil, nil
	}

	ws, err := getConnection(ctx, prop.ControllerID)
	if err != nil {
		return []shared.RawDevice{}, fmt.Errorf("error getting websocket: %s", err.Error())
	}
	defer ws.Close()

	devices, err := getRawDevices(ws)
	if err != nil {
		return []shared.RawDevice{}, fmt.Errorf("error getting raw devices: %s", err.Error())
	}

	return devices, nil
}

func (d *DeviceController) RemoveLockCode(ctx context.Context, prop shared.Property, device shared.Device, code string) error {
	if prop.ControllerID == "" {
		return fmt.Errorf("property doesn't have a controller ID")
	}

	ws, err := getConnection(ctx, prop.ControllerID)
	if err != nil {
		return fmt.Errorf("error getting websocket: %s", err.Error())
	}
	defer ws.Close()

	lockCodes, item, err := wsGetLockCodesForDevice(ws, device.RawDevice.ID)
	if err != nil {
		return fmt.Errorf("error getting lock codes: %s", err.Error())
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

			return lockCodes, item, nil
		}
	}

	return []shared.RawDeviceLockCode{}, wsItem{}, fmt.Errorf("couldn't find lock codes")
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
