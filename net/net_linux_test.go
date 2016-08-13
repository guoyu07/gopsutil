package net

import (
	"os"
	"syscall"
	"testing"

	"github.com/percona/gopsutil/internal/common"
)

func TestGetProcInodesAll(t *testing.T) {
	if os.Getenv("CIRCLECI") == "true" {
		t.Skip("Skip CI")
	}

	root := common.HostProc("")
	v, err := getProcInodesAll(root)
	if err != nil {
		t.Fatalf("cannot get process inodes: %s", err)
	}
	if len(v) == 0 {
		t.Error("inodes list is empty")
	}
}

type AddrTest struct {
	IP    string
	Port  int
	Error bool
}

func TestDecodeAddress(t *testing.T) {

	addr := map[string]AddrTest{
		"0500000A:0016": AddrTest{
			IP:   "10.0.0.5",
			Port: 22,
		},
		"0100007F:D1C2": AddrTest{
			IP:   "127.0.0.1",
			Port: 53698,
		},
		"11111:0035": AddrTest{
			Error: true,
		},
		"0100007F:BLAH": AddrTest{
			Error: true,
		},
		"0085002452100113070057A13F025401:0035": AddrTest{
			IP:   "2400:8500:1301:1052:a157:7:154:23f",
			Port: 53,
		},
		"00855210011307F025401:0035": AddrTest{
			Error: true,
		},
	}

	for src, dst := range addr {
		family := syscall.AF_INET
		if len(src) > 13 {
			family = syscall.AF_INET6
		}
		addr, err := decodeAddress(uint32(family), src)
		if dst.Error {
			if err == nil {
				t.Errorf("error decoding address: %s\nsrc:%#v", err, src)
			}
		} else {
			if err != nil {
				t.Errorf("error decoding address: %s\nsrc:%#v", err, src)
			}
			if dst.IP != addr.IP {
				t.Errorf("IP addresses are different. %v != %v\nsrc: %#v", dst.IP, addr.IP, src)
			}
			if dst.Port != int(addr.Port) {
				t.Errorf("ports don't match. %v != %v", dst.Port, addr.Port)
			}
		}
	}
}

func TestReverse(t *testing.T) {
	src := []byte{0x01, 0x02, 0x03}
	expect := []byte{0x03, 0x02, 0x01}
	got := Reverse(src)
	for i := 0; i < len(expect); i++ {
		if expect[i] != got[i] {
			t.Errorf("reverse is wrong. Got: %v, expected: %v", got, expect)
			break
		}
	}
}
