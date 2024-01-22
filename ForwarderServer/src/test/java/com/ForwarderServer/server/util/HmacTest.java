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
		String username = "aditya";
		String accessToken = CryptoHandler.getAccessToken("adityafa8f5481591bae0cb439671abecb76eb2473f457dc445fe197f869dcc159ed37", "0abfb985fe6a95cfdce9429f42e43c9c012387470b8bd352f0ee2bca0d7ba39e6bd58231c87a0d7d4c375c8c0cc38d37d1ccbafe173bcef432885d339dfb820b28fcb5db591500a44e7edea426077e454eb8f4db89233b7061fa78f0d5529d04251eea4d24ae65f0cbfd854aaf186cad2316f47492a1acc001e641d86f68e8f82dbd0d33fa6c8b86ff1cdf1bb5fa601aa263c89c443cb581266825f93b626afb9bbcf603149585bd2d273fa747730ceb1b1c4565c5d12741c91e22eacd702fe5d722978cd0d0eb970b72199276ee1f85afb54b2232d4dbb82f8cf3e4c6e90f4652a337b3c20dcadd95b736bc0960234735c1b90a77e142a8313a8a9058fde15b");
		
		boolean isTokenEqual = CryptoHandler.verifyAccessToken(username, "fa8f5481591bae0cb439671abecb76eb2473f457dc445fe197f869dcc159ed37", accessToken, "0abfb985fe6a95cfdce9429f42e43c9c012387470b8bd352f0ee2bca0d7ba39e6bd58231c87a0d7d4c375c8c0cc38d37d1ccbafe173bcef432885d339dfb820b28fcb5db591500a44e7edea426077e454eb8f4db89233b7061fa78f0d5529d04251eea4d24ae65f0cbfd854aaf186cad2316f47492a1acc001e641d86f68e8f82dbd0d33fa6c8b86ff1cdf1bb5fa601aa263c89c443cb581266825f93b626afb9bbcf603149585bd2d273fa747730ceb1b1c4565c5d12741c91e22eacd702fe5d722978cd0d0eb970b72199276ee1f85afb54b2232d4dbb82f8cf3e4c6e90f4652a337b3c20dcadd95b736bc0960234735c1b90a77e142a8313a8a9058fde15b");
		
		System.out.println(isTokenEqual);
		
		assert isTokenEqual : "Wrong Equality";
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
