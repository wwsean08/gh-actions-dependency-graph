package scanner

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScanner_isRepoSusceptibleToRepoJacking(t *testing.T) {
	scan := Scanner{}
	type testRun struct {
		Repo     string
		APIData  string
		Expected bool
	}

	testData := []testRun{
		{
			Repo:     "team-TAU/tau",
			APIData:  "{\"id\":338929421,\"node_id\":\"MDEwOlJlcG9zaXRvcnkzMzg5Mjk0MjE=\",\"name\":\"tau\",\"full_name\":\"Team-TAU/tau\",\"private\":false,\"owner\":{\"login\":\"Team-TAU\",\"id\":86118136}}",
			Expected: false,
		},
		{
			Repo:     "FiniteSingularity/tau",
			APIData:  "{\"id\":338929421,\"node_id\":\"MDEwOlJlcG9zaXRvcnkzMzg5Mjk0MjE=\",\"name\":\"tau\",\"full_name\":\"Team-TAU/tau\",\"private\":false,\"owner\":{\"login\":\"Team-TAU\",\"id\":86118136}}",
			Expected: true,
		},
	}

	for _, run := range testData {
		repo := new(ghRepo)
		err := json.Unmarshal([]byte(run.APIData), repo)
		require.NoError(t, err)
		assert.Equal(t, run.Expected, scan.isRepoSusceptibleToRepoJacking(run.Repo, repo))
	}
}
