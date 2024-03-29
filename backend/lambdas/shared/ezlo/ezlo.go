package ezlo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mlock/lambdas/shared"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type authIdentity struct {
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

type authResponse struct {
	Identity          string `json:"Identity"`
	IdentitySignature string `json:"IdentitySignature"`
	ServerEvent       string `json:"Server_Event"`
	ServerEventAlt    string `json:"Server_Event_Alt"`
	ServerAccount     string `json:"Server_Account"`
	ServerAccountAlt  string `json:"Server_Account_Alt"`
}

type authData struct {
	Identity authIdentity
	Response authResponse
}

type device struct {
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

type deviceResponse struct {
	//EngineStatus: "0"
	//FK_Branding: "1"
	//HasAlarmPanel: "0"
	//HasWifi: "0"
	//LinuxFirmware: 1
	//MacAddress: "..."
	NMAControllerStatus int `json:"NMAControllerStatus"`
	//NMAUuid: "..."
	PKDevice    string `json:"PK_Device"`
	ServerRelay string `json:"Server_Relay"`
	//UI: "4"
	//public_key_android: ""
	//public_key_ios: ""
}

type devicesResponse struct {
	Devices []device `json:"Devices"`
}

type wsDeviceListResponse struct {
	API   string `json:"api"`
	Error *struct {
		Code        int    `json:"code"`
		Data        string `json:"data"`
		Description string `json:"description"`
	} `json:"error"`
	ID     string                     `json:"id"`
	Method string                     `json:"method"`
	Result wsDeviceListResponseResult `json:"result"`
}
type wsDeviceListResponseResult struct {
	Devices []wsDeviceListResponseResultDevice `json:"devices"`
	Sender  wsDeviceListResponseResultSender   `json:"sender"`
}

type wsDeviceListResponseResultDevice struct {
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

type wsDeviceListResponseResultSender struct {
	ConnID string `json:"conn_id"`
	Type   string `json:"type"`
}

type wsItemsListResponse struct {
	API   string `json:"api"`
	Error *struct {
		Code        int    `json:"code"`
		Data        string `json:"data"`
		Description string `json:"description"`
	} `json:"error"`
	ID     string                    `json:"id"`
	Method string                    `json:"method"`
	Result wsItemsListResponseResult `json:"result"`
}

type wsItemsListResponseResult struct {
	Items []wsItem `json:"items"`
}

type wsLogInRequest struct {
	Method string               `json:"method"`
	ID     string               `json:"id"`
	Params wsLogInRequestParams `json:"params"`
}

type wsLogInRequestParams struct {
	MMSAuth    string `json:"MMSAuth"`
	MMSAuthSig string `json:"MMSAuthSig"`
}

type wsResponse struct {
	Error *struct {
		Code        int    `json:"code"`
		Data        string `json:"data"`
		Description string `json:"description"`
	} `json:"error"`
	Method string `json:"method"`
	ID     string `json:"id"`
}

type wsRegisterRequest struct {
	Method  string                  `json:"method"`
	ID      string                  `json:"id"`
	JSONRPC string                  `json:"jsonrpc"`
	Params  wsRegisterRequestParams `json:"params"`
}

type wsRegisterRequestParams struct {
	Serial string `json:"serial"`
}

const ENHANCED_LOGGING = false

func getRawDevices(ws *websocket.Conn) ([]shared.RawDevice, error) {
	deviceListResp, err := wsDeviceList(ws)
	if err != nil {
		return []shared.RawDevice{}, fmt.Errorf("error getting device list: %s", err.Error())
	}

	itemsByDevice, err := wsItemsByDevice(ws)
	if err != nil {
		return []shared.RawDevice{}, fmt.Errorf("error getting items by device: %s", err.Error())
	}

	result := []shared.RawDevice{}
	for _, d := range deviceListResp.Result.Devices {
		status := shared.DeviceStatusOffline
		deviceInitialized := d.DeviceTypeID != "0_0_0"
		if d.Reachable && deviceInitialized && d.Status != "broken" {
			status = shared.DeviceStatusOnline
		}

		rd := shared.RawDevice{
			Battery: shared.RawDeviceBattery{
				BatteryPowered: d.Reachable && d.BatteryPowered,
			},
			Category:     d.Category,
			DeviceTypeID: d.DeviceTypeID,
			ID:           d.ID,
			Name:         d.Name,
			Status:       status,
		}

		for _, item := range itemsByDevice[d.ID] {
			if item.Name == "battery" {
				if rd.Battery.Level, err = item.getBatteryLevel(); err != nil {
					return []shared.RawDevice{}, fmt.Errorf("error getting battery: %s", err.Error())
				}
			} else if item.Name == "user_codes" {
				if rd.LockCodes, err = item.getLockCodes(); err != nil {
					return []shared.RawDevice{}, fmt.Errorf("error getting lock codes: %s", err.Error())
				}
			}
		}

		result = append(result, rd)
	}

	return result, nil
}

func authenticate(ctx context.Context, username string, password string) (authData, error) {
	// TODO: get a list of the endpoints and use a random one?
	// "vera-us-oem-account12.mios.com"
	url := fmt.Sprintf("https://vera-us-oem-account11.mios.com/autha/auth/username/%s", username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return authData{}, fmt.Errorf("error creating request: %s", err.Error())
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

	resp, err := client.Do(req)
	if err != nil {
		return authData{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return authData{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	body := authResponse{}
	if err := json.Unmarshal(respBody, &body); err != nil {
		return authData{}, fmt.Errorf("error unmarshalling body: %s", err.Error())
	}

	identityString, err := base64.StdEncoding.DecodeString(body.Identity)
	if err != nil {
		return authData{}, fmt.Errorf("error decoding identity: %s", err.Error())
	}

	identity := authIdentity{}
	if err := json.Unmarshal(identityString, &identity); err != nil {
		return authData{}, fmt.Errorf("error decoding identity: %s", err.Error())
	}

	return authData{
		Identity: identity,
		Response: body,
	}, nil
}

func getDevice(ctx context.Context, ad authData, d device) (deviceResponse, error) {
	// TODO: we probably need to use the same domain that we used to auth.
	url := fmt.Sprintf("https://vera-us-oem-account11.mios.com/device/device/device/%s", d.PKDevice)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return deviceResponse{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	req.Header.Set("mmsAuth", ad.Response.Identity)
	req.Header.Set("mmsAuthSig", ad.Response.IdentitySignature)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return deviceResponse{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return deviceResponse{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	body := deviceResponse{}
	if err := json.Unmarshal(respBody, &body); err != nil {
		return deviceResponse{}, fmt.Errorf("error unmarshalling body: %s; error: %s", string(respBody), err.Error())
	}

	return body, nil
}

func getDeviceByID(ctx context.Context, ad authData, controllerID string) (device, error) {
	devices, err := getDevices(ctx, ad)
	if err != nil {
		return device{}, fmt.Errorf("error getting devices: %s", err.Error())
	}

	for _, d := range devices {
		if d.PKDevice == controllerID {
			return d, nil
		}
	}

	return device{}, fmt.Errorf("unable to find device: %+v", devices)
}
func getDevices(ctx context.Context, ad authData) ([]device, error) {
	// TODO: we probably need to use the same domain that we used to auth.
	url := fmt.Sprintf("https://vera-us-oem-account11.mios.com/account/account/account/%d/devices", ad.Identity.PKAccount)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []device{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	req.Header.Set("mmsAuth", ad.Response.Identity)
	req.Header.Set("mmsAuthSig", ad.Response.IdentitySignature)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return []device{}, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []device{}, fmt.Errorf("error reading body: %s", err.Error())
	}

	body := devicesResponse{}
	if err := json.Unmarshal(respBody, &body); err != nil {
		return []device{}, fmt.Errorf("error unmarshalling body: %s; url: %s; authData: %+v; error: %s", string(respBody), url, ad, err.Error())
	}

	return body.Devices, nil
}

func wsSendCommand(ws *websocket.Conn, id string, request interface{}, outResponse interface{}) error {
	jsonReq, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal: %s", err.Error())
	}

	if ENHANCED_LOGGING {
		fmt.Printf("sending: %s\n", string(jsonReq))
	}

	if err := ws.WriteMessage(websocket.TextMessage, jsonReq); err != nil {
		return fmt.Errorf("write: %s", err.Error())
	}

	for {
		_, jsonResp, err := ws.ReadMessage()
		if err != nil {
			return fmt.Errorf("read: %s", err.Error())
		}

		resp := wsResponse{}
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

		if ENHANCED_LOGGING {
			fmt.Printf("received: %s\n", string(jsonResp))
		}

		return nil
	}
}

func wsLogIn(ws *websocket.Conn, ar authResponse) error {
	id := fmt.Sprintf("loginUserMios.%s", uuid.New())
	err := wsSendCommand(
		ws,
		id,
		wsLogInRequest{
			Method: "loginUserMios",
			ID:     id,
			Params: wsLogInRequestParams{
				MMSAuth:    ar.Identity,
				MMSAuthSig: ar.IdentitySignature,
			},
		},
		&struct{}{},
	)
	if err != nil {
		return fmt.Errorf("error sending command: %s", err.Error())
	}

	return nil
}

func wsRegisterHub(ws *websocket.Conn, hubSerialNumber string) error {
	id := fmt.Sprintf("register.%s", uuid.New())
	err := wsSendCommand(
		ws,
		id,
		wsRegisterRequest{
			Method: "register",
			ID:     id,
			Params: wsRegisterRequestParams{
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

func wsDeviceList(ws *websocket.Conn) (wsDeviceListResponse, error) {
	id := fmt.Sprintf("hub.devices.list.%s", uuid.New())
	resp := wsDeviceListResponse{}
	err := wsSendCommand(
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
		return wsDeviceListResponse{}, fmt.Errorf("error sending command: %s", err.Error())
	}

	return resp, nil
}

func wsItemsByDevice(ws *websocket.Conn) (map[string][]wsItem, error) {
	itemsByDevice := map[string][]wsItem{}

	id := fmt.Sprintf("hub.items.list.%s", uuid.New())
	resp := wsItemsListResponse{}
	err := wsSendCommand(
		ws,
		id,
		struct {
			Method string   `json:"method"`
			ID     string   `json:"id"`
			Params struct{} `json:"params"`
		}{
			Method: "hub.items.list",
			ID:     id,
		},
		&resp,
	)
	if err != nil {
		return itemsByDevice, fmt.Errorf("error sending command: %s", err.Error())
	}

	for _, item := range resp.Result.Items {
		itemsByDevice[item.DeviceID] = append(itemsByDevice[item.DeviceID], item)
	}

	return itemsByDevice, nil
}
