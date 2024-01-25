import sys
import os.path
import curses

from client import Client
from connection import Connection
from receiver import Receiver
from history_handler import HistoryHandler
from call_handler import CallHandler

from ui.cli.controller.writer import Writer
from ui.cli.controller.page import ScreenHandler
from ui.cli.controller.keyboard.shortcut_handler import ShortcutHandler
from ui.cli.page.home import HomePage
from ui.cli.page.chat_room import ChatRoomPage

CHAT_HISTORY_DATA = None
CALL_INFO = None
ERR = None
CHAT_HISTORY_DATA_FILEPATH = os.path.join(os.path.dirname(os.path.abspath(sys.argv[0])), "data/history.txt")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Please specified the username!")
        print("Usage python main.py [username]")
        sys.exit(-1)
    
    scr = curses.initscr()
    curses.raw()
    username = sys.argv[1]

    # read history data
    CHAT_HISTORY_DATA = HistoryHandler.read_history_data(CHAT_HISTORY_DATA_FILEPATH)
    
    # creating all necessary page
    home_page = HomePage(CHAT_HISTORY_DATA, username=username)
    chat_room_page = ChatRoomPage(CHAT_HISTORY_DATA)
    # IncomingCallPage = IncomingCallPage(CALL_INFO)
    # OnCallPage = OnCallPage(CALL_INFO)
    
    pages = [home_page, chat_room_page] # list of pages
    
    # initiate page loader with home page
    shortcut_handler = ShortcutHandler()
    screen_handler = ScreenHandler(home_page, shortcut_handler)
    
    # create a writer
    writer = Writer(page_loader=screen_handler)
    
    conn = Connection(ui_handler=writer)
    client = Client(username, buffer=writer.job, conn=conn)
    call_handler = CallHandler()
    
    # # create a writer thread
    # writer_t = Writer(name="writer_thread", args=(client, ), daemon=True)
    # writer_t.start()
        
    # connect and register channel
    try:
        conn.connect_to_socket(host_ip="127.0.0.1", port=6565)
    except Exception as e:
        print("Cant connect to server!")
        print(e)
        sys.exit(-1)        
        
    receiver = Receiver(buffer=client.message_buffer, conn=conn, call_handler=call_handler, ui_handler=writer)
    receiver.start()

    rc = client.check_in()
    
    if rc == 1:
        print("Login Failed. Please try again!")
        sys.exit(1)
        
    action = ""

    rc = client.fetch()
    client.flush_buffer(CHAT_HISTORY_DATA)
    rc = client.ready()
    # sys.exit(1)
    try:
        if rc == 0:
            writer.start()
            writer.update()
            writer.expecting_input = True
            while rc != 1:
                # action = input("")  
                action = screen_handler.input()
                
                if not conn.running:
                    break
                
                rc = screen_handler.get_loaded_page().handle_input(
                    user_input = action,
                    writer= writer,
                    page_loader=screen_handler,
                    pages=pages,
                    client=client
                )
                     
        else:
            conn.running = False
            print("Login failed. Please try again!")
           
    # except KeyboardInterrupt:
    #     print("CTRL + C is pressed")
    #     writer.set_err_signal(err=type(KeyboardInterrupt).__name__, terminate=True)
    
    except Exception as e:
        print(e)
        writer.set_err_signal(err=type(e).__name__, terminate=True)
        
    finally:
        curses.nocbreak()
        scr.keypad(False)
        curses.echo()
        curses.endwin()
        curses.noraw()
        
        conn.running = False
        writer.expecting_input = False
        
        writer.terminate()
        if call_handler.check_process_status():
            call_handler.force_stop()
            
        HistoryHandler.write_history_data(CHAT_HISTORY_DATA_FILEPATH, CHAT_HISTORY_DATA)        