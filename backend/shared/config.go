package shared

import (
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
		return fmt.Errorf("error loading .env file from %s", dir)
	}
	return nil
}

func GetConfigUnsafe(name string) string {
	v, err := GetConfig(name)
	if err != nil {
		fmt.Printf("error getting config: \"%s\"\n", err.Error())
	}
	return v
}

func GetConfig(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("can't find config for \"%s\"", name)
	}
	return val, nil
}
