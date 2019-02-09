package sources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFractal(t *testing.T) {
	s, err := NewFractal()
	require.NotNil(t, s)
	require.NoError(t, err)
}
