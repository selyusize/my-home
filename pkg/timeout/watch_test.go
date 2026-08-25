package timeout

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/selyusize/my-home/pkg/bitrix"
)

type fakeShift struct {
	status    bitrix.TimeManStatus
	statusErr error
	openErr   error
	closeErr  error
	opens     int
	closes    int
}

func (f *fakeShift) TimeMan(context.Context) (*bitrix.TimeMan, error) {
	if f.statusErr != nil {
		return nil, f.statusErr
	}
	return &bitrix.TimeMan{Status: f.status}, nil
}

func (f *fakeShift) TimeManOpen(context.Context) error {
	if f.openErr != nil {
		return f.openErr
	}
	f.opens++
	f.status = bitrix.TimeManOpened
	return nil
}

func (f *fakeShift) TimeManClose(context.Context) error {
	if f.closeErr != nil {
		return f.closeErr
	}
	f.closes++
	f.status = bitrix.TimeManClosed
	return nil
}

func TestOpenIfNeededOpensClosedShift(t *testing.T) {
	t.Parallel()
	shift := &fakeShift{status: bitrix.TimeManClosed}
	opened, err := openIfNeeded(t.Context(), shift)
	if err != nil {
		t.Fatal(err)
	}
	if !opened || shift.opens != 1 {
		t.Fatalf("opened=%v opens=%d", opened, shift.opens)
	}
}

func TestOpenIfNeededClosesExpiredThenOpens(t *testing.T) {
	t.Parallel()
	shift := &fakeShift{status: bitrix.TimeManExpired}
	opened, err := openIfNeeded(t.Context(), shift)
	if err != nil {
		t.Fatal(err)
	}
	if !opened || shift.closes != 1 || shift.opens != 1 {
		t.Fatalf("opened=%v closes=%d opens=%d", opened, shift.closes, shift.opens)
	}
}

func TestOpenIfNeededSkipsOpenedShift(t *testing.T) {
	t.Parallel()
	shift := &fakeShift{status: bitrix.TimeManOpened}
	opened, err := openIfNeeded(t.Context(), shift)
	if err != nil {
		t.Fatal(err)
	}
	if opened || shift.opens != 0 {
		t.Fatalf("opened=%v opens=%d", opened, shift.opens)
	}
}

func TestCloseIfNeededClosesOpenShift(t *testing.T) {
	t.Parallel()
	shift := &fakeShift{status: bitrix.TimeManOpened}
	closed, err := closeIfNeeded(t.Context(), shift)
	if err != nil {
		t.Fatal(err)
	}
	if !closed || shift.closes != 1 {
		t.Fatalf("closed=%v closes=%d", closed, shift.closes)
	}
}

func TestCloseIfNeededSkipsClosedShift(t *testing.T) {
	t.Parallel()
	shift := &fakeShift{status: bitrix.TimeManClosed}
	closed, err := closeIfNeeded(t.Context(), shift)
	if err != nil {
		t.Fatal(err)
	}
	if closed || shift.closes != 0 {
		t.Fatalf("closed=%v closes=%d", closed, shift.closes)
	}
}

func TestRunOpensAfterNine(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("Khabarovsk", 10*3600)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	shift := &fakeShift{status: bitrix.TimeManClosed}
	now := time.Date(2026, 8, 24, 9, 1, 0, 0, loc)
	got := ""
	err := run(ctx, shift, loc, func() time.Time { return now }, func(ctx context.Context, until time.Time) error {
		want := time.Date(2026, 8, 24, 18, 0, 0, 0, loc)
		if !until.Equal(want) {
			t.Fatalf("sleep until %s want %s", until, want)
		}
		if shift.opens != 1 {
			t.Fatalf("opens=%d", shift.opens)
		}
		if got != actionOpened {
			t.Fatalf("notify=%q", got)
		}
		cancel()
		return ctx.Err()
	}, func(action string) { got = action })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}

func TestRunWaitsUntilNine(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("Khabarovsk", 10*3600)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	shift := &fakeShift{status: bitrix.TimeManClosed}
	now := time.Date(2026, 8, 24, 8, 0, 0, 0, loc)
	err := run(ctx, shift, loc, func() time.Time { return now }, func(ctx context.Context, until time.Time) error {
		want := time.Date(2026, 8, 24, 9, 0, 0, 0, loc)
		if !until.Equal(want) {
			t.Fatalf("sleep until %s want %s", until, want)
		}
		if shift.opens != 0 {
			t.Fatal("should not open before 09:00")
		}
		cancel()
		return ctx.Err()
	}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}

func TestRunClosesAfterSix(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("Khabarovsk", 10*3600)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	shift := &fakeShift{status: bitrix.TimeManOpened}
	now := time.Date(2026, 8, 24, 18, 1, 0, 0, loc)
	got := ""
	err := run(ctx, shift, loc, func() time.Time { return now }, func(ctx context.Context, until time.Time) error {
		want := time.Date(2026, 8, 25, 9, 0, 0, 0, loc)
		if !until.Equal(want) {
			t.Fatalf("sleep until %s want %s", until, want)
		}
		if shift.closes != 1 {
			t.Fatalf("closes=%d", shift.closes)
		}
		if got != actionClosed {
			t.Fatalf("notify=%q", got)
		}
		cancel()
		return ctx.Err()
	}, func(action string) { got = action })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}

func TestRunWaitsUntilSixWhenOpened(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("Khabarovsk", 10*3600)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	shift := &fakeShift{status: bitrix.TimeManOpened}
	now := time.Date(2026, 8, 24, 10, 0, 0, 0, loc)
	err := run(ctx, shift, loc, func() time.Time { return now }, func(ctx context.Context, until time.Time) error {
		want := time.Date(2026, 8, 24, 18, 0, 0, 0, loc)
		if !until.Equal(want) {
			t.Fatalf("sleep until %s want %s", until, want)
		}
		if shift.opens != 0 || shift.closes != 0 {
			t.Fatal("should not change a running shift")
		}
		cancel()
		return ctx.Err()
	}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}

func TestRunRetriesOnError(t *testing.T) {
	t.Parallel()
	loc := time.FixedZone("Khabarovsk", 10*3600)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	now := time.Date(2026, 8, 24, 18, 0, 0, 0, loc)
	shift := &fakeShift{status: bitrix.TimeManOpened, closeErr: errors.New("bitrix down")}
	err := run(ctx, shift, loc, func() time.Time { return now }, func(ctx context.Context, until time.Time) error {
		want := now.Add(retryAfter)
		if !until.Equal(want) {
			t.Fatalf("sleep until %s want %s", until, want)
		}
		cancel()
		return ctx.Err()
	}, nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err=%v", err)
	}
}
