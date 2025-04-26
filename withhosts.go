package withhosts

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/caddyserver/caddy/v2/caddyconfig"
)

type FromHostsAdapter struct{}

func init() {
	caddyconfig.RegisterAdapter("caddyfile_withhosts", FromHostsAdapter{})
}

func (FromHostsAdapter) Adapt(raw []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	localwarnings := []caddyconfig.Warning{}

	vars := map[string]string{}

	lines := strings.Split(string(raw), "\n")
	var transformed []string

	for _, line := range lines {
		if strings.HasPrefix(line, "@") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				vars[strings.Trim(fields[0], "@")] = fields[1]
			}
		} else if strings.HasPrefix(line, "from_hosts") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				hostsPath := vars["hosts_filepath"]
				domains, err := getMatchingDomains(hostsPath, fields[1])
				if err != nil {
					return nil, nil, fmt.Errorf("matching error: %v", err)
				}

				etcPath := vars["etc_hosts_filepath"]
				localHosts, _ := getEtcHostsEntries(etcPath)

				for _, domain := range domains {
					if !slices.Contains[[]string](localHosts, domain) {
						localwarnings = append(localwarnings, caddyconfig.Warning{
							Message: "/etc/hosts is missing " + domain,
						})
					}
				}

				newLine := strings.Join(append(domains, fields[2:]...), " ")
				transformed = append(transformed, newLine)
			}
		} else {
			transformed = append(transformed, line)
		}
	}

	adapted, warnings, err := caddyconfig.GetAdapter("caddyfile").Adapt([]byte(strings.Join(transformed, "\n")), options)

	warnings = append(warnings, localwarnings...)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to adapt expanded Caddyfile: %v", err)
	}

	return adapted, warnings, nil
}

func getMatchingDomains(path, tag string) ([]string, error) {
	// contents, err := os.ReadFile("/etc/caddy/hosts")
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read hosts: %v", err)
	}

	var matches []string
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == tag {
			matches = append(matches, fields[0])
		}
	}

	return matches, nil
}

func getEtcHostsEntries(path string) ([]string, error) {
	etcHosts := []string{}
	// contents, err := os.ReadFile("/hostmachine/hosts")
	contents, err := os.ReadFile(path)
	if err != nil {
		return etcHosts, err
	}

	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				etcHosts = append(etcHosts, fields[1])
			}
		}
	}

	return etcHosts, nil
}
