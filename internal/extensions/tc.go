package extensions

import (
	"fmt"
	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"golang.org/x/sys/unix"
	"net"
	"port-traffic-control/internal/configs"
)

func NewTC(config *configs.TCConfig) (tc_ *tc.Tc, err error) {

	iface, err := net.InterfaceByName(config.InterfaceName)
	if err != nil {
		err = fmt.Errorf("failed to obtain the interface, Error=%v", err)
		return
	}

	tc_, err = tc.Open(&tc.Config{})
	if err != nil {
		err = fmt.Errorf("failed to open TC, Error=%v", err)
		return
	}

	rootQdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(iface.Index),
			Handle:  core.BuildHandle(0xFFFF, 0),
			Parent:  tc.HandleRoot,
		},
		Attribute: tc.Attribute{
			Kind: "htb",
			Htb: &tc.Htb{
				Init: &tc.HtbGlob{
					Rate2Quantum: 10,
				},
			},
		},
	}
	if err = tc_.Qdisc().Add(&rootQdisc); err != nil {
		_ = tc_.Close()
		err = fmt.Errorf("failed to add Qdisc, Error=%v", err)
		return
	}

	return

}
