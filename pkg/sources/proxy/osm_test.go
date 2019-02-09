package proxy_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"AGES/pkg/sources/proxy"
)

func TestNewOSM(t *testing.T) {
	s, err := proxy.NewOSM("", 30*time.Second)
	require.NotNil(t, s)
	require.NoError(t, err)
}
