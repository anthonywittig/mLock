package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func buildLambdaBinary(srcDirectory string, buildDirectory string) error {
	cmd := exec.Command("go", "build", "-o", buildDirectory, fmt.Sprintf("%s/main.go", srcDirectory))
	cmd.Dir = srcDirectory
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command, %v", err)
	}
	return nil
}

func cpIfExists(src string, dest string) error {
	if _, err := os.Stat(src); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("error checking if file exists, %v", err)
		} else {
			return nil
		}
	}

	cmd := exec.Command("cp", src, dest)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command, %v", err)
	}

	return nil
}

func createLambdaZip(buildDirectory string) error {
	args := []string{"-r", "function.zip"}

	files, err := ioutil.ReadDir(buildDirectory)
	if err != nil {
		return fmt.Errorf("unable to read directory, %v", err)
	}
	for _, file := range files {
		args = append(args, file.Name())
	}

	cmd := exec.Command("zip", args...)
	cmd.Dir = buildDirectory
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command, %v", err)
	}
	return nil
}

func mkDir(directory string) error {
	if _, err := os.Stat(directory); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("error checking if directory exists, %v", err)
		} else {
			err := os.Mkdir(directory, os.ModePerm)
			if err != nil {
				return fmt.Errorf("unable to create directory, %v", err)
			}
		}
	}
	return nil
}

func rmDir(directory string) error {
	if _, err := os.Stat(directory); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("error checking if directory exists, %v", err)
		} else {
			// Directory does not exist, nothing to do.
			return nil
		}
	}
	if err := os.RemoveAll(directory); err != nil {
		return fmt.Errorf("unable to remove directory, %v", err)
	}
	return nil
}
