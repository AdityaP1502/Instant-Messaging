CREATE TABLE account(
	account_id SERIAL PRIMARY KEY,
	username CHAR(64) NOT NULL
	name CHAR(64) NOT NULL,
	password CHAR(64) NOT NULL,
	email CHAR(64) NOT NULL,
	is_active bool NOT NULL
);

CREATE TABLE invalid_token(
	token_id SERIAL PRIMARY KEY,
	account_id INT references account(account_id),
	token VARCHAR(1024) NOT NULL,
	type CHAR(16) NOT NULL,
	expired_at TIMESTAMP NOT NULL
);

CREATE TABLE user_otp(
	account_id INT references account(account_id),
	otp CHAR(16) NOT NULL,
	last_resend TIMESTAMP NOT NULL,
	marked_for_deletion bool NOT NULL
);
