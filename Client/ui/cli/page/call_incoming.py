from time import time
from ui.cli.page.page import BasePage

class IncomingCallPage(BasePage):
  def __init__(self, username):
    self.username = username
    self._ = list(map(lambda x : x.format(username), 
                 ["You got a call from {}.", "You got a call from {}..", "You got a call from {}..."]))
    self._ctr = 0
    self.start_time = time()
  
  def get_content(self) -> str:
    t = time() - self.start_time
    content = self._[self._ctr]
    content += "\n00:{:02d}".format(int(30 - t))
    self._ctr = (self._ctr + 1) % 3
    return content
  
  def get_prompt(self) -> str:
    prompt = "Do you want to accept?"
    return prompt
  
#   def get_page(self) -> str:
#     return self.__get_content() + '\n' + self.__get_prompt()

  def handle_input(self, *, user_input='', audio_client=None, **kwargs) -> int:
    if audio_client is None:
      return 1
    
    user_input = user_input.upper()
    
    if user_input == "Y":
      audio_client.accept_connection()
      audio_client.set_state(1)
      
    elif user_input == "N":
      audio_client.declined_connection()
      audio_client.set_state(-1)
    
  