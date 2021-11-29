#!/bin/bash
# sawtooth peer list --url http://rest-api-1:8008
#
#
#
#
sawtooth_auction_client addItem "11111" "T-shirt" "nice" "a" "v"
sleep 1
sawtooth_auction_client addItem "22222" "Jacket" "nice" "a" "v"
sleep 1
sawtooth_auction_client addItem "33333" "Bike" "nice" "a" "v"
sleep 1
sawtooth_auction_client addItem "44444" "Car" "nice" "a" "v"
sleep 4
sawtooth_auction_client placeBid "11111" "12" "Ted" "152" "timestamp"
sleep 1
sawtooth_auction_client placeBid "11111" "11"  "Joe" "150" "timestamp"
sleep 1
sawtooth_auction_client placeBid "22222" "21" "Oliver" "19" "timestamp"
sleep 1
sawtooth_auction_client placeBid "22222" "22" "Elon" "18" "timestamp"
sleep 1
sawtooth_auction_client placeBid "22222" "23" "Mary" "15" "timestamp"
sleep 1
sawtooth_auction_client placeBid "44444" "41" "Mary" "1232" "timestamp"
sleep 1
sawtooth_auction_client placeBid "33333" "31" "Ted" "550" "timestamp"
sleep 1

sawtooth_auction_client show "11111"
sawtooth_auction_client show "22222"
sawtooth_auction_client show "33333"
sawtooth_auction_client show "44444"
