package com.ForwarderServer.server.audio;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.DatagramChannel;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;

import com.ForwarderServer.server.util.CryptoHandler;

import io.github.cdimascio.dotenv.Dotenv;

public class AudioServer {

	private static final int BUFFER_SIZE = 4096; // bytes

	private static AudioChannels channels;
	private static Map<String, SocketAddress> waitingForConnection;

	private static DatagramChannel serverSocket;
	private static String SUPER_SECRET_KEY;
	private static String MODE;
	private static Dotenv dotenv;

	public static void main(String[] args) throws IOException, InvalidKeyException {
		// SET SECRET KEY

		// SUPER_SECRET_KEY = System.getenv("FORWARD_SERVER_SUPER_SECRET_KEY");

		String IPaddress;
		int port;

		dotenv = Dotenv.load();

		SUPER_SECRET_KEY = dotenv.get("SECRET_KEY");
		MODE = dotenv.get("MODE");

		// Network configuration
		IPaddress = MODE.equals("DEV") ? dotenv.get("IP_A1_DEV") : dotenv.get("IP_A1_PROD");
		port = Integer.parseInt(dotenv.get("PORT_A1"));

		Selector selector = Selector.open();
		serverSocket = DatagramChannel.open();

		serverSocket.bind(new InetSocketAddress(IPaddress, port));
		serverSocket.configureBlocking(false);
		serverSocket.register(selector, SelectionKey.OP_READ);

		waitingForConnection = new HashMap<>();
		channels = new AudioChannels();

		while (true) {
			selector.select();
			Set<SelectionKey> selectedKeys = selector.selectedKeys();
			Iterator<SelectionKey> iter = selectedKeys.iterator();

			while (iter.hasNext()) {
				SelectionKey key = iter.next();

				if (key.isReadable()) {
					processIncomingRequest(key);
				}

				iter.remove();
			}
		}
	}

	private static UDPSocketData getSocketData(DatagramChannel s) throws IOException {
		ByteBuffer buffer = ByteBuffer.allocate(BUFFER_SIZE);

		SocketAddress a = s.receive(buffer);
		buffer.flip();
		byte[] message = new byte[buffer.remaining()];
		buffer.get(message);

		return new UDPSocketData(a, message);
	}

	private static int readActionInfo(byte[] words, int start) {
		// assuming the data is sorted in big endian order
		int bitmask = 0xFF;

		int firstByte = (words[start] & bitmask) << 24;
		int secondByte = (words[start + 1] & bitmask) << 16;
		int thirdByte = (words[start + 2] & bitmask) << 8;
		int fourthByte = words[start + 3] & bitmask;

		return firstByte + secondByte + thirdByte + fourthByte;
	}

