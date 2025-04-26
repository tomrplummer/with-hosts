package caddyvar

import "strings"

type CaddyVar struct {
	Name  string
	Value string
}

func New(line string) CaddyVar {
	return parse(line)
}

func parse(line string) CaddyVar {
	fields := strings.Fields(line)

	return CaddyVar{
		Name:  strings.Trim(fields[0], "@"),
		Value: fields[1],
	}
}

func IsCaddyVar(line string) bool {
	if len(line) < 4 {
		return false
	}

	fields := strings.Fields(line)
	if !(len(fields) >= 2) || !strings.HasPrefix(fields[0], "@") {
		return false
	}

	return true
}
