class HomePage():
  def __init__(self, data, username):
    self.data = data
    self.username = username
  
  def __get_content(self) -> str:
    content = "{}, welcome to Aditya Messenger\nChat:\n".format(self.username)

    # get all chat recent data 
    # do something with data
    for (username, chat_history) in self.data.items():
      content += "{}\n{} : {}\n".format(username, chat_history[-1][0], chat_history[-1][2])
      
    return content

  def __get_prompt(self) -> str:
    prompt = "Who do you want to chat:"
    return prompt

  def get_page(self) -> str:
    content = self.__get_content()
    prompt = self.__get_prompt()
    
    return content + "\n" + prompt
  

    
  