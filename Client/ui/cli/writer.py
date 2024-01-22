from threading import Thread
from multiprocessing import Queue

class Writer():
    def __init__(self, page_loader) -> None:
        self.state = 0
        self.page_loader = page_loader        
        self.job = Queue(maxsize=10)
        self.__thread = Thread(target=self.__update_page)
        self.err = None
        self.expecting_input = False
        
    def start(self) -> None:
        self.__thread.start()

    def join(self, timeout: float | None = None) -> None:
        return self.__thread.join()
    
    def update(self):
        self.state = 0

    def set_err_signal(self, err, terminate=False):
        self.err = err
        self.state = 2
        if terminate:
            self.state = 6
    
    def terminate(self):
        self.state = 2
        
    def __update_page(self):
        # while True:
        #     # update the page every 200 ms
        #     self.state = 1
        #     sleep(0.02)
        #     self.state = 0
            
        #     self.page_loader.update()
        
        while True:
            # only issue page_loader to update when received an update signal
            # via self.state
            
            # state = 0  : There is a signal to update
            # state = 1  : No signal to update, only takes job
            # state = -1 : err signal is raised
            if self.state > 2:
                print("Exception is raised")
                print(self.err)

                if self.state == 6:
                    if self.expecting_input:
                        print("Press ENTER to Exit!")
                        
                    break
            
            if self.state == 2:
                break
            
            if self.state == 0:
                self.page_loader.update()
                self.state = 1
            
            while self.state == 1 and not self.job.empty():
                sender, timestamp, message = self.job.get()
                entry = self.page_loader.get_loaded_page().data.get(sender, None)
                # t = dateutil.parser.isoparse(timestamp).strftime('%m/%d %H:%M')
                if entry == None:
                    self.page_loader.get_loaded_page().data[sender] = [[sender, timestamp, message]]

                else:
                    entry.append([sender, timestamp, message])
                    
                self.page_loader.update()
        
        print("Writer Thread stopped gracefully")                
                
                    
                
            

            
