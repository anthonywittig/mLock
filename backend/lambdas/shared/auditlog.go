package shared

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID      uuid.UUID       `json:"id"`
	Entries []AuditLogEntry `json:"entries"`
}

type AuditLogEntry struct {
	CreatedAt time.Time `json:"createdAt"`
	Log       string    `json:"log"`
}
