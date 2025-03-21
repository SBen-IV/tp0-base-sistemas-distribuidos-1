import logging

from common.client_socket import ClientSocket
from common.protocol_translator import ProtocolTranslator
from common.bet_info import BetInfo
from common.utils import Bet, store_bets

class ClientHandler():
    def __init__(self, client_socket: ClientSocket):
        self._client_socket = client_socket
        self._translator = ProtocolTranslator()

    def run(self):
        try:
            # First receive client id
            self._manage_client_id()
            # Then wait for client to send the Bet amount and bytes mount
            bet_info = self._manage_bet_info()
            # Wait for the client to send the amount of bets and bytes
            bet = self._manage_bet(bet_info)
            # Send OK and close the client
            # Store bet
            bets = [bet]
            store_bets(bets=bets)
            logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")

        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            self._client_socket.close()

    def _manage_client_id(self):
        # TODO: Modify the receive to avoid short-reads
        msg = self._client_socket.recv(2)
        logging.debug(f"Received message {msg} with length {len(msg)}")
        
        client_id = self._translator.translate_cli_id(msg)
        addr = self._client_socket.getpeername()
        logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {client_id}')
        
        # TODO: Modify the send to avoid short-writes
        # self._client_socket.send("{}\n".format("OK").encode('utf-8'))
        self._client_socket.send("{}".format("OK").encode('utf-8'))

        self._client_id = client_id


    def _manage_bet_info(self) -> BetInfo:
        msg = self._client_socket.recv(8)
        logging.debug(f"Received message {msg} with length {len(msg)}")

        bet_info = self._translator.translate_bet_info(msg)

        self._client_socket.send("{}".format("OK").encode('utf-8'))

        return bet_info
    
    def _manage_bet(self, bet_info: BetInfo) -> Bet:
        msg = self._client_socket.recv(bet_info._bytes_amount)
        logging.debug(f"Received message {msg} with length {len(msg)}")

        bet = self._translator.translate_bet(msg, self._client_id)

        self._client_socket.send("{}".format("OK").encode('utf-8'))

        return bet

    def stop(self):
        self._client_socket.close()