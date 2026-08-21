package dl

import (
	"regexp"
	"strings"
	"unicode"
)

const (
	composeServiceLabel = "com.docker.compose.service"

	serviceTraefik   = "traefik"
	servicePortainer = "portainer"
	serviceMail      = "mail"
)

type knownService struct {
	name    string
	aliases []string
}

var infrastructure = []knownService{
	{name: serviceTraefik, aliases: []string{"traefik"}},
	{name: servicePortainer, aliases: []string{"portainer"}},
	{name: serviceMail, aliases: []string{"mail", "mailhog", "mailcatcher", "mailpit"}},
}

var versionTokenRe = regexp.MustCompile(`(?i)v?(\d+\.\d+(?:\.\d+)?)`)

func defaultServices() []ServiceState {
	out := make([]ServiceState, 0, len(infrastructure))
	for _, svc := range infrastructure {
		out = append(out, ServiceState{Name: svc.name})
	}
	return out
}

// DetectServices maps docker containers onto the three dl infrastructure services.
func DetectServices(hints []ContainerHint) []ServiceState {
	out := defaultServices()
	for i, svc := range infrastructure {
		out[i] = matchService(hints, svc)
	}
	return out
}

// IsTraefikRunning reports whether the core traefik service is up.
func IsTraefikRunning(services []ServiceState) bool {
	for _, svc := range services {
		if svc.Name == serviceTraefik {
			return svc.IsRunning
		}
	}
	return false
}

func matchService(hints []ContainerHint, svc knownService) ServiceState {
	state := ServiceState{Name: svc.name}
	for _, hint := range hints {
		if !hintMatches(hint, svc.aliases) {
			continue
		}
		state.IsPresent = true
		if isRunning(hint.State) {
			state.IsRunning = true
			return state
		}
	}
	return state
}

func hintMatches(hint ContainerHint, aliases []string) bool {
	if label := strings.TrimSpace(hint.Labels[composeServiceLabel]); label != "" {
		for _, alias := range aliases {
			if strings.EqualFold(label, alias) {
				return true
			}
		}
	}

	names := make([]string, 0, 1+len(hint.Names))
	if hint.Name != "" {
		names = append(names, hint.Name)
	}
	names = append(names, hint.Names...)
	for _, name := range names {
		for _, alias := range aliases {
			if tokenMatches(name, alias) {
				return true
			}
		}
	}

	base := imageBase(hint.Image)
	for _, alias := range aliases {
		if tokenMatches(base, alias) {
			return true
		}
	}
	return false
}

func tokenMatches(value, alias string) bool {
	value = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(value), "/"))
	alias = strings.ToLower(alias)
	if value == "" || alias == "" {
		return false
	}
	if value == alias {
		return true
	}
	for _, token := range strings.FieldsFunc(value, isNameSeparator) {
		if token == alias {
			return true
		}
	}
	return false
}

func isNameSeparator(r rune) bool {
	return r == '-' || r == '_' || r == '.' || unicode.IsSpace(r)
}

func imageBase(image string) string {
	img := strings.ToLower(strings.TrimSpace(image))
	if i := strings.LastIndex(img, "/"); i >= 0 {
		img = img[i+1:]
	}
	if i := strings.IndexAny(img, ":@"); i >= 0 {
		img = img[:i]
	}
	return img
}

func isRunning(state string) bool {
	return strings.EqualFold(strings.TrimSpace(state), "running")
}

func isUpdateAvailable(current, latest string) bool {
	c := versionToken(current)
	l := versionToken(latest)
	if c == "" || l == "" {
		return false
	}
	return c != l
}

func versionToken(s string) string {
	match := versionTokenRe.FindStringSubmatch(s)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}
