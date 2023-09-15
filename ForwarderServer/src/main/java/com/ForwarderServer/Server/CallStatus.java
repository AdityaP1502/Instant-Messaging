package com.ForwarderServer.Server;

public enum CallStatus {
	INCOMINGCALL(ResponseHeader.INCOMINGCALL), 
	CALLACCEPTED(ResponseHeader.CALLACCEPTED), 
	CALLDECLINED(ResponseHeader.CALLDECLINED), 
	CALLTERMINATE(ResponseHeader.CALLTERMINATE),
	CALLTIMEOUT(ResponseHeader.CALLTIMEOUT);

	private ResponseHeader header;

	private CallStatus(ResponseHeader header) {
		this.header = header;
	}

	public ResponseHeader getHeader() {
		return header;
	}
}
