package com.ForwarderServer.server;

import java.security.InvalidKeyException;
import com.ForwarderServer.server.util.CryptoHandler;


public class CallPayload {
	private static String generateChannelAllocPayload(String username, String key, String IPaddr, int port) throws InvalidKeyException
	{
		String[] token = CryptoHandler.getSessionToken(username);
		
		String accessToken = CryptoHandler.getAccessToken(token[0], key);
		
		String networkAddress = String.format("%s:%d", IPaddr, port);
		String payload = String.format("token=%s;salt=%s;NetworkAddress=%s", accessToken, token[1], networkAddress);
		
		return payload;
	}
	
	public static String generateCallRequestPayload(String uid, String requestPayload, String username, String key, String IPaddr, int port) throws InvalidKeyException
	{
		String channelAllocPayload = generateChannelAllocPayload(username, key, IPaddr, port);
		return String.format("uid=%s;%s;%s", uid, requestPayload, channelAllocPayload);
	}
}
