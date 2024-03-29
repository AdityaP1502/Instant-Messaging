import socket
import errno
from time import sleep
from time import time
from threading import Thread
import struct

class AudioReceiver():
    # PACKET FORMAT
    START_BUFF_SIZE = 8 
    BUFF_SIZE = 2048 + 8
    
    def __init__(self, conn, f, client, logger = None):
        self._stop = False
        self._t = Thread(target=self.receive, daemon=True)
        self._conn = conn
        self._logger = logger
        
        if self._logger != None:
            self._logger.register_event(name="UDP Latency")
            
        self.output_f = f
        self.client = client
        
    # def get(self):
    #     return self._q.get()
            
    def start(self):
        self._stop = False 
        self._t.start()
        
    def stop(self):
        self._stop = True
        try:
            self._t.join(timeout=10)
        except:
            return
    
    def terminate(self):
        self.stop()
        self._t.join(timeout=1)
     
    def receive(self):
        data = b""
        buff_size = self.START_BUFF_SIZE
        
        while not self._stop:
            try:
                s_time = time()
                # packet = self._conn.socket.recv(self.BUFF_SIZE)
                
                # while (len(data) < self.BUFF_SIZE):
                #     data += packet
                #     packet = self._conn.socket.recv(self.BUFF_SIZE)
                
                # frame = data[:self.BUFF_SIZE]
                # data = data[self.BUFF_SIZE:]
                
                packet, _ = self._conn.socket.recvfrom(buff_size)
                channel, frame_id, frame = packet[:4], packet[4:8], packet[8:]
                
                channel_int = struct.unpack(">I", channel)[0]
                
                if channel_int == 0x8fffffff:
                    channel_args = frame_id
                    self.client.set_channel(channel_args)
                    print("Received Channel : {}".format(struct.unpack(">I", channel_args)[0]))
                    buff_size = self.BUFF_SIZE
                    continue
                
                if channel_int == 0xaaaaaaaa:
                    # both connection has been established 
                    # start to send audio data
                    self.client.set_state(1)
                    continue
                    
                # TODO: channel_closed because user declined
                
                if channel_int == 0x90000000:
                    self.client.set_state(-1)
                    continue
                    
                # TODO: channel_closed because user didn't answer
                if channel_int == 0x90000001:
                    self.client.set_state(-2)
                    continue
                    
                frame_id = struct.unpack(">I", frame_id)[0]
                print("Received : frame - {}".format(frame_id))
                self.output_f(frame, frame_id)
                e_time = time()
                
                if self._logger:
                    self._logger.emit("UDP Latency", "{} ms".format((e_time - s_time) * 1000))
                
            except socket.error as e:
                err = e.args[0]
                if err == errno.EAGAIN or err == errno.EWOULDBLOCK:
                    sleep(1)
                    continue
                
        print("Closing receiver Thread")

    
    
    
            
            
