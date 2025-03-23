import logging

from common.client_socket import ClientSocket
from common.bet_info import BetInfo
from common.utils import Bet, store_bets

from protocol.agency_id import AGENCY_ID_MESSAGE_LEN, AgencyID
from protocol.bet_info import BET_INFO_MESSAGE_LEN, BetInfoProtocol
from protocol.bet import BetProtocol
from protocol.ok_message import OkMessage

from enum import Enum


class ProtocolState(Enum):
    """
    State of the protocol used to communicate with the client
    """
    AgencyIdentification = 1
    RecvBets = 2
    Fin = 3


class ClientHandler():
    """
    Handles communication with the Agency connected to client_socket
    """
    def __init__(self, client_socket: ClientSocket):
        self._client_socket = client_socket
        self._state = ProtocolState.AgencyIdentification

    def run(self):
        """
        Manage communication with the Agency
        """
        try:
            while self._state != ProtocolState.Fin:
                if self._state == ProtocolState.AgencyIdentification:
                    # First receive client id
                    self._manage_client_id()

                    self._state = ProtocolState.RecvBets

                elif self._state == ProtocolState.RecvBets:
                    # Then wait for client to send the Bet amount and bytes mount
                    bet_info = self._manage_bet_info()

                    # Wait for the client to send the amount of bets and bytes
                    # Send OK and close the client
                    bet = self._manage_bet(bet_info)
                    
                    # Store bet
                    store_bets(bets=[bet])
                    
                    logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")
                    
                    self._state = ProtocolState.Fin

        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | exception: {e}")
        finally:
            self._client_socket.close()

    def _manage_client_id(self):
        msg = self._recv_msg(AGENCY_ID_MESSAGE_LEN)

        agency_id = AgencyID.from_bytes(msg)
        
        addr = self._client_socket.getpeername()
        logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {agency_id}')
        
        self._send_ok()

        self._agency_id = agency_id


    def _manage_bet_info(self) -> BetInfo:
        msg = self._recv_msg(BET_INFO_MESSAGE_LEN)

        bet_info = BetInfoProtocol.from_bytes(msg)

        self._send_ok()

        return bet_info
    
    def _manage_bet(self, bet_info: BetInfo) -> Bet:
        msg = self._recv_msg(bet_info._bytes_amount)

        bet = BetProtocol.from_bytes(msg, self._agency_id)

        self._send_ok()

        return bet
    
    def _recv_msg(self, bytes_amount):
        msg = self._client_socket.recv(bytes_amount)
        logging.debug(f"Received message {msg} with length {len(msg)}")

        return msg

    def _send_ok(self):
        ok = OkMessage()
        buf, size = ok.encode()

        self._client_socket.send(buf, size)

    def stop(self):
        self._client_socket.close()