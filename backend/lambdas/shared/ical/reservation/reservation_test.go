package reservation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_getCheckinTimestamp_nonMidnightUTC(t *testing.T) {
	in, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.Nil(t, err)

	_, err = getCheckinTimestamp(in)
	assert.NotNil(t, err)
}

func Test_getCheckinTimestamp_happyPath(t *testing.T) {
	in, err := time.Parse(time.RFC3339, "2006-01-02T00:00:00Z")
	assert.Nil(t, err)

	actual, err := getCheckinTimestamp(in)
	assert.Nil(t, err)

	expected, err := time.Parse(time.RFC3339, "2006-01-02T16:00:00-07:00")
	assert.Nil(t, err)

	assert.Equal(t, expected.String(), actual.String())
}

func Test_getCheckoutTimestamp_nonMidnightUTC(t *testing.T) {
	in, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	assert.Nil(t, err)

	_, err = getCheckoutTimestamp(in)
	assert.NotNil(t, err)
}

func Test_getCheckoutTimestamp_happyPath(t *testing.T) {
	in, err := time.Parse(time.RFC3339, "2006-01-02T00:00:00Z")
	assert.Nil(t, err)

	actual, err := getCheckoutTimestamp(in)
	assert.Nil(t, err)

	expected, err := time.Parse(time.RFC3339, "2006-01-02T11:00:00-07:00")
	assert.Nil(t, err)

	assert.Equal(t, expected.String(), actual.String())
}
