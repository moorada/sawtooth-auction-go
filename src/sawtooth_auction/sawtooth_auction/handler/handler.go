/**
 * Copyright 2017-2018 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ------------------------------------------------------------------------------
 */

package handler

import (
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/processor_pb2"
	"github.com/moorada/sawtooth-sdk-go/examples/auction_go/src/sawtooth_auction/auction_payload"
	"github.com/moorada/sawtooth-sdk-go/examples/auction_go/src/sawtooth_auction/auction_state"
	"strings"
	"time"
)

var logger *logging.Logger = logging.Get()

type AuctionHandler struct {
}

const layoutDate = "2006-01-02T15:04:05.000Z"

func (self *AuctionHandler) FamilyName() string {
	return "auction"
}

func (self *AuctionHandler) FamilyVersions() []string {
	return []string{"1.0"}
}

func (self *AuctionHandler) Namespaces() []string {
	return []string{auction_state.Namespace}
}

func (self *AuctionHandler) Apply(request *processor_pb2.TpProcessRequest, context *processor.Context) error {
	// The auction player is defined as the signer of the transaction, so we unpack
	// the transaction header to obtain the signer's public key, which will be
	// used as the player's identity.
	header := request.GetHeader()
	user := header.GetSignerPublicKey()

	// The payload is sent to the transaction processor as bytes (just as it
	// appears in the transaction constructed by the transactor).  We unpack
	// the payload into an auctionPayload struct so we can access its fields.
	payload, err := auction_payload.FromBytes(request.GetPayload())
	if err != nil {
		return err
	}

	auctionState := auction_state.NewAuctionState(context)

	switch payload.Action {
	case "addItem":
		ok, err := validateAddBid(auctionState, payload.IdItem)
		if err != nil {
			return err
		}
		if ok {
			auction := &auction_state.AuctionTable{
				Item: auction_state.Item{
					ID:          payload.IdItem,
					Name:        payload.NameItem,
					Description: payload.Description,
					PostTime:    payload.PostTime,
					ExpiryTime:  payload.ExpiryTime,
				},
			}
			displayCreate(payload, user)
			return auctionState.AddItem(payload.IdItem, auction)
		} else {
			return &processor.InvalidTransactionError{Msg: "Auction not valid"}
		}
	case "placeBid":
		ok, err := validateBid(auctionState, payload, user)
		if err != nil {
			return err
		}
		auction, err := auctionState.GetAuction(payload.IdItem)
		if err != nil {
			return err
		}
		if ok {
			bid := auction_state.Bid{
				BidderName: payload.BidderName,
				ID:         payload.BidId,
				Amount:     payload.Amount,
				Timestamp:  payload.Timestamp,
			}
			auction.Bids = append(auction.Bids, bid)
			displayBid(payload, user)
			return auctionState.PlaceBid(payload.IdItem, auction)
		} else {
			return &processor.InvalidTransactionError{Msg: "Bid not valid"}
		}
	default:
		return &processor.InvalidTransactionError{
			Msg: fmt.Sprintf("Invalid Action : '%v'", payload.Action)}
	}
}

func validateAddBid(auctionState *auction_state.AuctionState, name string) (bool, error) {
	auction, err := auctionState.GetAuction(name)
	if err != nil {
		return false, err
	}
	if auction != nil {
		return false, &processor.InvalidTransactionError{Msg: "Auction already exists"}
	}
	return true, nil
}

func validateBid(auctionState *auction_state.AuctionState, payload *auction_payload.AuctionPayload, signer string) (bool, error) {
	auction, err := auctionState.GetAuction(payload.IdItem)
	if err != nil {
		return false, err
	}
	if auction == nil {
		return false, &processor.InvalidTransactionError{Msg: "PlaceBid requires an existing auction"}
	}
	tNow := time.Now()

	tMin, err := time.Parse(layoutDate, auction.Item.PostTime)
	if err != nil {
		return false, &processor.InvalidTransactionError{Msg: "Error parse date PostTime :"+auction.Item.PostTime}
	}

	tMaxBidder, err := time.Parse(layoutDate, payload.Timestamp)
	if err != nil {
		return false, &processor.InvalidTransactionError{Msg: "Error parse date Timestamp :"+payload.Timestamp}
	}
	tMax, err := time.Parse(layoutDate, auction.Item.ExpiryTime)
	if err != nil {
		return false, &processor.InvalidTransactionError{Msg: "Error parse date ExpiryTime :"+auction.Item.ExpiryTime}
	}
	if tNow.After(tMaxBidder) {
		return false, &processor.InvalidTransactionError{Msg: "After Max time bidder"}
	}
	if tNow.After(tMax) {
		return false, &processor.InvalidTransactionError{Msg: "Auction expired"}
	}
	if tNow.Before(tMin) {
		return false, &processor.InvalidTransactionError{Msg: "The auction has not yet started, tnow:"+tNow.String()+", tMin:"+tMin.String()}
	}
	return true, nil
}


func displayCreate(payload *auction_payload.AuctionPayload, signer string) {
	s := fmt.Sprintf("+ User %s created auction %s +", signer[:6], payload.IdItem)
	sLength := len(s)
	border := "+" + strings.Repeat("-", sLength-2) + "+"
	fmt.Println(border)
	fmt.Println(s)
	fmt.Println(border)
}

func displayBid(payload *auction_payload.AuctionPayload, signer string) {
	s := fmt.Sprintf("+ BidPlaced by %s, amount: %s +", payload.BidderName, payload.Amount)
	sLength := len(s)
	border := "+" + strings.Repeat("-", sLength-2) + "+"
	fmt.Println(border)
	fmt.Println(s)
	fmt.Println(border)
}

