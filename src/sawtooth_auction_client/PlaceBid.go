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
)

type PlaceBid struct {
	Args struct {
		//Name  string `positional-arg-name:"name" required:"true" description:"Name of key to set"`
		//Value string `positional-arg-name:"value" required:"true" description:"Amount to set"`
		IdItem     string `positional-arg-name:"IdItem" required:"true" description:"Name of key to set"`
		BidId      string `positional-arg-name:"BidId" required:"true" description:"Name of key to set"`
		BidderName string `positional-arg-name:"BidderName" required:"true" description:"Name of key to set"`
		Amount     string `positional-arg-name:"Amount" required:"true" description:"Name of key to set"`
		//Timestamp  string `positional-arg-name:"Timestamp" required:"true" description:"Name of key to set"`
	} `positional-args:"true"`
	Url     string `long:"url" description:"Specify URL of REST API"`
	Keyfile string `long:"keyfile" description:"Identify file containing user's private key"`
	Wait    uint   `long:"wait" description:"Set time, in seconds, to wait for transaction to commit"`
}

func (args *PlaceBid) Name() string {
	return "placeBid"
}

func (args *PlaceBid) KeyfilePassed() string {
	return args.Keyfile
}

func (args *PlaceBid) UrlPassed() string {
	return args.Url
}

func (args *PlaceBid) Register(parent *flags.Command) error {
	_, err := parent.AddCommand(args.Name(), "Sets an intkey value", "Sends an intkey transaction to set <name> to <value>.", args)
	if err != nil {
		return err
	}
	return nil
}

func (args *PlaceBid) Run() error {
	// Construct client

	wait := args.Wait

	auctionClient, err := GetClient(args, true)
	if err != nil {
		return err
	}
	_, err = auctionClient.PlaceBid(args.Args.IdItem, args.Args.BidId, args.Args.BidderName, args.Args.Amount /*args.Args.Timestamp, */, wait)
	return err
}
