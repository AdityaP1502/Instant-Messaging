
import sys
import signal
import traceback
import getopt
import os
import os.path
import functools

from sys import argv
from time import sleep
from threading import Thread, Timer

# Add parent directory into PYTHONPATH
path = os.path.join(os.path.dirname(os.path.abspath(sys.argv[0])), "../../")
sys.path.append(path)

from audio.mixer import Mixer
from audio.audio_receiver import AudioReceiver
from audio.audio_connection import AudioConnection
from audio.audio_client import AudioClient

from ui.cli.page.page import PageLoader, UserInputHandler
from ui.cli.page.calling import CallPage
from ui.cli.page.call_incoming import IncomingCallPage
from ui.cli.page.on_call import OnCallPage
from ui.cli.page.call_declined import DeclinedPage
from ui.cli.page.call_timeout import TimeoutPage
from ui.cli.writer import Writer



# from audio.logger import Logger

# Configuration variables
MODE = None  # 0 : Caller ; 1 : Receiver
ACCESS_TOKEN = None
SALT = None
SENDER_USERNAME = None
RECIPIENT_USERNAME = None
IP_ADDERSS = None
PORT = None


# TODO: Handle receiver input (such as accept declined and timeout)

def sigterm_handler(active_thread):
    print("SIGTERM Signal Received. Terminating all thread...")
    for thread in active_thread:
        thread.terminate()
    print("Closed")

def sigint_handler(active_thread):
    print("CTRL + C event is sent")
    
    for thread in active_thread:
        thread.terminate()
        
    input()

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

    print(opts)
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
                print("Error: Invalid Port number. Exiting...")
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
                traceback.print_exc(e)
                sys.exit(1)

        if opt == "--sender":
            if SENDER_USERNAME == None:
                SENDER_USERNAME = arg
                c += 1
                continue
            
            print("Error: Assigning sender twice. Exiting...")
            sys.exit(2)

        if opt == "--recipient":
            if RECIPIENT_USERNAME == None:
                RECIPIENT_USERNAME = arg
                c += 1
                continue
            
            print("Error: Assigning recipient twice. Exiting...")
            sys.exit(2)

        if opt == "--token":
            try:
                print("token:", arg)
                ACCESS_TOKEN = bytes.fromhex(arg)
                c += 1
            except:
                print("Error: Cannot decode token. Exiting...")
                sys.exit(4)

        if opt == "--salt":
            try:
                print("salt:", arg)
                SALT = bytes.fromhex(arg)
                c += 1
            except:
                print("Error: Cannot decode salt. Exiting...")
                sys.exit(4)

        if opt == "--caller":
            if MODE != None:
                print("Error: Assigning mode option twice. Exiting...")
                sys.exit(2)

            MODE = 0
            c += 1

        if opt == "--receiver":
            if MODE != None:
                print("Error: Assigning mode option twice. Exiting...")
                sys.exit(2)

            MODE = 1
            c += 1

    if c != 7:
        print("Error: Invalid number of options. Exiting...")
        sys.exit(3)

class Updater():
    def __init__(self, writer, aud):
        self._writer = writer
        self._client = aud
        self._stop = False
        self._thread = Thread(target=self.run)
        
    def run(self):
        while (not self._stop) and self._client.get_state() >= 0:
            sleep(1)
            self._writer.update()
            
        return

    def start(self):
        self._thread.start()
        
    def terminate(self):
        self._stop = True
        self._thread.join()
    
