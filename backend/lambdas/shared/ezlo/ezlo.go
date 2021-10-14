package ezlo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type AuthResponse struct {
	Identity          string `json:"Identity"`
	IdentitySignature string `json:"IdentitySignature"`
	ServerEvent       string `json:"Server_Event"`
	ServerEventAlt    string `json:"Server_Event_Alt"`
	ServerAccount     string `json:"Server_Account"`
	ServerAccountAlt  string `json:"Server_Account_Alt"`
}

type WSRequest struct {
	Method string   `json:"method"`
	ID     string   `json:"id"`
	Params struct{} `json:"params"`
}

type WSDeviceListResponse struct {
	API    string                     `json:"api"`
	Error  *string                    `json:"error"` // Not sure if this is really a string.
	ID     string                     `json:"id"`
	Method string                     `json:"method"`
	Result WSDeviceListResponseResult `json:"result"`
}
type WSDeviceListResponseResult struct {
	Devices []WSDeviceListResponseResultDevice `json:"devices"`
	Sender  WSDeviceListResponseResultSender   `json:"sender"`
}

type WSDeviceListResponseResultDevice struct {
	ID             string   `json:"_id"`
	BatteryPowered bool     `json:"batteryPowered"`
	Category       string   `json:"category"`
	DeviceTypeID   string   `json:"deviceTypeId"`
	GatewayID      string   `json:"gatewayId"`
	Info           struct{} `json:"info"`
	Name           string   `json:"name"`
	ParentDeviceID string   `json:"parentDeviceId"`
	Persistent     bool     `json:"persistent"`
	Reachable      bool     `json:"reachable"`
	Ready          bool     `json:"ready"`
	RoomID         string   `json:"roomId"`
	Security       string   `json:"security"`
	Status         string   `json:"status"`
	Subcategory    string   `json:"subcategory"`
	Type           string   `json:"type"`
}

type WSDeviceListResponseResultSender struct {
	ConnID string `json:"conn_id"`
	Type   string `json:"type"`
}

type WSLogInRequest struct {
	Method string               `json:"method"`
	ID     string               `json:"id"`
	Params WSLogInRequestParams `json:"params"`
}

type WSLogInRequestParams struct {
	MMSAuth    string `json:"MMSAuth"`
	MMSAuthSig string `json:"MMSAuthSig"`
}

type WSRegisterRequest struct {
	Method  string                  `json:"method"`
	ID      string                  `json:"id"`
	JSONRPC string                  `json:"jsonrpc"`
	Params  WSRegisterRequestParams `json:"params"`
}

type WSRegisterRequestParams struct {
	Serial string `json:"serial"`
}

