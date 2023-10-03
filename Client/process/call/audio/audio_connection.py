import socket

class AudioConnection():    
    def __init__(self, ip, port) -> None:
        self.socket : socket.socket = None
        self.ip = ip
        self.port = port
        
    def connect_to_socket_tcp(self):
        try:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        except socket.error as err:
            return -1, "socket creation failed with error {}".format(err)


        self.socket.connect((self.ip, self.port))
        self.socket.bind()
        self.socket.setblocking(0)
        
        return 0, ""
    
    def connect_to_socket_udp(self):
        try:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        except socket.error as err:
            return -1, "socket creation failed with error {}".format(err)


        self.socket.connect((self.ip, self.port))
        self.socket.setblocking(0)
        return 0, ""
    
    def send(self, data : bytes):
        self.socket.send(data)
    
    def terminate_connection(self):
        self.socket.close()
    
    