import threading
import logging

from common.national_lottery import NationalLottery
from common.client_socket import ClientSocket
from common.client_handler import ClientHandler

from multiprocessing import Process

class ClientManager():
    def __init__(self, clients_amount: int):
        self._national_lottery = NationalLottery(clients_amount)
        self._clients = []
        
    def add_client(self, client_socket: ClientSocket):
        client = ClientHandler(client_socket, self._national_lottery)

        # p = Process(target=client.run)

        logging.debug("Starting new client")
        # p.start()

        # self._clients.append((client, p))

        t = threading.Thread(target=client.run)

        t.start()

        self._clients.append((client, t))

    def stop(self):
        for (client, p) in self._clients:
            client.stop()
            p.join()