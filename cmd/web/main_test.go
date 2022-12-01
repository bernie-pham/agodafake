package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	_, err := run()
	require.NoError(t, err)
}
