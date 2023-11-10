package scanner

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/wwsean08/actions-dependency-graph/pkg/action"
)

type Scanner struct {
	ScanNodeVersionEOL bool
	RepoJackingScan    bool
	LatestTagScan      bool
}

type Results struct {
	Action                   string
	NodeVersionEOL           string
	RepoJackingPossible      string
	RepoJackingComment       string
	LatestTagIsCurrentLatest bool
	LatestTagTag             string
	LatestTagError           string
}

func NewDefaultScanner() *Scanner {
	return &Scanner{
		ScanNodeVersionEOL: true,
		RepoJackingScan:    true,
		LatestTagScan:      true,
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

	if s.RepoJackingScan {
		rjResult := s.RepoJackingCheck(action.Repo)
		if rjResult.err != nil {
			results.RepoJackingPossible = "unknown"
			results.RepoJackingComment = rjResult.err.Error()
		} else if rjResult.susceptible {
			results.RepoJackingPossible = "true"
			results.RepoJackingComment = fmt.Sprintf("Repo jacking possible.  %s has been moved and is now accessible via %s", rjResult.original, rjResult.new)
		} else {
			results.RepoJackingPossible = "false"
		}
	}

	if s.LatestTagScan {
		var err error
		results.LatestTagIsCurrentLatest, results.LatestTagTag, err = s.CheckActionVersionIsLatest(action)
		if err != nil {
			results.LatestTagError = err.Error()
		}
	}

	return results, errs
}

func (s *Scanner) FormatResults(results *Results) string {
	sb := strings.Builder{}
	_, _ = sb.WriteString("Node Version EOL: ")
	if s.ScanNodeVersionEOL {
		sb.WriteString(fmt.Sprintf("%s\n", results.NodeVersionEOL))
	} else {
		sb.WriteString("Not Scanned\n")
	}

	sb.WriteString("Repo Jacking Possible: ")
	if s.RepoJackingScan {
		sb.WriteString(fmt.Sprintf("%s\n", results.RepoJackingPossible))
		if results.RepoJackingComment != "" {
			sb.WriteString(fmt.Sprintf("%s\n", results.RepoJackingComment))
		}
	} else {
		sb.WriteString("Not Scanned\n")
	}

	sb.WriteString(fmt.Sprintf("Is Referencing Latest Tag: %t\n", results.LatestTagIsCurrentLatest))
	if results.LatestTagTag != "" {
		sb.WriteString(fmt.Sprintf("The latest tag for %s is %s\n", results.Action, results.LatestTagTag))
	}
	if results.LatestTagError != "" {
		sb.WriteString(fmt.Sprintf("Error getting latest tag: %s\n", results.LatestTagError))
	}

	return sb.String()
}
