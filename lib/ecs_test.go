package lib

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_max(t *testing.T) {
	require.Equal(t, max(2, 1), int64(2))
	require.Equal(t, max(1, 2), int64(2))
	require.Equal(t, max(-2, 1), int64(1))
	require.Equal(t, max(1, -2), int64(1))
	require.Equal(t, max(math.MaxInt, -2), int64(math.MaxInt))
	require.Equal(t, max(math.MinInt, -2), int64(-2))
}

func Test_min(t *testing.T) {
	require.Equal(t, min(2, 1), int64(1))
	require.Equal(t, min(1, 2), int64(1))
	require.Equal(t, min(-2, 1), int64(-2))
	require.Equal(t, min(1, -2), int64(-2))
	require.Equal(t, min(math.MaxInt, -2), int64(-2))
	require.Equal(t, min(math.MinInt, -2), int64(math.MinInt))
}
