package shared

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevice_Marshaing_SimpleCases(t *testing.T) {
	d := Device{}

	j, err := json.Marshal(d)
	assert.Nil(t, err)
	// Lame...
	assert.Contains(
		t,
		string(j),
		"\"managedLockCodes\":null",
	)

	d.SortManagedLockCodes()
	j, err = json.Marshal(d)
	assert.Nil(t, err)
	assert.Contains(
		t,
		string(j),
		"\"managedLockCodes\":[]",
	)
}
