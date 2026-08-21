package bitrix_test

import (
	"errors"
	"testing"

	"github.com/selyusize/my-home/pkg/bitrix"
)

func TestComposeWebhook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		domain      string
		webhook     string
		wantPortal  string
		wantWebhook string
		wantErr     error
	}{
		{
			name:        "full url ignores domain",
			domain:      "other.bitrix24.ru",
			webhook:     "https://company.bitrix24.ru/rest/12/secretcode",
			wantPortal:  "https://company.bitrix24.ru",
			wantWebhook: "https://company.bitrix24.ru/rest/12/secretcode/",
		},
		{
			name:        "user and code",
			domain:      "company.bitrix24.ru",
			webhook:     "12/secretcode",
			wantPortal:  "https://company.bitrix24.ru",
			wantWebhook: "https://company.bitrix24.ru/rest/12/secretcode/",
		},
		{
			name:        "rest path",
			domain:      "https://company.bitrix24.ru/",
			webhook:     "rest/12/secretcode",
			wantPortal:  "https://company.bitrix24.ru",
			wantWebhook: "https://company.bitrix24.ru/rest/12/secretcode/",
		},
		{
			name:        "secret defaults to user 1",
			domain:      "company.bitrix24.ru",
			webhook:     "secretcode",
			wantPortal:  "https://company.bitrix24.ru",
			wantWebhook: "https://company.bitrix24.ru/rest/1/secretcode/",
		},
		{
			name:    "missing webhook",
			domain:  "company.bitrix24.ru",
			wantErr: bitrix.ErrMissingWebhook,
		},
		{
			name:    "missing portal",
			webhook: "secretcode",
			wantErr: bitrix.ErrMissingPortal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			portal, webhookURL, err := bitrix.ComposeWebhook(tt.domain, tt.webhook)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err=%v want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if portal != tt.wantPortal {
				t.Fatalf("portal=%q want %q", portal, tt.wantPortal)
			}
			if webhookURL != tt.wantWebhook {
				t.Fatalf("webhook=%q want %q", webhookURL, tt.wantWebhook)
			}
		})
	}
}
