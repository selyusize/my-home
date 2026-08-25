package timeout

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/selyusize/my-home/pkg/bitrix"
)

const (
	retryAfter   = 5 * time.Minute
	actionOpened = "opened"
	actionClosed = "closed"
)

// Shift is the TimeMan surface used to start and finish a workday.
type Shift interface {
	TimeMan(ctx context.Context) (*bitrix.TimeMan, error)
	TimeManOpen(ctx context.Context) error
	TimeManClose(ctx context.Context) error
}

type sleeper func(ctx context.Context, until time.Time) error

// Watch opens at 09:00 Khabarovsk if idle and closes at 18:00 if still running.
func Watch(ctx context.Context, shift Shift, notify func(string)) error {
	loc, err := khabarovsk()
	if err != nil {
		return err
	}
	log.Printf("timeout: watching 09:00–18:00 %s", loc)
	return run(ctx, shift, loc, time.Now, sleepUntil, notify)
}

func run(
	ctx context.Context,
	shift Shift,
	loc *time.Location,
	now func() time.Time,
	sleep sleeper,
	notify func(string),
) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := step(ctx, shift, loc, now, sleep, notify); err != nil {
			return err
		}
	}
}

func step(
	ctx context.Context,
	shift Shift,
	loc *time.Location,
	now func() time.Time,
	sleep sleeper,
	notify func(string),
) error {
	moment := now()
	if pastClose(moment, loc) {
		return apply(ctx, shift, closeIfNeeded, actionClosed, nextOpen(moment, loc), moment, sleep, notify)
	}
	if inOpenWindow(moment, loc) {
		return apply(ctx, shift, openIfNeeded, actionOpened, nextClose(moment, loc), moment, sleep, notify)
	}
	return sleep(ctx, nextOpen(moment, loc))
}

func apply(
	ctx context.Context,
	shift Shift,
	act func(context.Context, Shift) (bool, error),
	action string,
	until time.Time,
	moment time.Time,
	sleep sleeper,
	notify func(string),
) error {
	changed, err := act(ctx, shift)
	if err != nil {
		log.Printf("timeout: %v", err)
		return sleep(ctx, moment.Add(retryAfter))
	}
	if changed && notify != nil {
		notify(action)
	}
	return sleep(ctx, until)
}

func openIfNeeded(ctx context.Context, shift Shift) (bool, error) {
	tm, err := shift.TimeMan(ctx)
	if err != nil {
		return false, fmt.Errorf("status: %w", err)
	}
	if !needsOpen(tm.Status) {
		log.Printf("timeout: shift is %s", tm.Status)
		return false, nil
	}
	if tm.Status == bitrix.TimeManExpired {
		if err := shift.TimeManClose(ctx); err != nil {
			return false, fmt.Errorf("close expired: %w", err)
		}
	}
	if err := shift.TimeManOpen(ctx); err != nil {
		return false, fmt.Errorf("open: %w", err)
	}
	log.Printf("timeout: opened shift")
	return true, nil
}

func closeIfNeeded(ctx context.Context, shift Shift) (bool, error) {
	tm, err := shift.TimeMan(ctx)
	if err != nil {
		return false, fmt.Errorf("status: %w", err)
	}
	if !needsClose(tm.Status) {
		log.Printf("timeout: shift is %s", tm.Status)
		return false, nil
	}
	if err := shift.TimeManClose(ctx); err != nil {
		return false, fmt.Errorf("close: %w", err)
	}
	log.Printf("timeout: closed shift")
	return true, nil
}

func sleepUntil(ctx context.Context, until time.Time) error {
	delay := time.Until(until)
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
