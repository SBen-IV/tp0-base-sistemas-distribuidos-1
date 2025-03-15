FROM ubuntu:noble

RUN apt update && apt upgrade --yes && apt install --no-install-recommends --yes netcat-traditional

# COPY ./validar-echo-server.sh .

ENTRYPOINT ["/bin/bash"]