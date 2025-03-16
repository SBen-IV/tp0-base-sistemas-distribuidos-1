#!/bin/bash

MSG_SENT="Hola"
NETWORK=tp0_testing_net

docker build -f ./netcat/Dockerfile -t "nc-client:latest" .

MSG_RESPONSE=$(docker run --rm --network ${NETWORK} -e MSG_SENT=${MSG_SENT} -v ./netcat/nc.sh:/nc.sh nc-client:latest "/nc.sh")

if test "${MSG_RESPONSE}" = "${MSG_SENT}"; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi