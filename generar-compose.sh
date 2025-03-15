#!/bin/bash

echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"

source .venv/bin/activate

python3 mi-generador.py $1 $2
