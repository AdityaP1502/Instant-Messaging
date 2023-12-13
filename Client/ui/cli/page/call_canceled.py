class CallCanceled():
  def __init__(self, username, mode) -> None:
    self.username = username
    self.mode = mode
    
  def __get_content(self):
    if self.mode == 0:
      return f"Your call with {self.username} has been terminated because of the user cancel the call"
    
    return f"Your call with {self.username} has been cancelled"
  
  def __get_prompt(self):
    return "Press ENTER to exit"
  
  def get_page(self):
    return self.__get_content() + '\n' + self.__get_prompt()