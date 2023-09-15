from time import sleep
from threading import Thread, Event

from ui.cli.page.call_incoming import IncomingCallPage
from ui.cli.page.calling import CallPage
from ui.cli.page.on_call import OnCallPage

from audio.mixer import Mixer

class TimeoutHandler():
  def __init__(self) -> None:
    self._t = Thread(target=self.run)
    self._stop = Event()
    self._timeout_s = 30
  
  def start(self):
    self._t.start()
  
  def run(self):
    while not self.stopped():
      self._stop.wait(self._timeout_s)
      break
    
  def stop(self):
    self._stop.set()
  
  def stopped(self):
    return self._stop.is_set()
  
  def is_alive(self):
    return self._t.is_alive()
  
class CallHandler():
  def __init__(self, ui_handler, page_loader, client) -> None:
    self._ui_handler = ui_handler
    self._page_loader = page_loader
    self._client = client
    self._username = None
    self._state = -1
    self._prev_page = None
    self._t = None
    self.mixer = Mixer(client=client)
  
  def _start(self, username, target):
    self._username = username
    self._prev_page = self._page_loader.get_loaded_page()
    self._t = Thread(target=target, daemon=True)
    self._t.start()
  
  def get_caller(self):
    return self._username

  def set_state(self, state):
    self._state = state
  
  def wait_call_response(self):
    calling_page = CallPage(self._username)
    self._page_loader.load_new_page(calling_page)
    
    while self._state == 0:
      self._ui_handler.update()
      sleep(1)
      
    self.process_caller_state()
      
  def wait_user_response(self):
    call_page = IncomingCallPage(self._username)
    
    # load incoming call page
    self._page_loader.load_new_page(call_page)
    
    waiter = TimeoutHandler()
    waiter.start()
    
    self._ui_handler.update()

    while waiter.is_alive():
      self._ui_handler.update()
      sleep(1)
      
      if self._state == 0:
        continue
      
      waiter.stop()
      break
    
    self.process_incoming_state()
  
  def process_caller_state(self):
    if self._state == 3:
      self.start_calling()
      
    # go back to previous page
    self._page_loader.load_new_page(self._prev_page)
    self._ui_handler.update()
    # reset state
    self.clear()
    
  def process_incoming_state(self):
    if self._state == 1:
      self._client.timeout_call(self._username)
      
    elif self._state == 2:
      self._client.decline_call(self._username)
      
    elif self._state == 3:
      self._client.accept_call(self._username)
      self.start_calling()
      
    # go back to previous page
    self._page_loader.load_new_page(self._prev_page)
    self._ui_handler.update()
    
    # reset state
    self.clear()
    
  def initiate_call(self, username):
    self._state = 0
    self._start(username=username, target=self.wait_call_response)
    self._client.init_call(self._username)
    
  def start_incoming_call(self, username):
    self._state = 0
    self._start(username=username, target=self.wait_user_response)
  
  def clear(self):
    self._username = None
    self._prev_page = None
    self._state = -1
    self._t = None
    
  def start_calling(self):
    on_call_page = OnCallPage(self._username)
    self._page_loader.load_new_page(on_call_page)
    self.mixer.start(self._username)
    
    while (self._state != 4 and self._state != 5 and not self.mixer._stop):
      # self._ui_handler.update()
      sleep(2)
      
    self.mixer.terminate()
      