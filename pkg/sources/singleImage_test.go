package sources

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSingleImage(t *testing.T) {
	s, err := NewSingleImage("")
	require.Nil(t, s)
	require.Error(t, err)

	s, err = NewSingleImage("pipe.jpg")
	require.Nil(t, s)
	require.Error(t, err)
}
