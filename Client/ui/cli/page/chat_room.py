from math import ceil
import dateutil.parser

class ChatRoomPage():
  def __init__(self, data) -> None:
    # data is a dictionary consist of all of the chat history data
    self.data = data
    self.recpt = None
    self.history = None

  def __from_history_to_text_box(self, history):
    sender, timestamp, message = history
    box = ""
    lpad = ""
    
    if sender == "you":
      lpad = " " * 60
    
    t = dateutil.parser.isoparse(timestamp).strftime('%m/%d %H:%M')
    box += "{}{} {}\n".format(lpad, t, sender)
    _ = ceil(len(message)/40)
      
    for i in range(_):
      box += "{}{}\n".format(lpad, message[40 * i:min(40 * (i + 1), len(message))])
    
    return box
      
  def __get_content(self) -> str:
    header = "You are currently messaging {}\n".format(self.recpt)
    content = ""
    
    for history in self.history:
      content += self.__from_history_to_text_box(history)
      
    return header + content
  
  def __get_prompt(self) -> str:
    return "What message do you want to send:"  
    
  def set_recipient_username(self, username):
    self.recpt = username
    
    if x := self.data.get(username, None) == None:
      self.data[username] = []
    
    self.history = self.data[username]
    
  def get_page(self) -> str:
    content = self.__get_content()
    prompt = self.__get_prompt()
    
    return content + "\n" + prompt
  
    
  
  