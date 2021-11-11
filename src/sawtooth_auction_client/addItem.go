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

type AddItem struct {
	Args struct {
		//Name  string `positional-arg-name:"name" required:"true" description:"Name of key to set"`
		//Value string `positional-arg-name:"value" required:"true" description:"Amount to set"`
		IdItem      string `positional-arg-name:"IdItem" required:"true" description:"Name of key to set"`
		NameItem    string `positional-arg-name:"NamItem" required:"true" description:"Name of key to set"`
		Description string `positional-arg-name:"Description" required:"true" description:"Name of key to set"`
		PostTime    string `positional-arg-name:"PostTime" required:"true" description:"Name of key to set"`
		ExpiryTime  string `positional-arg-name:"ExpiryTime" required:"true" description:"Name of key to set"`
		//BidderName  string
		//Timestamp   string
		//Amount      string
		//BidId       string
	} `positional-args:"true"`
	Url     string `long:"url" description:"Specify URL of REST API"`
	Keyfile string `long:"keyfile" description:"Identify file containing user's private key"`
	Wait    uint   `long:"wait" description:"Set time, in seconds, to wait for transaction to commit"`
}

func (args *AddItem) Name() string {
	return "addItem"
}

func (args *AddItem) KeyfilePassed() string {
	return args.Keyfile
}

func (args *AddItem) UrlPassed() string {
	return args.Url
}

func (args *AddItem) Register(parent *flags.Command) error {
	_, err := parent.AddCommand(args.Name(), "Sets an intkey value", "Sends an intkey transaction to set <name> to <value>.", args)
	if err != nil {
		return err
	}
	return nil
}

func (args *AddItem) Run() error {
	// Construct client

	wait := args.Wait

	auctionClient, err := GetClient(args, true)
	if err != nil {
		return err
	}
	_, err = auctionClient.AddItem(args.Args.IdItem, args.Args.NameItem, args.Args.Description, args.Args.PostTime, args.Args.ExpiryTime, wait)
	return err
}
