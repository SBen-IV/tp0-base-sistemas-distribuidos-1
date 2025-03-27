import logging

from common.client_socket import ClientSocket
from common.national_lottery import NationalLottery

from model.bet_info import BetInfo
from model.operation import Operation

from protocol.agency_id import AGENCY_ID_MESSAGE_LEN, AgencyID
from protocol.bet_info import BET_INFO_MESSAGE_LEN, BetInfoProtocol
from protocol.operation import OPERATION_MESSAGE_LEN, OperationProtocol
from protocol.client_message import OK_CLIENT_MESSAGE_LEN, ClientMessage, ClientResponseMessage
from protocol.bets import BetsProtocol
from protocol.ok_message import OkMessage
from protocol.err_message import ErrMessage
from protocol.betfrombytes_error import BetFromBytesError
from protocol.winners_info_message import WinnersInfoMessageProtocol
from protocol.winners_message import WinnersMessageProtocol
from protocol.not_available_message import NotAvailableMessage
from protocol.response_message import ResponseMessage

from enum import Enum


class ProtocolState(Enum):
    """
    State of the protocol used to communicate with the client
    """
    AgencyIdentification = 1
    WaitingOperation = 2
    RecvBets = 3
    Fin = 4
    GetWinners = 5


class ClientHandler():
    """
    Handles communication with the Agency connected to client_socket
    """
    def __init__(self, client_socket: ClientSocket, national_lottery: NationalLottery):
        self._client_socket = client_socket
        self._state = ProtocolState.AgencyIdentification
        self._national_lottery = national_lottery

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
                elif self._state == ProtocolState.GetWinners:
                    logging.info("action: sorteo | result: success")
                    self._manage_get_winners()

                    self._state = ProtocolState.WaitingOperation

        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        except Exception as e:
            logging.error(f"action: receive_message | result: fail | exception: {e}")
        finally:
            self.stop()

    def _manage_client_id(self):
        # First receive client id
        msg = self._recv_msg(AGENCY_ID_MESSAGE_LEN)

        agency_id = AgencyID.from_bytes(msg)
        
        addr = self._client_socket.getpeername()
        logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {agency_id}')
        
        self._send_msg(OkMessage())

        self._agency_id = agency_id
        self._national_lottery.add_agency(agency_id)

    def _manage_operation(self):
        msg = self._recv_msg(OPERATION_MESSAGE_LEN)

        operation = OperationProtocol.from_bytes(msg)

        logging.debug(f"Received operation: {operation} from {self._agency_id}")

        if operation == Operation.Bet:
            self._state = ProtocolState.RecvBets
        elif operation == Operation.NoMoreBets:
            self._national_lottery.mark_no_more_bets(self._agency_id)
        elif operation == Operation.Draw:
            if self._national_lottery.can_draw():
              self._state = ProtocolState.GetWinners
            else:
              self._send_msg(NotAvailableMessage())
              return
        else:
            self._state = ProtocolState.Fin

        self._send_msg(OkMessage())

    def _manage_bet_batch(self):
        # Then wait for client to send the Bet amount and bytes mount
        bet_info = self._manage_bet_info()

        # Wait for the client to send the amount of bets and bytes
        # Send the appropriate message
        self._manage_bets(bet_info)
        

    def _manage_bet_info(self) -> BetInfo:
        msg = self._recv_msg(BET_INFO_MESSAGE_LEN)

        bet_info = BetInfoProtocol.from_bytes(msg)

        self._send_msg(OkMessage())

        return bet_info
    
    def _manage_bets(self, bet_info: BetInfo):
        msg = self._recv_msg(bet_info._bytes_amount)

        try:
            bets = BetsProtocol.from_bytes(msg, self._agency_id, bet_info._bets_amount)
        
            # Store bets
            self._national_lottery.store_bets(bets)
            
            logging.info(f"action: apuesta_recibida | result: success | cantidad: {bet_info._bets_amount}")

            self._send_msg(OkMessage())
        except (BetFromBytesError, ValueError) as e:
            self._send_msg(ErrMessage())
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {bet_info._bets_amount}")
            raise e # Propagate the error to close the connection
            
    def _manage_get_winners(self):
        winners = self._national_lottery.get_winners(self._agency_id)

        logging.debug(f"Winners {len(winners)} for {self._agency_id}")

        winners_buf, winners_buf_size = WinnersMessageProtocol.to_bytes(winners)

        winners_info_buf, winners_info_buf_size = WinnersInfoMessageProtocol.to_bytes(winners_buf_size)

        self._client_socket.send(winners_info_buf, winners_info_buf_size)

        msg = self._recv_msg(OK_CLIENT_MESSAGE_LEN)

        if ClientMessage.from_bytes(msg, OK_CLIENT_MESSAGE_LEN) != ClientResponseMessage.OkMessage:
            raise Exception("Client failed unexpectedly")

        self._client_socket.send(winners_buf, winners_buf_size)

        msg = self._recv_msg(OK_CLIENT_MESSAGE_LEN)
        
        if ClientMessage.from_bytes(msg, OK_CLIENT_MESSAGE_LEN) != ClientResponseMessage.OkMessage:
            raise Exception("Client failed unexpectedly")

    def _recv_msg(self, bytes_amount):
        msg = self._client_socket.recv(bytes_amount)
        logging.debug(f"Received message {msg} with length {len(msg)}")

        return msg

    def _send_msg(self, msg: ResponseMessage):
        buf, size = msg.encode()

        self._client_socket.send(buf, size)

    def stop(self):
        self._client_socket.close()