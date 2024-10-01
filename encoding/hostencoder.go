package encoding

import (
	"fmt"
	"strings"
)

type HostnameEncoder struct {
	Zone string
}

const maxNamePartSize = 63
const maxNameSize = 255

var ErrNameSize = fmt.Errorf("name exceeds maximum size, might not be resolved")

func (e HostnameEncoder) EncodeToHostname(data string) (string, error) {
	res := strings.Builder{}

	lastPos := len(data)

	for x := 0; x < lastPos; x += maxNamePartSize {
		res.WriteString(data[x:min(x+maxNamePartSize, lastPos)])
		res.WriteRune('.')
	}

	res.WriteString(e.Zone)

	name := res.String()
	var err error

	if len(name) > maxNameSize {
		err = ErrNameSize
	}

	return name, err
}

func (e HostnameEncoder) DecodeHostname(hname string) string {
	return strings.Replace(
		strings.TrimSuffix(
			strings.TrimSuffix(
				hname,
				e.Zone,
			),
			strings.ToUpper(e.Zone),
		),
		".",
		"",
		-1,
	)
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
