package com.ForwarderServer.server.audio;
import java.net.SocketAddress;

public class UDPSocketData {
	private SocketAddress a;
	private byte[] data;
	
	public UDPSocketData(SocketAddress a, byte[] data)
	{
		this.a = a;
		this.data = data;
	}
	
	public byte[] getData()
	{
		return data;
	}
	
	public SocketAddress getSocketAddress()
	{
		return a;
	}
}
