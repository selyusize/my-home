package bitrix

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	b24 "github.com/bitrix24/b24gosdk"
)

// Client talks to Bitrix24 REST through an inbound webhook.
type Client struct {
	mu      sync.Mutex
	sdk     *b24.Client
	portal  string
	webhook string
}

// New creates an unconfigured Bitrix24 client.
func New() *Client {
	return &Client{}
}

// SetWebhook configures the inbound webhook from portal + token/URL fields.
func (c *Client) SetWebhook(domain, webhook string) error {
	portal, webhookURL, err := ComposeWebhook(domain, webhook)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sdk = b24.NewClient(webhookURL)
	c.portal = portal
	c.webhook = webhookURL
	return nil
}

// Logout clears the in-memory webhook.
func (c *Client) Logout() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sdk = nil
	c.portal = ""
	c.webhook = ""
}

// IsAuthenticated reports whether a webhook is configured in memory.
func (c *Client) IsAuthenticated() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sdk != nil
}

// Portal returns the configured portal origin.
func (c *Client) Portal() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.portal
}

// WebhookURL returns the composed inbound webhook URL.
func (c *Client) WebhookURL() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.webhook
}

func (c *Client) core() (*b24.Core, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.sdk == nil {
		return nil, "", ErrNotConfigured
	}
	return c.sdk.Core(), c.portal, nil
}

func (c *Client) call(ctx context.Context, method string, dest any) error {
	raw, err := c.callRaw(ctx, method)
	if err != nil {
		return err
	}
	if dest == nil {
		return nil
	}
	if err := json.Unmarshal(raw, dest); err != nil {
		return fmt.Errorf("bitrix: %s: decode: %w", method, err)
	}
	return nil
}

func (c *Client) callRaw(ctx context.Context, method string) ([]byte, error) {
	core, _, err := c.core()
	if err != nil {
		return nil, err
	}
	raw, err := core.CallJSON(ctx, method, map[string]any{})
	if err != nil {
		return nil, fmt.Errorf("bitrix: %s: %w", method, err)
	}
	return raw, nil
}

type rawUser struct {
	ID            b24.ID `json:"ID"`
	Name          string `json:"NAME"`
	LastName      string `json:"LAST_NAME"`
	SecondName    string `json:"SECOND_NAME"`
	Email         string `json:"EMAIL"`
	PersonalPhoto string `json:"PERSONAL_PHOTO"`
	WorkPosition  string `json:"WORK_POSITION"`
}

// Profile returns the current user from user.current.
func (c *Client) Profile(ctx context.Context) (*Profile, error) {
	data, err := c.callRaw(ctx, "user.current")
	if err != nil {
		return nil, err
	}
	_, portal, err := c.core()
	if err != nil {
		return nil, err
	}
	profile, err := DecodeProfile(data, portal)
	if err != nil {
		return nil, fmt.Errorf("bitrix: user.current: decode: %w", err)
	}
	return profile, nil
}

// durationSeconds accepts Bitrix REST durations as a number, numeric string, or HH:MM:SS.
type durationSeconds int

func (n *durationSeconds) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*n = 0
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		v, err := parseDurationSeconds(s)
		if err != nil {
			return err
		}
		*n = durationSeconds(v)
		return nil
	}
	var v json.Number
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	f, err := v.Float64()
	if err != nil {
		return err
	}
	*n = durationSeconds(f)
	return nil
}

func parseDurationSeconds(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	if !strings.Contains(raw, ":") {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return 0, err
		}
		return int(v), nil
	}

	parts := strings.Split(raw, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return 0, fmt.Errorf("invalid duration %q", raw)
	}
	nums := make([]int, len(parts))
	for i, part := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return 0, err
		}
		nums[i] = v
	}
	if len(nums) == 2 {
		return nums[0]*60 + nums[1], nil
	}
	return nums[0]*3600 + nums[1]*60 + nums[2], nil
}

type rawTimeMan struct {
	Status    string          `json:"STATUS"`
	TimeStart string          `json:"TIME_START"`
	Duration  durationSeconds `json:"DURATION"`
	TimeLeaks durationSeconds `json:"TIME_LEAKS"`
}

// TimeMan returns the current workday status.
func (c *Client) TimeMan(ctx context.Context) (*TimeMan, error) {
	data, err := c.callRaw(ctx, "timeman.status")
	if err != nil {
		return nil, err
	}
	tm, err := DecodeTimeMan(data)
	if err != nil {
		return nil, fmt.Errorf("bitrix: timeman.status: decode: %w", err)
	}
	return tm, nil
}

// TimeManOpen starts or resumes the workday.
func (c *Client) TimeManOpen(ctx context.Context) error {
	return c.call(ctx, "timeman.open", nil)
}

// TimeManPause pauses the workday.
func (c *Client) TimeManPause(ctx context.Context) error {
	return c.call(ctx, "timeman.pause", nil)
}

// TimeManClose finishes the workday.
func (c *Client) TimeManClose(ctx context.Context) error {
	return c.call(ctx, "timeman.close", nil)
}

func mapProfile(user rawUser, portal string) *Profile {
	parts := make([]string, 0, 3)
	for _, part := range []string{user.Name, user.SecondName, user.LastName} {
		part = strings.TrimSpace(part)
		if part != "" {
			parts = append(parts, part)
		}
	}
	name := strings.Join(parts, " ")
	if name == "" {
		name = user.Email
	}
	id := user.ID.Int64()
	return &Profile{
		ID:        id,
		Name:      name,
		Email:     strings.TrimSpace(user.Email),
		Position:  strings.TrimSpace(user.WorkPosition),
		AvatarURL: strings.TrimSpace(user.PersonalPhoto),
		PortalURL: portal,
		PageURL:   personalPageURL(portal, id),
	}
}

func mapTimeMan(raw rawTimeMan) *TimeMan {
	return &TimeMan{
		Status:    normalizeTimeManStatus(raw.Status),
		TimeStart: raw.TimeStart,
		Duration:  int(raw.Duration),
		TimeLeaks: int(raw.TimeLeaks),
	}
}

// DecodeTimeMan maps a Bitrix timeman.status payload.
func DecodeTimeMan(data []byte) (*TimeMan, error) {
	var raw rawTimeMan
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return mapTimeMan(raw), nil
}

// DecodeProfile maps a Bitrix user.current payload.
func DecodeProfile(data []byte, portal string) (*Profile, error) {
	var user rawUser
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return mapProfile(user, portal), nil
}

func normalizeTimeManStatus(status string) TimeManStatus {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "OPENED":
		return TimeManOpened
	case "PAUSED":
		return TimeManPaused
	case "CLOSED":
		return TimeManClosed
	case "EXPIRED":
		return TimeManExpired
	default:
		return TimeManUnknown
	}
}
