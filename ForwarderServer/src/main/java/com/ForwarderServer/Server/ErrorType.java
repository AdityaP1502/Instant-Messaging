package com.ForwarderServer.Server;

public enum ErrorType {
	InvalidMessageFormat("Missing Header or Payload"),
	InvalidRequestTypeEntryName("Request Type field name is invalid"),
	InvalidUID("Request UID isn't valid"),
	InvalidRequestType("Request Type is unrecognized"), 
	InvalidPayloadFormat("Payload format is invalid"), 
	InvalidCheckInPayloadEntryName("Entry name in the payload is invalid"), 
	InvalidCheckInUsername("Username in the field already been used and not match with this channel"),
	InvalidCheckInProtocol("This channel has been registered to the server"),
	InvalidFetchProtocol("Username in payload hasn't registered to the channel"), 
	InvalidReadyProtocol("Username in payload hasn't registered to the channel"),
	InvalidSendMessagePayloadLength("Payload length is invalid"),
	InvalidSendMessagePayloadEntryName("Entry name in the payload is invalid"),
	InvalidSenderUsername("Username in the field is not match with the channel"),
	InvalidRecipientUsername("Recipient username is invalid or isn't online"), 
	InvalidInitiateCallPacketFormat("Packet format is invalid"), 
	ServerError("Request failed to be processed. Server Error");
	private final String errorDescription;
	
	private ErrorType(String errDescription) {
		this.errorDescription = errDescription;
	}
	
	public String getErrorDescription() {
		return errorDescription;
	}
	
}
