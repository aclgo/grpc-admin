package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFromFile(t *testing.T) {

	cfg := Config{}

	err := cfg.LoadFromFile(".")
	require.NoError(t, err)
	require.NotNil(t, cfg)
}
