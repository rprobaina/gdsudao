#!/bin/bash
echo "Inserindo dados das Estações Meteorológicas"
./estacoes
echo "Dados de Estações Meteorológicas inseridos em: $(date)" >> /home/ric/gdsudao.log

echo "Inserindo dados das Normais Climatológicas"
./normais
echo "Dados de Normais Climatológicas inseridos em: $(date)" >> /home/ric/gdsudao.log