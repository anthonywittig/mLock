package ezlo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func HW(ctx context.Context, username string, password string) (string, error) {
	url := fmt.Sprintf("https://vera-us-oem-account11.mios.com/autha/auth/username/%s", username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %s", err.Error())
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
		return "", fmt.Errorf("error doing request: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading body: %s", err.Error())
	}

	return string(respBody), nil
}
