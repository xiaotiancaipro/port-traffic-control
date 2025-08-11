package extensions

import (
	"fmt"
	"net"
	"port-traffic-control/internal/configs"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"golang.org/x/sys/unix"
)

func NewTC(config *configs.TCConfig) (tc_ *TC, err error) {

	iface, err := net.InterfaceByName(config.InterfaceName)
	if err != nil {
		err = fmt.Errorf("failed to obtain interface '%s': %v", config.InterfaceName, err)
		return
	}
	if iface.Flags&net.FlagUp == 0 {
		err = fmt.Errorf("interface '%s' is not up", config.InterfaceName)
		return
	}

	connect, err := tc.Open(&tc.Config{})
	if err != nil {
		err = fmt.Errorf("failed to open TC connection: %v", err)
		return
	}

	rootHandle := core.BuildHandle(0x1, 0x0)

	exists := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(iface.Index),
			Handle:  rootHandle,
			Parent:  tc.HandleRoot,
		},
	}
	_ = connect.Qdisc().Delete(&exists)

	root := &tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(iface.Index),
			Handle:  rootHandle,
			Parent:  tc.HandleRoot,
			Info:    0,
		},
		Attribute: tc.Attribute{
			Kind: "htb",
			Htb: &tc.Htb{
				Init: &tc.HtbGlob{
					Version:      config.HTBVersion,
					Rate2Quantum: config.Rate2Quantum,
				},
			},
		},
	}
	err = connect.Qdisc().Replace(root)
	if err != nil {
		_ = connect.Close()
		err = fmt.Errorf("failed to create HTB qdisc: %v", err)
		return
	}

	tc_ = &TC{
		TC_:        connect,
		Iface:      iface,
		ObjectRoot: root,
		HandleRoot: rootHandle,
	}
	return

}

func (tc_ *TC) CloseTC() error {
	if tc_ == nil {
		return nil
	}
	if err := tc_.TC_.Qdisc().Delete(tc_.ObjectRoot); err != nil {
		return fmt.Errorf("failed to delete Qdisc: %v", err)
	}
	if err := tc_.TC_.Close(); err != nil {
		return fmt.Errorf("failed to close TC: %v", err)
	}
	return nil
}
