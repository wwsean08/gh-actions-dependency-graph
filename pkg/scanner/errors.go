package scanner

import "errors"

var (
	NotApplicableError          = errors.New("scanner is not applicable to to this action")
	UnexpectedResponseCodeError = errors.New("received an unexpected http response code")
)
