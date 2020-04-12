package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	assert.Equal(t, 0, Abs(0))
	assert.Equal(t, 0.0, Abs(0.0))
	assert.Equal(t, 1, Abs(-1))
	assert.Equal(t, 1.1, Abs(-1.1))
}
func TestMin(t *testing.T) {
	assert.EqualValues(t, -1, Min(0, 0.1, 0, -1.0, -1, 1000.87896))
}
func TestMax(t *testing.T) {
	assert.EqualValues(t, 1000.87896, Max(0, 0.1, 0, -1, 1000.87896))
}
