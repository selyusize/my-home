package bitrix_test

import (
	"errors"
	"testing"

	"github.com/selyusize/my-home/pkg/bitrix"
)

func TestDecodeProfile(t *testing.T) {
	t.Parallel()

	profile, err := bitrix.DecodeProfile([]byte(`{
		"ID":"12",
		"NAME":"Ivan",
		"LAST_NAME":"Petrov",
		"EMAIL":"ivan@example.com",
		"WORK_POSITION":"Dev",
		"PERSONAL_PHOTO":"https://img"
	}`), "https://company.bitrix24.ru")
	if err != nil {
		t.Fatal(err)
	}
	if profile.Name != "Ivan Petrov" {
		t.Fatalf("name=%q", profile.Name)
	}
	if profile.PageURL != "https://company.bitrix24.ru/company/personal/user/12/" {
		t.Fatalf("page=%q", profile.PageURL)
	}
	if profile.Position != "Dev" {
		t.Fatalf("position=%q", profile.Position)
	}
}

func TestClientNotConfigured(t *testing.T) {
	t.Parallel()

	client := bitrix.New()
	if client.IsAuthenticated() {
		t.Fatal("expected empty client")
	}
	if err := client.TimeManOpen(t.Context()); !errors.Is(err, bitrix.ErrNotConfigured) {
		t.Fatalf("err=%v", err)
	}
}
