package com.ForwarderServer.Server;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.SelectionKey;
import java.nio.channels.Selector;
import java.nio.channels.ServerSocketChannel;
import java.nio.channels.SocketChannel;
import java.util.Arrays;
import java.util.HashMap;
import java.util.Iterator;
import java.util.Map;
import java.util.Set;

public class MessageForwarderServer {
	// clientTable is the list of user
	// that are online / connected to the server
	private static Map<String, RegisteredChannel> clientTable;
	private static Map<SocketAddress, String> clientIPAddressTable;
	private static Map<String, String> callTable;

	private static final int BUFFER_SIZE = 1048576;
	private static final char EOF = '\n';
	private static String SUPER_SECRET_KEY;

	public static void main(String[] args) throws IOException {
		
		SUPER_SECRET_KEY = System.getenv("FORWARD_SERVER_SUPER_SECRET_KEY");
		
		if (SUPER_SECRET_KEY == null)
		{
			System.out.println("SUPER SECRET KEY environment variables hasn't been set!");
			System.exit(-1);
		}
		
		Selector selector = Selector.open();
		ServerSocketChannel serverSocket = ServerSocketChannel.open();
		
		serverSocket.bind(new InetSocketAddress("localhost", 8080));
		serverSocket.configureBlocking(false);
		serverSocket.register(selector, SelectionKey.OP_ACCEPT);

		ByteBuffer buffer = ByteBuffer.allocate(BUFFER_SIZE);

		clientTable = new HashMap<>();
		clientIPAddressTable = new HashMap<>();
		callTable = new HashMap<>();

		while (true) {
			selector.select();
			Set<SelectionKey> selectedKeys = selector.selectedKeys();
			Iterator<SelectionKey> iter = selectedKeys.iterator();

			while (iter.hasNext()) {
				SelectionKey key = iter.next();

				if (key.isAcceptable()) {
					register(selector, serverSocket);
				}

				if (key.isReadable()) {
					processIncomingRequest(buffer, key);
				}

				iter.remove();
			}
		}
	}

	private static void forceCloseConnection(ByteBuffer buffer, SocketChannel client) throws IOException {

		String username = clientIPAddressTable.getOrDefault(client.getRemoteAddress(), null);

		// username has already been registered
		if (username != null) {
			// abort a call connection if in call
			String currentlyInCallWith = callTable.getOrDefault(username, null);

			if (currentlyInCallWith != null) {
				SocketChannel otherClient = clientTable.get(currentlyInCallWith).getRegisteredChannel();
				try {
					abortCall(buffer, otherClient);

				} catch (Exception e) {
					// the otherClient also abort the connection abruptly
					otherClient.close();
					clientTable.remove(currentlyInCallWith);
					clientIPAddressTable.remove(otherClient.getRemoteAddress());

				} finally {
					// remove from the call table
					callTable.remove(currentlyInCallWith);
					callTable.remove(username);
				}

			}

			clientTable.remove(username);
			clientIPAddressTable.remove(client.getRemoteAddress());

			System.out.println("Removing " + username + " from registered channel");
			System.out.println(username + " is offline");
		}

		System.out.println(client.getRemoteAddress() + " removed");

		client.close();
		buffer.clear();

	}

	private static void putStringToBuffer(ByteBuffer buffer, String message) throws IOException {
		buffer.clear();
		if (message.length() > BUFFER_SIZE) {
			buffer.put(message.substring(0, BUFFER_SIZE - 1).getBytes("UTF-8"));
			return;
		}

		buffer.put(message.getBytes("UTF-8"));

	}

	private static void writeResponseToChannel(ByteBuffer buffer, SocketChannel client, String response)
			throws IOException {

		// put the string into the buffer and
		// write the buffer into the channel
		System.out.println("Writing " + response + " to " + client.getRemoteAddress());
		putStringToBuffer(buffer, response + EOF);

		try {
			buffer.flip();
			client.write(buffer);
			buffer.clear();

			buffer.put(new byte[BUFFER_SIZE]);
			buffer.clear();
		} catch (IOException e) {
			if (!client.isConnected()) {
				forceCloseConnection(buffer, client);
			}
		}
	}

