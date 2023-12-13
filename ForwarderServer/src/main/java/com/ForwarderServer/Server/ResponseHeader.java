package com.ForwarderServer.server;

public enum ResponseHeader {
	ERROR("responsetype=ERROR"),
	FETCH("responsetype=FETCH"), 
	SENDMESSAGE("responsetype=MESSAGE"),
	INCOMINGCALL("responsetype=INCOMING_CALL"),
	CHANNELALLOCATION("responsetype=CHANNEL_ALLOCATION"),
	CALLACCEPTED("responsetype=CALL_ACCEPTED"),
	CALLDECLINED("responsetype=CALL_DECLINED"), 
	CALLTIMEOUT("responsetype=CALL_TIMEOUT"), 
	CALLTERMINATE("responsetype=CALL_TERMINATE"),
	CALLABORT("responsetype=CALL_ABORT"), 
	CALLBUSY("responsetype=CALL_BUSY"),
	SENDAUDIO("responsetype=AUDIO"), 
	OK("responsetype=OK");
	
	private final String header;
	
	ResponseHeader(String header) {
		this.header = header;
	}
	
	public String getHeader() {
		return header;
	}
	
}
