package sources

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewWMSProxy(t *testing.T) {
	s, err := NewWMSProxy("", 30*time.Seconds)
	require.Error(t, err)
}
