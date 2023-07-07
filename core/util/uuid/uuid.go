package uuid

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

func NewUuid() string {
	id, _ := uuid.NewV4()
	return strings.ReplaceAll(id.String(), "-", "")
}
