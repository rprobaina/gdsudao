#!/bin/bash

echo "Estação mais proxima"
curl localhost:8082/estacao/maisproxima/-54.013292/-31.347801

echo ""
echo "Normais"
curl localhost:8082/normais/BAGE

echo ""
echo "Diarios"
curl localhost:8082/diarios/A827/2020-09-10/2020-09-15

echo ""
echo "Previsões"
curl localhost:8082/previsoes/A827/2020-09-10