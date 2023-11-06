package scanner

import (
	"github.com/wwsean08/actions-dependency-graph/pkg/action"
)

// IsNodeVersionEOL checks the version of node and verifies that it is not EOL
func (s *Scanner) IsNodeVersionEOL(action *action.Action) (bool, error) {
	if !action.IsJavascript() {
		return false, NotApplicableError
	}
	version, err := action.GetNodeVersion()
	if err != nil {
		return false, err
	}

	return version < 20, nil
}
