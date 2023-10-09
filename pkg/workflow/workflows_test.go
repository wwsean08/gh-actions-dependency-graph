package workflow

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseWorkflowReturnsSimpleWorkflow(t *testing.T) {
	workflow, err := ParseWorkflow("testdata/simple.yaml")
	require.NoError(t, err)
	require.NotNil(t, workflow)

	// validate the jobs are parsed at the top level
	assert.Len(t, workflow.Jobs, 1)
	assert.Contains(t, workflow.Jobs, "test")

	// validate the jobs have their steps
	assert.Len(t, workflow.Jobs["test"].Steps, 3)

	// Validate the jobs are what we expect
	assert.Equal(t, "actions/checkout@v4", *workflow.Jobs["test"].Steps[0].Uses)
	assert.Equal(t, "actions/setup-go@v4", *workflow.Jobs["test"].Steps[1].Uses)
	assert.Nil(t, workflow.Jobs["test"].Steps[2].Uses)
	assert.Equal(t, "no-uses", *workflow.Jobs["test"].Steps[2].Id)
	assert.Equal(t, "Run tests", *workflow.Jobs["test"].Steps[2].Name)
}

func TestParseWorkflowReturnsErrorOnFileNotFound(t *testing.T) {
	workflow, err := ParseWorkflow("testdata/non-existent-file.404")
	require.Error(t, err)
	require.Nil(t, workflow)
}
