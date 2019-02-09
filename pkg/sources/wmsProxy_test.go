package sources

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewWMSProxy(t *testing.T) {
	s, err := NewWMSProxy("", 30*time.Second)
	require.NotNil(t, s)
	require.NoError(t, err)
}
