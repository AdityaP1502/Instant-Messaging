
class CallPage():
  def __init__(self, username):
    self.username = username
    self._ = list(map(lambda x : x.format(username), 
                 ["Calling {}.", "Calling {}..", "Calling {}..."]))
    self.ctr = 0
    
  def __get_content(self) -> str:
    content = self._[self.ctr]
    self.ctr += 1
    self.ctr = self.ctr % 3
    
    return content

  def __get_prompt(self) -> str:
    prompt = "Do you want to cancel calling {}".format(self.username)
    return prompt
  
  def get_page(self) -> str:
    return self.__get_content() + '\n' + self.__get_prompt()