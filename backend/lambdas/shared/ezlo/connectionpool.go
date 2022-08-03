package ezlo

import (
	"context"
	"fmt"
	"log"
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

	// Connecting to the ws should be quick.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	ws, err := cp.connect(ctx, controllerID)
	if err != nil {
		return nil, fmt.Errorf("error connecting: %s", err.Error())
	}

	cp.connectionByControllerID[controllerID] = ws
	return ws, nil
}

func (cp *ConnectionPool) connect(ctx context.Context, controllerID string) (*websocket.Conn, error) {
	ad, err := getAuthData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting auth data: %s", err.Error())
	}

	device, err := getDeviceByID(ctx, ad, controllerID)
	if err != nil {
		return nil, fmt.Errorf("error getting device by ID (%s): %s", controllerID, err.Error())
	}

	deviceResponse, err := getDevice(ctx, ad, device)
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

	if err := wsLogIn(ws, ad.Response); err != nil {
		return nil, fmt.Errorf("login: %s", err.Error())
	}

	if err := wsRegisterHub(ws, device.PKDevice); err != nil {
		return nil, fmt.Errorf("register: %s", err.Error())
	}

	return ws, nil
}
