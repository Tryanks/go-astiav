package astiav_test

import (
	"testing"

	"github.com/asticode/go-astiav"
	"github.com/stretchr/testify/require"
)

func TestIntReadWrite(t *testing.T) {
	is := []uint8{1, 2, 3, 4, 5, 6, 7, 8}
	require.Equal(t, uint32(0), astiav.RL32([]byte{}))
	require.Equal(t, uint32(0x4030201), astiav.RL32(is))
	require.Equal(t, uint32(0), astiav.RL32WithOffset([]byte{}, 4))
	require.Equal(t, uint32(0x8070605), astiav.RL32WithOffset(is, 4))
}