	private static String makeResponseMessage(String header, String payload) {
		return header + ";" + "payload=" + payload;
	}

	private static void handleError(ByteBuffer buffer, SocketChannel client, ErrorMessage err, String uid, int action)
			throws IOException {

		ResponseHeader header = ResponseHeader.ERROR;
		String payload = "status=" + err.getErrorMessage() + ";uid=" + uid;
		String response = makeResponseMessage(header.getHeader(), payload);
		writeResponseToChannel(buffer, client, response);
		switch (action) {
		case 1:
			// terminate connection
			forceCloseConnection(buffer, client);
		case 3:
			// terminate connection
			// client hasn't checked in yet
			client.close();
			break;
		}
	}

	private static void sendOKMessage(ByteBuffer buffer, SocketChannel client, String uid) throws IOException {
		ResponseHeader header = ResponseHeader.OK;
		String payload = "uid=" + uid;
		String response = makeResponseMessage(header.getHeader(), payload);
		writeResponseToChannel(buffer, client, response);
	}

	private static ErrorMessage isRequestValid(String uidEntryName, String uidValue, String reqTypeEntryName,
			String payloadEntryName) {
		if (!reqTypeEntryName.equals("reqtype")) {
			return new ErrorMessage(ErrorType.InvalidRequestTypeEntryName,
					"expected 'reqtype' got '" + reqTypeEntryName + "'");
		}

		if (!uidEntryName.equals("uid")) {
			return new ErrorMessage(ErrorType.InvalidRequestTypeEntryName, "expected 'uid' got '" + uidEntryName + "'");
		}

		try {
			int uidValueInt = Integer.parseInt(uidValue);
			if (uidValueInt < 0) {
				return new ErrorMessage(ErrorType.InvalidUID,
						"Expected uid value to be greater than 0. Received " + uidValueInt);
			}
		} catch (NumberFormatException e) {
			return new ErrorMessage(ErrorType.InvalidUID, "Expected uid value to be an integer. Received " + uidValue);
		}

		if (!payloadEntryName.equals("payload")) {
			return new ErrorMessage(ErrorType.InvalidRequestTypeEntryName,
					"expected 'payload' got '" + payloadEntryName + "'");
		}

		return null;
	}

