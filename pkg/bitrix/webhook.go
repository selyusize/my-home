package bitrix

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var userCodeRe = regexp.MustCompile(`^\d+/[A-Za-z0-9]+$`)

// ComposeWebhook builds a portal origin and inbound webhook URL.
//
// domain is a host or https URL. webhook may be a full incoming webhook URL,
// "rest/{userId}/{code}", "{userId}/{code}", or a secret (user id defaults to 1).
func ComposeWebhook(domain, webhook string) (portal, webhookURL string, err error) {
	webhook = strings.TrimSpace(webhook)
	if webhook == "" {
		return "", "", ErrMissingWebhook
	}

	if strings.HasPrefix(webhook, "http://") || strings.HasPrefix(webhook, "https://") {
		parsed, err := url.Parse(webhook)
		if err != nil || parsed.Host == "" {
			return "", "", fmt.Errorf("%w: %s", ErrInvalidWebhook, webhook)
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return "", "", fmt.Errorf("%w: %s", ErrInvalidWebhook, webhook)
		}
		portal = parsed.Scheme + "://" + parsed.Host
		return portal, ensureTrailingSlash(parsed.String()), nil
	}

	portal, err = normalizePortal(domain)
	if err != nil {
		return "", "", err
	}

	path := strings.Trim(webhook, "/")
	switch {
	case strings.HasPrefix(strings.ToLower(path), "rest/"):
		webhookURL = portal + "/" + path
	case userCodeRe.MatchString(path):
		webhookURL = portal + "/rest/" + path
	default:
		if path == "" {
			return "", "", ErrMissingWebhook
		}
		webhookURL = portal + "/rest/1/" + path
	}
	return portal, ensureTrailingSlash(webhookURL), nil
}

func normalizePortal(domain string) (string, error) {
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return "", ErrMissingPortal
	}
	if !strings.Contains(domain, "://") {
		domain = "https://" + domain
	}
	parsed, err := url.Parse(domain)
	if err != nil || parsed.Host == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingPortal, domain)
	}
	scheme := parsed.Scheme
	if scheme != "http" && scheme != "https" {
		return "", fmt.Errorf("%w: %s", ErrMissingPortal, domain)
	}
	return scheme + "://" + parsed.Host, nil
}

func ensureTrailingSlash(raw string) string {
	if strings.HasSuffix(raw, "/") {
		return raw
	}
	return raw + "/"
}

func personalPageURL(portal string, userID int64) string {
	if portal == "" || userID <= 0 {
		return portal
	}
	return fmt.Sprintf("%s/company/personal/user/%d/", strings.TrimRight(portal, "/"), userID)
}
