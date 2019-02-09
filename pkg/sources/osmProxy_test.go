package sources

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewOSMProxy(t *testing.T) {
	s, err := NewOSMProxy("", 30*time.Second)
	require.NotNil(t, s)
	require.NoError(t, err)
}
