package bitrix

import (
	"fmt"
	"time"

	"github.com/selyusize/my-home/pkg/db"

	"gorm.io/gorm"
)

type bitrixCredential struct {
	ID        uint `gorm:"primaryKey"`
	Domain    string
	Webhook   string `gorm:"not null"`
	UpdatedAt time.Time
}

func (bitrixCredential) TableName() string {
	return "bitrix_credentials"
}

type CredentialStore struct {
	db db.Client
}

func NewCredentialStore(database db.Client) (*CredentialStore, error) {
	if database == nil {
		return nil, fmt.Errorf("bitrix credentials: database is required")
	}
	store := &CredentialStore{db: database}
	if err := database.DB().AutoMigrate(&bitrixCredential{}); err != nil {
		return nil, fmt.Errorf("bitrix credentials: migrate: %w", err)
	}
	return store, nil
}

func (s *CredentialStore) Load() (domain, webhook string, err error) {
	var row bitrixCredential
	if err := s.db.DB().Limit(1).Find(&row).Error; err != nil {
		return "", "", fmt.Errorf("bitrix credentials: load: %w", err)
	}
	return row.Domain, row.Webhook, nil
}

func (s *CredentialStore) Save(domain, webhook string) error {
	var row bitrixCredential
	if err := s.db.DB().Limit(1).Find(&row).Error; err != nil {
		return fmt.Errorf("bitrix credentials: load: %w", err)
	}
	row.Domain = domain
	row.Webhook = webhook
	if row.ID == 0 {
		if err := s.db.DB().Create(&row).Error; err != nil {
			return fmt.Errorf("bitrix credentials: save: %w", err)
		}
		return nil
	}
	if err := s.db.DB().Save(&row).Error; err != nil {
		return fmt.Errorf("bitrix credentials: save: %w", err)
	}
	return nil
}

func (s *CredentialStore) Delete() error {
	if err := s.db.DB().Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&bitrixCredential{}).Error; err != nil {
		return fmt.Errorf("bitrix credentials: delete: %w", err)
	}
	return nil
}
