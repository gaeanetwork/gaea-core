#!/bin/bash

version=1.2

cd $FABRIC_CFG_PATH

echo "start deploy tee_data..."
echo "packet..."
./peer chaincode package teepack.out -n tee_data -v $version -s -S -p github.com/gaeanetwork/gaea-core/smartcontract/fabric/chaincode/tee/core
mkdir $HOME/chaincodes/tee

echo "move..."
mv -fv teepack.out $HOME/chaincodes/tee/teepack.out

echo "install..."
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode install $HOME/chaincodes/tee/teepack.out

echo "upgrade..."
CORE_PEER_MSPCONFIGPATH=crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp ./peer chaincode upgrade -C syschannel -n tee_data -v $version -c '{"Args":[]}' -o orderer.rabbit.com:7050
