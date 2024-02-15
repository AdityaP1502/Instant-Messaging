package unauthenticated

type errorType string

const (
	UserMarkedInActive errorType = "user_marked_inactive"
	InvalidCredentials errorType = "invalid_credentials"
)
