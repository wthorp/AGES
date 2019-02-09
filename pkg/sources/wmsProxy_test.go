package sources

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewWMSProxy(t *testing.T) {
	s, err := NewWMSProxy("", 30*time.Second)
	require.Nil(t, s)
	require.Error(t, err)
}
