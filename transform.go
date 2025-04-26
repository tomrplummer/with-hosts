package withhosts

import (
	"os"
	"slices"
	"strings"

	"github.com/caddyserver/caddy/v2/caddyconfig"
)

const FROM_HOSTS_IDENTIFIER_TEXT = "from_hosts"

func isFromHostsIdentifier(line string) bool {
	// min - FROM_HOSTS_IDENTIFIER_TEXT 1 {
	if len(line) < (len(FROM_HOSTS_IDENTIFIER_TEXT) + 4) {
		return false
	}

	fields := strings.Fields(line)

	if len(fields) < 3 {
		return false
	}

	return fields[0] == FROM_HOSTS_IDENTIFIER_TEXT
}

func parseTag(line string) (string, []string) {
	fields := strings.Fields(line)

	return fields[1], fields[2:]
}

func parseEtcHostsEntry(line string) (string, bool) {
	if len(line) > 0 && !strings.HasPrefix(line, "#") {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			// etcHosts = append(etcHosts, fields[1])
			return fields[1], true
		}
	}

	return "", false
}

func parseHostsEntry(line, tag string) (string, bool) {
	line = strings.TrimSpace(line)
	if len(line) == 0 || strings.HasPrefix(line, "#") {
		return "", false
	}

	fields := strings.Fields(line)
	if len(fields) >= 2 && fields[1] == tag {
		return fields[0], true
	}

	return "", false
}

func getMatchingDomains(path, tag string) []string {
	// contents, err := os.ReadFile("/etc/caddy/hosts")
	matches := []string{}
	contents, err := os.ReadFile(path)
	if err != nil {
		return matches
	}

	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		if entry, ok := parseHostsEntry(line, tag); ok {
			matches = append(matches, entry)
		}
	}

	return matches
}

func getEtcHostsEntries(path string) []string {
	etcHosts := []string{}

	contents, err := os.ReadFile(path)
	if err != nil {
		return etcHosts
	}

	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		if entry, ok := parseEtcHostsEntry(line); ok {
			etcHosts = append(etcHosts, entry)
		}
	}

	return etcHosts
}

func warnHostFilesMismatch(localHosts, domains []string) []caddyconfig.Warning {
	warnings := []caddyconfig.Warning{}
	for _, domain := range domains {
		if !slices.Contains[[]string](localHosts, domain) {
			warnings = append(warnings, caddyconfig.Warning{
				Message: "/etc/hosts is missing " + domain,
			})
		}
	}

	return warnings
}
