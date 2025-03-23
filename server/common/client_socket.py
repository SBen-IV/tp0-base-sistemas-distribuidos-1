from socket import socket, SHUT_RDWR

from common.client_disconnected_exception import ClientDisconnectedException


class ClientSocket():
    def __init__(self, socket: socket):
        self._socket = socket

    def send(self, buffer, size: int):
        total_bytes_sent = 0

        while total_bytes_sent < size:
            bytes_sent = self._socket.send(buffer[total_bytes_sent:size])

            total_bytes_sent += bytes_sent

        return total_bytes_sent

    def recv(self, size: int):
        total_bytes_received = 0
        data = b''

        while total_bytes_received < size:
            new_data = self._socket.recv((size - total_bytes_received))

            if not new_data:
                raise ClientDisconnectedException()

            total_bytes_received += len(new_data)
            data += new_data


        return data    

    def getpeername(self):
        return self._socket.getpeername()

    def close(self):
        self._socket.shutdown(SHUT_RDWR)