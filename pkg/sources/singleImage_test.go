package sources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSingleImage(t *testing.T) {
	s, err := SingleImage("")
	require.Error(t, err)
}
