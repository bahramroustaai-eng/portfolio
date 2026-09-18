package portfo

import (
	"portfolio/internal/apperr"
	"strings"
)

var (
	ErrInvalidItemTypeName = apperr.New(apperr.CodeInvalid, "item type name must not be empty")
	ErrItemTypeConflict    = apperr.New(apperr.CodeConflict, "item type already exists")
)

type ItemType struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func NewItemTypeName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrInvalidItemTypeName
	}
	return name, nil
}
