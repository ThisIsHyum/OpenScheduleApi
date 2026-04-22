package weekdays

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseWeekday(t *testing.T) {
	w, exists := ParseWeekday("monday")
	assert.True(t, exists)
	assert.Equal(t, time.Monday, w)

	w, exists = ParseWeekday("tuesday")
	assert.True(t, exists)
	assert.Equal(t, time.Tuesday, w)

	w, exists = ParseWeekday("WednesDay")
	assert.True(t, exists)
	assert.Equal(t, time.Wednesday, w)

	w, exists = ParseWeekday("WDay")
	assert.False(t, exists)
	assert.NotEqual(t, time.Wednesday, w)
}
