package scanner

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/wwsean08/actions-dependency-graph/pkg/action"
)

type Scanner struct {
	ScanNodeVersionEOL bool
}

type Results struct {
	Action         string
	NodeVersionEOL string
}

func NewDefaultScanner() *Scanner {
	return &Scanner{
		ScanNodeVersionEOL: true,
	}
}

func (s *Scanner) Scan(action *action.Action) (results *Results, errs []error) {
	results = new(Results)
	results.Action = fmt.Sprintf("%s/%s@%s", action.Repo, action.Path, action.Ref)
	if action.Path == "" {
		results.Action = fmt.Sprintf("%s@%s", action.Repo, action.Ref)
	}
	if s.ScanNodeVersionEOL {
		eol, err := s.IsNodeVersionEOL(action)
		if err != nil && !errors.Is(err, NotApplicableError) {
			errs = append(errs, errors.Wrap(err, fmt.Sprintf("failed to determine if node version is EOL for %s/%s@%s", action.Repo, action.Path, action.Ref)))
		} else if errors.Is(err, NotApplicableError) {
			results.NodeVersionEOL = "Not Applicable"
		} else {
			results.NodeVersionEOL = fmt.Sprintf("%t", eol)
		}
	}
	return results, nil
}

func (s *Scanner) FormatResults(results *Results) string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("Node Version EOL: ")
	if s.ScanNodeVersionEOL {
		sb.WriteString(fmt.Sprintf("%s\n", results.NodeVersionEOL))
	} else {
		sb.WriteString("Not Scanned\n")
	}
	return sb.String()
}
