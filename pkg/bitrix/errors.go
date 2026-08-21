package bitrix

import "errors"

var (
	ErrNotConfigured  = errors.New("bitrix: webhook is not configured")
	ErrMissingPortal  = errors.New("bitrix: portal is required")
	ErrMissingWebhook = errors.New("bitrix: webhook is required")
	ErrInvalidWebhook = errors.New("bitrix: invalid webhook")
)
