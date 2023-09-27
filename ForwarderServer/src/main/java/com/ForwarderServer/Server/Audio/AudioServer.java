package com.ForwarderServer.server.audio;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.DatagramChannel;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;

public class AudioServer {
	
	private static final int BUFFER_SIZE = 4096; // bytes

	private static AudioChannels channels;
	private static Map<String, SocketAddress> waitingForConnection;
	
	private static DatagramChannel serverSocket;
	private static String SUPER_SECRET_KEY;
	
	public static void main(String[] args) throws IOException {
		// SET SECRET KEY 
		SUPER_SECRET_KEY = System.getenv("FORWARD_SERVER_SUPER_SECRET_KEY");
		
		if (SUPER_SECRET_KEY == null)
		{
			System.out.println("SUPER SECRET KEY environment variables hasn't been set!");
			System.exit(-1);
		}
		
		System.out.println(SUPER_SECRET_KEY);
		
		Selector selector = Selector.open();
		serverSocket = DatagramChannel.open();

		serverSocket.bind(new InetSocketAddress("localhost", 8080));
		serverSocket.configureBlocking(false);
		serverSocket.register(selector, SelectionKey.OP_READ);

		waitingForConnection = new HashMap<>();
		channels = new AudioChannels();
		
		// TODO: connect to main server 
		
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
	
	
	private static UDPSocketData getSocketData(DatagramChannel s) throws IOException
	{
		ByteBuffer buffer = ByteBuffer.allocate(BUFFER_SIZE);
		
		SocketAddress a = s.receive(buffer);
		buffer.flip();
		byte[] message = new byte[buffer.remaining()];
	    buffer.get(message);
	    
	    return new UDPSocketData(a, message);
	}
	
	private static int readActionInfo(byte[] words, int start)
	{
		// assuming the data is sorted in big endian order
		int bitmask = 0xFF;
		
		int firstByte = (words[start] & bitmask) << 24;
		int secondByte = (words[start + 1] & bitmask) << 16;
		int thirdByte = (words[start + 2] & bitmask) << 8;
		int fourthByte = words[start + 3] & bitmask;
		
		
		return firstByte + secondByte + thirdByte + fourthByte;
	}
	
	private static void processIncomingRequest(SelectionKey key) throws IOException
	{
		// TODO: Verify HMAC
		// TODO: Messages the main server that the client already connected
		
		ByteBuffer buffer = ByteBuffer.allocate(2048 + 8);
		long s_time = System.currentTimeMillis();
		DatagramChannel client = (DatagramChannel) key.channel();
		
		
		UDPSocketData a = getSocketData(client);
		byte[] message = a.getData();
		
		// check action flag
		
		int action = readActionInfo(message, 0);
		System.out.println("Action Info : " + action);
		
		switch (action)
		{
			case 0x8fffffff:
				// START PACKET
				// String HMAC = new String(Arrays.copyOfRange(message, 4, 4 + HMAC_LENGTH));
				byte senderUsernameLength = message[36];
				String senderUsername = new String(Arrays.copyOfRange(message, 37, 37 + senderUsernameLength), StandardCharsets.UTF_8);
				
				byte recipientUsernameLength = message[69]; // :)))
				
				String recipientUsername = new String(Arrays.copyOfRange(message, 70, 70 + recipientUsernameLength), StandardCharsets.UTF_8);
				
				// check if the recipient is waiting for connection with the sender
				SocketAddress recipientAddress = waitingForConnection.getOrDefault(recipientUsername, null);
				
				System.out.println(senderUsername + " trying to connect to " + recipientUsername);
				
				if (recipientAddress == null)
				{
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
			
			case 0xffffffff:
				// TERMINATE
				client.close();
				
			default:
				// AUDIO PACKET
				
				int channel = action;
				SocketAddress otherClientAddr = channels.getChannelAddress(channel);
				
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
