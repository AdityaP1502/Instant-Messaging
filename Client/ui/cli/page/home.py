from ui.cli.page.page import BasePage
class HomePage(BasePage):
  def __init__(self, data, username):
    self.data = data
    self.username = username
  
  def get_content(self) -> str:
    content = "{}, welcome to Aditya Messenger\nChat:\n".format(self.username)

    # get all chat recent data 
    # do something with data
    for (username, chat_history) in self.data.items():
      content += "{}\n{} : {}\n".format(username, chat_history[-1][0], chat_history[-1][2])
      
    return content

  def get_prompt(self) -> str:
    prompt = "Who do you want to chat:"
    return prompt

#   def get_page(self) -> str:
#     content = self.__get_content()
#     prompt = self.__get_prompt()
    
#     return content + "\n" + prompt
  def handle_input(self, *, user_input='', writer=None, page_loader=None, pages=None, **kwargs) -> int:
    user_input_parsed = user_input.split(" ", maxsplit=1)
    
    if len(user_input_parsed) != 2:
      user_input_parsed.append("")
    
    command, args = user_input_parsed
    
    if command.upper() == "BACK":
      #TODO: Add extra dialog to confirm the user that they will exit the program
      return 1
    
    elif command.upper() == "CHAT":
      # User want to chat someone
      # args =  username the user want to chat with
      # set the username to chat_room pages
      try:
        pages[1].set_recipient_username(args)
      except Exception as e:
        print("Error")
        print(e)
        return 1
    
      page_loader.load_new_page(pages[1])
      writer.update()
    
    else:
      # Invalid Command 
      return 1
    
  

    
  