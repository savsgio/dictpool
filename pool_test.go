package dictpool

import "testing"

func TestAcquireDict(t *testing.T) {
	d := AcquireDict()

	if d == nil {
		t.Error("nil dict")
	}
}

func TestReleaseDict(t *testing.T) {
	d := AcquireDict()
	d.Set("key", "value")

	ReleaseDict(d)

	if len(d.D) > 0 {
		t.Error("the dict has not been reseted")
	}
}
