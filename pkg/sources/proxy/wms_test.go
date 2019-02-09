package proxy_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"AGES/pkg/sources/proxy"
)

func TestNewWMS(t *testing.T) {
	s, err := proxy.NewWMS("", "JPEG", 30*time.Second)
	require.NotNil(t, s)
	require.NoError(t, err)
}
