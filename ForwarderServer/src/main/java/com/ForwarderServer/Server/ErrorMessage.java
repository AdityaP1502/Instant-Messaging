package com.ForwarderServer.Server;

public class ErrorMessage {
	private final ErrorType err;
	private final String message;

	public ErrorMessage(ErrorType err, String message) {
		this.err = err;
		this.message = message;
	}

	public String getErrorMessage() {
		if (!message.equals(""))
			return err.name() + ':' + err.getErrorDescription() + "." + message;
		
		return err.name() + ':' + err.getErrorDescription();
	}
}
