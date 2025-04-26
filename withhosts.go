package withhosts

import (
	"fmt"
	"strings"

	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/tomrplummer/with-hosts/caddyvar"
)

type FromHostsAdapter struct{}

func init() {
	caddyconfig.RegisterAdapter("caddyfile_withhosts", FromHostsAdapter{})
}

func (FromHostsAdapter) Adapt(raw []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	vars := map[string]string{}

	localWarnings := []caddyconfig.Warning{}

	lines := strings.Split(string(raw), "\n")
	var transformed []string

	for _, line := range lines {
		if caddyvar.IsCaddyVar(line) {
			cvar := caddyvar.New(line)
			vars[cvar.Name] = cvar.Value
		} else if isFromHostsIdentifier(line) {
			tag, rest := parseTag(line)

			domains := getMatchingDomains(vars["hosts_filepath"], tag)
			localHosts := getEtcHostsEntries(vars["etc_hosts_filepath"])

			localWarnings = warnHostFilesMismatch(localHosts, domains)

			newLine := strings.Join(append(domains, rest...), " ")
			transformed = append(transformed, newLine)
		} else {
			transformed = append(transformed, line)
		}
	}

	adapted, warnings, err := caddyconfig.GetAdapter("caddyfile").Adapt([]byte(strings.Join(transformed, "\n")), options)

	warnings = append(warnings, localWarnings...)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to adapt expanded Caddyfile: %v", err)
	}

	return adapted, warnings, nil
}
