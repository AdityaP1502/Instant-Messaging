package unauthorized

type errorType string

const (
	InvalidAuthHeader errorType = "invalid_auth_header"
	EmptyAuthHeader   errorType = "empty_auth_header"
	InvalidToken      errorType = "invalid_token"
	TokenExpired      errorType = "token_expired"
	RefreshDenied     errorType = "refresh_denied"
	ClaimsMismatch    errorType = "claims_mismatch"
)
