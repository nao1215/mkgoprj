package cmdinfo

import (
	"fmt"
)

const (
	name    = "ubume"
	version = "2.0.0"
)

// Version return ubume command version.
func Version() string {
	return fmt.Sprintf("%s version %s (under Apache License version 2.0)",
		Name(), version)
}

// Name return command name.
func Name() string {
	return name
}
