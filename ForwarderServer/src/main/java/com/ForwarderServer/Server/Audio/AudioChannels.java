package com.ForwarderServer.Server.Audio;

import java.util.ArrayList;
import java.net.SocketAddress;

public class AudioChannels {
	private static final int CHANNEL_MAX_SIZE = 0x7fffffff; // 2147483647 (max 32 bit signed int)
	private static final int CHANNEL_BASE_SIZE = 0x0000ffff; // 65535
	
	private ArrayList<SocketAddress> channels;
	private ArrayList<Integer> freeChannels;
	private int channelCurrMaxSize;
	private int freeChannelsSize;
	private boolean isMaxSizeLimitReached;
	
	public AudioChannels()
	{
		channelCurrMaxSize = CHANNEL_BASE_SIZE; 
		freeChannelsSize = CHANNEL_BASE_SIZE;
		channels = new ArrayList<>(CHANNEL_BASE_SIZE);
		isMaxSizeLimitReached = false;
		
		freeChannels = new ArrayList<>(CHANNEL_BASE_SIZE);
		fillChannels(0, CHANNEL_BASE_SIZE);
		
	}
	
	private void fillChannels(int start, int end)
	{
		for (int i = start; i < end; i++)
		{
			freeChannels.add(i);
		}
		
		freeChannelsSize += (end - start) + 1;
	}
	
	private void returnChannel(int index)
	{
		freeChannels.add(index);
		freeChannelsSize++;
	}
	
	public int allocate(SocketAddress a)
	{
		if (freeChannelsSize == 0 && !isMaxSizeLimitReached)
		{
			int t = channelCurrMaxSize;
			channelCurrMaxSize = 2 * channelCurrMaxSize;
			fillChannels(t, channelCurrMaxSize);
			
			if (channelCurrMaxSize == CHANNEL_MAX_SIZE)
			{
				isMaxSizeLimitReached = true;
				return -1;
			}
		}
		
		int index = freeChannels.get(0);
		channels.add(index, a);
		freeChannels.remove(0);
		
		return index;
	}
	
	public SocketAddress getChannelAddress(int channelIndex)
	{
		return channels.get(channelIndex);
	}
	
	
	public void close(int index)
	{
		returnChannel(index);
	}
}
