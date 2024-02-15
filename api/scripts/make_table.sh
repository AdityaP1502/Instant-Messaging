CREATE TABLE account(
	account_id SERIAL PRIMARY KEY,
	username VARCHAR(64) NOT NULL,
	name VARCHAR(64) NOT NULL,
	password vARCHAR(80) NOT NULL,
    password_salt VARCHAR(64) NOT NULL,
	email CHAR(64) NOT NULL,
	is_active VARCHAR(6) NOT NULL
);

CREATE TABLE invalid_token(
	token_id SERIAL PRIMARY KEY,
	username VARCHAR(64) references account(username) ON DELETE CASCADE,
	token VARCHAR(1024) NOT NULL,
	type CHAR(16) NOT NULL,
	expired_at TIMESTAMP NOT NULL
);

CREATE TABLE user_otp(
	username  VARCHAR(64) references account(username) ON DELETE CASCADE,
    otp_confirmation_id varchar(512) NOT NULL, 
	otp CHAR(6) NOT NULL,
	last_resend TIMESTAMP NOT NULL,
	marked_for_deletion VARCHAR(6) NOT NULL
);
