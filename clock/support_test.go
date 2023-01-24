package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartDate(t *testing.T) {
	d := time.Date(2021, 6, 22, 11, 33, 40, 1233, time.UTC)
	result := StartDate(d)
	assert.Equal(t, 2021, result.Year(), "Year")
	assert.Equal(t, time.Month(6), result.Month(), "Month")
	assert.Equal(t, 22, result.Day(), "Day")
	assert.Equal(t, 0, result.Hour(), "Hour")
	assert.Equal(t, 0, result.Minute(), "Minute")
	assert.Equal(t, 0, result.Second(), "Second")
	assert.Equal(t, 0, result.Nanosecond(), "Nanosecond")
}

func TestEndDate(t *testing.T) {
	d := time.Date(2021, 6, 22, 11, 33, 40, 1233, time.UTC)
	result := EndDate(d)
	assert.Equal(t, 2021, result.Year(), "Year")
	assert.Equal(t, time.Month(6), result.Month(), "Month")
	assert.Equal(t, 22, result.Day(), "Day")
	assert.Equal(t, 23, result.Hour(), "Hour")
	assert.Equal(t, 59, result.Minute(), "Minute")
	assert.Equal(t, 59, result.Second(), "Second")
	assert.Equal(t, 999999999, result.Nanosecond(), "Nanosecond")
}
