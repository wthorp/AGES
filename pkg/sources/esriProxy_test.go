package sources

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewESRIProxy(t *testing.T) {
	s, err := NewESRIProxy("", 30*time.Seconds)
	require.Error(t, err)
}
