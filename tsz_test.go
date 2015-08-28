package tsz

import (
	"testing"
	"time"
)

func TestExampleEncoding(t *testing.T) {

	// Example from the paper
	t0, _ := time.ParseInLocation("Jan _2 2006 15:04:05", "Mar 24 2015 02:00:00", time.Local)
	tunix := uint32(t0.Unix())

	s := New(tunix)

	tunix += 62
	s.Push(tunix, 12)

	tunix += 60
	s.Push(tunix, 12)

	tunix += 60
	s.Push(tunix, 24)

	// extra tests

	// floating point masking/shifting bug
	tunix += 60
	s.Push(tunix, 13)

	tunix += 60
	s.Push(tunix, 24)

	// delta-of-delta sizes
	tunix += 300 // == delta-of-delta of 240
	s.Push(tunix, 24)

	tunix += 900 // == delta-of-delta of 600
	s.Push(tunix, 24)

	tunix += 900 + 2050 // == delta-of-delta of 600
	s.Push(tunix, 24)

	s.Finish()

	it := s.Iter()

	tunix = uint32(t0.Unix())
	want := []struct {
		t uint32
		v float64
	}{
		{tunix + 62, 12},
		{tunix + 122, 12},
		{tunix + 182, 24},

		{tunix + 242, 13},
		{tunix + 302, 24},

		{tunix + 602, 24},
		{tunix + 1502, 24},
		{tunix + 4452, 24},
	}

	for _, w := range want {
		if !it.Next() {
			t.Fatalf("Next()=false, want true")
		}
		tt, vv := it.Values()
		if w.t != tt || w.v != vv {
			t.Errorf("Values()=(%v,%v), want (%v,%v)\n", tt, vv, w.t, w.v)
		}
	}

	if it.Next() {
		t.Fatalf("Next()=true, want false")
	}

	if err := it.Err(); err != nil {
		t.Errorf("it.Err()=%v, want nil", err)
	}
}
