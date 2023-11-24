package shared

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Unit struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	PropertyID        uuid.UUID `json:"propertyId"`
	RemotePropertyURL string    `json:"remotePropertyUrl"`
	UpdatedBy         string    `json:"updatedBy"`
}

func (u *Unit) GetRemotePropertyID() int {
	// RemotePropertyURL is of the form:
	// https://dashboard.hostaway.com/listing/211374
	if u.RemotePropertyURL == "" {
		return -1
	}
	split := strings.Split(u.RemotePropertyURL, "/")
	id := split[len(split)-1]
	intID, err := strconv.Atoi(id)
	if err != nil {
		return -1
	}
	return intID
}
