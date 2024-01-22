import socket
import sys
from time import sleep

class Connection():
    PENDING_REQUEST_MAXSIZE = 256  # the maximum amount of pending request

    def __init__(self, ui_handler) -> None:
        self.socket = None
        self.ui_handler = ui_handler
        self.is_secure = False
        self.running = False  # the state of the connection
        self.has_registered = False
        self.request_status = [
            "" for i in range(Connection.PENDING_REQUEST_MAXSIZE + 1)
        ]
        self.request_state = [
            0 for i in range(Connection.PENDING_REQUEST_MAXSIZE + 1)
        ]
        self.request_number = -1

    def send(self, data: str, in_bytes=False):
        # send a request/data to the server
        # encrypt if use secure_connection

        self.request_number = (self.request_number +
                               1) % Connection.PENDING_REQUEST_MAXSIZE

        if self.request_state[self.request_number] == 1:
            # wait until the request is finished
            while self.request_state[self.request_number] == 1:
                continue
        
        _ = "uid={};".format(self.request_number)
        _ = _.encode()
        
        if not in_bytes:
            data = data.encode()
        
        data  = _ + data + '\n'.encode()
        
        if self.is_secure:
            pass
        
        try:
            self.socket.send(data)
            
        except BaseException as e:
                sleep(0.1)
                self.ui_handler.set_err_signal(e, terminate=1)
                self.running = False
                
                return "ERROR"
                 
        status = self.__wait_for_response(self.request_number)
        
        return status

    def fast_send(self, data : bytes):
        self.socket.send(data)
           
    def __wait_for_response(self, uid: int):
        
        while self.request_status[uid] == "":
            if not self.running:
                return "ERROR"
            continue

        status = self.request_status[uid]
        self.request_status[uid] = ""
        self.request_state[uid] = 0

        return status

    def connect_to_socket(self, host_ip: str, port: int):
        """ Establish a socket connection to the server. 
            Raised a socketerror connection. 
        """
        # create a socket

        try:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        except socket.error as err:
            return -1, "socket creation failed with error {}".format(err)

        self.socket.connect((host_ip, port))
        self.socket.setblocking(0)
        self.running = True

    def secure_connection(self):
        """Establish TLS connection to the server
    """
        pass

    def terminate_connection(self, username: str):
        """close the connection the server
    """

        message = "reqtype=TERMINATE;payload=username={}".format(username)
        status = self.send(message)

        if status == "OK":
            self.socket.close()
            self.running = False
            return 0

        return 1
