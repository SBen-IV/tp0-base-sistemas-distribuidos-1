#!/bin/bash

MSG_SENT="Hola"
MSG_RESPONSE=$(echo ${MSG_SENT} | nc server 12345)

if [[ "${MSG_RESPONSE}" == "${MSG_SENT}" ]]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
