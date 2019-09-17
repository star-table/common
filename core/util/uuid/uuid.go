package uuid

import (
	uuid "github.com/satori/go.uuid"
	"strings"
)

func NewUuid() string {
	id := uuid.NewV4().String()

	return strings.ReplaceAll(id, "-", "")
}
