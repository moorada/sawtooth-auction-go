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
	"fmt"
	"github.com/evilsocket/islazy/tui"
	"github.com/jessevdk/go-flags"
	"os"
	"strconv"
	"strings"
	"time"
)

type Show struct {
	Args struct {
		Name string `positional-arg-name:"name" required:"true" description:"Name of key to show"`
	} `positional-args:"true"`
	Url string `long:"url" description:"Specify URL of REST API"`
}

func (args *Show) Name() string {
	return "show"
}

func (args *Show) KeyfilePassed() string {
	return ""
}

func (args *Show) UrlPassed() string {
	return args.Url
}

func (args *Show) Register(parent *flags.Command) error {
	_, err := parent.AddCommand(args.Name(), "Displays the specified intkey value", "Shows the value of the key <name>.", args)
	if err != nil {
		return err
	}
	return nil
}

func (args *Show) Run() error {
	// Construct client
	name := args.Args.Name
	auctionClient, err := GetClient(args, false)
	if err != nil {
		return err
	}
	value, err := auctionClient.Show(name)
	if err != nil {
		return err
	}
	PrintAuction(value)
	return nil
}

func PrintAuction(data string) {

	var rows [][]string

	for _, str := range strings.Split(data, "|") {

		tables := strings.Split(str, "!")

		itemSlice := strings.Split(tables[0], ",")

		s := fmt.Sprintf("+ Auction of the item %s with ID: %s created +", itemSlice[1], itemSlice[0])
		sLength := len(s)
		border := "+" + strings.Repeat("-", sLength-2) + "+"
		fmt.Println(border)
		fmt.Println(s)
		fmt.Println(border)

		fmt.Println("Description item:", itemSlice[2], ", PostTime:", itemSlice[3], ", ExpiryTime:", itemSlice[4])

		//if time.Now().After()

		higherAmount := -1
		pWinner := ""

		if tables[1] != "" {
			bids := strings.Split(tables[1], ";")
			for _, bid := range bids {
				bidSlice := strings.Split(bid, ",")
				amountInt, err := strconv.Atoi(bidSlice[2])
				if err != nil {
					fmt.Println(bidSlice[2])
					panic(err)
				}
				if higherAmount < 0 {
					higherAmount = amountInt
					pWinner = bidSlice[1]
				} else if amountInt > higherAmount {
					higherAmount = amountInt
					pWinner = bidSlice[1]
				}
				bidSlice[2] = bidSlice[2] + " $"
				rows = append(rows, bidSlice)
			}
		}

		expiry, err := time.Parse(layoutDate, itemSlice[4])
		if err != nil {
			panic(err)
		}
		tNow := time.Now()
		if tNow.After(expiry) {
			fmt.Println("The winner is:", pWinner)
		}
	}

	columns := []string{
		"Id Bid", "Name Bidder", "Amount", "Timestamp",
	}

	tui.Table(os.Stdout, columns, rows)

}
