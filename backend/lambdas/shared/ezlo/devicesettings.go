package ezlo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type wsDeviceSetting struct {
	ID       string `json:"_id"`
	DeviceID string `json:"deviceId"`
	Label    struct {
		Text string `json:"text"`
	} `json:"label"`
}

type wsDeviceSettingsListResponse struct {
	API   string `json:"api"`
	Error *struct {
		Code        int    `json:"code"`
		Data        string `json:"data"`
		Description string `json:"description"`
	} `json:"error"`
	ID     string                             `json:"id"`
	Method string                             `json:"method"`
	Result wsDeviceSettingsListResponseResult `json:"result"`
}

type wsDeviceSettingsListResponseResult struct {
	Settings []wsDeviceSetting `json:"settings"`
}

type wsSetDeviceSettingResponse struct {
	Error *struct {
		Code        int    `json:"code"`
		Data        string `json:"data"`
		Description string `json:"description"`
	} `json:"error"`
	ID     string `json:"id"`
	Method string `json:"method"`
}

func wsGetDeviceSettings(ws *websocket.Conn, deviceID string) ([]wsDeviceSetting, error) {
	method := "hub.device.settings.list"
	id := fmt.Sprintf("%s.%s", method, uuid.New())
	resp := wsDeviceSettingsListResponse{}
	type params struct{}
	err := wsSendCommand(
		ws,
		id,
		struct {
			Method string `json:"method"`
			ID     string `json:"id"`
			Params params `json:"params"`
		}{
			Method: method,
			ID:     id,
			Params: params{},
		},
		&resp,
	)
	if err != nil {
		return nil, fmt.Errorf("error sending command: %s", err.Error())
	}

	results := []wsDeviceSetting{}
	for _, setting := range resp.Result.Settings {
		if setting.DeviceID == deviceID {
			results = append(results, setting)
		}
	}

	return results, nil
}

func wsSetDeviceSetting(ws *websocket.Conn, settingID string, value string) error {
	method := "hub.device.setting.value.set"
	id := fmt.Sprintf("%s.%s", method, uuid.New())
	resp := wsSetDeviceSettingResponse{}

	type params struct {
		ID    string `json:"_id"`
		Value string `json:"value"`
	}

	err := wsSendCommand(
		ws,
		id,
		struct {
			Method string `json:"method"`
			ID     string `json:"id"`
			Params params `json:"params"`
		}{
			Method: method,
			ID:     id,
			Params: params{
				ID:    settingID,
				Value: value,
			},
		},
		&resp,
	)
	if err != nil {
		return fmt.Errorf("error sending command: %s", err.Error())
	}

	return nil
}
