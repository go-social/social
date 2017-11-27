package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/go-social/social/providers"
)

func providerError(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*anaconda.ApiError); ok {
		// twitter errors: https://dev.twitter.com/docs/error-codes-responses
		if len(e.Decoded.Errors) > 0 {
			apiErr := e.Decoded.Errors[0] // oh anaconda..

			switch apiErr.Code {
			case anaconda.TwitterErrorCouldNotAuthenticate:
				return providers.ErrAuthFailed
			case anaconda.TwitterErrorDoesNotExist:
				return providers.ErrInvalidQuery
			case anaconda.TwitterErrorAccountSuspended:
				return providers.ErrBadAccount
			case anaconda.TwitterErrorRateLimitExceeded:
				return providers.ErrHitRateLimit
			case anaconda.TwitterErrorInvalidToken:
				return providers.ErrInvalidToken
			case anaconda.TwitterErrorOverCapacity:
				return providers.ErrProviderDown
			case anaconda.TwitterErrorInternalError:
				return providers.ErrProviderDown
			case anaconda.TwitterErrorCouldNotAuthenticateYou:
				return providers.ErrAuthFailed
			case anaconda.TwitterErrorStatusIsADuplicate:
				return providers.ErrWritingPost
			case anaconda.TwitterErrorBadAuthenticationData:
				return providers.ErrAuthFailed
			case anaconda.TwitterErrorUserMustVerifyLogin:
				return providers.ErrMustReauth
			default:
				return providers.ErrUnknown.Err(e)
			}
		} else {
			// It seems that twitter doesn't return a "code" on the unauthorized error
			if e.StatusCode == 401 {
				return providers.ErrUnauthorizedQuery
			} else {
				return providers.ErrUnknown.Err(err)
			}
		}
	} else {
		return providers.ErrUnknown.Err(err)
	}
}
