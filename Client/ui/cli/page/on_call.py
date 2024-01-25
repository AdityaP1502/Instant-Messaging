from time import time
from ui.cli.page.page import BasePage
class OnCallPage(BasePage):
  def __init__(self, username):
    self.username = username
    self.start_time = time()
  
  def get_content(self) -> str:
    new_time = int((time() - self.start_time)) # in seconds
    h, m, s = new_time // 3600, new_time % 3600 // 60, new_time % 3600 % 60
    
    content = [
        "In Call With {}".format(self.username), 
        "{:02d}:{:02d}:{:02d}".format(h, m, s)
    ]
    
    return content

  def get_prompt(self) -> str:
    return ""
  
  def get_header(self):
    return ""
#   def get_page(self) -> str:
#     content = self.__get_content()
#     prompt = self.__get_prompt()
    
#     return content
  
  def handle_input(self, *, user_input='', **kwargs) -> int:
    user_input_parsed = user_input.split(" ", maxsplit=1)
    
    if len(user_input_parsed) != 2:
      user_input_parsed.append("")
    
    command, args = user_input_parsed
    
    if command.upper() == "HANG":
      return 1