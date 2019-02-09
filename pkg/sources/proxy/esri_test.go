package proxy_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"AGES/pkg/sources/proxy"
)

func TestNewESRI(t *testing.T) {
	s, err := proxy.NewESRI("", "JPEG", 30*time.Second)
	require.NotNil(t, s)
	require.NoError(t, err)
}
