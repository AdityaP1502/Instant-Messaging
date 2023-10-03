import struct
import queue
import pyaudio
import traceback

from threading import Thread
from multiprocessing import Value
from time import time

class Mixer():
    def __init__(self, send, logger=None) -> None:
        self._recorder_t = Thread(target=self.record)
        self._player_t = Thread(target=self.play)
        self._stream_p : pyaudio._Stream = None
        self._stream_r : pyaudio._Stream = None
        self._format = pyaudio.paInt16
        self._channels = 1
        self._rate = 5000
        self._chunk = 1024 
        self._stop_record = False
        self._stop_play = False
        self._send = send
        self._buffer = queue.Queue(maxsize=10000)
        self._expected_frame = Value("i")
        self.running = False
        
        self._logger = logger
        if logger != None:
            self._logger.register_event(name="Recorder Latency")
            self._logger.register_event(name="Player Fetch Latency")
            self._logger.register_event(name="Dropped Frame")
        
        self._prev_length = 0
         
    def terminate(self):
        print("Terminating Mixer...")
        self._stop_record = True
        self._stop_play = True
        self.running = False
        
        self._buffer.put(None)
        
    def append_audio(self, audio, frame_id):
        with self._expected_frame.get_lock():
            if frame_id >= self._expected_frame.value:
                self._buffer.put(audio)
                
            elif frame_id < self._expected_frame.value and self._logger != None:
                self._logger.emit("Dropped Frame", "[WARNING] Dropped frame {}".format(frame_id))
            
    def mute_record(self):
        self._stop_record = True
        self._recorder_t.join(timeout=10)
    
    def mute_play(self):
        self._stop_play = True
        self._buffer.put(None)
        self._recorder_t.join(timeout=10)
    
    def is_alive(self):
        return self._player_t.is_alive() or self._recorder_t.is_alive()
          
    def record(self): 
        frame_id = 0  
        try:
            while not self._stop_record:
                s_time = time()
                data = self._stream_r.read(self._chunk)
                self._send(data, struct.pack(">I", frame_id))
                frame_id += 1
                e_time = time()
                
                if self._logger != None:
                    self._logger.emit("Recorder Latency","{} ms".format((e_time - s_time) * 1000))
                
                # sleep(4 * self._chunk / self._rate)
                        
        except:
            print("Error Occured")
            traceback.print_exc()
            self.terminate()
            
        self._stream_r.close()
        self._stream_r = None
        print("Closing Recording Thread")
    
    def play(self):
        try:
            while not self._stop_play:
                s_time = time()
                try:
                    audio = self._buffer.get(timeout=self._chunk / self._rate)
                except queue.Empty: 
                    with self._expected_frame.get_lock():  
                        self._expected_frame.value += 1
                    continue
                
                if audio == None:
                    continue
                        
                e_time = time()
                
                if self._logger != None:
                    self._logger.emit("Player Fetch Latency", "{} ms".format((e_time - s_time) * 1000))
                    
                self._stream_p.write(audio, self._chunk)
                
        except:
            print("Error Occured")
            traceback.print_exc()
            self.terminate()
                        
        self._stream_p.close()
        self._stream_p = None
        print("Closing Player Thread")
         
    def start(self):
        p = pyaudio.PyAudio()
        self.running = True
        self._stop_record = False
        self._stop_play = False
        self._stream_p = p.open(format=self._format, 
                            channels=self._channels,
                            rate=self._rate, 
                            output=True, 
                            frames_per_buffer=self._chunk)
        
        self._stream_r = p.open(format=self._format, 
                            channels=self._channels,
                            rate=self._rate, 
                            input=True, 
                            frames_per_buffer=self._chunk)
        
        self._player_t.start()
        self._recorder_t.start()
        
        
        