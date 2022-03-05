package ezlo

import (
	"context"
	"fmt"
	"log"
	mshared "mlock/shared"
	"net/url"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectionPool struct {
	connectionByControllerID map[string]*websocket.Conn
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{
		connectionByControllerID: map[string]*websocket.Conn{},
	}
}

func (cp *ConnectionPool) Close() {
	for _, c := range cp.connectionByControllerID {
		if err := c.Close(); err != nil {
			log.Printf("error while closing connection: %s", err.Error())
		}
	}
}

func (cp *ConnectionPool) GetConnection(ctx context.Context, controllerID string) (*websocket.Conn, error) {
	if ws, ok := cp.connectionByControllerID[controllerID]; ok {
		return ws, nil
	}

	ws, err := cp.connect(ctx, controllerID)
	if err != nil {
		return nil, fmt.Errorf("error connecting: %s", err.Error())
	}

	cp.connectionByControllerID[controllerID] = ws
	return ws, nil
}

func (cp *ConnectionPool) connect(ctx context.Context, controllerID string) (*websocket.Conn, error) {
	username, err := mshared.GetConfig("EZLO_USERNAME")
	if err != nil {
		return nil, fmt.Errorf("error getting username: %s", err.Error())
	}

	password, err := mshared.GetConfig("EZLO_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("error getting password: %s", err.Error())
	}

	authData, err := authenticate(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("error authenticating: %s", err.Error())
	}

	device, err := getDeviceByID(ctx, authData, controllerID)
	if err != nil {
		return nil, fmt.Errorf("error getting devices: %s", err.Error())
	}

	deviceResponse, err := getDevice(ctx, authData, device)
	if err != nil {
		return nil, fmt.Errorf("error getting device: %s", err.Error())
	}

	r, err := regexp.Compile("wss://(.*):443")
	if err != nil {
		return nil, fmt.Errorf("error compiling regex: %s", err.Error())
	}

	wsURLs := r.FindStringSubmatch(deviceResponse.ServerRelay)
	if c := len(wsURLs); c != 2 {
		return nil, fmt.Errorf("unexpected match count: %d", c)
	}

	u := url.URL{Scheme: "wss", Host: wsURLs[1], Path: ""}

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err.Error())
	}

	if err := ws.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		return nil, fmt.Errorf("setting read deadline: %s", err.Error())
	}

	if err := wsLogIn(ws, authData.Response); err != nil {
		return nil, fmt.Errorf("login: %s", err.Error())
	}

	if err := wsRegisterHub(ws, device.PKDevice); err != nil {
		return nil, fmt.Errorf("register: %s", err.Error())
	}

	return ws, nil
}
