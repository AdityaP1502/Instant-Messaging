from abc import ABC, abstractmethod

class BasePage(ABC):
    @abstractmethod
    def get_content(self) -> list[str]:
        pass
    
    @abstractmethod
    def get_prompt(self) -> str:
        pass
    
    @abstractmethod
    def get_header(self) -> str:
        pass
    
    @abstractmethod
    def handle_input(self, *, user_input='', **kwargs) -> int:
        """Handle input that user entered

        Args:
            user_input (str, optional): User command. Defaults to ''.

        Returns:
            int: Return code. return 1 to exit the main loop. 
        """
        pass
    