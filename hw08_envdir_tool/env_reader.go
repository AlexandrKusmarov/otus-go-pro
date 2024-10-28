package main

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	listFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	env := make(Environment)
	for _, file := range listFiles {
		if file.IsDir() {
			continue
		}
		envVal := EnvValue{
			Value:      "",
			NeedRemove: true,
		}
		if strings.Contains(file.Name(), "=") {
			return nil, errors.New("error: file name contains '='")
		}
		fData, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		envVal.NeedRemove = len(fData) == 0
		scanner := bufio.NewScanner(strings.NewReader(string(fData)))
		if scanner.Scan() {
			envVal.Value = scanner.Text()
			envVal.Value = strings.TrimRight(envVal.Value, " \t")
			envVal.Value = strings.ReplaceAll(envVal.Value, "\x00", "\n")
		}
		env[file.Name()] = envVal
	}

	return env, nil
}
