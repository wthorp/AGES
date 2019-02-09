package sources

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewESRIProxy(t *testing.T) {
	s, err := NewESRIProxy("", 30*time.Second)
	require.NotNil(t, s)
	require.NoError(t, err)
}
