#!/bin/bash
go build -o spital-server
docker build -t spital-server .
rm spital-server
