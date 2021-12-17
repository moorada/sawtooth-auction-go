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

package main

import (
	"github.com/jessevdk/go-flags"
	"math/rand"
	"strconv"
	"time"
)

type WorkLoad struct {
	Args struct {
		WaitTime int `positional-arg-name:"waitTime" required:"true" description:"wait time"`
		Quantity int `positional-arg-name:"quantity" required:"true" description:"quantity"`
		//Value string `positional-arg-name:"value" required:"true" description:"Amount to set"`
		//IdItem     string `positional-arg-name:"IdItem" required:"true" description:"Name of key to set"`
		//BidId      string `positional-arg-name:"BidId" required:"true" description:"Name of key to set"`
		//BidderName string `positional-arg-name:"BidderName" required:"true" description:"Name of key to set"`
		//Amount     string `positional-arg-name:"Amount" required:"true" description:"Name of key to set"`
		//Timestamp  string `positional-arg-name:"Timestamp" required:"true" description:"Name of key to set"`
	} `positional-args:"true"`
	Url     string `long:"url" description:"Specify URL of REST API"`
	Keyfile string `long:"keyfile" description:"Identify file containing user's private key"`
	Wait    uint   `long:"wait" description:"Set time, in seconds, to wait for transaction to commit"`
}

func (args *WorkLoad) Name() string {
	return "workload"
}

func (args *WorkLoad) KeyfilePassed() string {
	return args.Keyfile
}

func (args *WorkLoad) UrlPassed() string {
	return args.Url
}

func (args *WorkLoad) Register(parent *flags.Command) error {
	_, err := parent.AddCommand(args.Name(), "Workload", "Workload", args)
	if err != nil {
		return err
	}
	return nil
}

func (args *WorkLoad) Run() error {
	// Construct client

	wait := args.Wait

	auctionClient, err := GetClient(args, true)
	if err != nil {
		return err
	}
	idItem := "idItem"
	nameItem := "nameItem"
	description := "description"

	bidId := "bidId"
	bidderName := "bidderName"

	for i := 0; i < args.Args.Quantity; i++ {

		time.Sleep(time.Millisecond * time.Duration(rand.Intn(args.Args.WaitTime)))
		expiryTime := time.Now().Add(time.Minute * 30)
		postTime := time.Now()
		number := strconv.Itoa(i)

		_, err = auctionClient.AddItem(idItem+number, nameItem+number, description+number, postTime.Format(layoutDate), expiryTime.Format(layoutDate), wait)
		if err != nil {
			return err
		}
	}

	for i := 0; i < args.Args.Quantity; i++ {
		numberItem := strconv.Itoa(i)
		for j := 0; j < args.Args.Quantity; j++ {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(args.Args.WaitTime)))
			numberBid := strconv.Itoa(j)
			amount := strconv.Itoa(rand.Intn(100))

			//timeStamp := time.Now().Add(time.Second * 3)
			_, err = auctionClient.PlaceBid(idItem+numberItem, bidId+numberBid, bidderName+numberBid, amount /*timeStamp.Format(layoutDate),*/, wait)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
