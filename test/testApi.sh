#!/bin/bash

### Info ### ### ###
#	
#	Script de teste e colehata de estatísticas da API GD Sudão
#
#	Autor: Ricardo Robaina
#
#	Data de atualização: 23 de fevereiro de 2021
#
### ### ### ### ####


# Estação mais próxima
echo "--------- Estação mais próxima ----------"

echo "Bagé"
time curl 'http://localhost:8082/estacao/maisproxima/-31.347801/-54.013292'
echo ""

echo "São Luiz Gonzaga"
time curl 'http://localhost:8082/estacao/maisproxima/-28.417113/-54.962403'
echo ""

echo "Cambará do Sul"
time curl 'http://localhost:8082/estacao/maisproxima/-29.049125/-50.149636'
echo ""


# Normais

echo "--------- Normais ----------"

echo "Bagé"
time curl 'http://localhost:8082/normais/BAGE'
echo ""

echo "São Luiz Gonzaga"
time curl 'http://localhost:8082/normais/SAO%20LUIZ%20GONZAGA'
echo ""

echo "Cambará do Sul"
time curl 'http://localhost:8082/normais/CAMBARA%20DO%20SUL'
echo ""



# Dados diários

echo "--------- Diários ----------"

echo "Bagé"
time curl 'http://localhost:8082/diarios/A827/2020-10-01/2020-10-31'
echo ""

echo "São Luiz Gonzaga"
time curl 'http://localhost:8082/diarios/A852/2020-10-01/2020-10-31'
echo ""

echo "Cambará do Sul"
time curl 'http://localhost:8082/diarios/A897/2020-10-01/2020-10-31'
echo ""


# Dados previsões

echo "--------- Previsões ----------"

echo "Bagé"
time curl 'http://localhost:8082/previsoes/A827/2020-10-02'
echo ""

echo "São Luiz Gonzaga"
time curl 'http://localhost:8082/previsoes/A852/2020-10-02'
echo ""

echo "Cambará do Sul"
time curl 'http://localhost:8082/previsoes/A897/2020-10-02'
echo ""


# Soma térmica capim-sudão BRS-Estribo

echo "--------- Soma térmica capim-sudão BRS-Estribo ----------"

echo "Bagé"
time curl 'http://localhost:8082/gdsudao/A827/2020-10-01/2020-10-31'
echo ""

echo "São Luiz Gonzaga"
time curl 'http://localhost:8082/gdsudao/A852/2020-10-01/2020-10-31'
echo ""

echo "Cambará do Sul"
time curl 'http://localhost:8082/gdsudao/A897/2020-10-01/2020-10-31'
echo ""


# Soma térmica capim-sudão BRS-Estribo

echo "--------- Soma térmica ----------"

echo "Bagé"
time curl 'http://localhost:8082/grausdia/A827/15.0/2020-10-01/2020-10-31'
echo ""

echo "São Luiz Gonzaga"
time curl 'http://localhost:8082/grausdia/A852/15.0/2020-10-01/2020-10-31'
echo ""

echo "Cambará do Sul"
time curl 'http://localhost:8082/grausdia/A897/15.0/2020-10-01/2020-10-31'
echo ""


# Próximo pastejo

echo "--------- Próximo pastejo ----------"

echo "Bagé"
time curl 'http://localhost:8082/proximoPastejo/A827/2020-10-01/1'
echo ""

echo "São Luiz Gonzaga"
time curl 'http://localhost:8082/proximoPastejo/852/2020-10-01/1'
echo ""

echo "Cambará do Sul"
time curl 'http://localhost:8082/proximoPastejo/A897/2020-10-01/1'
echo ""


# Número estimado de pastejos

echo "--------- Número estimado de pastejos ----------"

echo "Bagé"
time curl 'http://localhost:8082/pastejos/A827/2020-10-01/2020-10-31'
echo ""

echo "São Luiz Gonzaga"
time curl 'http://localhost:8082/pastejos/A852/2020-10-01/2020-10-31'
echo ""

echo "Cambará do Sul"
time curl 'http://localhost:8082/pastejos/A897/2020-10-01/2020-10-31'
echo ""