func Authenticate(ctx context.Context, username string, password string) (AuthResponse, error) {
	// TODO: get a list of the endpoints and use a random one.
	url := fmt.Sprintf("https://vera-us-oem-account11.mios.com/autha/auth/username/%s", username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	q := req.URL.Query()
	q.Add("SHA1Password", password)
	q.Add("SHA1PasswordCS", password)
	q.Add("PK_Oem", "1")
	q.Add("TokenVersion", "2")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	fmt.Printf("%s %s\n", req.Host, req.URL.Path)

	resp, err := client.Do(req)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return AuthResponse{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	body := AuthResponse{}
	json.Unmarshal(respBody, &body)

	return body, nil
}

func WSLogIn(ws *websocket.Conn, authResponse AuthResponse) error {
	logInReq, err := json.Marshal(
		WSLogInRequest{
			Method: "loginUserMios",
			ID:     "loginUser",
			Params: WSLogInRequestParams{
				MMSAuth:    authResponse.Identity,
				MMSAuthSig: authResponse.IdentitySignature,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("marshal: %s", err.Error())
	}

	if err := ws.WriteMessage(websocket.TextMessage, logInReq); err != nil {
		return fmt.Errorf("write: %s", err.Error())
	}

	_, _, err = ws.ReadMessage()
	if err != nil {
		return fmt.Errorf("read: %s", err.Error())
	}

	// TODO: confirm `message` has expected response ID (probably want a random piece to it).

	return nil
}

func WSRegisterHub(ws *websocket.Conn, hubSerialNumber string) error {
	registerReq, err := json.Marshal(
		WSRegisterRequest{
			Method: "register",
			ID:     "register",
			Params: WSRegisterRequestParams{
				Serial: hubSerialNumber,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("marshal: %s", err.Error())
	}

	if err := ws.WriteMessage(websocket.TextMessage, registerReq); err != nil {
		return fmt.Errorf("write: %s", err.Error())
	}

	_, _, err = ws.ReadMessage()
	if err != nil {
		return fmt.Errorf("read: %s", err.Error())
	}

	// TODO: confirm `message` has expected response ID (probably want a random piece to it).

	return nil
}

func WSDeviceList(ws *websocket.Conn) (WSDeviceListResponse, error) {
	id := fmt.Sprintf("hub.devices.list.%s", uuid.New())
	registerReq, err := json.Marshal(
		WSRequest{
			Method: "hub.devices.list",
			ID:     id,
		},
	)
	if err != nil {
		return WSDeviceListResponse{}, fmt.Errorf("marshal: %s", err.Error())
	}

	if err := ws.WriteMessage(websocket.TextMessage, registerReq); err != nil {
		return WSDeviceListResponse{}, fmt.Errorf("write: %s", err.Error())
	}

	_, message, err := ws.ReadMessage()
	if err != nil {
		return WSDeviceListResponse{}, fmt.Errorf("read: %s", err.Error())
	}

	// TODO: confirm `message` has expected response ID (probably want a random piece to it).
	// We should probably check the response ID before unmarshalling it.

	resp := WSDeviceListResponse{}
	if err := json.Unmarshal(message, &resp); err != nil {
		return WSDeviceListResponse{}, fmt.Errorf("unmarshal: %s", err.Error())
	}

	return resp, nil
}

func X(context context.Context, authResponse AuthResponse, hubSerialNumber string) (string, error) {
	// TODO: get a list of the endpoints and use a random one.
	wsURL := "nma-server7-ui-cloud.ezlo.com"
	u := url.URL{Scheme: "wss", Host: wsURL, Path: ""}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("dial: %s", err.Error())
	}
	defer ws.Close()

	if err := WSLogIn(ws, authResponse); err != nil {
		return "", fmt.Errorf("login: %s", err.Error())
	}

	if err := WSRegisterHub(ws, hubSerialNumber); err != nil {
		return "", fmt.Errorf("register: %s", err.Error())
	}

	resp, err := WSDeviceList(ws)
	if err != nil {
		return "", fmt.Errorf("device list: %s", err.Error())
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("marshal: %s", err.Error())
	}

	return string(jsonResp), nil
}

/*
	Unexpected receive:
	{
		"id": "ui_broadcast",
		"msg_id": "616891e2939a9314217c08ea",
		"msg_subclass": "hub.device.updated",
		"result": {
			"_id": "61576aad939a9313ca42350c",
			"reachable": true,
			"serviceNotification": false,
			"syncNotification": false
		}
	}

	Unexpected receive:
	{
		"id": "ui_broadcast",
		"msg_id": "6168969c939a9314217c08eb",
		"msg_subclass": "hub.item.updated",
		"result": {
			"_id": "61576aae939a9313ca423514",
			"deviceCategory": "door_lock",
			"deviceId": "61576aad939a9313ca42350c",
			"deviceName": "1A",
			"deviceSubcategory": "",
			"name": "battery",
			"notifications": null,
			"roomName": "",
			"serviceNotification": false,
			"syncNotification": false,
			"userNotification": false,
			"value": 99,
			"valueFormatted": "99",
			"valueType": "int"
		}
	}







	Sends:
	{
		"method": "loginUserMios",
		"id": "loginUser",
		"params": {
			"MMSAuth": authResponse.Identity,
			"MMSAuthSig": authResponse.IdentitySignature,
		}
	}
	Receives:
	{
		"id": "loginUser",
		"method": "loginUserMios",
		"error": null,
		"result": {}
	}

	Sends:
	{
		"method": "register",
		"id": "register",
		"jsonrpc": "2.0",
		"params": {
			"serial": "92001809"
		}
	}
	Receives:
	{
		"id": "register",
		"method": "registered",
		"error": null,
		"result": {}
	}

	Sends:
	{"method":"hub.devices.list","id":"1634235234608","params":{}}

	Receives:
	{
		"api": "1.0",
		"error": null,
		"id": "1634235234608",
		"method": "hub.devices.list",
		"result": {
			"devices": [
				{
					"_id": "606704e200000015eb0782eb",
					"batteryPowered": false,
					"category": "siren",
					"deviceTypeId": "e550_siren",
					"gatewayId": "606704e200000015eb0782ea",
					"info": {},
					"name": "Controller Siren",
					"parentDeviceId": "",
					"persistent": true,
					"reachable": true,
					"ready": true,
					"roomId": "",
					"security": "high",
					"status": "idle",
					"subcategory": "",
					"type": "siren"
				},
				{
					"_id": "61576aad939a9313ca42350c",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_25409_20548",
					"firmware": [
						{
							"id": "us.59.28950.0",
							"version": "113.22"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "3.42",
						"hardware": "0",
						"manufacturer": "Schlage",
						"model": "BE469NX",
						"protocol": "zwave",
						"zwave.node": "13",
						"zwave.smartstart": "no"
					},
					"name": "1A",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "middle",
					"status": "idle",
					"subcategory": "",
					"type": "doorlock"
				},
				{
					"_id": "61576f2f939a9313ca423533",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_1_1128",
					"firmware": [
						{
							"id": "us.59.18064.0",
							"version": "3.3"
						},
						{
							"id": "us.59.18065.1",
							"version": "11.0"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "6.3",
						"hardware": "3",
						"manufacturer": "Schlage",
						"model": "Unknown",
						"protocol": "zwave",
						"zwave.node": "17",
						"zwave.smartstart": "no"
					},
					"name": "1B",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "high",
					"status": "idle",
					"subcategory": "",
					"type": "device"
				},
				{
					"_id": "6159fc8e939a9313ca423558",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_1_1128",
					"firmware": [
						{
							"id": "us.59.18064.0",
							"version": "3.3"
						},
						{
							"id": "us.59.18065.1",
							"version": "11.0"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "6.3",
						"hardware": "3",
						"manufacturer": "Schlage",
						"model": "Unknown",
						"protocol": "zwave",
						"zwave.node": "22",
						"zwave.smartstart": "no"
					},
					"name": "ZC3 (In Box)",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": true,
					"ready": true,
					"roomId": "",
					"security": "high",
					"status": "idle",
					"subcategory": "",
					"type": "device"
				},
				{
					"_id": "6159ff37939a9313ca42356c",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_1_1128",
					"firmware": [
						{
							"id": "us.59.18064.0",
							"version": "3.3"
						},
						{
							"id": "us.59.18065.1",
							"version": "11.0"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "6.3",
						"hardware": "3",
						"manufacturer": "Schlage",
						"model": "Unknown",
						"protocol": "zwave",
						"zwave.node": "24",
						"zwave.smartstart": "no"
					},
					"name": "ZC2 (In Box)",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": true,
					"ready": true,
					"roomId": "",
					"security": "high",
					"status": "idle",
					"subcategory": "",
					"type": "device"
				},
				{
					"_id": "615a01ce939a9313ca423580",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_25409_20548",
					"firmware": [
						{
							"id": "us.59.28950.0",
							"version": "113.22"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "3.42",
						"hardware": "0",
						"manufacturer": "Schlage",
						"model": "BE469NX",
						"protocol": "zwave",
						"zwave.node": "25",
						"zwave.smartstart": "no"
					},
					"name": "LV4 (In Box)",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": true,
					"ready": true,
					"roomId": "",
					"security": "middle",
					"status": "idle",
					"subcategory": "",
					"type": "doorlock"
				},
				{
					"_id": "6165d1b8939a9313e502a895",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_25409_20548",
					"firmware": [
						{
							"id": "us.59.28950.0",
							"version": "113.22"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "3.42",
						"hardware": "0",
						"manufacturer": "Schlage",
						"model": "BE469NX",
						"protocol": "zwave",
						"zwave.node": "32",
						"zwave.smartstart": "no"
					},
					"name": "2A",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "middle",
					"status": "idle",
					"subcategory": "",
					"type": "doorlock"
				},
				{
					"_id": "6165d271939a9313e502a8b3",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_25409_20548",
					"firmware": [
						{
							"id": "us.59.28950.0",
							"version": "113.22"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "3.42",
						"hardware": "0",
						"manufacturer": "Schlage",
						"model": "BE469NX",
						"protocol": "zwave",
						"zwave.node": "33",
						"zwave.smartstart": "no"
					},
					"name": "2B",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "middle",
					"status": "idle",
					"subcategory": "",
					"type": "doorlock"
				},
				{
					"_id": "6165d418939a9313e502a8d1",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_25409_20548",
					"firmware": [
						{
							"id": "us.59.28950.0",
							"version": "113.22"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "3.42",
						"hardware": "0",
						"manufacturer": "Schlage",
						"model": "BE469NX",
						"protocol": "zwave",
						"zwave.node": "34",
						"zwave.smartstart": "no"
					},
					"name": "4B",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "middle",
					"status": "idle",
					"subcategory": "",
					"type": "doorlock"
				},
				{
					"_id": "6165d546939a9313e502a8ef",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_25409_20548",
					"firmware": [
						{
							"id": "us.59.28950.0",
							"version": "113.22"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "3.42",
						"hardware": "0",
						"manufacturer": "Schlage",
						"model": "BE469NX",
						"protocol": "zwave",
						"zwave.node": "35",
						"zwave.smartstart": "no"
					},
					"name": "3A",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "middle",
					"status": "idle",
					"subcategory": "",
					"type": "doorlock"
				},
				{
					"_id": "6165d5f7939a9313e502a90d",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_25409_20548",
					"firmware": [
						{
							"id": "us.59.28950.0",
							"version": "113.22"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "3.42",
						"hardware": "0",
						"manufacturer": "Schlage",
						"model": "BE469NX",
						"protocol": "zwave",
						"zwave.node": "36",
						"zwave.smartstart": "no"
					},
					"name": "3B",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "middle",
					"status": "idle",
					"subcategory": "",
					"type": "doorlock"
				},
				{
					"_id": "6165d7b3939a9313e502a92b",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_25409_20548",
					"firmware": [
						{
							"id": "us.59.28950.0",
							"version": "113.22"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "3.42",
						"hardware": "0",
						"manufacturer": "Schlage",
						"model": "BE469NX",
						"protocol": "zwave",
						"zwave.node": "37",
						"zwave.smartstart": "no"
					},
					"name": "5B",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "middle",
					"status": "idle",
					"subcategory": "",
					"type": "doorlock"
				},
				{
					"_id": "6165d8d8939a9313e502a949",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_1_1128",
					"firmware": [
						{
							"id": "us.59.18064.0",
							"version": "3.3"
						},
						{
							"id": "us.59.18065.1",
							"version": "11.0"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "6.3",
						"hardware": "3",
						"manufacturer": "Schlage",
						"model": "Unknown",
						"protocol": "zwave",
						"zwave.node": "38",
						"zwave.smartstart": "no"
					},
					"name": "5A",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "high",
					"status": "idle",
					"subcategory": "",
					"type": "device"
				},
				{
					"_id": "6165da82939a9313e502a95d",
					"batteryPowered": true,
					"category": "door_lock",
					"deviceTypeId": "59_25409_20548",
					"firmware": [
						{
							"id": "us.59.28950.0",
							"version": "113.22"
						}
					],
					"gatewayId": "606704e8939a9315eb0782f3",
					"info": {
						"firmware.stack": "3.42",
						"hardware": "0",
						"manufacturer": "Schlage",
						"model": "BE469NX",
						"protocol": "zwave",
						"zwave.node": "39",
						"zwave.smartstart": "no"
					},
					"name": "6B",
					"parentDeviceId": "",
					"persistent": false,
					"reachable": false,
					"ready": true,
					"roomId": "",
					"security": "middle",
					"status": "idle",
					"subcategory": "",
					"type": "doorlock"
				}
			]
		},
		"sender": {
			"conn_id": "2aab83bc-c2c8-4dd8-9a1e-ff7858753a4a",
			"type": "ui"
		}
	}
*/
