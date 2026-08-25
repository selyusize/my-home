package timeout

import (
	"context"
	"errors"
	"log"

	"github.com/selyusize/my-home/internal/bitrix"
	pktimeout "github.com/selyusize/my-home/pkg/timeout"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const timeManEvent = "bitrix:timeman"

type TimeoutService struct {
	bitrix *bitrix.BitrixService
}

func NewTimeoutService(bitrixService *bitrix.BitrixService) *TimeoutService {
	return &TimeoutService{bitrix: bitrixService}
}

func (s *TimeoutService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	go s.watch(ctx)
	return nil
}

func (s *TimeoutService) watch(ctx context.Context) {
	if s == nil || s.bitrix == nil {
		return
	}
	err := pktimeout.Watch(ctx, s.bitrix, emitTimeMan)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("timeout: %v", err)
	}
}

func emitTimeMan(action string) {
	app := application.Get()
	if app == nil {
		return
	}
	app.Event.Emit(timeManEvent, action)
}
