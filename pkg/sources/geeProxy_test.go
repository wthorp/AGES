package sources

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewGEEProxy(t *testing.T) {
	s, err := NewGEEProxy("", 30*time.Seconds)
	require.Error(t, err)
}
