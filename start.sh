#!/bin/bash
docker-compose -f sawtooth-default.yaml down
docker build -t sawtooth-auction-tp-go .
docker-compose -f sawtooth-default.yaml up