class Watcher():
    def __init__(self, aud: AudioClient, mixer: Mixer, page_loader: PageLoader, mode : int, writer : Writer = None) -> None:
        self.aud = aud
        self.mix = mixer
        self.pg = page_loader
        self.mode = mode
        self._timeout_s = 30 + 100 if mode == 0 else 30 # 30 seconds + 100 seconds grace period 
        self._stop = False
        self._timeout = False
        self._thread = Thread(target=self.watch_connection_state)
        self._ui_handler = writer
        
    def start(self):
        self._thread.start()
        
    def terminate(self):
        self._stop = True
        self._thread.join(timeout=1)

    def timeout_handler(self):
        if self.aud.get_state() != 0 or self._stop:
            return
        
        # if self._ui_handler != None:
        #     self._ui_handler.pause()
            
        # print("Reach Timeout!...")
        self._timeout = True
        
        self.aud.set_state(-2)
        
        if self.mode == 1:
            # sent timeout packet to server
            self.aud.connection_timeout()
        
        # load timeout page
        timeout = TimeoutPage(username=RECIPIENT_USERNAME, mode=self.mode)
        self.pg.load_new_page(timeout)
        
        # if self._ui_handler != None:
        #     self._ui_handler.un_pause()
            
    def watch_connection_state(self):
        timer = Timer(self._timeout_s, self.timeout_handler)
        
        timer.start()
        
        while not self.aud.get_state() != 0 or self._stop:
            continue
        
        timer.cancel()
        
        if timer.finished:
            if self._timeout == True:
                return
             
        if self._stop:
            return

        if self.aud.get_state() == -1:
            # load call declined page
            declined = DeclinedPage(username=RECIPIENT_USERNAME, mode=self.mode)
            self.pg.load_new_page(declined)
            return
        
        elif self.aud.get_state() == -2:
            timeout = TimeoutPage(username=RECIPIENT_USERNAME, mode=self.mode)
            self.pg.load_new_page(timeout)
            return
        
        
        print("Opening mic and player!")
        
        self.mix.start()
        
        if MODE == 0:
            call_page = OnCallPage(RECIPIENT_USERNAME)

        else:
            call_page = OnCallPage(SENDER_USERNAME)

        self.pg.load_new_page(call_page)
        
        return
         
if __name__ == "__main__":
    parse_option()
    
    print(SENDER_USERNAME, RECIPIENT_USERNAME,
          IP_ADDERSS, PORT, ACCESS_TOKEN, MODE)
    
    # Initialize Object
    if MODE == 0:
        on_call_page = CallPage(RECIPIENT_USERNAME)
        page_loader = PageLoader(on_call_page)

    elif MODE == 1:
        incoming_call_page = IncomingCallPage(RECIPIENT_USERNAME)
        page_loader = PageLoader(incoming_call_page)
        
    conn = AudioConnection(ip=IP_ADDERSS, port=PORT)
    
    aud = AudioClient(username=SENDER_USERNAME, conn=conn,
                      rcpt_username=RECIPIENT_USERNAME, token=ACCESS_TOKEN, salt=SALT)
    
    ## Init thread class     
    mixer = Mixer(send=aud.send_audio)
    
    writer = Writer(page_loader=page_loader)
    
    watcher = Watcher(aud, mixer, page_loader, MODE, writer=writer)
    
    rec = AudioReceiver(conn=conn, f=mixer.append_audio, client=aud)
    
    updater = Updater(writer=writer, aud=aud)
    
    # Init signal watcher for sigterm to terminate used thread
    signal.signal(signal.SIGTERM, functools.partial(
        sigterm_handler, [mixer, writer, watcher, rec, updater]))
    
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
    updater.start()
    watcher.start()
    
    rc = 0
    
    while rc == 0:
        action = input()

        # Do some stuff here
        
        if page_loader.current_active_page() == "ONCALLPAGE":
            rc = UserInputHandler.process_on_call_input(action)

        elif page_loader.current_active_page() == "INCOMINGCALLPAGE":
            rc = UserInputHandler.process_incoming_call_input(aud, action)

        elif page_loader.current_active_page() == "CALLPAGE":
            rc = UserInputHandler.process_calling_input(action)

        elif page_loader.current_active_page() == "DECLINEDPAGE":
            rc = UserInputHandler.process_declined_input(action)
            
        elif page_loader.current_active_page() == "TIMEOUTPAGE":
            rc = UserInputHandler.process_timeout_input(action)
            
        ## writer.update()

    aud.terminate_channel()
    
    for t in [mixer, writer, watcher, rec, updater]:
        t.terminate()
    
      
