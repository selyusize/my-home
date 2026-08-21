package settings

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/selyusize/my-home/pkg/db"
)

const (
	maxOrderKeyLength  = 64
	maxOrderLength     = 32
	maxOrderItemLength = 128
)

type uiSetting struct {
	Key       string `gorm:"primaryKey;size:64"`
	Value     string `gorm:"type:text;not null"`
	UpdatedAt time.Time
}

func (uiSetting) TableName() string {
	return "ui_settings"
}

type SettingsService struct {
	db db.Client
}

func NewSettingsService(database db.Client) (*SettingsService, error) {
	if database == nil {
		return nil, fmt.Errorf("settings: database is required")
	}
	service := &SettingsService{db: database}
	if err := database.AutoMigrate(&uiSetting{}); err != nil {
		return nil, fmt.Errorf("settings: migrate: %w", err)
	}
	return service, nil
}

func (s *SettingsService) GetOrder(key string) ([]string, error) {
	key, err := normalizeOrderKey(key)
	if err != nil {
		return nil, err
	}

	var row uiSetting
	if err := s.db.DB().Where("key = ?", key).Limit(1).Find(&row).Error; err != nil {
		return nil, fmt.Errorf("settings: get %s: %w", key, err)
	}
	if row.Key == "" {
		return []string{}, nil
	}

	order, err := parseOrder(row.Value)
	if err != nil {
		return nil, fmt.Errorf("settings: parse %s: %w", key, err)
	}
	return order, nil
}

func (s *SettingsService) SetOrder(key string, order []string) error {
	key, err := normalizeOrderKey(key)
	if err != nil {
		return err
	}
	if err := validateOrder(order); err != nil {
		return err
	}

	normalized := make([]string, len(order))
	for i, id := range order {
		normalized[i] = strings.TrimSpace(id)
	}

	payload, err := json.Marshal(normalized)
	if err != nil {
		return fmt.Errorf("settings: encode %s: %w", key, err)
	}

	row := uiSetting{Key: key, Value: string(payload)}
	if err := s.db.DB().Save(&row).Error; err != nil {
		return fmt.Errorf("settings: save %s: %w", key, err)
	}
	return nil
}

func normalizeOrderKey(key string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", fmt.Errorf("settings: key is required")
	}
	if utf8.RuneCountInString(key) > maxOrderKeyLength {
		return "", fmt.Errorf("settings: key is too long")
	}
	for _, r := range key {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == ':' || r == '-' || r == '_' {
			continue
		}
		return "", fmt.Errorf("settings: key contains invalid characters")
	}
	return key, nil
}

func validateOrder(order []string) error {
	if len(order) > maxOrderLength {
		return fmt.Errorf("settings: order is too long")
	}
	seen := make(map[string]struct{}, len(order))
	for _, id := range order {
		id = strings.TrimSpace(id)
		if id == "" {
			return fmt.Errorf("settings: order item is empty")
		}
		if utf8.RuneCountInString(id) > maxOrderItemLength {
			return fmt.Errorf("settings: order item is too long")
		}
		if _, exists := seen[id]; exists {
			return fmt.Errorf("settings: order contains duplicates")
		}
		seen[id] = struct{}{}
	}
	return nil
}

func parseOrder(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return []string{}, nil
	}
	var order []string
	if err := json.Unmarshal([]byte(raw), &order); err != nil {
		return nil, err
	}
	if err := validateOrder(order); err != nil {
		return nil, err
	}
	if order == nil {
		return []string{}, nil
	}
	return order, nil
}
