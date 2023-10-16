package common

import (
	"fmt"
	"strings"
)

type Step struct {
	Uses *string `yaml:"uses"`
	Name *string `yaml:"name"`
	Id   *string `yaml:"id"`
}

// GetRepo parses out the repo from the uses statement if there is one.
func (s *Step) GetRepo() string {
	if s.Uses == nil || *s.Uses == "" {
		return ""
	}
	uses := *s.Uses

	// this is a local action within the current repo
	if strings.HasPrefix(uses, "./") {
		return ""
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

// GetGitRef gets the commitish (sha, branch, or tag) for an action
// if there is a "uses" statement for a remote repo
func (s *Step) GetGitRef() string {
	if s.Uses == nil || *s.Uses == "" {
		return ""
	}
	if strings.ContainsRune(*s.Uses, '@') {
		return strings.Split(*s.Uses, "@")[1]
	}
	return ""
}

// ParseUses parses the step's uses block and returns the git repo, the repo path,
// the git ref, and an error if something goes wrong parsing it.
func (s *Step) ParseUses() (repo string, path string, ref string, err error) {
	if s.Uses == nil || *s.Uses == "" {
		err = fmt.Errorf("uses statement is null or blank, nothing to parse")
		return
	}

	repo = s.GetRepo()
	path = s.GetRepoPath()
	ref = s.GetGitRef()
	return
}
