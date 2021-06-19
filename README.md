# GD Sudao
### Sobre:

GD Sudão é um Sistema de Informação de Gestão Agropacuária (FMIS) voltado ao manejo de capim-sudão, desenvolvido como projeto de mestrado de Ricardo Robaina, em uma parceria entre a Universidade Federal do Pampa e a Embrapa Pecuária Sul.



### Arquitetura

![](gdsudao-arquitetura.png)



### Compilação

1. Clone o repositoório

   ```bash
   $ git clone https://github.com/rprobaina/gdsudao
   $ cd gdsudao
   ```

2. Instale as dependências

   ```bash
   $ ./dependencias/instala_dependencias.sh
   ```

3. Instalar a linguagem de programação Go

   ```bash
   $ sudo dnf install go
   ```

4. Compilar e executar (Desenvolvimento e Testes)

   ```bash
   $ go run main.go
   ```

5. Compilar para um arquivo executável (Implantação)

   ```bash
   $ go build main.go <nome_do_programa>
   ```

   

## Implantação

Os passos a seguir devem ser executados em um sistema operacional linux, preferencnialmente na distribuição Fedora 32+.

#### Instalação do banco de dados

1. Habilitar o repositório

   ```bash
   $ sudo vi /etc/yum.repos.d/mongodb-org-4.4.repo
   
   ## Adicionar o seguinte contepudo nesse arquivo: ##
   
   [mongodb-org-4.4]
   name=MongoDB Repository
   baseurl=https://repo.mongodb.org/yum/redhat/8Server/mongodb-org/4.4/x86_64/
   gpgcheck=1
   enabled=1
   gpgkey=https://www.mongodb.org/static/pgp/server-4.4.asc
   ```

2. Instalar os pacotes

   ```bash
   $sudo dnf install -y mongodb-org
   ```

3. Instalar o SGBD

   https://downloads.mongodb.com/compass/mongodb-compass-1.26.1.x86_64.rpm



#### Instalar o Módulo de coleta de dados

1. Copiar os arquivos executáveis para ``/home/bin``

   ```bash
   $ cd ppgcap-recursos
   $ mkdir -p $HOME/bin/gdsudao
   $ cp api diarios estacoes previsoes normais $HOME/bin/gdsudao
   ```

2. Configurar a crontab

   ```bash
   $ cd ppgcap-recursos
   $ sudo cp crontab.bk /etc/crontab
   ```

   

#### Restaurar o banco de dados

1. Crie o banco de dados

   ```bash
   $ cd ppgcap-recursos
   $ ./criar-db.sh
   Inserindo dados das Estações Meteorológicas
   Database connection successfully started!
   Document successfully inserted!  ObjectID("60ce59d801ed03754bf0da49")
   Document successfully inserted!  ObjectID("60ce59da01ed03754bf0da4a")
   ....
   (Espere a inserção terminar)
   ```

2. Use o Mongodb Compass

   1. Conecte localmente
   2. Acesse o banco de dados **gdsudao**
   3. Clique no botão verde **CREATE COLLECTION**
   4. Crie a coleção (por exemplo: **normais**)
   5. Clique no botão verde **ADD DATA**, e então em **Inport File**. Em seguiga selecione o arquivo ``ppgcap-recursos/normais.json``
   6. Repita esse processo com as coleções **diarios** e **previsões**



#### Executar

1. Iniciar a API

   ```bash
   $ $HOME/bin/gdsuduao/api
   ```

2. Acessar a API

   Acesse http://localhost:8082/ em um navegador

   ```
    GD Sudão Application Programming Interface
   
   Essa API compõe o FMIS GD Sudão, densenvolvido em parceiria pela Universidade Federal do Pampa e Embrapa Pecuária Sul.
   Como usar:
   
       Retorna dados da estação meteorológica mais próxima a um ponto:
           https://localhost:8080/estacao/maisproxima/{latitude}/{longitude} 
   
       Retorna dados de normais climatológicas de uma estação:
           https://localhost:8080/normais/{nomeEstacao} 
   
       Retorna dados diários de uma estação:
           https://localhost:8080/diarios/{codigoINMET}/{dataInicial}/{dataFinal} 
   
       Retorna dados de previsão do tempo de uma estação:
           https://localhost:8080/previsoes/{codigoINMET}/{dataAtual} 
   
       Retorna a soma térmica para o capim sudão BRS-Stribo e a proproção de dados utililizados:
           https://localhost:8080/gdsudao/{codigoINMET}/{dataInicial}/{dataFinal} 
   
       Retorna a soma térmica para determinada temperatura basal e a proproção de dados utililizados:
           https://localhost:8080/somatermica/{codigoINMET}/{temperaturaBasal}/{dataInicial}/{dataFinal} 
   
       Retorna a data do próximo pastejo do Capim-Sudão BRS-Stribo:
           https://localhost:8080/proximoPastejo/{codigoINMET}/{dataInicial}/{numeroPastejos} 
   
       Retorna o número estimados de pastejos de uma região:
           https://localhost:8080/pastejos/{codigoINMET}/{dataInicial}/{dataFinal} 
   ```

   