	private static void processIncomingRequest(ByteBuffer buffer, SelectionKey key) throws IOException {
		ErrorMessage err;
		RequestType req;

		SocketChannel client = (SocketChannel) key.channel();

		try {
			client.read(buffer);
		} catch (IOException e) {
			System.out.println("Something is wrong!");
			System.out.println(e.getMessage());

			forceCloseConnection(buffer, client);
			return;
		}
		
		String clientRequest = new String(buffer.array()).trim();
		
		if (clientRequest.length() > 100) {
			System.out.println(client.getRemoteAddress() + " send:" + clientRequest.substring(0, 100));
		}
		else
		{
			System.out.println(client.getRemoteAddress() + " send:" + clientRequest);
		}
		
		// Request format:
		// reqtype=####;payload=###

		String[] parsedRequest = clientRequest.split(";", 3);

		if (parsedRequest.length == 0) {
			System.out.println("Client close the connection? Closing the connection.");
			forceCloseConnection(buffer, client);
			return;
		}

		if (parsedRequest.length == 1) {
			handleError(buffer, client, new ErrorMessage(ErrorType.InvalidMessageFormat, "Received " + clientRequest),
					"-1", 0);
			return;
		}

		String uid = parsedRequest[0];
		String reqType = parsedRequest[1];
		String payload = parsedRequest[2];

		// Parsed array will have
		// the 0th index as the entry name and the 1st index as the value
		String[] parsedUIDRequestField = uid.split("=", 2);
		String[] parsedRequestField = reqType.split("=", 2);
		String[] parsedPayloadField = payload.split("=", 2);

		err = isRequestValid(parsedUIDRequestField[0], parsedUIDRequestField[1], parsedRequestField[0],
				parsedPayloadField[0]);

		if (!(err == null)) {
			handleError(buffer, client, err, uid, 0);
			return;
		}

		try {
			req = RequestType.valueOf(parsedRequestField[1]);
			switch (req) {
			case CHECKIN:
				checkIn(buffer, client, parsedPayloadField[1], parsedUIDRequestField[1]);
				break;

			case FETCH:
				fetchMessage(buffer, client, parsedPayloadField[1], parsedUIDRequestField[1]);
				break;

			case READY:
				processReadyRequest(buffer, client, parsedPayloadField[1], parsedUIDRequestField[1]);
				break;

			case SENDMESSAGE:
				sendMessage(buffer, client, parsedPayloadField[1], parsedUIDRequestField[1]);
				break;

			case SENDAUDIO:
				sendAudio(buffer, client, parsedPayloadField[1], parsedUIDRequestField[1]);
				break;

			case INITIATECALL:
				callHandler(buffer, client, parsedPayloadField[1], CallStatus.INCOMINGCALL, parsedUIDRequestField[1]);
				break;

			case ACCEPTCALL:
				callHandler(buffer, client, parsedPayloadField[1], CallStatus.CALLACCEPTED, parsedUIDRequestField[1]);
				break;

			case DECLINECALL:
				callHandler(buffer, client, parsedPayloadField[1], CallStatus.CALLDECLINED, parsedUIDRequestField[1]);
				break;

			case TIMEOUTCALL:
				callHandler(buffer, client, parsedPayloadField[1], CallStatus.CALLTIMEOUT, parsedUIDRequestField[1]);
				break;

			case TERMINATE:
				terminateConnection(buffer, client, parsedPayloadField[1], parsedUIDRequestField[1]);
				break;

			default:
				break;
			}
		} catch (IllegalArgumentException | NullPointerException e) {
			// send error message
			handleError(buffer, client, new ErrorMessage(ErrorType.InvalidRequestType, "Received " + parsedRequestField[1]),
					parsedUIDRequestField[1], 0);
		}

	}

	private static void fetchMessage(ByteBuffer buffer, SocketChannel client, String payload, String uid)
			throws IOException {

		String[] parsedUsernameField = payload.split("=", 2);
		String username = parsedUsernameField[1];

		// check for format ?
		// i guess later

		if (clientTable.getOrDefault(username, null) == null) {
			handleError(buffer, client, new ErrorMessage(ErrorType.InvalidFetchProtocol, ""), uid, 0);
			return;
		}

		// get history
		String userHistory = HistoryReader.readUserHistory(username);
		String responsePayload = "uid=" + uid + ";history=" + userHistory;
		String response = makeResponseMessage(ResponseHeader.FETCH.getHeader(), responsePayload);

		writeResponseToChannel(buffer, client, response);
	}

	private static void processReadyRequest(ByteBuffer buffer, SocketChannel client, String payload, String uid)
			throws IOException {
		String[] parsedUsernameField = payload.split("=", 2);
		String username = parsedUsernameField[1];

		if (clientTable.getOrDefault(username, null) == null) {
			handleError(buffer, client, new ErrorMessage(ErrorType.InvalidReadyProtocol, ""), uid, 0);
			return;
		}

		// remove user history
		HistoryReader.deleteUserHistory(username);

		// send an ok message
		sendOKMessage(buffer, client, uid);
	}

	private static ErrorMessage isSenderContentValid(String entryName, String username, SocketChannel sender) {
		if (!entryName.equals("sdr")) {
			return new ErrorMessage(ErrorType.InvalidSendMessagePayloadEntryName,
					"Expected 'sdr' got '" + entryName + "'");
		}

		if (!clientTable.get(username).getRegisteredChannel().equals(sender)) {
			return new ErrorMessage(ErrorType.InvalidSenderUsername, "Received " + username);
		}

		return null;
	}

	private static ErrorMessage isRecipientContentValid(String entryName, String username) {
		if (!entryName.equals("rcpt")) {
			return new ErrorMessage(ErrorType.InvalidSendMessagePayloadEntryName,
					"Expected 'rcpt' got '" + entryName + "'");
		}

		// must check to the server API
		// if a username is valid or not

		// if (clientTable.getOrDefault(username, null) == null) {
		// return new ErrorMessage(ErrorType.InvalidRecipientUsername, "Received " +
		// username);
		// }

		return null;
	}

