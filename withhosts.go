package withhosts

import (
	"fmt"
	"os"
	"strings"

	// "github.com/caddyserver/caddy/caddyhttp"

	"github.com/caddyserver/caddy/v2/caddyconfig"
)

type FromHostsAdapter struct{}

func init() {
	caddyconfig.RegisterAdapter("caddyfile_withhosts", FromHostsAdapter{})
}
func (FromHostsAdapter) Adapt(raw []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	lines := strings.Split(string(raw), "\n")
	var transformed []string
	for _, line := range lines {
		if strings.HasPrefix(line, "from_hosts") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				domains, err := getMatchingDomains(fields[1])
				if err != nil {
					return nil, nil, fmt.Errorf("matching error: %v", err)
				}

				newLine := strings.Join(append(domains, fields[2:]...), " ")
				transformed = append(transformed, newLine)
			}
		} else {
			transformed = append(transformed, line)
		}
	}

	fmt.Println(strings.Join(transformed, "\n"))

	adapted, warnings, err := caddyconfig.GetAdapter("caddyfile").Adapt([]byte(strings.Join(transformed, "\n")), options)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to adapt expanded Caddyfile: %v", err)
	}

	return adapted, warnings, nil
}

func getMatchingDomains(tag string) ([]string, error) {
	contents, err := os.ReadFile("./hosts")
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
