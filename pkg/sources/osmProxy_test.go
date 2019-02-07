package sources

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewOSMProxy(t *testing.T) {
	s, err := NewOSMProxy("", 30*time.Seconds)
	require.Error(t, err)
}
