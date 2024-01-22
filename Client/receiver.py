import os
import os.path
import sys
import socket
import errno

from time import sleep
from threading import Thread
from history_handler import HistoryHandler

class Receiver():
    EOF = "\n"
    def __init__(self, buffer, conn, ui_handler, call_handler, buffer_size=4096):
        self._ui_handler = ui_handler
        self._call_handler = call_handler
        self._buffer_size = buffer_size
        self._conn = conn
        self._t = Thread(target=self.receive, daemon=True)
        self.buffer = buffer

    def start(self) -> None:
        return self._t.start()

    def join(self, timeout: float | None = None) -> None:
        return self._t.join(timeout)

    def receive(self):
        # where the payload size is fixed
        
        while not self._conn.has_registered or self._conn.running:
            try:
                msg = self._conn.socket.recv(self._buffer_size)
                incoming_response = msg.decode('UTF-8')
                # current_prompt = self.client.get_current_prompt()

                # if current_prompt == '':
                #   break

                packets = incoming_response.split(Receiver.EOF)

                for packet in packets[:-1]:
                    self.process_incoming_response(packet)
                    
            except socket.error as e:
                err = e.args[0]
                if err == errno.EAGAIN or err == errno.EWOULDBLOCK:
                    sleep(1)
                    continue
                else:
                    self._ui_handler.set_err_signal(e, terminate=1)
                    self._conn.running = False
                    break

                # sys.stdout.write('\rGetting {}\n{}'.format(msg, current_prompt))

    def __parse_message(self, packet: str):
        """pass the packet that is received
        """

        a = packet.split(";", 1)
        try:
            restype, payload = a
        except Exception as e:
            print(e)
            raise Exception
                 
        _, restype_value = restype.split("=", 1)
        _, payload_value = payload.split("=", 1)

        return restype_value, payload_value

    def process_incoming_response(self, packet: str):
        """process response
        """
        restype, payload = self.__parse_message(packet)
        
        if restype == None:
            return
            
        uid_integer = -1
        
        if restype == "OK":
            _, uid_value = payload.split("=", 1)
            try:
                uid_integer = int(uid_value)
                self._conn.request_status[uid_integer] = "OK"
                self._conn.request_state[uid_integer] = 0

            except ValueError as e:
                self._conn.request_status[-1] = "Cannot Get UID in {}".format(restype)
                self._ui_handler.set_err_signal(e)
                return
            
            except BaseException as e:
                self._conn.request_status[uid_integer] = e
                self._conn.request_state[uid_integer] = 0
                self._ui_handler.set_err_signal(e)
                return

        elif restype == "MESSAGE":
            parsed_response = payload.split(";", 3)
            data = []

            for field in parsed_response:
                _, value = field.split("=", 1)
                data.append(value)

            del data[1]
            self.buffer.put(data)

        elif restype == "FETCH":
            parsed_response = payload.split(";", 2)
            data = []
            
            for field in parsed_response:
                _, value = field.split("=", 1)
                data.append(value)
                
            try:
                # uid, length, raw_messages = data
                uid, raw_messages = data

                uid_integer = int(uid)
                
                
                if raw_messages == "":
                    self._conn.request_status[uid_integer] = "OK"
                    self._conn.request_state[uid_integer] = 0
                    return
                
                # each message is seperated by a newline                
                messages = raw_messages.split('|')
                
                # print(len(messages), length)
                # if len(messages) != length:
                #     self.conn.request_status[uid_integer] = "Invalid length"
                #     self.conn.request_state[uid_integer] = 0
                #     return

                for message in messages:
                    message_data = HistoryHandler.parse_history(message)
                    self.buffer.put(message_data)

                self._conn.request_status[uid_integer] = "OK"
                self._conn.request_state[uid_integer] = 0
            
            except ValueError as e:
                self._conn.request_status[-1] = "Cannot Get UID in {}".format(restype)
                self._ui_handler.set_err_signal(e)
                return
                
            except BaseException as e:
                self._conn.request_status[uid_integer] = e
                self._conn.request_state[uid_integer] = 0
                self._ui_handler.set_err_signal(e)
                return
        
        # elif restype == "AUDIO":
        #     parsedResponse = payload.split(";", 3)
        #     data = []

        #     for field in parsedResponse:
        #         _, value = field.split("=", 1)
        #         data.append(value)

        #     _, __, audio = data
            
        #     audio = audio.encode()
            
        #     if len(audio) != 1024:
        #         self.audio_temp = audio
        #         return
                
        #     self.call_handler.mixer.append_audio(audio)
        elif restype == "CHANNEL_ALLOCATION" or restype == "INCOMING_CALL":
            parsed_response = payload.split(";", 6)
            data = []

            for field in parsed_response:
                _, value = field.split("=", 1)
                data.append(value)
                
            uid, sender, recipient, token, salt, network_address = data
            
            if restype == "CHANNEL_ALLOCATION":
                try:
                    uid_integer = int(uid)
                    self._conn.request_status[uid_integer] = "OK"
                    self._conn.request_state[uid_integer] = 0 
                except ValueError:
                    pass
            
            self._call_handler.spawn_process(restype, sender, recipient, token, salt, network_address)
            
        # elif restype == "INCOMING_CALL":
        #     parsedResponse = payload.split(";", 4)
        #     data = []

        #     for field in parsedResponse:
        #         _, value = field.split("=", 1)
        #         data.append(value)

        #     caller, _ = data
        #     # TODO: Spawn call process with the correct option

        elif restype == "CALL_TIMEOUT":
            parsed_response = payload.split(";", 2)
            data = []

            for field in parsed_response:
                _, value = field.split("=", 1)
                data.append(value)

            # self._call_handler.set_state(1)
        
        elif restype == "CALL_DECLINED":
            parsed_response = payload.split(";", 2)
            data = []

            for field in parsed_response:
                _, value = field.split("=", 1)
                data.append(value)
                
            # self._call_handler.set_state(2)
        
        elif restype == "CALL_ACCEPTED":
            parsed_response = payload.split(";", 2)
            data = []

            for field in parsed_response:
                _, value = field.split("=", 1)
                data.append(value)

            # self._call_handler.set_state(3)
            
        elif restype == "CALL_TERMINATE":
            parsed_response = payload.split(";", 2)
            data = []

            for field in parsed_response:
                _, value = field.split("=", 1)
                data.append(value)

            # self._call_handler.set_state(4)
            
        elif restype == "CALL_ABORT":
            # parsedResponse = payload.split(";", 2)
            # data = []

            # for field in parsedResponse:
            #     _, value = field.split("=", 1)
            #     data.append(value)
            # _, value = payload.split("=", 1)
            # try:
            #     uid = int(value)
            #     self._conn.request_status[uid] = "ERROR"
            #     self._conn.request_state[uid] = 0 
            # except ValueError:
            #     pass
            
            if self._call_handler.check_process_status():
              self._call_handler.force_stop()
            
            
            # self._call_handler.set_state(5)
                       
        elif restype == "ERROR":
            error_field, uid_field = payload.split(";", 2)

            _, uid = uid_field.split("=", 2)
            _, err = error_field.split("=", 2)

            try:
                uid_integer = int(uid)
                self._conn.request_status[uid_integer] = err
                self._conn.request_state[uid_integer] = 0
                self._ui_handler.set_err_signal(err)
                
            except ValueError as e:
                self._conn.request_status[-1] = "Cannot Get UID in {}".format(restype)
                self._ui_handler.set_err_signal(e)
