
import sys
import signal
import traceback
import getopt
import os
import os.path
import functools

from sys import argv
from time import sleep
from threading import Thread

# Add parent directory into PYTHONPATH
path = os.path.join(os.path.dirname(os.path.abspath(sys.argv[0])), "../../")
sys.path.append(path)

from ui.cli.page.calling import CallPage
from ui.cli.page.call_incoming import IncomingCallPage
from ui.cli.page.on_call import OnCallPage
from ui.cli.page.page import PageLoader, UserInputHandler
from ui.cli.writer import Writer

from audio.mixer import Mixer
from audio.audio_receiver import AudioReceiver
from audio.audio_connection import AudioConnection
from audio.audio_client import AudioClient
# from audio.logger import Logger

# Configuration variables
MODE = None  # 0 : Caller ; 1 : Receiver
ACCESS_TOKEN = None
SALT = None
SENDER_USERNAME = None
RECIPIENT_USERNAME = None
IP_ADDERSS = None
PORT = None


def sigterm_handler(active_thread):
    print("SIGTERM Signal Received. Terminating all thread...")
    for thread in active_thread:
        # Each class must have terminate method
        thread.terminate()
    print("Closed")


def parse_option():
    global MODE, ACCESS_TOKEN, SALT
    global SENDER_USERNAME, RECIPIENT_USERNAME, IP_ADDERSS, PORT

    c = 0
    try:
        opts, args = getopt.getopt(argv[1:], shortopts="s:r:", longopts=[
                                   "ip=", "port=", "network=", "sender=", "recipient=", "token=", "salt=", "caller", "receiver"])
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

                if IP_ADDERSS != None and PORT != None:
                    sys.exit(2)

                IP_ADDERSS, PORT = _
                PORT = int(PORT)

                c += 2
            except Exception as e:
                traceback.print_exc()
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
                print(arg)
                ACCESS_TOKEN = bytes.fromhex(arg)
                c += 1
            except:
                sys.exit(4)

        if opt == "--salt":
            try:
                print(arg)
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

    if c != 7:
        sys.exit(3)

class Watcher():
    def __init__(self, aud: AudioClient, mixer: Mixer, page_loader: PageLoader) -> None:
        self.aud = aud
        self.mix = mixer
        self.pg = page_loader
        self._stop = False
        self._thread = Thread(target=self.watch_connection_state, daemon=True)

    def start(self):
        self._thread.start()
        
    def terminate(self):
        self._stop = True
        self._thread.join(timeout=1)

    def watch_connection_state(self):
        while not (self.aud.is_connection_ready and self._stop):
            continue

        if self._stop:
            return

        self.mixer.start()
        if MODE == 0:
            call_page = OnCallPage(RECIPIENT_USERNAME)

        else:
            call_page = OnCallPage(SENDER_USERNAME)

        self.page_loader.load_new_page(call_page)
        return


if __name__ == "__main__":
    parse_option()
    
    print(SENDER_USERNAME, RECIPIENT_USERNAME,
          IP_ADDERSS, PORT, ACCESS_TOKEN, MODE)
    
    # Initialize Object
    conn = AudioConnection(ip=IP_ADDERSS, port=PORT)
    
    aud = AudioClient(username=SENDER_USERNAME, conn=conn,
                      rcpt_username=RECIPIENT_USERNAME, token=ACCESS_TOKEN, salt=SALT)
    
    if MODE == 0:
        on_call_page = OnCallPage(RECIPIENT_USERNAME)
        page_loader = PageLoader(on_call_page)

    elif MODE == 1:
        incoming_call_page = IncomingCallPage(SENDER_USERNAME)
        page_loader = PageLoader(incoming_call_page)
    
    ## Init thread class     
    mixer = Mixer(send=aud.send_audio)
    
    watcher = Watcher(aud, mixer, page_loader)
    
    writer = Writer(page_loader=page_loader)
    
    rec = AudioReceiver(conn=conn, f=mixer.append_audio, client=aud)
    
    # Init signal watcher for sigterm to terminate used thread
    signal.signal(signal.SIGTERM, functools.partial(
        sigterm_handler, [mixer, writer, watcher, rec]))

    # start connection
    ret, _ = conn.connect_to_socket_udp()
    if ret != 0:
        print(_)
        sys.exit(1)
    
    rec.start()
    aud.register_channel()
    
    while not aud.is_channel_set():
        continue
    
    writer.start() 
    watcher.start()
    
    rc = 0
    
    while rc == 0:
        action = input()

        if page_loader.current_active_page == "ONCALLPAGE":
            rc = UserInputHandler.process_on_call_input(action)

        elif page_loader.current_active_page() == "INCOMINGCALLPAGE":
            rc = UserInputHandler.process_incoming_call_input(action)

        elif page_loader.current_active_page == "CALLPAGE":
            rc = UserInputHandler.process_calling_input(action)

        ## writer.update()