	private static void sendMessage(ByteBuffer buffer, SocketChannel sender, String payload, String uid)
			throws IOException {
		// Parse clientMessage
		// message format :
		// sdr=######;rcpt=#####;msg=#######
		ErrorMessage err;

		String[] parsedMessage = payload.split(";", 4);
		System.out.println(Arrays.toString(parsedMessage));

		if (parsedMessage.length != 4) {
			handleError(buffer, sender, new ErrorMessage(ErrorType.InvalidSendMessagePayloadLength,
					"Expected 4 got " + Integer.toString(parsedMessage.length)), uid, 0);
			return;
		}

		String[] senderField = parsedMessage[0].split("=");
		err = isSenderContentValid(senderField[0], senderField[1], sender);
		if (!(err == null)) {
			handleError(buffer, sender, err, uid, 0);
			return;
		}

		String[] recipientField = parsedMessage[1].split("=");
		err = isRecipientContentValid(recipientField[0], recipientField[1]);
		if (!(err == null)) {
			handleError(buffer, sender, err, uid, 0);
			return;
		}

		if (clientTable.getOrDefault(recipientField[1], null) == null) {
			// store the message to the local storage
			String message = parsedMessage[3].split("=")[1];
			String ts = parsedMessage[2].split("=")[1];

			HistoryReader.writeUserHistory(recipientField[1], senderField[1], message, ts);
			sendOKMessage(buffer, sender, uid);
			return;
		}

		forward(buffer, recipientField[1], payload, ResponseHeader.SENDMESSAGE);
		sendOKMessage(buffer, sender, uid);
	}

	private static void sendAudio(ByteBuffer buffer, SocketChannel sender, String payload, String uid)
			throws IOException {
		ErrorMessage err;

		String[] parsedMessage = payload.split(";", 3);
		// System.out.println(Arrays.toString(parsedMessage));

		if (parsedMessage.length != 3) {
			handleError(buffer, sender, new ErrorMessage(ErrorType.InvalidSendMessagePayloadLength,
					"Expected 3 got " + Integer.toString(parsedMessage.length)), uid, 0);
			return;
		}

		String[] senderField = parsedMessage[0].split("=");
		err = isSenderContentValid(senderField[0], senderField[1], sender);
		if (!(err == null)) {
			handleError(buffer, sender, err, uid, 0);
			return;
		}

		String[] recipientField = parsedMessage[1].split("=");
		err = isRecipientContentValid(recipientField[0], recipientField[1]);
		if (!(err == null)) {
			handleError(buffer, sender, err, uid, 0);
			return;
		}
		
		if (clientTable.getOrDefault(recipientField[1], null) == null) {
			return;
		}
		
		forward(buffer, recipientField[1], payload, ResponseHeader.SENDAUDIO);
		// sendOKMessage(buffer, sender, uid);
	}

	private static void callHandler(ByteBuffer buffer, SocketChannel sender, String payload, CallStatus status,
			String uid) throws IOException {
		ErrorMessage err;

		ResponseHeader header = status.getHeader();

		String[] parsedMessage = payload.split(";", 2);
		if (parsedMessage.length != 2) {
			handleError(buffer, sender, new ErrorMessage(ErrorType.InvalidInitiateCallPacketFormat,
					"Expected 2 Field in the packet got " + Integer.toString(parsedMessage.length)), uid, 0);
			return;
		}

		String[] senderField = parsedMessage[0].split("=");
		err = isSenderContentValid(senderField[0], senderField[1], sender);
		if (!(err == null)) {
			handleError(buffer, sender, err, uid, 0);
			return;
		}

		String[] recipientField = parsedMessage[1].split("=");
		err = isRecipientContentValid(recipientField[0], recipientField[1]);
		if (!(err == null)) {
			handleError(buffer, sender, err, uid, 0);
			return;
		}

		String senderUsername = senderField[1];
		String recipientUsername = recipientField[1];
		
		if (clientTable.getOrDefault(recipientField[1], null) == null) {
			header = ResponseHeader.CALLABORT;
			// send the response to the recipient
			forward(buffer, senderField[1], payload, header);
			sendOKMessage(buffer, sender, uid);		
		}
		
		switch (status) {
		case INCOMINGCALL:
			// check if user already in call
			if (callTable.getOrDefault(recipientUsername, null) != null) {
				// send an abort message
				header = ResponseHeader.CALLABORT;
				break;
			}

			callTable.put(senderUsername, recipientUsername);
			callTable.put(recipientUsername, senderUsername);

			break;

		case CALLDECLINED:

		case CALLTIMEOUT:

		case CALLTERMINATE:
			callTable.remove(senderUsername);
			callTable.remove(recipientUsername);
			break;

		default:
			break;
		}

		// send the response to the recipient
		forward(buffer, recipientField[1], payload, header);

		sendOKMessage(buffer, sender, uid);
	}

