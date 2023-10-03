import os
import platform
import datetime

class PageLoader():
  def __init__(self, page) -> None:
    self.__loaded_page = page
 
  def current_active_page(self) -> str:
    return type(self.__loaded_page).__name__.upper()
  
  def show(self) -> None:
    page = self.__loaded_page.get_page()
    print(page)
    
  def clear(self) -> None:
    if platform.system() == "Windows":
      os.system(command="cls")
      
  def update(self) -> None:
    self.clear()
    self.show()
    
  def load_new_page(self, new_page) -> None:
    self.__loaded_page = new_page
  
  def get_loaded_page(self) :
    return self.__loaded_page
  
class UserInputHandler():
  @staticmethod
  def process_home_input(writer, user_input : str, page_loader : PageLoader, pages):
    user_input_parsed = user_input.split(" ", maxsplit=1)
    
    if len(user_input_parsed) != 2:
      user_input_parsed.append("")
    
    command, args = user_input_parsed
    
    if command.upper() == "BACK":
      # TODO: Add extra dialog to confirm the user that they will exit the program
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
  
  @staticmethod
  def process_chat_input(client, writer, user_input : str, page_loader : PageLoader, pages):
    user_input_parsed = user_input.split(" ", maxsplit=1)
    
    if len(user_input_parsed) != 2:
      user_input_parsed.append("")
    
    command, args = user_input_parsed
    
    if command.upper() == "BACK":
      # Go back to home page
      page_loader.load_new_page(pages[0])
      writer.update()
      return 0
    
    elif command.upper() == "MESSAGE":
      # user input would be the message they want to sent to x
      recpt_username = page_loader.get_loaded_page().recpt

      client.send_message(recipient=recpt_username, message=args)

      timestamp = datetime.datetime.now().isoformat()
      # add the new messages to the history data
      page_loader.get_loaded_page().history.append(["you", timestamp, args])
      writer.update()
      # update the pages
      return 0
    
    elif command.upper() == "CALL":
      # TODO: Spawned a new process
      recpt_username = page_loader.get_loaded_page().recpt
    
    else:
      # invalid command
      return 1
  
  @staticmethod
  def process_calling_input(user_input : str):
    user_input_parsed = user_input.split(" ", maxsplit=1)
    
    if len(user_input_parsed) != 2:
      user_input_parsed.append("")
    
    command, args = user_input_parsed
    
    if command.upper() == "CANCEL":
      return 1
    
    return 0
    
  @staticmethod
  def process_incoming_call_input(user_input : str):
    user_input = user_input.upper()
    
    if user_input == "Y":
      return 0
      
    elif user_input == "N":
      return 1
  
  @staticmethod
  def process_on_call_input(user_input : str):
    user_input_parsed = user_input.split(" ", maxsplit=1)
    
    if len(user_input_parsed) != 2:
      user_input_parsed.append("")
    
    command, args = user_input_parsed
    
    if command.upper() == "HANG":
      return 1
    
  @staticmethod
  def process_user_input(user_input : str, curr_page, page_loader : PageLoader, history):
    NotImplemented

