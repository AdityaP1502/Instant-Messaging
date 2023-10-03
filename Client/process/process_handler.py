import os
import psutil
import subprocess

from multiprocessing import Process
from time import sleep

class ProcessHandler():
    def __init__(self, cmd, shortopts=[], longopts=[]):
        shortopts = map(lambda x: "=".join(x), shortopts)
        longopts = map(lambda x: "=".join(x), longopts)
        
        shortopts_str = " ".join(shortopts)
        longopts_str = " ".join(longopts)

        self._parent_pid = os.getpid() # whose pid should i use?
        self._cmd = "{} {} {}".format(cmd, shortopts_str, longopts_str)
        self._state = 0
        self._process = Process(target=self._run_child, daemon=True)
    
    def run(self):
        self._process.start()
    
    def join(self):
        self._process.join()
        
    def force_terminate(self):
        self._state = 1
            
    def _run_child(self):
        parent = psutil.Process(pid=self._parent_pid)
        child = psutil.Popen(self._cmd, creationflags=subprocess.CREATE_NEW_CONSOLE)
        
        try:
            # Run while parent and the child process is active
            while parent.status() == psutil.STATUS_RUNNING and child.status() == psutil.STATUS_RUNNING:
                if (self._state) == 1:
                    break
                sleep(1)
                
        except psutil.NoSuchProcess:
            # Raised when either parent or child is dead
            pass
            
        finally:
            try:
                child.terminate()
            except psutil.NoSuchProcess:
                pass
        return
    
    
    
    
        
        
        
        