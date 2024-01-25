import dateutil.parser
import datetime

from math import ceil
from ui.cli.page.page import BasePage

class ChatRoomPage(BasePage):
  def __init__(self, data) -> None:
    # data is a dictionary consist of all of the chat history data
    self.data = data
    self.recpt = None
    self.history = None

  def __from_history_to_text_box(self, history):
    sender, timestamp, message = history
    lpad = ""
    
    if sender == "you":
      lpad = " " * 60
    
    t = dateutil.parser.isoparse(timestamp).strftime('%m/%d %H:%M')
    box = ["{}{} {}".format(lpad, t, sender)]
    _ = ceil(len(message)/40)
      
    for i in range(_):
      box.append("{}{}".format(lpad, message[40 * i:min(40 * (i + 1), len(message))]))
    
    return box
      
  def get_header(self):
    return "You are currently messaging {}".format(self.recpt)

  def get_content(self) -> str:
    content = []
    
    for history in self.history:
      content += self.__from_history_to_text_box(history)
      
    return content
  
  def get_prompt(self) -> str:
    return "What message do you want to send:"  
    
  def set_recipient_username(self, username):
    self.recpt = username
    
    if x := self.data.get(username, None) == None:
      self.data[username] = []
    
    self.history = self.data[username]
    
#   def get_page(self) -> str:
#     content = self.__get_content()
#     prompt = self.__get_prompt()
    
#     return content + "\n" + prompt
  
  def handle_input(self, *, user_input='', client=None, writer=None, page_loader=None, pages=None, **kwargs) -> int:
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
      recpt_username = page_loader.get_loaded_page().recpt
      client.init_call(recipient=recpt_username)
      return 0
    
    else:
      # invalid command
      return 1
    
  
  