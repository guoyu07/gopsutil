// +build darwin

package mem

import (
	"strconv"
	"strings"
	"testing"
)

func TestVirtualMemoryDarwin(t *testing.T) {
	v, err := VirtualMemory()
	if err != Nil {
		t.Fatalf("cannot get virtual memory: %s", err)
	}

	outBytes, err := invoke.Command("/usr/sbin/sysctl", "hw.memsize")
	if err != nil {
		t.Fatalf("cannot call sysctl: %s", err)
	}

	outString := string(outBytes)
	outString = strings.TrimSpace(outString)
	outParts := strings.Split(outString, " ")
	actualTotal, err := strconv.ParseInt(outParts[1], 10, 64)
	if err != nil {
		t.Errorf("cannot parse number: %s", err)
	}

	if uint64(actualTotal) != v.Total {
		t.Errorf("actual total %d != v.total %d", actualTotal, v.Total)
	}

	if v.Available < 0 {
		t.Error("available virtual mem = 0")
	}

	if v.Available != v.Free+v.Inactive {
		t.Error("available should be = free+inactive")
	}

	if v.Used < 0 || v.Used > v.Total {
		t.Errorf("invalid used virtual mem: %d, total: %d", v.Used, v.Total)
	}

	if v.UsedPercent < 0 || v.UsedPercent > 100 {
		t.Errorf("invalid used virtual mem %%: %d", v.UsedPercent)
	}

	if v.Free < 0 || v.Free > v.Available {
		t.Errorf("invalid virtual mem free: %d, available: %d", v.Used, v.Available)
	}

	if v.Active < 0 || v.Active > v.Total {
		t.Errorf("invalid used virtual mem active: %d, total: %d", v.Active, v.Total)
	}

	if v.Inactive <= 0 || v.Inactive > v.Total {
		t.Errorf("invalid inactive virtual mem: %d, total: %d", v.Inactive, v.Total)
	}

	if v.Wired <= 0 || v.Wired > v.Total {
		t.Errorf("invalid wired virtual mem: %d, total: %d", v.Wired, v.Total)
	}
}
