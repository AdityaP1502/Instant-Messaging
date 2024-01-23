from ui.cli.page.page import BasePage
class CallPage(BasePage):
  def __init__(self, username):
    self.username = username
    self._ = list(map(lambda x : x.format(username), 
                 ["Calling {}.", "Calling {}..", "Calling {}..."]))
    self.ctr = 0
    
  def get_content(self) -> str:
    content = self._[self.ctr]
    self.ctr += 1
    self.ctr = self.ctr % 3
    
    return content

  def get_prompt(self) -> str:
    prompt = "Do you want to cancel calling {}".format(self.username)
    return prompt
  
#   def get_page(self) -> str:
#     return self.__get_content() + '\n' + self.__get_prompt()

  def handle_input(self, *, user_input='', **kwargs) -> int:
    user_input_parsed = user_input.split(" ", maxsplit=1)
    
    if len(user_input_parsed) != 2:
      user_input_parsed.append("")
    
    command, args = user_input_parsed
    
    if command.upper() == "HANG":
      return 1
    
    return 0
  
  