package saOrm

import (
	"gorm.io/gorm"
	"strings"
)

func (o *DB) IsError(err error) bool {
	if err == nil || err == gorm.ErrRecordNotFound || strings.Contains(err.Error(), "no row") {
		return false
	}
	return true
}