	private static void forward(ByteBuffer buffer, String recipientUsername, String payload, ResponseHeader header)
			throws IOException {

		SocketChannel recChannel = clientTable.get(recipientUsername).getRegisteredChannel();
		String response = makeResponseMessage(header.getHeader(), payload);
		writeResponseToChannel(buffer, recChannel, response);
	}

	private static void register(Selector selector, ServerSocketChannel serverSocket) throws IOException {

		SocketChannel client = serverSocket.accept();
		client.configureBlocking(false);
		client.register(selector, SelectionKey.OP_READ);
	}

	private static ErrorMessage isCheckInValid(String entryName, String username, SocketChannel client) {
		if (!entryName.equals("username")) {
			return new ErrorMessage(ErrorType.InvalidPayloadFormat, "expected 'username' received " + entryName);
		}

		RegisteredChannel clientWithSameUsername = clientTable.getOrDefault(username, null);
		if ((clientWithSameUsername != null) && clientWithSameUsername.getRegisteredChannel().equals(client)) {

			return new ErrorMessage(ErrorType.InvalidCheckInProtocol,
					"Last check-in:" + clientWithSameUsername.getTimestamp());
		}

		else if ((clientWithSameUsername != null) && !clientWithSameUsername.getRegisteredChannel().equals(client)) {

			return new ErrorMessage(ErrorType.InvalidCheckInUsername, "received " + username);
		}
		return null;

	}

	private static void checkIn(ByteBuffer buffer, SocketChannel client, String payload, String uid)
			throws IOException {
		ErrorMessage err;

		String[] checkInField = payload.split("=");

		if (checkInField.length < 2) {
			handleError(buffer, client,
					new ErrorMessage(ErrorType.InvalidCheckInPayloadEntryName, "received " + payload), uid, 3);
			return;
		}

		err = isCheckInValid(checkInField[0], checkInField[1], client);
		if (!(err == null)) {
			handleError(buffer, client, err, uid, 3);
			return;
		}

		clientTable.put(checkInField[1], new RegisteredChannel(client));
		clientIPAddressTable.put(client.getRemoteAddress(), checkInField[1]);

		sendOKMessage(buffer, client, uid);
	}

	private static void abortCall(ByteBuffer buffer, SocketChannel client) throws IOException {
		ResponseHeader header = ResponseHeader.CALLABORT;
		String payload = "uid=-1";
		String response = makeResponseMessage(header.getHeader(), payload);
		writeResponseToChannel(buffer, client, response);
	}

	private static void terminateConnection(ByteBuffer buffer, SocketChannel client, String payload, String uid)
			throws IOException {
		String[] splitMessage = payload.split("=");

		if (!splitMessage[0].equals("username")) {
			handleError(buffer, client, new ErrorMessage(ErrorType.InvalidCheckInPayloadEntryName,
					"expected username got " + splitMessage[0]), uid, 0);
		}

		sendOKMessage(buffer, client, uid);

		client.close();
		clientTable.remove(splitMessage[1]);
		clientIPAddressTable.remove(client.getRemoteAddress());
		System.out.print(splitMessage[1] + " is offline");
	}
}
