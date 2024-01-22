import sys
import os.path

from client import Client
from connection import Connection
from receiver import Receiver
from history_handler import HistoryHandler
from call_handler import CallHandler

from ui.cli.writer import Writer
from ui.cli.page.page import PageLoader, UserInputHandler
from ui.cli.page.home import HomePage
from ui.cli.page.chat_room import ChatRoomPage
from ui.cli.page.call_incoming import IncomingCallPage
from ui.cli.page.on_call import OnCallPage

CHAT_HISTORY_DATA = None
CALL_INFO = None
CHAT_HISTORY_DATA_FILEPATH = os.path.join(os.path.dirname(os.path.abspath(sys.argv[0])), "data/history.txt")
ERR = None

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Please specified the username!")
        print("Usage python main.py [username]")
        sys.exit(-1)
    
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
    
    page_loader = PageLoader(home_page)
    
    # create a writer
    writer = Writer(page_loader=page_loader)
    
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
            writer.expecting_input = True
            while rc != 1:
                action = input("")  
                
                if not conn.running:
                    break
                
                if page_loader.current_active_page() == "HOMEPAGE":
                    rc = UserInputHandler.process_home_input(writer, action, page_loader, pages)

                elif page_loader.current_active_page() == "CHATROOMPAGE":
                    rc = UserInputHandler.process_chat_input(client, writer, action, page_loader, pages)
                     
        else:
            conn.running = False
            print("Login failed. Please try again!")
           
    except KeyboardInterrupt:
        print("CTRL + C is pressed")
      
    finally:
        conn.running = False
        writer.expecting_input = False
        writer.set_err_signal(err=KeyboardInterrupt.__name__, terminate=True)
        
        if call_handler.check_process_status():
            call_handler.force_stop()
            
        HistoryHandler.write_history_data(CHAT_HISTORY_DATA_FILEPATH, CHAT_HISTORY_DATA)        
                   
    # while (action != "TERMINATE"):
    #     client.state.value = 0
    #     sys.stdout.write("\rinput action:")
    #     action = input("").upper()

    #     if not conn.running:
    #         break

    #     if action == "SENDMESSAGE":
    #         client.state.value = 1
    #         sys.stdout.write("\rrecipient:")
    #         recipient = input("")

    #         if not conn.running:
    #             break

    #         client.state.value = 2
    #         sys.stdout.write("\rmessage:")
    #         message = input("")

    #         if not conn.running:
    #             break

    #         client.send_message(recipient=recipient,
    #                             message=message,
    #                             conn=conn)

    #     elif action == "TERMINATE":
    #         print("Terminating connection!")

    #     else:
    #         print("Wrong action")
    #         continue
