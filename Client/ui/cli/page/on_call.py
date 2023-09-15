from time import time
class OnCallPage():
  def __init__(self, username):
    self.username = username
    self.start_time = time()
  
  def __get_content(self) -> str:
    new_time = int((time() - self.start_time)) # in seconds
    h, m, s = new_time // 3600, new_time % 3600 // 60, new_time % 3600 % 60
    
    content = "In Call With {}".format(self.username)
    content += "\n{:02d}:{:02d}:{:02d}".format(h, m, s)
    
    return content

  def __get_prompt(self) -> str:
    return ""
  
  def get_page(self) -> str:
    content = self.__get_content()
    prompt = self.__get_prompt()
    
    return content
  