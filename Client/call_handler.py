import os
import os.path
import sys

from process.process_handler import ProcessHandler

class CallHandler():
    def __init__(self) -> None:
        self._p_handler = None
    
    def spawn_process(self, restype, sender, recipient, token, salt, network_address):
        if restype != "CHANNEL_ALLOCATION":
            sender, recipient = recipient, sender
        
        print(sender, recipient)
        shortopts = [
                ["-s", sender],
                ["-r", recipient],
            ]
            
        longopts = [
            ["--token", token],
            ["--salt", salt],
            ["--network", network_address],
        ]
            
        longopts.append(["--caller", ""] if restype == "CHANNEL_ALLOCATION" else ["--receiver", ""])
            
        path = os.path.join(os.path.dirname(os.path.abspath(sys.argv[0])), "process", "call", "main.py")
        
        # Use this handler when trying to debug, the terminal will stay open when an error occured  
        # self._p_handler = ProcessHandler(cmd="cmd /K python {}".format(path), shortopts=shortopts, longopts=longopts)
        
        self._p_handler = ProcessHandler(cmd="python {}".format(path), shortopts=shortopts, longopts=longopts)
        self._p_handler.run()
    
    def check_process_status(self):
        if self._p_handler == None:
            return False
        
        return self._p_handler.is_spawned_process_alive()
    
    def force_stop(self):
        self._p_handler.force_terminate()
        self._p_handler.join()
        
    
    
    