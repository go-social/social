package facebook

import (
	"errors"
	"net/url"
	"strings"

	"github.com/go-social/social/providers"
	fb "github.com/huandu/facebook"
)

var errNodeTypePage = errors.New("Tried accessing a nonexisting field on a node type Page")

// Facebook error codes
// https://developers.facebook.com/docs/graph-api/using-graph-api/v2.2#errors
func providerError(err error) error {
	if err == nil {
		return nil
	}

	// Most probably oauth2 error.
	if e, ok := err.(*url.Error); ok {
		if strings.Contains(strings.ToLower(e.Error()), "unauthorized") {
			return providers.ErrAuthFailed
		}
		return providers.ErrUnknown.Err(err)
	}

	if e, ok := err.(*fb.Error); ok {
		switch e.Code {
		case 1, 2:
			return providers.ErrProviderDown

		case 4, 17, 341:
			return providers.ErrHitRateLimit

		case 10:
			return providers.ErrMustReauth

		case 100:
			// This happens when the user tries to get an e-mail from a page.
			// "(#100) Tried accessing nonexisting field (email) on node type (Page)"
			if strings.Contains(e.Error(), "Page") {
				return errNodeTypePage
			}
		case 102, 190:
			switch e.ErrorSubcode {
			case 458, 459, 460:
				return providers.ErrMustReauth
			case 463:
				return providers.ErrExpiredToken
			case 464:
				return providers.ErrBadAccount
			default:
				return providers.ErrInvalidToken
			}

		case 506:
			return providers.ErrDuplicatePost

		case 803:
			return providers.ErrUsernameSearch
		}
	}

	return providers.ErrUnknown.Err(err)
}
