package bitrix

import (
	"context"
	"fmt"
	"log"

	pkgbitrix "github.com/selyusize/my-home/pkg/bitrix"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type BitrixService struct {
	client *pkgbitrix.Client
	creds  *CredentialStore
}

func NewBitrixService(creds *CredentialStore) (*BitrixService, error) {
	return &BitrixService{client: pkgbitrix.New(), creds: creds}, nil
}

func (s *BitrixService) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	if s.creds == nil {
		return nil
	}
	domain, webhook, err := s.creds.Load()
	if err != nil {
		return err
	}
	if webhook == "" {
		return nil
	}
	if err := s.client.SetWebhook(domain, webhook); err != nil {
		return s.dropStored("invalid stored webhook", err)
	}
	if _, err := s.client.Profile(ctx); err != nil {
		return s.dropStored("stored webhook is invalid", err)
	}
	return nil
}

func (s *BitrixService) dropStored(reason string, cause error) error {
	log.Printf("bitrix: %s, dropping session: %v", reason, cause)
	s.client.Logout()
	if s.creds == nil {
		return nil
	}
	if err := s.creds.Delete(); err != nil {
		log.Printf("bitrix: drop webhook: %v", err)
	}
	return nil
}

func (s *BitrixService) IsAuthenticated() bool {
	return s.client.IsAuthenticated()
}

func (s *BitrixService) SetWebhook(ctx context.Context, domain, webhook string) (*pkgbitrix.Profile, error) {
	if err := s.client.SetWebhook(domain, webhook); err != nil {
		return nil, err
	}
	profile, err := s.client.Profile(ctx)
	if err != nil {
		s.client.Logout()
		return nil, err
	}
	if s.creds != nil {
		if err := s.creds.Save(s.client.Portal(), s.client.WebhookURL()); err != nil {
			s.client.Logout()
			return nil, fmt.Errorf("bitrix: persist webhook: %w", err)
		}
	}
	return profile, nil
}

func (s *BitrixService) Profile(ctx context.Context) (*pkgbitrix.Profile, error) {
	return s.client.Profile(ctx)
}

func (s *BitrixService) TimeMan(ctx context.Context) (*pkgbitrix.TimeMan, error) {
	return s.client.TimeMan(ctx)
}

func (s *BitrixService) TimeManOpen(ctx context.Context) error {
	return s.client.TimeManOpen(ctx)
}

func (s *BitrixService) TimeManPause(ctx context.Context) error {
	return s.client.TimeManPause(ctx)
}

func (s *BitrixService) TimeManClose(ctx context.Context) error {
	return s.client.TimeManClose(ctx)
}

func (s *BitrixService) Logout() error {
	s.client.Logout()
	if s.creds == nil {
		return nil
	}
	return s.creds.Delete()
}
