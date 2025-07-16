package extensions

import (
	"github.com/florianl/go-tc"
	"gorm.io/gorm"
)

type Extensions struct {
	Database *gorm.DB
	TC       *TC
}

type TC struct {
	TC_  *tc.Tc
	Root *tc.Object
}
