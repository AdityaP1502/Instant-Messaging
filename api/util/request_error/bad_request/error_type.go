package badrequest

type errorType string

const (
	MissingParameter     errorType = "missing_parameter"
	HeaderValueMistmatch errorType = "header_value_mismatch"
	UsernameExists       errorType = "username_exists"
	EmailExists          errorType = "email_exists"
	UsernameInvalid      errorType = "username_invalid"
	PasswordWeak         errorType = "password_weak"
	EmailInvalid         errorType = "invalid_email"
	PayloadInvalid       errorType = "invalid_payload"
	OTPInvalid           errorType = "invalid_otp"
)
