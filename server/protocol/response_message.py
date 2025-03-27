class ResponseMessage():
    def __init__(self, data):
        self._data = data

    def encode(self) -> tuple[bytes, int]:
        data_as_bytes = self._data.encode("utf-8")
        data_len = len(data_as_bytes)

        return data_as_bytes, data_len