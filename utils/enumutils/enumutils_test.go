package enumutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalText(t *testing.T) {
	eb := NewEnumBase("TEST", map[int]string{
		1: "ONE",
		2: "TWO",
	})
	val, err := eb.UnmarshalText([]byte("ONE"))
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = eb.UnmarshalText([]byte("UNKNOWN"))
	assert.Equal(t, 0, val)
	assert.Equal(t, "unknown TEST: UNKNOWN", err.Error())
}
