package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	e, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	fmt.Println(os.LookupEnv("SHELL"))
	code := RunCmd([]string{`/usr/local/bin/bash.exe`, "-c", "echo arg=1"}, e)
	require.Equal(t, 0, code, "Exit code should be 0")
}
