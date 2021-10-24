package ezlo

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Test_HW(t *testing.T) {
	assert.Nil(t, loadConfig())

	resp, err := X(context.Background(), os.Getenv("USERNAME"), os.Getenv("PASSWORD"), os.Getenv("HUB_SERIAL_NUMBER"))
	assert.Nil(t, err)

	assert.Equal(t, "ahh", resp)
}

func loadConfig() error {
	if err := godotenv.Load(".env.test"); err != nil {
		return fmt.Errorf("error loading .env file: %s", err.Error())
	}
	return nil
}
