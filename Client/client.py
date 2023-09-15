import datetime

from multiprocessing import Value, Queue
from connection import Connection
class Client():
    def __init__(self, username: str, buffer : Queue, conn : Connection) -> None:
        self._conn = conn
        self._audio_conn = None 
        self._username = username
        self._state = Value("i", -1)
        self.message_buffer = buffer
            
    def get_state(self):
        """get the state of the client

        Returns:
            state : int
        """

        with self._state.get_lock():
            return self._state.value

    def get_current_prompt(self):
        """return the current prompt

        Returns:
            str : current prompt for the user input
        """
        state = self.get_state()

        if state == 0:
            return 'input action:'

        elif state == 1:
            return 'recipient:'

        elif state == 2:
            return 'message:'

        return ''
        
    def check_in(self):
        """Register the user username to the socket
        """

        message = "reqtype=CHECKIN;payload=username={}".format(self._username)

        status = self._conn.send(message)

        if status == "OK":
            self.has_registered = True
            return 0

        else:
            return 1
    
    def ready(self):
        message = "reqtype=READY;payload=username={}".format(self._username)

        status = self._conn.send(message)

        if status == "OK":
            self.has_registered = True
            return 0

        else:
            return 1
            
    def fetch(self):        
        message = "reqtype=FETCH;payload=username={}".format(self._username)
        
        status = self._conn.send(message)
        
        if status == "OK":
            return 0
        
        else:
            # TODO: Create a retry protocol
            return 1
        
    def send_message(self, message: str, recipient: str):
        """Send a message to recipient via the connection made with the server

        Args:
            message (str): Message to be sent
            recipient (str): the recipient of the message
            conn (Connection) : Connection to the channel
        """
        ts = datetime.datetime.now().isoformat()
        
        message = "reqtype=SENDMESSAGE;payload=sdr={};rcpt={};timestamp={};message={}".format(
            self._username, recipient, ts, message)

        self._conn.send(message)

    def send_audio(self, audio : bytes, recipient : str):
        message = "uid=0;reqtype=SENDAUDIO;payload=sdr={};rcpt={};message=".format(self._username, 
                                                                             recipient)
        message = message.encode()
        message += audio
        message += '\n'.encode()
        
        self._conn.fast_send(message)
        
    def init_call(self, recipient: str):
        """Send a call notification to the recipient via the server

        Args:
            recipient (str): the user that want to be called
        """

        message = "reqtype=INITIATECALL;payload=sdr={};rcpt={}".format(self._username, recipient)
        self._conn.send(message)

    def accept_call(self, recipient: str):
        """_summary_

        Args:
            recipient (str): _description_
        """
        message = "reqtype=ACCEPTCALL;payload=sdr={};rcpt={}".format(self._username, recipient)
        self._conn.send(message)
    
    def decline_call(self, recipient : str):
        message = "reqtype=DECLINECALL;payload=sdr={};rcpt={}".format(self._username, recipient)
        self._conn.send(message)
        
    def timeout_call(self, recipient : str):
        message = "reqtype=TIMEOUTCALL;payload=sdr={};rcpt={}".format(self._username, recipient)
        self._conn.send(message)
          
    def end_call(self, recipient: str, conn : Connection):
        """end the call to the recipient

        Args:
            recipient (str): _description_
        """
        
        message = "reqtype=TERMINATECALL;payload=sdr={};rcpt={}".format(self._username, recipient)
        self._conn.send(message)

    def end_connection(self):
        """
        End the connection
        """
    
    def flush_buffer(self, chat_data : dict):
        """Emptied buffer data and store it in chat data

        Args:
            chat_data (dict): chat history data. To store data in chat_data
            the key is the chat_room name and the value is a list of chat history daya
            with format [sender, timestamp, message]
        """
        
        # flush chat data from buffer to chat_history
        while not self.message_buffer.empty():
            # data format : [sender, timestamp, message]
            data = self.message_buffer.get()
            
            if (x := chat_data.get(data[0], None)) == None:
                chat_data[data[0]] = [data]
                return 0
                
            x.append(data)
            
        return 0