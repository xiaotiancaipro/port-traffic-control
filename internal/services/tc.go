package services

import (
	"fmt"

	"github.com/florianl/go-tc"
	"github.com/florianl/go-tc/core"
	"golang.org/x/sys/unix"
)

func (tcs *TCService) CreateParentClass(minor uint32, rate uint32) error {
	rate = tcs.MbpsToBps(rate)
	parentClass := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(tcs.Iface.Index),
			Handle:  core.BuildHandle(0x1, minor),
			Parent:  tcs.HandleRoot,
		},
		Attribute: tc.Attribute{
			Kind: "htb",
			Htb: &tc.Htb{
				Parms: &tc.HtbOpt{
					Rate: tc.RateSpec{
						Rate:      rate,
						CellLog:   3,
						Linklayer: 1,
						Overhead:  26,
						Mpu:       64,
						CellAlign: 0,
					},
					Ceil: tc.RateSpec{
						Rate:      rate,
						CellLog:   3,
						Linklayer: 1,
						Overhead:  26,
						Mpu:       64,
						CellAlign: 0,
					},
					Buffer:  uint32(1),
					Cbuffer: uint32(1),
				},
			},
		},
	}
	if err := tcs.TC.Class().Add(&parentClass); err != nil {
		err = fmt.Errorf("create parent class failed, Error=%v", err)
		tcs.Log.Error(err)
		return err
	}
	return nil
}

func (tcs *TCService) CreateChildClass(parentMinor uint32, childMinor uint32, rate uint32, rateCeil uint32) error {

	rate = tcs.MbpsToBps(rate)
	rateCeil = tcs.MbpsToBps(rateCeil)

	childClass := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(tcs.Iface.Index),
			Handle:  core.BuildHandle(parentMinor, childMinor),
			Parent:  core.BuildHandle(0x1, parentMinor), // TODO
		},
		Attribute: tc.Attribute{
			Kind: "htb",
			Htb: &tc.Htb{
				Parms: &tc.HtbOpt{
					Rate: tc.RateSpec{
						Rate: rate,
					},
					Ceil: tc.RateSpec{
						Rate:      rateCeil,
						CellLog:   3,
						Linklayer: 1,
						Overhead:  26,
						Mpu:       64,
						CellAlign: 0,
					},
					Buffer:  uint32(1),
					Cbuffer: uint32(1),
					Quantum: 1500,
				},
			},
		},
	}
	if err := tcs.TC.Class().Add(&childClass); err != nil {
		err = fmt.Errorf("create child class failed, Error=%v", err)
		tcs.Log.Error(err)
		return err
	}

	target := uint32(5000)
	limit := uint32(10240)
	quantum := uint32(1500)
	fqQdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(tcs.Iface.Index),
			Handle:  core.BuildHandle(0x0, 0x0),
			Parent:  core.BuildHandle(parentMinor, childMinor), // TODO
		},
		Attribute: tc.Attribute{
			Kind: "fq_codel",
			FqCodel: &tc.FqCodel{
				Target:  &target,
				Limit:   &limit,
				Quantum: &quantum,
			},
		},
	}
	if err := tcs.TC.Qdisc().Add(&fqQdisc); err != nil {
		err = fmt.Errorf("failed to add fq_codel queue discipline, Error=%v", err)
		tcs.Log.Error(err)
		return err
	}

	return nil

}

func (tcs *TCService) DeleteChildClass(parentMinor uint32, childMinor uint32) error {

	// First remove the fq_codel qdisc attached to the child class
	fqQdisc := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(tcs.Iface.Index),
			Handle:  core.BuildHandle(0x0, 0x0),
			Parent:  core.BuildHandle(parentMinor, childMinor),
		},
	}
	if err := tcs.TC.Qdisc().Delete(&fqQdisc); err != nil {
		// Log and continue. If qdisc is missing, deleting class might still succeed.
		tcs.Log.Warningf("failed to delete fq_codel qdisc for %x:%x, Error=%v", parentMinor, childMinor, err)
	}

	class := tc.Object{
		Msg: tc.Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: uint32(tcs.Iface.Index),
			Handle:  core.BuildHandle(parentMinor, childMinor),
			Parent:  core.BuildHandle(0x1, parentMinor),
		},
	}
	if err := tcs.TC.Class().Delete(&class); err != nil {
		err = fmt.Errorf("delete child class failed, Error=%v", err)
		tcs.Log.Error(err)
		return err
	}
	return nil
}

func (tcs *TCService) MbpsToBps(mbps uint32) uint32 {
	return mbps * 1_000_000 / 8
}
