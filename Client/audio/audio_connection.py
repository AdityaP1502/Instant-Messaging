import socket

class AudioConnection():
    HOST = "localhost"
    PORT = 8080
    
    def __init__(self) -> None:
        self.socket : socket.socket = None
        
    def connect_to_socket_tcp(self):
        try:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        except socket.error as err:
            return -1, "socket creation failed with error {}".format(err)

        try:
            host_ip = socket.gethostbyname("localhost")
        except socket.error as err:
            return -1, "There was an error resolving the hostname"
        
        self.socket.connect((host_ip, self.PORT))
        self.socket.bind()
        self.socket.setblocking(0)
        return 0, ""
    
    def connect_to_socket_udp(self):
        try:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        except socket.error as err:
            return -1, "socket creation failed with error {}".format(err)

        try:
            host_ip = socket.gethostbyname("localhost")
        except socket.error as err:
            return -1, "There was an error resolving the hostname"
        
        self.socket.connect((host_ip, self.PORT))
        self.socket.setblocking(0)
        return 0, ""
    
    def send(self, data : bytes):
        self.socket.send(data)
    
    def terminate_connection(self):
        self.socket.close()
    
    