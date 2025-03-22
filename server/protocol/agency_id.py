AGENCY_ID_MESSAGE_LEN = 2

class AgencyID():
    @staticmethod
    def from_bytes(msg):
        return str(int.from_bytes(msg[0:AGENCY_ID_MESSAGE_LEN], 'big'))