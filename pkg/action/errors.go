package action

// NoDependenciesError represents an error which can generally be ignored
// (though could be useful information), when an action is not dependent
// on other actions
type NoDependenciesError struct {
	msg string
}

// Error implementing the error interface
func (e NoDependenciesError) Error() string {
	return e.msg
}
