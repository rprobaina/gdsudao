#!/bin/bash
echo "Inserindo dados de Previsão do Tempo"
./previsoes
echo "Dados de Previsão do Tempo inseridos em: $(date)" >> /home/ric/gdsudao.log


echo "Inserindo dados de Medicoes Diárias"
./diarios
echo "Dados de Dados Diários inseridos em: $(date)" >> /home/ric/gdsudao.log
