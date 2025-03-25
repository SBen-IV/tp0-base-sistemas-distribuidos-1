import socket
import logging
import signal

from common.client_handler import ClientHandler
from common.client_socket import ClientSocket
from common.national_lottery import NationalLottery
from common.client_manager import ClientManager


class Server:
    def __init__(self, port, listen_backlog, clients_amount: int):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_is_running = True
        self._client_manager = ClientManager(clients_amount)

        signal.signal(signal.SIGTERM, self.__stop)


    def __stop(self, sig, frame):
        logging.info("SIGTERM received")
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        
        self._client_manager.stop()

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        while self._server_is_running:
            try:
                # Create a new ClientSocket()
                # Pass it to a ClientManager instance
                client_socket = self.__accept_new_connection()
                self._client_manager.add_client(ClientSocket(client_socket))
            except OSError:
                logging.info("Server socket closed")
                self._server_is_running = False

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
