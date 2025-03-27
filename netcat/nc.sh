#!/bin/bash

echo ${MSG_SENT} | nc server 12345 -q 5
