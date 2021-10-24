package ezlo

import (
	"context"
	"fmt"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Test_GetDevices(t *testing.T) {
	assert.Nil(t, loadConfig())

	resp, err := GetDevices(context.Background())
	assert.Nil(t, err)

	assert.Equal(t, "ahh", resp)
}

func loadConfig() error {
	if err := godotenv.Load(".env.test"); err != nil {
		return fmt.Errorf("error loading .env file: %s", err.Error())
	}
	return nil
}
