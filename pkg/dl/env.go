package dl

import (
	"os"
	"path/filepath"
	"strings"
)

func execEnv(isolatedHome, binDir, dockerConfig string, parent []string) []string {
	env := filterEnv(parent, "HOME", "PATH", "DOCKER_CONFIG")
	path := envValue(parent, "PATH")
	env = append(env,
		"HOME="+isolatedHome,
		"PATH="+binDir+string(os.PathListSeparator)+path,
	)
	if dockerConfig != "" {
		env = append(env, "DOCKER_CONFIG="+dockerConfig)
	}
	return env
}

func filterEnv(parent []string, drop ...string) []string {
	skip := make(map[string]struct{}, len(drop))
	for _, key := range drop {
		skip[key] = struct{}{}
	}

	out := make([]string, 0, len(parent))
	for _, entry := range parent {
		key, _, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		if _, found := skip[key]; found {
			continue
		}
		out = append(out, entry)
	}
	return out
}

func envValue(parent []string, key string) string {
	for i := len(parent) - 1; i >= 0; i-- {
		name, value, ok := strings.Cut(parent[i], "=")
		if ok && name == key {
			return value
		}
	}
	return ""
}

func dockerConfigDir(parent []string, userHome string) string {
	if value := envValue(parent, "DOCKER_CONFIG"); value != "" {
		return value
	}
	if userHome == "" {
		return ""
	}
	return filepath.Join(userHome, ".docker")
}
