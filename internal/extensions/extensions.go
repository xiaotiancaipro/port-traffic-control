package extensions

import (
	"fmt"
	"port-traffic-control/internal/configs"
)

func New(config *configs.Configuration) (ext *Extensions, err error) {

	db, err := NewDB(config.Database)
	if err != nil {
		err = fmt.Errorf("failed to initialize the database, Error=%v", err)
		return
	}

	tc_, err := NewTC(config.TC)
	if err != nil {
		err = fmt.Errorf("failed to initialize the TC, Error=%v", err)
		return
	}

	ext = &Extensions{
		Database: db,
		TC:       tc_,
	}

	return

}

func (ext *Extensions) Close() {
	_ = ext.TC.Close()
}
