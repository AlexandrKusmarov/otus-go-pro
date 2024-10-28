package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	e, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	env, _ := os.LookupEnv("SHELL")
	fmt.Println("SHELL: " + env)
	fmt.Println("ALL ENV" + strings.Join(os.Environ(), ","))
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	code := RunCmd([]string{`/usr/local/bin/bash.exe`, "-c", "echo arg=1"}, e)
	require.Equal(t, 0, code, "Exit code should be 0")
}
