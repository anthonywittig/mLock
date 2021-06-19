package hab

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mlock/shared"
	"mlock/shared/protos/messaging"
	"net/http"
	"time"
)

func ProcessCommand(ctx context.Context, in *messaging.HabCommand) (*messaging.OnPremHabCommandResponse, error) {
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	endpoint, err := shared.GetConfig("OPEN_HAB_ENDPOINT")
	if err != nil {
		return nil, fmt.Errorf("error getting endpoint: %s", err.Error())
	}
	url := endpoint + in.Request.Path

	req, err := http.NewRequestWithContext(ctx, in.Request.Method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err.Error())
	}

	auth, err := shared.GetConfig("OPEN_HAB_AUTH")
	if err != nil {
		return nil, fmt.Errorf("error getting auth: %s", err.Error())
	}
	encAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+encAuth)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %s", err.Error())
	}

	return &messaging.OnPremHabCommandResponse{
		Description: "HAB command response",
		HabCommand:  in,
		Response:    body,
	}, nil
}
