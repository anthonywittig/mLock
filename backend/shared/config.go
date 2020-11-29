package shared

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	dir := path.Dir(ex)
	if err := godotenv.Load(dir + "/.env"); err != nil {
		return errors.New("Error loading .env file")
	}
	return nil
}

func GetConfig(name string) string {
	val := os.Getenv(name)
	if val == "" {
		fmt.Printf("can't find config for \"%s\"\n", name)
	}
	return val
}
