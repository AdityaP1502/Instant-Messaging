import sys
from audio_client import AudioClient
from audio_connection import AudioConnection
from audio_receiver import AudioReceiver
from mixer import Mixer
from logger import Logger

from time import sleep

from random import randint
if __name__ == "__main__":
    
    if len(sys.argv) < 3:
        print("Usage: python test.py [username] [recipient username]")
        sys.exit(1)
    
    try:    
        username, recipient_username = sys.argv[1:]
        
        conn = AudioConnection()
        ret, _ = conn.connect_to_socket_udp()
        
        if ret != 0:
            print(_)
            sys.exit(1)
        
        aud = AudioClient(username=username, conn=conn, rcpt_username=recipient_username)
        logger = Logger()
        logger.start()
        mixer = Mixer(send=aud.send_audio, logger=logger)
        rec = AudioReceiver(conn=conn, f=mixer.append_audio, logger=logger, client=aud)
        
        key = bytes(bytearray([randint(0, 7) for i in range(32)]))
        aud.set_keys(key)
        rec.start()
        aud.register_channel()
        
        while aud._channel == None:
            continue
        
        sleep(5)
        
        mixer.start(recpt=recipient_username)
        mixer.mute_play()
        
        while mixer.is_alive():
            continue
        
    except KeyboardInterrupt:
        print("CTRL + C is pressed...")
        
    finally:
        mixer.terminate()
        rec.stop()
        logger.terminate()
        
        
    # aud.terminate_channel()
    
    