package ezlo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type AuthIdentity struct {
	// There are lot of other fields that we get back.
	PKAccount int `json:"PK_Account"`
	//"Expires":1635128723,
	//"Generated":1635042323,
	//"PermissionsEnabled":[
	//  1,
	//  2,
	//  {
	//	"PK":1653,
	//	"Arguments":"[9,10,11,12,13]"
	//  }
	//],
	//"PermissionsDisabled":[],
	//"Version":2,
	//"PK_AccountType":5,
	//"PK_AccountChild":0,
	//"PK_Account_Parent":2375,
	//"PK_Account_Parent2":1,
	//"PK_Account_Parent3":0,
	//"PK_Account_Parent4":0,
	//"PK_App":0,
	//"PK_Oem":1,
	//"PK_Oem_User":"",
	//"PK_PermissionRole":10,
	//"PK_User":2928592,
	//"PK_Server_Auth":1,
	//"PK_Server_Account":5,
	//"PK_Server_Event":53,
	//"Server_Auth":"vera-us-oem-autha11.mios.com",
	//"Seq":6359588,
	//"Username":"anthony.wittig"
}

type AuthResponse struct {
	Identity          string `json:"Identity"`
	IdentitySignature string `json:"IdentitySignature"`
	ServerEvent       string `json:"Server_Event"`
	ServerEventAlt    string `json:"Server_Event_Alt"`
	ServerAccount     string `json:"Server_Account"`
	ServerAccountAlt  string `json:"Server_Account_Alt"`
}

type AuthData struct {
	Identity AuthIdentity
	Response AuthResponse
}

type Device struct {
	Blocked         int    `json:"Blocked"`
	DeviceAssigned  string `json:"DeviceAssigned"`
	MACAddress      string `json:"MacAddress"`
	PKDevice        string `json:"PK_Device"`
	PKDeviceSubType string `json:"PK_DeviceSubType"`
	PKDeviceType    string `json:"PK_DeviceType"`
	PKInstallation  string `json:"PK_Installation"`
	ServerDevice    string `json:"Server_Device"`
	ServerDeviceAlt string `json:"Server_Device_Alt"`
}

type DeviceResponse struct {
	//EngineStatus: "0"
	//FK_Branding: "1"
	//HasAlarmPanel: "0"
	//HasWifi: "0"
	//LinuxFirmware: 1
	//MacAddress: "..."
	//NMAControllerStatus: 1
	//NMAUuid: "..."
	//PK_Device: "..."
	ServerRelay string `json:"Server_Relay"`
	//UI: "4"
	//public_key_android: ""
	//public_key_ios: ""
}

type DevicesResponse struct {
	Devices []Device `json:"Devices"`
}

