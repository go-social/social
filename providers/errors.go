package providers

import "fmt"

// TODO: redo the errors.. use pkg/errors
// review box errors thing or other stuff in render..

// TODO: We need iota-style error codes, so we can catch'em in the below layers easily.
var (
	ErrUnknownProviderID = &Error{Code: 1, Msg: "unknown provider id"}

	// Authorization
	ErrNoCredentials = &Error{Code: 1000, Msg: "missing provider credentials, please connect your social account first"}
	ErrAuthFailed    = &Error{Code: 1001, Msg: "provider authorization failed, please re-connect your social account"}
	ErrInvalidToken  = &Error{Code: 1002, Msg: "invalid provider token, please re-connect your social account"}
	ErrExpiredToken  = &Error{Code: 1003, Msg: "expired provider token, please re-connect your social account"}
	ErrHitRateLimit  = &Error{Code: 1004, Msg: "hit token rate limit"}
	ErrBadAccount    = &Error{Code: 1005, Msg: "disabled account"}
	ErrMustReauth    = &Error{Code: 1006, Msg: "authentication error, please re-connect your social account"}
	ErrGetUser       = &Error{Code: 1007, Msg: "unable to fetch user profile"}
	ErrEmptyCode     = &Error{Code: 1008, Msg: "empty code in callback"}

	// Queries
	ErrInvalidQuery      = &Error{Code: 2000, Msg: "invalid request query"}
	ErrNoQueryAccess     = &Error{Code: 2001, Msg: "provider does not have access for this query"}
	ErrInvalidAsset      = &Error{Code: 2002, Msg: "invalid asset"}
	ErrWritingPost       = &Error{Code: 2003, Msg: "unable to post asset"}
	ErrDuplicatePost     = &Error{Code: 2004, Msg: "duplicate post"}
	ErrUsernameSearch    = &Error{Code: 2005, Msg: "provided doesn't allow @username searches, @page (brand) searches work"}
	ErrUnauthorizedQuery = &Error{Code: 2006, Msg: "user unauthorized to make this query"}

	// Everything else
	ErrUnknown        = &Error{Code: 5000, Msg: "unknown provider error"}
	ErrProviderDown   = &Error{Code: 5001, Msg: "provider is down"}
	ErrUnsupported    = &Error{Code: 5002, Msg: "unsupported operation"}
	ErrNotImplemented = &Error{Code: 5003, Msg: "not implemented"}
	ErrInvalidContent = &Error{Code: 5004, Msg: "empty title and url provided"}
)

// Provider-specific error
type Error struct {
	err  error  // the original error
	Code int    // provider error code
	Msg  string // provider error string
}

func (e *Error) Error() string {
	s := fmt.Sprintf("%d - %s", e.Code, e.Msg)
	if e.err != nil {
		return s + ": " + e.err.Error()
	}
	return s
}

func (e *Error) Err(err error) error {
	return &Error{err: err, Code: e.Code, Msg: e.Msg}
}
