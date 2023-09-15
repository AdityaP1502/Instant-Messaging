
import struct

class AudioClient():   
    REGISTER_PACKET_SIZE_B = 102
    
    
    def __init__(self, conn, username, rcpt_username) -> None:
        self._conn = conn
        
        self.sdr = username
        self.rcpt = rcpt_username
        
        self._key = None
        self._sdr_b = username.encode()
        self._len_sdr_b = struct.pack("B", len(self._sdr_b))
        
        self._sdr_b = self._len_sdr_b + self._sdr_b + b"0" * (32 - len(self._sdr_b))
        
        self._rcpt_b = rcpt_username.encode()
        self._len_rcpt_b = struct.pack("B", len(self._rcpt_b))
        self._rcpt_b = self._len_rcpt_b + self._rcpt_b + b"0" * (32 - len(self._rcpt_b))
        self._channel = None
    
    def set_keys(self, key):
        self._key = key
    
    def set_channel(self, channel : bytes):
        self._channel = channel
        
    def register_channel(self):
        packet = b"\x8f\xff\xff\xff" + self._key + self._sdr_b + self._rcpt_b
        print("Sending {}".format(packet))
        assert len(packet) == self.REGISTER_PACKET_SIZE_B, "Packet length must be 102 bytes, received {}".format(len(packet))
        self._conn.send(packet)
    
    def terminate_channel(self):
        packet = self._key + b'\xFF' + self._sdr_b + self._rcpt_b + b'\x00' * 2048
        assert len(packet) == self.PACKET_SIZE_B, "Packet length must be 2147 bytes, received {}".format(len(packet))
        self._conn.send(packet)
        
    def send_audio(self, audio_data : bytes, frame_id : bytes):
        assert len(audio_data) == 2048, "Audio data must be 2048 bytes long (1024 chunks with 16 bit integer). Received {}".format(len(audio_data))
        
        packet = self._channel + frame_id + audio_data
        print("Sending {} bytes to server".format(len(packet)))
        
        self._conn.send(packet)
            
        

    
    