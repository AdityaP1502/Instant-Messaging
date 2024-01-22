
import struct

class AudioClient():   
    REGISTER_PACKET_SIZE_B = 134
    
    def __init__(self, conn, username, rcpt_username, salt, token) -> None:
        self._conn = conn
        
        self.sdr = username
        self.rcpt = rcpt_username
        
        self._token = token
        self._salt = salt
        
        self._sdr_b = username.encode()
        self._len_sdr_b = struct.pack("B", len(self._sdr_b))
        
        self._sdr_b = self._len_sdr_b + self._sdr_b + b"0" * (32 - len(self._sdr_b))
        
        self._rcpt_b = rcpt_username.encode()
        self._len_rcpt_b = struct.pack("B", len(self._rcpt_b))
        self._rcpt_b = self._len_rcpt_b + self._rcpt_b + b"0" * (32 - len(self._rcpt_b))
        self._channel = None
        self._state = 0
    
    def set_channel(self, channel : bytes):
        self._channel = channel
        
    def is_channel_set(self):
        if self._channel == None:
            return False
        
        return True
    
    def set_state(self, state):
        self._state = state
    
    def get_state(self):
        return self._state
        
    def register_channel(self):
        packet = b"\x8f\xff\xff\xff" + self._token + self._salt + self._sdr_b + self._rcpt_b
        print("Sending {}".format(packet))
        assert len(packet) == self.REGISTER_PACKET_SIZE_B, "Packet length must be 134 bytes, received {}".format(len(packet))
        self._conn.send(packet)
    
    def declined_connection(self):
        packet = b"\x90\x00\x00\x00" + self._channel
        self._conn.send(packet)
    
    def accept_connection(self):
        packet = b"\xaa\xaa\xaa\xaa" + self._channel
        self._conn.send(packet)
        
    def connection_timeout(self):
        packet = b"\x90\x00\x00\x01" + self._channel
        self._conn.send(packet)
    
    def declined_connection(self):
        packet = b"\x90\x00\x00\x00" + self._channel
        self._conn.send(packet)
    
    def accept_connection(self):
        packet = b"\xaa\xaa\xaa\xaa" + self._channel
        self._conn.send(packet)
          
    def terminate_channel(self):
        packet = b"\xff\xff\xff\xff" + self._channel
        self._conn.send(packet)
        
    def send_audio(self, audio_data : bytes, frame_id : bytes):
        # assert len(audio_data) == 2048, "Audio data must be 2048 bytes long (1024 chunks with 16 bit integer). Received {}".format(len(audio_data))
        
        packet = self._channel + frame_id + audio_data
        print("Sending {} bytes to server".format(len(packet)))
        
        self._conn.send(packet)
            
        

    
    