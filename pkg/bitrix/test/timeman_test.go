package bitrix_test

import (
	"testing"

	"github.com/selyusize/my-home/pkg/bitrix"
)

func TestDecodeTimeManStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected bitrix.TimeManStatus
	}{
		{name: "opened", input: `{"STATUS":"OPENED"}`, expected: bitrix.TimeManOpened},
		{name: "paused", input: `{"STATUS":"paused"}`, expected: bitrix.TimeManPaused},
		{name: "expired", input: `{"STATUS":"EXPIRED"}`, expected: bitrix.TimeManExpired},
		{name: "empty", input: `{"STATUS":""}`, expected: bitrix.TimeManUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := bitrix.DecodeTimeMan([]byte(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			if got.Status != tt.expected {
				t.Fatalf("status=%s want %s", got.Status, tt.expected)
			}
		})
	}
}

func TestDecodeTimeManDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{name: "string seconds", input: `{"STATUS":"OPENED","TIME_START":"2026-08-21T01:00:00+10:00","DURATION":"3600"}`, expected: 3600},
		{name: "clock zero", input: `{"STATUS":"CLOSED","DURATION":"00:00:00"}`, expected: 0},
		{name: "clock hours", input: `{"STATUS":"OPENED","DURATION":"01:02:03"}`, expected: 3723},
		{name: "clock minutes", input: `{"STATUS":"PAUSED","DURATION":"12:30"}`, expected: 750},
		{name: "number seconds", input: `{"STATUS":"OPENED","DURATION":90}`, expected: 90},
		{name: "empty string", input: `{"STATUS":"CLOSED","DURATION":""}`, expected: 0},
		{name: "null", input: `{"STATUS":"CLOSED","DURATION":null}`, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := bitrix.DecodeTimeMan([]byte(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			if got.Duration != tt.expected {
				t.Fatalf("duration=%d want %d", got.Duration, tt.expected)
			}
		})
	}
}

func TestDecodeTimeManLeaks(t *testing.T) {
	t.Parallel()

	got, err := bitrix.DecodeTimeMan([]byte(`{"STATUS":"PAUSED","DURATION":"01:00:00","TIME_LEAKS":"00:17:53"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.Duration != 3600 {
		t.Fatalf("duration=%d", got.Duration)
	}
	if got.TimeLeaks != 17*60+53 {
		t.Fatalf("leaks=%d", got.TimeLeaks)
	}
}
