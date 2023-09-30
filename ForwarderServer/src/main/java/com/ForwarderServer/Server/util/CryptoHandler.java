package com.ForwarderServer.server.util;

import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

public class CryptoHandler {
	private static final String ALGORITHM = "HmacSHA256";
	
	public static byte[] hexStringToByteArray(String hexString) {
		if ((hexString.length() & 1) == 1) {
			throw new IllegalArgumentException("Hex string length must be a multiple of 2");
		}
		
		byte f, g;
		
		int length = (int) hexString.length() / 2;
		byte[] bytearray = new byte[length];
		int c = 0;
		
		for (int i = 0; i < length; i++)
		{
			f = (byte) Character.digit(hexString.charAt(c), 16);
			g = (byte) Character.digit(hexString.charAt(c + 1), 16);
			bytearray[i] = (byte) ((f << 4) + g);
			c += 2;
		}
		
		return bytearray;
	}
	
	public static String bytesToHexString(byte[] bytearray) {
		int t;
		StringBuilder hex = new StringBuilder();
		
		for (int i = 0; i < bytearray.length; i++)
		{
			t = ((int) bytearray[bytearray.length - (i + 1)]) & 0xff;
			hex.insert(0, Integer.toHexString(t));
			
			if (t < 16)
			{
				// toHexString omit leading zero
				hex.insert(0, "0");
			}
		}
		
		return hex.toString();
	}
	
	public static String getAccessToken(String sessionKey, String key) throws InvalidKeyException {
		/**
		 * Get access token used to allocate channel in audio server
		 * Algorithm used is SHA-256
		 */
		
		byte[] keyByteArray = hexStringToByteArray(key);
		SecretKeySpec secretkeyspec = new SecretKeySpec(keyByteArray, ALGORITHM);
		
		try {
			Mac mac = Mac.getInstance(ALGORITHM);
			mac.init(secretkeyspec);
			return bytesToHexString(mac.doFinal(sessionKey.getBytes()));
		} catch (NoSuchAlgorithmException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		return "";
	}
}
