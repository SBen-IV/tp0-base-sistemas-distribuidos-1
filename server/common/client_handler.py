import logging

from common.client_socket import ClientSocket
from common.utils import Bet, store_bets

from model.bet_info import BetInfo
from model.operation import Operation

from protocol.agency_id import AGENCY_ID_MESSAGE_LEN, AgencyID
from protocol.bet_info import BET_INFO_MESSAGE_LEN, BetInfoProtocol
from protocol.operation import OPERATION_MESSAGE_LEN, OperationProtocol
from protocol.bets import BetsProtocol
from protocol.ok_message import OkMessage
from protocol.err_message import ErrMessage
from protocol.betfrombytes_error import BetFromBytesError

from enum import Enum


class ProtocolState(Enum):
    """
    State of the protocol used to communicate with the client
    """
    AgencyIdentification = 1
    WaitingOperation = 2
    RecvBets = 3
    Fin = 4


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
                    self._manage_client_id()

                    self._state = ProtocolState.WaitingOperation

                elif self._state == ProtocolState.WaitingOperation:
                    self._manage_operation()

                elif self._state == ProtocolState.RecvBets:
                    self._manage_bet_batch()

                    self._state = ProtocolState.WaitingOperation

        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | exception: {e}")
        finally:
            self._client_socket.close()

    def _manage_client_id(self):
        # First receive client id
        msg = self._recv_msg(AGENCY_ID_MESSAGE_LEN)

        agency_id = AgencyID.from_bytes(msg)
        
        addr = self._client_socket.getpeername()
        logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {agency_id}')
        
        self._send_ok()

        self._agency_id = agency_id

    def _manage_operation(self):
        msg = self._recv_msg(OPERATION_MESSAGE_LEN)

        operation = OperationProtocol.from_bytes(msg)

        if operation == Operation.Bet:
            self._state = ProtocolState.RecvBets
        else:
            self._state = ProtocolState.Fin

        self._send_ok()

    def _manage_bet_batch(self):
        # Then wait for client to send the Bet amount and bytes mount
        bet_info = self._manage_bet_info()

        # Wait for the client to send the amount of bets and bytes
        # Send the appropriate message
        bets = self._manage_bets(bet_info)
        

    def _manage_bet_info(self) -> BetInfo:
        msg = self._recv_msg(BET_INFO_MESSAGE_LEN)

        bet_info = BetInfoProtocol.from_bytes(msg)

        self._send_ok()

        return bet_info
    
    def _manage_bets(self, bet_info: BetInfo) -> list[Bet]:
        msg = self._recv_msg(bet_info._bytes_amount)

        try:
            bets = BetsProtocol.from_bytes(msg, self._agency_id, bet_info._bets_amount)
        
            # Store bets
            store_bets(bets)
            
            logging.info(f"action: apuesta_recibida | result: success | cantidad: {bet_info._bets_amount}")

            self._send_ok()
        except (BetFromBytesError, ValueError) as e:
            self._send_error()
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {bet_info._bets_amount}")
            raise e # Propagate the error to close the connection
            
    
    def _recv_msg(self, bytes_amount):
        msg = self._client_socket.recv(bytes_amount)
        logging.debug(f"Received message {msg} with length {len(msg)}")

        return msg

    def _send_ok(self):
        ok = OkMessage()
        buf, size = ok.encode()

        self._client_socket.send(buf, size)

    def _send_error(self):
        err = ErrMessage()
        buf, size = err.encode()

        self._client_socket.send(buf, size)

    def stop(self):
        self._client_socket.close()