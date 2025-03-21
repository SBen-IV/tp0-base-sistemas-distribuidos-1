from socket import socket, SHUT_RDWR


class ClientSocket():
    def __init__(self, socket: socket):
        self._socket = socket

    def recv(self, size: int):
        return self._socket.recv(size)

    def send(self, buffer):
        return self._socket.send(buffer)

    def getpeername(self):
        return self._socket.getpeername()

    def close(self):
        self._socket.shutdown(SHUT_RDWR)