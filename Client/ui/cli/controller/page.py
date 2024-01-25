import curses
import functools

from ui.cli.controller.keyboard.shortcut_handler import ShortcutHandler
class InputWindow():
    def __init__(self, height : int, width : int, start_y : int,
                 start_x : int, handler : ShortcutHandler):
      
        self._handler = handler
        self._window = curses.newwin(height, width, start_y, start_x)
        self._window.keypad(True)
    
    def get_input(self):
        user_input = ""
        while True:
            c = self._window.getch()
            
            if c == 3:
                raise KeyboardInterrupt
            
            elif c >= 32 and c <= 126:
                user_input += chr(c)
                
            elif c == curses.KEY_ENTER or c == 10 or c == 13:
                # Enter key pressed, process the user input
                self._handler.handle_enter(user_input)
                self.clear()
                    
                break
                
            elif c == curses.KEY_BACKSPACE or c == 8:
                user_input = user_input[:-1]
                self.clear()
                self._window.addstr(0, 0, user_input)
                
            elif c == curses.KEY_UP:
                self._handler.handle_scroll_up(user_input)
                        
            elif c == curses.KEY_DOWN:
                self._handler.handle_scroll_down(user_input)
            
            elif c == 417:
                self._handler.handle_scroll_alt_a(user_input)
                            
            elif c == 418 :
                self._handler.handle_scroll_alt_b(user_input)
                         
        return user_input
    
    
    def clear(self):
        self._window.clear()
        self._window.refresh()

class Header():
    def __init__(self, width, start_y):
        self._window = curses.newwin(1, width, start_y, 0)
        self._width = width
        
    def set_header(self, string):
        # if a string length is greater than the width, then the string will be truncated
        header_string = string[:self._width - 1]
        self._window.addstr(0, 0, header_string)
        self._window.refresh()
        
    def clear(self):
        self._window.clear()
        self._window.refresh()
    
class ViewWindow():
    def __init__(self, height, width, start_y, start_x):
        self._x = start_x
        self._y = start_y
        self._w = width
        self._h = height
        self._min_history = 0
        self._window = curses.newwin(height, width, start_y, start_x)
        self._cursor_y = 0
        self._history = []
        self._scroll_length = 5
        self._skip_length = height // 2
        self._state = 0
        self._after_scroll = False
        self._max_retention = 100
    
    def clear_history(self):
      self._history = []
      self._min_history = 0
      self._cursor_y = 0
      
    def _add_to_history(self, data):
        if len(data) > self._max_retention:
            del data[0]
            
        self._history.append(data)
        
    def _print(self, string):
        if self._cursor_y >= self._h:
            self.clear()
            self._cursor_y = 0
            self.print(*self._history[-self._skip_length:], new=False)
            self._min_history += self._skip_length
            self._state = self._min_history
        
        self._window.addstr(self._cursor_y, 0, string)
        self._cursor_y += 1
        return
    
    def print(self, *usr_ins, new=True):
        # check if min history is at the latest
        if self._after_scroll and new:
            self._min_history = self._state
            self._cursor_y = 0
            self._window.clear()
            self.print(*self._history[self._state:], new=False)
            self._after_scroll = False
                       
        for x in usr_ins:
            while len(x) > self._w:
                x_truncate, x = x[:self._w - 1], x[self._w - 1:]
                self._print(x_truncate)
                if new:
                    self._add_to_history(x_truncate)
            
            if x != "":
                self._print(x)
                if new:
                    self._add_to_history(x)

        self._window.refresh()
    
    def scroll(self, direction):
        self._after_scroll = True
           
        if direction > 0:
            if self._min_history - self._scroll_length < 0:
                return
            
            self._min_history -= self._scroll_length
            
        else:
            if self._min_history + self._h >= len(self._history):
                return
            
            self._min_history += self._scroll_length 
        
        
        self._cursor_y = 0
        self.clear()
        self.print(*self._history[self._min_history: self._min_history + (self._h)], new=False)
        return
    
    def clear(self):
        self._window.clear()
        self._window.refresh()
class ScreenHandler():
  def __init__(self, page, shortcut_handler : ShortcutHandler) -> None:
    self._view_header = Header(100, 0)
    self._view_window = ViewWindow(10, 100, 1, 0)
    
    shortcut_handler.set_scroll_up_cb(
      self._scroll_up
    )
    
    shortcut_handler.set_scroll_down_cb(
      self._scroll_down
    )
    
    self._input_header = Header(100, 11)
    self._input_window = InputWindow(
        10,
        50,
        12,
        0,
        shortcut_handler,
    )
    
    self.__loaded_page = page

  
  def _scroll_up(self, user_input):
    return self._view_window.scroll(1)
  
  def _scroll_down(self, user_input):
    return self._view_window.scroll(-1)
    
  def current_active_page(self) -> str:
    return type(self.__loaded_page).__name__.upper()
  
  def show(self) -> None:
    page = self.__loaded_page
    
    header = page.get_header()
    content = page.get_content()
    prompt = page.get_prompt()
    
    self._view_header.set_header(header)
    self._view_window.clear_history()
    self._view_window.print(*content)
    
    self._input_header.set_header(prompt)
    
  def update(self) -> None:
    self._view_header.clear()
    self._view_window.clear()
    self._input_header.clear()
    
    self.show()
    
  def load_new_page(self, new_page) -> None:
    self.__loaded_page = new_page
  
  def get_loaded_page(self) :
    return self.__loaded_page

  def input(self):
    user_input = self._input_window.get_input()
    return user_input
    

