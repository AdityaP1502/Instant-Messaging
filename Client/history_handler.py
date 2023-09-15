class HistoryHandler():
    @staticmethod
    def parse_history(history: str):
        content = history.split(sep=",", maxsplit=2)
        username, timestamp, raw_messages = content
        messages = raw_messages[1: -1]  # discard the "" from the messages

        return [username, timestamp, messages]

    @staticmethod
    def read_history_data(filepath):
      dict = {}
      read_line = ""
      
      with open(filepath, mode="r") as f:
        read_line = f.readline()[:-1]
        
        while read_line != "END":
          if read_line == "START":
            chat_history = []
            username = f.readline()[:-1]
            read_line = f.readline()[:-1] # the first history data
            
            while read_line != "START" and read_line != "END":
              chat_history.append(HistoryHandler.parse_history(read_line))
              read_line = f.readline()[:-1]

              
            dict[username] = chat_history
            
      return dict
    
    @staticmethod
    def write_history_data(filepath, chat_data):
      content = ""
      formatter = lambda x: "{},{},\"{}\"".format(*x)
      
      with open(filepath, mode="w") as f:
        for (username, history) in chat_data.items():
          
          if len(history) == 0:
            continue
          
          content += "START\n{}\n".format(username)
          content += "\n".join(map(formatter, history))
          content += "\n"
        content += "END\n"
        f.write(content)
        
if __name__ == "__main__":
  HistoryHandler.read_history_data("./data/history.txt")