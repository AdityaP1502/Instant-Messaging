from ui.cli.page.page import BasePage
class CallCanceledPage(BasePage):
  def __init__(self, username, mode) -> None:
    self.username = username
    self.mode = mode

  def get_content(self):
    if self.mode == 0:
      return [f"Your call with {self.username} has been terminated because of the user cancel the call"]
    
    return [f"Your call with {self.username} has been cancelled"]
  
  def get_header(self):
    return ""
  
  def get_prompt(self):
    return "Press ENTER to exit"
  
  # def get_page(self):
  #   return self.__get_content() + '\n' + self.__get_prompt()

  def handle_input(self, *, user_input='', **kwargs) -> int:
    return 1
  