type WSDeviceListResponse struct {
	API   string `json:"api"`
	Error *struct {
		Code        int    `json:"code"`
		Data        string `json:"data"`
		Description string `json:"description"`
	} `json:"error"`
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

type WSResponse struct {
	Error *struct {
		Code        int    `json:"code"`
		Data        string `json:"data"`
		Description string `json:"description"`
	} `json:"error"`
	Method string `json:"method"`
	ID     string `json:"id"`
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

func authenticate(ctx context.Context, username string, password string) (AuthData, error) {
	// TODO: get a list of the endpoints and use a random one.
	// "vera-us-oem-account12.mios.com"
	url := fmt.Sprintf("https://vera-us-oem-account11.mios.com/autha/auth/username/%s", username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return AuthData{}, fmt.Errorf("error creating request: %s", err.Error())
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
		return AuthData{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return AuthData{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	body := AuthResponse{}
	if err := json.Unmarshal(respBody, &body); err != nil {
		return AuthData{}, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	identityString, err := base64.StdEncoding.DecodeString(body.Identity)
	if err != nil {
		return AuthData{}, fmt.Errorf("error decoding identity: %s", err.Error())
	}

	identity := AuthIdentity{}
	if err := json.Unmarshal(identityString, &identity); err != nil {
		return AuthData{}, fmt.Errorf("error decoding identity: %s", err.Error())
	}

	return AuthData{
		Identity: identity,
		Response: body,
	}, nil
}

func getDevice(ctx context.Context, authData AuthData, device Device) (DeviceResponse, error) {
	// TODO: we probably need to use the same domain that we used to auth.
	url := fmt.Sprintf("https://vera-us-oem-account11.mios.com/device/device/device/%s", device.PKDevice)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return DeviceResponse{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	req.Header.Set("mmsAuth", authData.Response.Identity)
	req.Header.Set("mmsAuthSig", authData.Response.IdentitySignature)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return DeviceResponse{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return DeviceResponse{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	body := DeviceResponse{}
	if err := json.Unmarshal(respBody, &body); err != nil {
		return DeviceResponse{}, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	return body, nil
}

func getDevices(ctx context.Context, authData AuthData) (Device, error) {
	// TODO: we probably need to use the same domain that we used to auth.
	url := fmt.Sprintf("https://vera-us-oem-account11.mios.com/account/account/account/%d/devices", authData.Identity.PKAccount)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Device{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	req.Header.Set("mmsAuth", authData.Response.Identity)
	req.Header.Set("mmsAuthSig", authData.Response.IdentitySignature)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return Device{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Device{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	body := DevicesResponse{}
	if err := json.Unmarshal(respBody, &body); err != nil {
		return Device{}, fmt.Errorf("error unmarshalling body: %s; url: %s; authData: %+v; error: %s", string(respBody), url, authData, err.Error())
	}

	// For now, we'll only work with a single device.
	if c := len(body.Devices); c != 1 {
		return Device{}, fmt.Errorf("unexpected device count: %d", c)
	}

	return body.Devices[0], nil
}

func WSLogIn(ws *websocket.Conn, authResponse AuthResponse) error {
	id := fmt.Sprintf("loginUserMios.%s", uuid.New())
	err := sendCommand(
		ws,
		id,
		WSLogInRequest{
			Method: "loginUserMios",
			ID:     id,
			Params: WSLogInRequestParams{
				MMSAuth:    authResponse.Identity,
				MMSAuthSig: authResponse.IdentitySignature,
			},
		},
		&struct{}{},
	)
	if err != nil {
		return fmt.Errorf("error sending command: %s", err.Error())
	}

	return nil
}

func WSRegisterHub(ws *websocket.Conn, hubSerialNumber string) error {
	id := fmt.Sprintf("register.%s", uuid.New())
	err := sendCommand(
		ws,
		id,
		WSRegisterRequest{
			Method: "register",
			ID:     id,
			Params: WSRegisterRequestParams{
				Serial: hubSerialNumber,
			},
		},
		&struct{}{},
	)
	if err != nil {
		return fmt.Errorf("error sending command: %s", err.Error())
	}

	return nil
}

func WSDeviceList(ws *websocket.Conn) (WSDeviceListResponse, error) {
	id := fmt.Sprintf("hub.devices.list.%s", uuid.New())
	resp := WSDeviceListResponse{}
	err := sendCommand(
		ws,
		id,
		struct {
			Method string   `json:"method"`
			ID     string   `json:"id"`
			Params struct{} `json:"params"`
		}{
			Method: "hub.devices.list",
			ID:     id,
		},
		&resp,
	)
	if err != nil {
		return WSDeviceListResponse{}, fmt.Errorf("error sending command: %s", err.Error())
	}

	return resp, nil
}

func X(ctx context.Context, username string, password string, hubSerialNumber string) (string, error) {
	authData, err := authenticate(ctx, username, password)
	if err != nil {
		return "", fmt.Errorf("error authenticating: %s", err.Error())
	}

	device, err := getDevices(ctx, authData)
	if err != nil {
		return "", fmt.Errorf("error getting devices: %s", err.Error())
	}

	deviceResponse, err := getDevice(ctx, authData, device)
	if err != nil {
		return "", fmt.Errorf("error getting device: %s", err.Error())
	}

	r, err := regexp.Compile("wss://(.*):443")
	if err != nil {
		return "", fmt.Errorf("error compiling regex: %s", err.Error())
	}

	wsURLs := r.FindStringSubmatch(deviceResponse.ServerRelay)
	if c := len(wsURLs); c != 2 {
		return "", fmt.Errorf("unexpected match count: %d", c)
	}

	u := url.URL{Scheme: "wss", Host: wsURLs[1], Path: ""}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("dial: %s", err.Error())
	}
	defer ws.Close()

	if err := ws.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return "", fmt.Errorf("setting read deadline: %s", err.Error())
	}

	if err := WSLogIn(ws, authData.Response); err != nil {
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

func sendCommand(ws *websocket.Conn, id string, request interface{}, outResponse interface{}) error {
	jsonReq, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal: %s", err.Error())
	}

	fmt.Printf("sending: %s\n", string(jsonReq))

	if err := ws.WriteMessage(websocket.TextMessage, jsonReq); err != nil {
		return fmt.Errorf("write: %s", err.Error())
	}

	for {
		_, jsonResp, err := ws.ReadMessage()
		if err != nil {
			return fmt.Errorf("read: %s", err.Error())
		}

		resp := WSResponse{}
		if err := json.Unmarshal(jsonResp, &resp); err != nil {
			return fmt.Errorf("unmarshal: %s", err.Error())
		}

		if resp.ID == "ui_broadcast" {
			// We don't care about these, try to get the next message.
			continue
		}

		if resp.ID != id {
			return fmt.Errorf("unexpected response ID: %s, expected: %s", resp.ID, id)
		}

		if resp.Error != nil {
			return fmt.Errorf("error in WS response: %+v", resp.Error)
		}

		if err := json.Unmarshal(jsonResp, &outResponse); err != nil {
			return fmt.Errorf("unmarshal: %s", err.Error())
		}

		fmt.Printf("received: %s\n", string(jsonResp))

		return nil
	}
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
