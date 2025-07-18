package extensions

import (
	"github.com/florianl/go-tc"
	"gorm.io/gorm"
	"net"
)

type Extensions struct {
	Database *gorm.DB
	TC       *TC
}

type TC struct {
	TC_        *tc.Tc
	Iface      *net.Interface
	ObjectRoot *tc.Object
	HandleRoot uint32
}
