class ShortcutHandler():
  def __init__(self) -> None:
    self._scroll_up_cb = self._base_fnc
    self._scroll_down_cb = self._base_fnc
    self._alt_a_cb = self._base_fnc
    self._alt_b_cb = self._base_fnc
    self._handle_enter_cb = self._base_fnc
  
  def _base_fnc(self, user_input):
    return None
  
  def set_scroll_up_cb(self, cb):
    self._scroll_up_cb = cb
  
  def set_scroll_down_cb(self, cb):
    self._scroll_down_cb = cb
  
  def set_alt_a_cb(self, cb):
    self._alt_a_cb = cb
  
  def set_alt_b_cb(self, cb):
    self._alt_b_cb = cb
    
  def set_handle_enter_cb(self, cb):
    self._handle_enter_cb = cb
  
  def handle_scroll_alt_a(self, user_input):
    return self._alt_a_cb(user_input)
  
  def handle_scroll_alt_b(self, user_input):
    return self._alt_b_cb(user_input)
  
  def handle_scroll_up(self, user_input):
    return self._scroll_up_cb(user_input)
  
  def handle_scroll_down(self, user_input):
    return self._scroll_down_cb(user_input)
  
  def handle_enter(self, user_input):
      return self._handle_enter_cb(user_input)