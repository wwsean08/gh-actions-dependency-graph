package scanner

import "errors"

var (
	NotApplicableError          = errors.New("scanner is not applicable to to this action")
	UnexpectedResponseCodeError = errors.New("received an unexpected http response code")
	LikelyShaError              = errors.New("reference is likely a sha-1 sum")
	LocalActionError            = errors.New("action is local, and so is not referred to by a git ref")
)
