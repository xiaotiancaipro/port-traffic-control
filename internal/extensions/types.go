package extensions

import (
	"github.com/florianl/go-tc"
	"gorm.io/gorm"
)

type Extensions struct {
	Database *gorm.DB
	TC       *tc.Tc
}
