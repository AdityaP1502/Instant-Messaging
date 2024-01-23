from ui.cli.page.page import BasePage
class TimeoutPage(BasePage):
  def __init__(self, username, mode) -> None:
    self.username = username
    self.mode = mode
    
  def get_content(self):
    if self.mode == 0:
      return f"Your call with {self.username} has been terminated because of the user didn't pick up"
    
    return f"Your call with {self.username} has been terminated because of inactivity"
  
  def get_prompt(self):
    return "Press ENTER to exit"
  
  def handle_input(self, *, user_input='', **kwargs) -> int:
    return 1




