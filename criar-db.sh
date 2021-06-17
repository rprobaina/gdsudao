#!/bin/bash
echo "Inserindo dados das Estações Meteorológicas"
./estacoes
echo "Dados de Estações Meteorológicas inseridos em: $(date)" >> $HOME/gdsudao.log

echo "Inserindo dados das Normais Climatológicas"
./normais
echo "Dados de Normais Climatológicas inseridos em: $(date)" >> $HOME/gdsudao.log
