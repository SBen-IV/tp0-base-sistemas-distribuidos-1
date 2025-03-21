import logging

from common.client_socket import ClientSocket

class ClientHandler():
    def __init__(self, client_socket: ClientSocket):
        self._client_socket = client_socket

    def run(self):
        try:
            # TODO: Modify the receive to avoid short-reads
            # msg = client_sock.recv(1024).rstrip().decode('utf-8')
            msg = self._client_socket.recv(2)
            logging.debug(f"Received message {int.from_bytes(msg, 'big')} with length {len(msg)}")
            addr = self._client_socket.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            
            # TODO: Modify the send to avoid short-writes
            self._client_socket.send("{}\n".format("OK").encode('utf-8'))
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            self._client_socket.close()


    def stop(self):
        self._client_socket.close()