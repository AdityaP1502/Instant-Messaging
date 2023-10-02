package com.ForwarderServer.server;

import java.security.InvalidKeyException;
import com.ForwarderServer.server.util.CryptoHandler;


public class CallPayload {
	public static String generateChannelAllocPayload(String uid, String username, String key, String IPaddr, int port) throws InvalidKeyException
	{
		String[] token = CryptoHandler.getSessionToken(username);
		String accessToken = CryptoHandler.getAccessToken(token[0], key);
		
		// Send the audio server IP for the server to connect to
		// TODO : For scaling purposes, integrate with cloud load balancer
		String networkAddress = String.format("%s:%d", IPaddr, port);
		String payload = String.format("uid=%s;token=%s;salt=%s;NetworkAddress=%s", uid, accessToken, token[1], networkAddress);
		
		return payload;
	}
	
	public static String generateCallRequestPayload(String uid, String requestPayload, String username, String key, String IPaddr, int port) throws InvalidKeyException
	{
		String channelAllocPayload = generateChannelAllocPayload(uid, username, key, IPaddr, port);
		return requestPayload + channelAllocPayload;
	}
}
