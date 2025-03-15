FROM ubuntu:noble

RUN apt update && apt upgrade --yes \
    && apt install --no-install-recommends --yes netcat-traditional

ENTRYPOINT ["/bin/bash"]