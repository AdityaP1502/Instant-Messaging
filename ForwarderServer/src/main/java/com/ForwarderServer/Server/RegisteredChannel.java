package com.ForwarderServer.Server;

import java.nio.channels.*;
import java.time.LocalDateTime;

public class RegisteredChannel {
	private final SocketChannel registeredChannel;
	private final LocalDateTime timestamp;
	
	public RegisteredChannel(SocketChannel channel){
		registeredChannel = channel;
		timestamp = LocalDateTime.now();
	}
	
	public SocketChannel getRegisteredChannel() {
		return registeredChannel;
	}
	
	public LocalDateTime getTimestamp() {
		return timestamp;
	}
	
	
	
	
	

}