	private static void processIncomingRequest(SelectionKey key) throws IOException, InvalidKeyException {
		// TODO: Messages the main server that the client already connected
		
		int channel;
		SocketAddress otherClientAddr;
		
		ByteBuffer buffer = ByteBuffer.allocate(2048 + 8);
		long s_time = System.currentTimeMillis();
		DatagramChannel client = (DatagramChannel) key.channel();

		UDPSocketData a = getSocketData(client);
		byte[] message = a.getData();

		// check action flag

		int action = readActionInfo(message, 0);
		System.out.println("Action Info : " + action);

		switch (action) {
		case 0x8fffffff:
			// START PACKET
			
			/*
			 * Packet Structure : 
			 * [0] 32 bit channel info | [4] 256 Bit HMAC-SHA256 | [36] 256 Bit salt | [68] 8 bit username
			 * length | [69] 256 bit sender username | [101] 8 bit username length | [102] 256 bit recipient
			 * username
			 */
			
			System.out.println(message.length);
			
			byte[] hmac = Arrays.copyOfRange(message, 4, 4 + 32);
			byte[] salt = Arrays.copyOfRange(message, 36, 36 + 32);
			
			String hmacHexString = CryptoHandler.bytesToHexString(hmac);
			String saltHexString = CryptoHandler.bytesToHexString(salt);
			
			System.out.println(hmacHexString);
			System.out.println(saltHexString);
			
			byte senderUsernameLength = message[68];
			
			System.out.println(message.length);
			
			String senderUsername = new String(Arrays.copyOfRange(message, 69, 69 + senderUsernameLength),
					StandardCharsets.UTF_8);

			System.out.println(senderUsername);
			
			byte recipientUsernameLength = message[101]; 

			String recipientUsername = new String(Arrays.copyOfRange(message, 102, 102 + recipientUsernameLength),
					StandardCharsets.UTF_8);
			
			// Verify Channel Allocation request
			boolean isVerified = CryptoHandler.verifyAccessToken(senderUsername, saltHexString, hmacHexString, SUPER_SECRET_KEY);
			System.out.println(isVerified);
			if (!isVerified)
			{
				// Ignore any unverified packet
				return;
			}
			
			// check if the recipient is waiting for connection with the sender
			SocketAddress recipientAddress = waitingForConnection.getOrDefault(recipientUsername, null);

			System.out.println(senderUsername + " trying to connect to " + recipientUsername);

			if (recipientAddress == null) {
				System.out.println(recipientUsername + " hasn't been registered yet");
				System.out.println("Adding " + senderUsername + " to the waiting list entry");
				waitingForConnection.put(senderUsername, a.getSocketAddress());
				return;
			}

			waitingForConnection.remove(recipientUsername);

			int senderChannels = channels.allocate(recipientAddress);
			int recipientChannels = channels.allocate(a.getSocketAddress());

			System.out.println("Allocating Channel  " + senderChannels + " to " + recipientAddress);
			System.out.println("Allocating channel " + recipientChannels + " to " + a.getSocketAddress());

			// forward the channels index to both the sender and recipient
			buffer.putLong(0x8fffffff00000000l + senderChannels);
			buffer.flip();
			serverSocket.send(buffer, a.getSocketAddress());
			buffer.clear();

			buffer.putLong(0x8fffffff00000000l + recipientChannels);
			buffer.flip();
			serverSocket.send(buffer, recipientAddress);
			buffer.clear();

			break;
		
		case 0x90000000:
			// CALL DECLINED
			channel = readActionInfo(message, 4);
			otherClientAddr = channels.getChannelAddress(channel);
			buffer.put(message);
			buffer.flip();
			serverSocket.send(buffer, otherClientAddr);
			buffer.clear();
			
			break;
			
		case 0xaaaaaaaa:
			// CALL ACCEPTED
			channel = readActionInfo(message, 4);
			otherClientAddr = channels.getChannelAddress(channel);
			buffer.put(message);
			buffer.flip();
			serverSocket.send(buffer, otherClientAddr);
			buffer.clear();
			
			break;
			
		case 0x90000001:
			// CALL TIMEOUT
			channel = readActionInfo(message, 4);
			otherClientAddr = channels.getChannelAddress(channel);
			buffer.put(message);
			buffer.flip();
			serverSocket.send(buffer, otherClientAddr);
			buffer.clear();
			
			break;
			
		case 0xffffffff:
			channel = readActionInfo(message, 4);
			otherClientAddr = channels.getChannelAddress(channel);
			buffer.put(message);
			buffer.flip();
			serverSocket.send(buffer, otherClientAddr);
			buffer.clear();
			break;


		default:
			// AUDIO PACKET

			channel = action;
			otherClientAddr = channels.getChannelAddress(channel);

			if (otherClientAddr == null) {
				return;
			}

			buffer.put(message);
			buffer.flip();
			serverSocket.send(buffer, otherClientAddr);
			buffer.clear();
			break;
		}

		long e_time = System.currentTimeMillis();
		System.out.println("Response Time : " + (e_time - s_time) + " ms");

	}
}
