#!/bin/bash

cd $FABRIC_CFG_PATH

## remove docker
docker ps -a | awk '{print $1}'| xargs docker rm 
docker images | grep tee | awk '{print $3}'| xargs docker rmi

## create and join channel
echo "Create syschannel"
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer channel create -c syschannel -f syschannel.tx -o orderer.rabbit.com:7050
echo "Join syschannel"
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer channel join -b ./syschannel.block -o orderer.rabbit.com:7050


## rabbit_02
echo "Install rabbit_02 chaincode"
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode install ~/chaincodes/rabbit_02pack.out

echo "Instantiate rabbit_02 chaincode"
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode instantiate -C syschannel -n rabbit_02 -v 1.0 -c '{"Args":[]}' -o orderer.rabbit.com:7050

## tee_data
echo "start deploy tee_data..."
echo "packet..."
./peer chaincode package teepack.out -n tee_data -v 1.0 -s -S -p github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/tee/core
mkdir $HOME/chaincodes/tee

echo "move..."
mv -fv teepack.out $HOME/chaincodes/tee/teepack.out

echo "install..."
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode install $HOME/chaincodes/tee/teepack.out

echo "instantiate..."
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode instantiate -C syschannel -n tee_data -v 1.0 -c '{"Args":[]}' -o orderer.rabbit.com:7050

## tee_exec
echo "start deploy tee_exec..."
echo "packet..."
./peer chaincode package teetaskpack.out -n tee_exec -v 1.0 -s -S -p github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/tee/task
mkdir $HOME/chaincodes/tee

echo "move..."
mv -fv teetaskpack.out $HOME/chaincodes/tee/teetaskpack.out

echo "install..."
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode install $HOME/chaincodes/tee/teetaskpack.out

echo "instantiate..."
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode instantiate -C syschannel -n tee_exec -v 1.0 -c '{"Args":["308187020100301306072a8648ce3d020106082a8648ce3d030107046d306b02010104202d130ea6dac76fcae718fbd20bf146643aa66fe6e5902975d2c5ed6ab3bcb5e2a144034200048f03f8321b00a4466f4bf4be51c91898cd50d8cc64c6ecf53e73443e348d5925a16f88c8952b78ebac2dc277a2cc54c77b4c3c07830f49629b689edf63086293", "048f03f8321b00a4466f4bf4be51c91898cd50d8cc64c6ecf53e73443e348d5925a16f88c8952b78ebac2dc277a2cc54c77b4c3c07830f49629b689edf63086293"]}' -o orderer.rabbit.com:7050
