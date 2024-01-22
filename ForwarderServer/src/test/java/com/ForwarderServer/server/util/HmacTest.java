package com.ForwarderServer.server.util;

import org.junit.jupiter.api.Test;

import java.security.InvalidKeyException;

import org.junit.jupiter.api.Assertions;

import com.ForwarderServer.server.util.CryptoHandler;

public class HmacTest {
	@Test
	public static void hexToBytesTest()
	{
		byte[] res = CryptoHandler.hexStringToByteArray("FFAAAF0F1F2AB3");
		byte[] expected = {-1, -86, -81, 15, 31, 42, -77};
		Assertions.assertArrayEquals(expected, res);
	}
	
	@Test
	public static void bytesToHexTest()
	{
		byte[] arg = {-1, -86, -81, 15, 31, 42, -77};
		String hex = CryptoHandler.bytesToHexString(arg);
		String expected = "FFAAAF0F1F2AB3";
		Assertions.assertEquals(expected.toLowerCase(), hex);
	}
	
	@Test
	public static void hmacTest()
	{
		String sessionKey = "Aditya1502";
		String key = "FF0A12";
		
		try {
			CryptoHandler.getAccessToken(sessionKey, key);
		} catch (InvalidKeyException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
			System.out.println("Failed");
			System.exit(1);
		}
	}
	
	@Test
	public static void verifyHmacTest() throws InvalidKeyException
	{
		String username = "Aditya";
		String[] sessionToken = CryptoHandler.getSessionToken(username);
		String accessToken = CryptoHandler.getAccessToken(sessionToken[0], "AAFFAAFF") + "AA";
		
		boolean isTokenEqual = CryptoHandler.verifyAccessToken(username, sessionToken[1], accessToken, "AAFFAAFF");
		
		assert !isTokenEqual : "Wrong Equality";
	}
	
	@Test 
	public static void getSessionTokenTest()
	{
		String token = CryptoHandler.getSessionToken("Aditya")[0];
		System.out.println(token);
	}
	
	public static void main(String[] args) throws InvalidKeyException
	{
		hexToBytesTest();
		System.out.println("Passed");
		bytesToHexTest();	
		System.out.println("Passed");
		hmacTest();
		System.out.println("Passed");
		getSessionTokenTest();
		System.out.println("Passed");
		verifyHmacTest();
		System.out.println("Passed");
	}
}
