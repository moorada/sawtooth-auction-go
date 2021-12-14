/**
 * Copyright 2018 Intel Corporation
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

package auction_state

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"sort"
	"strings"
)

var Namespace = hexdigest("auction")[:6]

type AuctionTable struct {
	Item Item
	Bids []Bid
}

type Bid struct {
	ID         string
	BidderName string
	Amount     string
	//Timestamp  string
}

type Item struct {
	ID          string
	Name        string
	Description string
	PostTime    string
	ExpiryTime  string
}

// AuctionState handles addressing, serialization, deserialization,
// and holding an addressCache of data at the address.
type AuctionState struct {
	context      *processor.Context
	addressCache map[string][]byte
}

// NewAuctionState constructs a new auctionState struct.
func NewAuctionState(context *processor.Context) *AuctionState {
	return &AuctionState{
		context:      context,
		addressCache: make(map[string][]byte),
	}
}

// GetAuction returns a game by it's name.
func (self *AuctionState) GetAuction(name string) (*AuctionTable, error) {
	auctions, err := self.loadAuctions(name)
	//	fmt.Println("Auctions", auctions)
	if err != nil {
		return nil, err
	}
	game, ok := auctions[name]
	if ok {
		return game, nil
	}
	return nil, nil
}

//// AddItem sets a game to it's name
func (self *AuctionState) AddItem(id string, auction *AuctionTable) error {
	auctions, err := self.loadAuctions(id)
	if err != nil {
		return err
	}
	auctions[id] = auction
	return self.storeAuctions(id, auctions)
}

//// AddItem sets a game to it's name
func (self *AuctionState) PlaceBid(id string, auction *AuctionTable) error {
	//fmt.Println("Inside PlaceBid", auction)
	auctions, err := self.loadAuctions(id)
	if err != nil {
		return err
	}
	auctions[id] = auction

	return self.storeAuctions(id, auctions)
}

func (self *AuctionState) loadAuctions(name string) (map[string]*AuctionTable, error) {
	address := makeAddress(name)

	data, ok := self.addressCache[address]
	if ok {
		if self.addressCache[address] != nil {
			return Deserialize(data)
		}
		return make(map[string]*AuctionTable), nil
	}

	results, err := self.context.GetState([]string{address})
	if err != nil {
		return nil, err
	}
	if len(string(results[address])) > 0 {
		self.addressCache[address] = results[address]
		return Deserialize(results[address])
	}
	self.addressCache[address] = nil
	auctions := make(map[string]*AuctionTable)
	return auctions, nil
}

func (self *AuctionState) storeAuctions(name string, auctions map[string]*AuctionTable) error {
	address := makeAddress(name)

	var names []string
	for id := range auctions {
		names = append(names, id)
	}

	sort.Strings(names)

	var g []*AuctionTable
	for _, id := range names {
		g = append(g, auctions[id])
	}

	data := Serialize(g)

	self.addressCache[address] = data

	_, err := self.context.SetState(map[string][]byte{
		address: data,
	})
	return err
}

func (self *AuctionState) deleteAuctions(name string) error {
	address := makeAddress(name)

	_, err := self.context.DeleteState([]string{address})
	return err
}

func Deserialize(data []byte) (map[string]*AuctionTable, error) {
	auctions := make(map[string]*AuctionTable)
	for _, str := range strings.Split(string(data), "|") {
		var auction = AuctionTable{}
		tables := strings.Split(str, "!")
		if len(tables) != 2 {
			return nil, &processor.InternalError{
				Msg: fmt.Sprintf("Malformed game data: '%v'", string(data))}
		}

		itemSlice := strings.Split(tables[0], ",")
		auction.Item = Item{ID: itemSlice[0], Name: itemSlice[1], Description: itemSlice[2], PostTime: itemSlice[3], ExpiryTime: itemSlice[4]}

		if tables[1] != "" {
			bids := strings.Split(tables[1], ";")
			for _, bid := range bids {
				bidSlice := strings.Split(bid, ",")
				auction.Bids = append(auction.Bids, Bid{ID: bidSlice[0], BidderName: bidSlice[1], Amount: bidSlice[2] /*Timestamp: bidSlice[3]*/})
			}
		}
		auctions[auction.Item.ID] = &auction

	}

	return auctions, nil
}

func Serialize(auctionTables []*AuctionTable) []byte {
	var buffer bytes.Buffer
	for i, at := range auctionTables {

		item := at.Item
		buffer.WriteString(item.ID)
		buffer.WriteString(",")
		buffer.WriteString(item.Name)
		buffer.WriteString(",")
		buffer.WriteString(item.Description)
		buffer.WriteString(",")
		buffer.WriteString(item.PostTime)
		buffer.WriteString(",")
		buffer.WriteString(item.ExpiryTime)
		buffer.WriteString("!")

		for j, bid := range at.Bids {
			buffer.WriteString(bid.ID)
			buffer.WriteString(",")
			buffer.WriteString(bid.BidderName)
			buffer.WriteString(",")
			buffer.WriteString(bid.Amount)
			buffer.WriteString(",")
			//buffer.WriteString(bid.Timestamp)
			if j+1 != len(at.Bids) {
				buffer.WriteString(";")
			}
		}

		if i+1 != len(auctionTables) {
			buffer.WriteString("|")
		}
	}
	return buffer.Bytes()
}

func makeAddress(name string) string {
	return Namespace + hexdigest(name)[:64]
}

func hexdigest(str string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	return strings.ToLower(hex.EncodeToString(hashBytes))
}
