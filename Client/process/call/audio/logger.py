from threading import Thread
import queue

class Logger():
    def __init__(self) -> None:
        self._t = Thread(target=self.run)
        self._log = {}
        self.job = queue.Queue(maxsize=100)
        self._stop = False
        
    def register_event(self, name):
        self._log[name] = ""
        print("Registered event {} to logger".format(name))
        
    def terminate(self):
        self._stop = True
        self.job.put(None)
        
    def emit(self, event_name, data):
        event = LogEvent(name=event_name, data=data)
        self.job.put(event)
    
    def start(self):
        self._stop = False
        self._t.start()
        
    def run(self):
        while not self._stop:
            event = self.job.get()
            
            if event == None:
                break
            
            try:
                log = self._log[event.name]
                log += event.data + "\n"
            except KeyError:
                print("Invalid event name {}, discarding event!".format(event.name))
                  
            filename = "log{}.txt".format(event.name.replace(" ", ""))
            
            with open(filename, "a") as f:
                f.write(log)
                
        print("Logger Thread Closed")
                
class LogEvent():
    def __init__(self, name, data) -> None:
        self.name = name
        self.data = data
        
            