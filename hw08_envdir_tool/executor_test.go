package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	e, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	code := RunCmd([]string{`$HOME/.bin`, "-c", "echo arg=1"}, e)
	require.Equal(t, 0, code, "Exit code should be 0")
}
