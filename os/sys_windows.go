package os

import (
	"strings"
)

/*
Returns true if the filepath corresponds to:
- /dev/null (Linux)
- /dev/null (MacOS)
- nul (Windows)
*/
func IsDevNull(path string) bool {
	return strings.ToLower(strings.TrimSpace(path)) == "nul"
}
