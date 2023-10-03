import sys
import getopt

from sys import argv 
from time import sleep
from threading import Thread

from audio.audio_client import AudioClient
from audio.audio_connection import AudioConnection
from audio.audio_receiver import AudioReceiver
from audio.mixer import Mixer

from ui.cli.writer import Writer
from ui.cli.page.page import PageLoader, UserInputHandler
from ui.cli.page.on_call import OnCallPage
from ui.cli.page.call_incoming import IncomingCallPage
from ui.cli.page.calling import CallPage

# from audio.logger import Logger

# Configuration variables
MODE = None  # 0 : Caller ; 1 : Receiver
ACCESS_TOKEN = None
SALT = None
SENDER_USERNAME = None
RECIPIENT_USERNAME = None
IP_ADDERSS = None
PORT = None

def parse_option():
    c = 0
    try:
        opts, args = getopt.getopt(argv[1:], shortopts="s:r:", longopts=["ip=, port=, network=, sender=, recipient=, token=, salt=, caller, receiver"])
    except getopt.GetoptError as err:
        print(err)
        print("Error : Invalid Argument")
        exit(1)

    for opt, arg in opts:
        if opt == "-s":
            SENDER_USERNAME = arg
            c += 1
        
        if opt == "-r":
            RECIPIENT_USERNAME = arg
            c += 1
                        
        if opt == "--ip":
            IP_ADDERSS = arg
            c += 1
            
        if opt == "--port":
            try:
                PORT = int(arg)
                c += 1
            except ValueError:
                sys.exit(1)
                
        if opt == "--network":
            try:
                _ = arg.split(":")
                
                if IP_ADDERSS !=  None and PORT != None:
                    sys.exit(2)
                    
                IP_ADDERSS, PORT = _
                PORT = int(PORT)
                
                c += 2
            except:
                sys.exit(1)
                
        if opt == "--sender":
            if SENDER_USERNAME == None:
                SENDER_USERNAME = arg
                c += 1
                continue
                
            sys.exit(2)
            
        if opt == "--recipient":
            if RECIPIENT_USERNAME == None:
                RECIPIENT_USERNAME = arg
                c += 1
                continue
                
            sys.exit(2)
            
        if opt == "--token":
            try:
                ACCESS_TOKEN = bytes.fromhex(arg)
                c += 1
            except:
                sys.exit(4)
                
        if opt == "--salt":
            try:
                SALT = bytes.fromhex(arg)
                c += 1
            except:
                sys.exit(4)
                
        if opt == "--caller":
            if MODE != None:
                sys.exit(2)
                
            MODE = 0
            c += 1
            
        if opt == "--receiver":
            if MODE != None:
                sys.exit(2)
                
            MODE = 1
            c += 1
                
    if c != 5:
        sys.exit(3)

def watch_connection_state(aud : AudioClient, mixer : Mixer, page_loader : PageLoader):
    while not aud.is_connection_ready:
        continue
    
    mixer.start()
    if MODE == 0:
        call_page = OnCallPage(RECIPIENT_USERNAME)
    
    else:
        call_page = OnCallPage(SENDER_USERNAME)
        
    page_loader.load_new_page(call_page)
    return

if __name__ == "__main__":
    parse_option()
    
    conn = AudioConnection(ip=IP_ADDERSS, port=PORT)
    ret, _ = conn.connect_to_socket_udp()
    
    if ret != 0:
        print(_)
        sys.exit(1)
    
    aud = AudioClient(username=SENDER_USERNAME, conn=conn, rcpt_username=RECIPIENT_USERNAME, token=ACCESS_TOKEN, salt=SALT)
    
    mixer = Mixer(send=aud.send_audio)
    rec = AudioReceiver(conn=conn, f=mixer.append_audio, client=aud)
    
    rec.start()
    aud.register_channel()
    
    while not aud.is_channel_set():
        continue
    
    if MODE == 0:
        on_call_page = OnCallPage(RECIPIENT_USERNAME)
        page_loader = PageLoader(on_call_page)
    
    elif MODE == 1:
        incoming_call_page = IncomingCallPage(SENDER_USERNAME)
        
    # create a writer
    writer = Writer(page_loader=page_loader)   
    writer.start()
    
    # TODO: When does we start to record the audio and send it over the channel? 
    watcher_r = Thread(target=watch_connection_state, args=(aud, mixer, page_loader, ), daemon=True)
    watcher_r.start()
    
    while rc == 0:
        action = input()
            
        if page_loader.current_active_page == "ONCALLPAGE":
            rc = UserInputHandler.process_on_call_input(action)
                    
        elif page_loader.current_active_page() == "INCOMINGCALLPAGE":
            rc = UserInputHandler.process_incoming_call_input(action)

        elif page_loader.current_active_page == "CALLPAGE":
            rc = UserInputHandler.process_calling_input(action)

        writer.update()
                    
    