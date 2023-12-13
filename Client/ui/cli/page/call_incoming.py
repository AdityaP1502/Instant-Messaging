from time import time
class IncomingCallPage():
  def __init__(self, username):
    self.username = username
    self._ = list(map(lambda x : x.format(username), 
                 ["You got a call from {}.", "You got a call from {}..", "You got a call from {}..."]))
    self._ctr = 0
    self.start_time = time()
  
  def __get_content(self) -> str:
    t = time() - self.start_time
    content = self._[self._ctr]
    content += "\n00:{:02d}".format(int(30 - t))
    self._ctr = (self._ctr + 1) % 3
    return content
  
  def __get_prompt(self) -> str:
    prompt = "Do you want to accept?"
    return prompt
  
  def get_page(self) -> str:
    return self.__get_content() + '\n' + self.__get_prompt()
    
    
  