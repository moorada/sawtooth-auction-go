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

package auction_payload

import (
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"strings"
)

type AuctionPayload struct {
	Action      string
	IdItem      string
	NameItem    string
	Description string
	PostTime    string
	ExpiryTime  string
	BidderName  string
	Timestamp   string
	Amount      string
	BidId       string
}

func FromBytes(payloadData []byte) (*AuctionPayload, error) {
	if payloadData == nil {
		return nil, &processor.InvalidTransactionError{Msg: "Must contain payload"}
	}

	parts := strings.Split(string(payloadData), ",")
	if len(parts) != 10 {
		return nil, &processor.InvalidTransactionError{Msg: "Payload is malformed"}
	}

	payload := AuctionPayload{}
	payload.Action = parts[0]
	payload.IdItem = parts[1]
	payload.NameItem = parts[2]
	payload.Description = parts[3]
	payload.PostTime = parts[4]
	payload.ExpiryTime = parts[5]
	payload.BidId = parts[6]
	payload.BidderName = parts[7]
	payload.Amount = parts[8]
	payload.Timestamp = parts[9]

	//TODO check everithing

	//if len(payload.Name) < 1 {
	//	return nil, &processor.InvalidTransactionError{Msg: "Name is required"}
	//}
	//
	//if len(payload.Action) < 1 {
	//	return nil, &processor.InvalidTransactionError{Msg: "Action is required"}
	//}

	//if payload.Action == "take" {
	//	space, err := strconv.Atoi(parts[2])
	//	if err != nil {
	//		return nil, &processor.InvalidTransactionError{
	//			Msg: fmt.Sprintf("Invalid Space: '%v'", parts[2])}
	//	}
	//	payload.Space = space
	//}
	//
	//if strings.Contains(payload.Name, "|") {
	//	return nil, &processor.InvalidTransactionError{
	//		Msg: fmt.Sprintf("Invalid Name (char '|' not allowed): '%v'", parts[2])}
	//}

	return &payload, nil
}
