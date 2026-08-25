package timeout

import (
	"testing"
	"time"

	"github.com/selyusize/my-home/pkg/bitrix"
)

func TestKhabarovsk(t *testing.T) {
	t.Parallel()
	loc, err := khabarovsk()
	if err != nil {
		t.Fatal(err)
	}
	if loc.String() != khabarovskZone {
		t.Fatalf("zone=%s want %s", loc, khabarovskZone)
	}
}

func TestInOpenWindow(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("Khabarovsk", 10*3600)
	cases := []struct {
		hour int
		want bool
	}{
		{8, false},
		{9, true},
		{17, true},
		{18, false},
		{0, false},
	}
	for _, tc := range cases {
		now := time.Date(2026, 8, 24, tc.hour, 0, 0, 0, loc)
		if got := inOpenWindow(now, loc); got != tc.want {
			t.Fatalf("hour %d: got %v want %v", tc.hour, got, tc.want)
		}
	}
}

func TestPastClose(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("Khabarovsk", 10*3600)
	cases := []struct {
		hour int
		want bool
	}{
		{17, false},
		{18, true},
		{19, true},
		{0, false},
	}
	for _, tc := range cases {
		now := time.Date(2026, 8, 24, tc.hour, 0, 0, 0, loc)
		if got := pastClose(now, loc); got != tc.want {
			t.Fatalf("hour %d: got %v want %v", tc.hour, got, tc.want)
		}
	}
}

func TestNextOpen(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("Khabarovsk", 10*3600)
	before := time.Date(2026, 8, 24, 8, 0, 0, 0, loc)
	got := nextOpen(before, loc)
	want := time.Date(2026, 8, 24, 9, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("before: got %s want %s", got, want)
	}

	atNine := time.Date(2026, 8, 24, 9, 0, 0, 0, loc)
	got = nextOpen(atNine, loc)
	want = time.Date(2026, 8, 25, 9, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("at nine: got %s want %s", got, want)
	}
}

func TestNextClose(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("Khabarovsk", 10*3600)
	morning := time.Date(2026, 8, 24, 10, 0, 0, 0, loc)
	got := nextClose(morning, loc)
	want := time.Date(2026, 8, 24, 18, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("morning: got %s want %s", got, want)
	}

	evening := time.Date(2026, 8, 24, 18, 0, 0, 0, loc)
	got = nextClose(evening, loc)
	want = time.Date(2026, 8, 25, 18, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Fatalf("evening: got %s want %s", got, want)
	}
}

func TestNeedsOpen(t *testing.T) {
	t.Parallel()
	if !needsOpen(bitrix.TimeManClosed) || !needsOpen(bitrix.TimeManExpired) {
		t.Fatal("idle shift should open")
	}
	if needsOpen(bitrix.TimeManOpened) || needsOpen(bitrix.TimeManPaused) || needsOpen(bitrix.TimeManUnknown) {
		t.Fatal("running or unknown should stay")
	}
}

func TestNeedsClose(t *testing.T) {
	t.Parallel()
	if !needsClose(bitrix.TimeManOpened) || !needsClose(bitrix.TimeManPaused) || !needsClose(bitrix.TimeManExpired) {
		t.Fatal("running shift should close")
	}
	if needsClose(bitrix.TimeManClosed) || needsClose(bitrix.TimeManUnknown) {
		t.Fatal("closed or unknown should stay")
	}
}
