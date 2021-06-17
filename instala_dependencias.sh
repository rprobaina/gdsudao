#!/bin/bash

# XLS decode
go get github.com/extrame/xls

# MongoDB drivers
go get go.mongodb.org/mongo-driver/mongo

# Decode XML for encoding='ISO-8859-1'
go get golang.org/x/net/html/charset

# Gorilla Mux (API)
go get -u github.com/gorilla/mux

# Mongo API
cp -r mongoapi/ /usr/lib/golang/src/

# Utils
cp -r utils/ /usr/lib/golang/src/
