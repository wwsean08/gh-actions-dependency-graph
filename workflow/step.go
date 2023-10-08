package workflow

import (
	"fmt"
	"strings"
)

// GetRepo parses out the repo from the uses statement if there is one.
func (s *Step) GetRepo() string {
	if s.Uses == nil || *s.Uses == "" {
		return ""
	}
	uses := *s.Uses

	// this is a local action within the current repo
	if strings.HasPrefix(uses, "./") {
		// TODO: Figure out how to handle this situation
	}

	uses = strings.Split(uses, "@")[0]
	split := strings.Split(uses, "/")
	return fmt.Sprintf("%s/%s", split[0], split[1])
}

// GetRepoPath gets the path within a repo where the action.yaml is stored,
// as a single repo can contain more than one action
func (s *Step) GetRepoPath() string {
	if s.Uses == nil || *s.Uses == "" {
		return ""
	}

	uses := *s.Uses
	// It's in the local repo, so give the path
	if strings.HasPrefix(uses, "./") {
		return uses
	}

	before, _, _ := strings.Cut(uses, "@")
	repoInfoSplit := strings.Split(before, "/")
	if len(repoInfoSplit) == 2 {
		return ""
	}
	return strings.Join(repoInfoSplit[2:], "/")
}

// GetCommitish gets the commitish (sha, branch, or tag) for an action
// if there is a "uses" statement for a remote repo
func (s *Step) GetCommitish() string {
	if s.Uses == nil || *s.Uses == "" {
		return ""
	}
	if strings.ContainsRune(*s.Uses, '@') {
		return strings.Split(*s.Uses, "@")[1]
	}
	return ""
}
