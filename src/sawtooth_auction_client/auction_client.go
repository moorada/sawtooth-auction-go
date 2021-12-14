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
	bytes2 "bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	cbor "github.com/brianolson/cbor_go"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/batch_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/transaction_pb2"
	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AuctionClient struct {
	url    string
	signer *signing.Signer
}

const layoutDate = "2006-01-02T15:04:05.000Z"

func NewAuctionClient(url string, keyfile string) (AuctionClient, error) {

	var privateKey signing.PrivateKey
	if keyfile != "" {
		// Read private key file
		privateKeyStr, err := ioutil.ReadFile(keyfile)
		if err != nil {
			return AuctionClient{},
				errors.New(fmt.Sprintf("Failed to read private key: %v", err))
		}
		// Get private key object
		privateKey = signing.NewSecp256k1PrivateKey(privateKeyStr)
	} else {
		privateKey = signing.NewSecp256k1Context().NewRandomPrivateKey()
	}
	cryptoFactory := signing.NewCryptoFactory(signing.NewSecp256k1Context())
	signer := cryptoFactory.NewSigner(privateKey)
	return AuctionClient{url, signer}, nil
}

func (auctionClient AuctionClient) AddItem(idItem string, nameItem string, description string, postTime string, expiryTime string, wait uint) (string, error) {

	//payload := []byte("addItem,item1,giocattolo,bellissimo,posttime,exp,alfredo,timestamp,400,bidid1")

	return auctionClient.sendTransaction(VERB_ADDITEM, idItem, nameItem, description, postTime, expiryTime, "", "", "", "", wait)
}

func (auctionClient AuctionClient) PlaceBid(idItem string, BidId string, bidderName string, amount string, timeStamp string, wait uint) (string, error) {

	//payload := []byte("addItem,item1,giocattolo,bellissimo,posttime,exp,alfredo,timestamp,400,bidid1")
	timeStamp1 := time.Now().Add(time.Minute)

	return auctionClient.sendTransaction(VERB_PLACEBID, idItem, "", "", "", "", BidId, bidderName, amount, timeStamp1.Format(layoutDate), wait)
}

func (auctionClient AuctionClient) List() ([]map[interface{}]interface{}, error) {

	// API to call
	apiSuffix := fmt.Sprintf("%s?address=%s",
		STATE_API, auctionClient.getPrefix())
	response, err := auctionClient.sendRequest(apiSuffix, []byte{}, "", "")
	if err != nil {
		return []map[interface{}]interface{}{}, err
	}

	var toReturn []map[interface{}]interface{}
	responseMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(response), &responseMap)
	if err != nil {
		return []map[interface{}]interface{}{},
			errors.New(fmt.Sprintf("Error reading response: %v", err))
	}
	encodedEntries := responseMap["data"].([]interface{})
	for _, entry := range encodedEntries {
		entryData, ok := entry.(map[interface{}]interface{})
		if !ok {
			return []map[interface{}]interface{}{},
				errors.New("Error reading entry data")
		}
		stringData, ok := entryData["data"].(string)
		if !ok {
			return []map[interface{}]interface{}{},
				errors.New("Error reading string data")
		}
		decodedBytes, err := base64.StdEncoding.DecodeString(stringData)
		if err != nil {
			return []map[interface{}]interface{}{},
				errors.New(fmt.Sprint("Error decoding: %v", err))
		}
		foundMap := make(map[interface{}]interface{})
		err = cbor.Loads(decodedBytes, &foundMap)
		if err != nil {
			return []map[interface{}]interface{}{},
				errors.New(fmt.Sprint("Error binary decoding: %v", err))
		}
		toReturn = append(toReturn, foundMap)
	}
	return toReturn, nil
}

func (auctionClient AuctionClient) Show(name string) (string, error) {
	apiSuffix := fmt.Sprintf("%s/%s", STATE_API, auctionClient.getAddress(name))
	response, err := auctionClient.sendRequest(apiSuffix, []byte{}, "", name)
	if err != nil {
		return "", err
	}
	//fmt.Println(response)
	responseMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(response), &responseMap)
	if err != nil {
		return "", errors.New(fmt.Sprint("Error reading response: %v", err))
	}
	data, ok := responseMap["data"].(string)
	if !ok {
		return "", errors.New("Error reading as string")
	}
	responseData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New(fmt.Sprint("Error decoding response: %v", err))
	}
	//responseFinal := make(map[interface{}]interface{})

	//err = cbor.Loads(responseData, &responseFinal)
	//if err != nil {
	//	return "", errors.New(fmt.Sprint("Error binary decoding: %v", err))
	//}
	return fmt.Sprintf("%v", string(responseData)), nil
}

func (auctionClient AuctionClient) getStatus(
	batchId string, wait uint) (string, error) {

	// API to call
	apiSuffix := fmt.Sprintf("%s?id=%s&wait=%d",
		BATCH_STATUS_API, batchId, wait)
	response, err := auctionClient.sendRequest(apiSuffix, []byte{}, "", "")
	if err != nil {
		return "", err
	}

	responseMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(response), &responseMap)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading response: %v", err))
	}
	entry :=
		responseMap["data"].([]interface{})[0].(map[interface{}]interface{})
	return fmt.Sprint(entry["status"]), nil
}

func (auctionClient AuctionClient) sendRequest(
	apiSuffix string,
	data []byte,
	contentType string,
	name string) (string, error) {

	// Construct URL
	var url string
	if strings.HasPrefix(auctionClient.url, "http://") {
		url = fmt.Sprintf("%s/%s", auctionClient.url, apiSuffix)
	} else {
		url = fmt.Sprintf("http://%s/%s", auctionClient.url, apiSuffix)
	}

	// Send request to validator URL
	var response *http.Response
	var err error
	if len(data) > 0 {
		response, err = http.Post(url, contentType, bytes2.NewBuffer(data))
	} else {
		response, err = http.Get(url)
	}
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("Failed to connect to REST API: %v", err))
	}
	if response.StatusCode == 404 {
		logger.Debug(fmt.Sprintf("%v", response))
		return "", errors.New(fmt.Sprintf("No such key: %s", name))
	} else if response.StatusCode >= 400 {
		return "", errors.New(
			fmt.Sprintf("Error %d: %s", response.StatusCode, response.Status))
	}
	defer response.Body.Close()
	reponseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading response: %v", err))
	}
	return string(reponseBody), nil
}

func (auctionClient AuctionClient) sendTransaction(
	verb string, idItem string, nameItem string, description string, postTime string, expiryTime string, idBid string, bidderName string, amount string, timestamp string, wait uint) (string, error) {

	// construct the payload information in CBOR format
	//payloadData := make(map[string]interface{})
	//payloadData["Action"] = verb
	//payloadData["Name"] = name
	//payloadData["Value"] = value
	payload := []byte(verb + "," + idItem + "," + nameItem + "," + description + "," + postTime + "," + expiryTime + "," + idBid + "," + bidderName + "," + amount + "," + timestamp)
	//payload, err := cbor.Dumps(payloadData)
	//if err != nil {
	//	return "", errors.New(fmt.Sprintf("Failed to construct CBOR: %v", err))
	//}

	// construct the address
	address := auctionClient.getAddress(idItem)
	fmt.Println("Address on the MerkleTree: ", address)
	// Construct TransactionHeader
	rawTransactionHeader := transaction_pb2.TransactionHeader{
		SignerPublicKey:  auctionClient.signer.GetPublicKey().AsHex(),
		FamilyName:       FAMILY_NAME,
		FamilyVersion:    FAMILY_VERSION,
		Dependencies:     []string{}, // empty dependency list
		Nonce:            strconv.Itoa(rand.Int()),
		BatcherPublicKey: auctionClient.signer.GetPublicKey().AsHex(),
		Inputs:           []string{address},
		Outputs:          []string{address},
		PayloadSha512:    Sha512HashValue(string(payload)),
	}
	transactionHeader, err := proto.Marshal(&rawTransactionHeader)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("Unable to serialize transaction header: %v", err))
	}

	// Signature of TransactionHeader
	transactionHeaderSignature := hex.EncodeToString(
		auctionClient.signer.Sign(transactionHeader))

	// Construct Transaction
	transaction := transaction_pb2.Transaction{
		Header:          transactionHeader,
		HeaderSignature: transactionHeaderSignature,
		Payload:         payload,
	}

	// Get BatchList
	rawBatchList, err := auctionClient.createBatchList(
		[]*transaction_pb2.Transaction{&transaction})
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("Unable to construct batch list: %v", err))
	}
	batchId := rawBatchList.Batches[0].HeaderSignature
	batchList, err := proto.Marshal(&rawBatchList)
	if err != nil {
		return "", errors.New(
			fmt.Sprintf("Unable to serialize batch list: %v", err))
	}

	if wait > 0 {
		waitTime := uint(0)
		startTime := time.Now()
		response, err := auctionClient.sendRequest(
			BATCH_SUBMIT_API, batchList, CONTENT_TYPE_OCTET_STREAM, idItem)
		if err != nil {
			return "", err
		}
		for waitTime < wait {
			status, err := auctionClient.getStatus(batchId, wait-waitTime)
			if err != nil {
				return "", err
			}
			waitTime = uint(time.Now().Sub(startTime))
			if status != "PENDING" {
				return response, nil
			}
		}
		return response, nil
	}

	return auctionClient.sendRequest(
		BATCH_SUBMIT_API, batchList, CONTENT_TYPE_OCTET_STREAM, idItem)
}

func (auctionClient AuctionClient) getPrefix() string {
	return Sha512HashValue(FAMILY_NAME)[:FAMILY_NAMESPACE_ADDRESS_LENGTH]
}

func (auctionClient AuctionClient) getAddress(name string) string {
	prefix := auctionClient.getPrefix()
	nameAddress := Sha512HashValue(name)[:FAMILY_VERB_ADDRESS_LENGTH]
	return prefix + nameAddress
}

func (auctionClient AuctionClient) createBatchList(
	transactions []*transaction_pb2.Transaction) (batch_pb2.BatchList, error) {

	// Get list of TransactionHeader signatures
	transactionSignatures := []string{}
	for _, transaction := range transactions {
		transactionSignatures =
			append(transactionSignatures, transaction.HeaderSignature)
	}

	// Construct BatchHeader
	rawBatchHeader := batch_pb2.BatchHeader{
		SignerPublicKey: auctionClient.signer.GetPublicKey().AsHex(),
		TransactionIds:  transactionSignatures,
	}
	batchHeader, err := proto.Marshal(&rawBatchHeader)
	if err != nil {
		return batch_pb2.BatchList{}, errors.New(
			fmt.Sprintf("Unable to serialize batch header: %v", err))
	}

	// Signature of BatchHeader
	batchHeaderSignature := hex.EncodeToString(
		auctionClient.signer.Sign(batchHeader))

	// Construct Batch
	batch := batch_pb2.Batch{
		Header:          batchHeader,
		Transactions:    transactions,
		HeaderSignature: batchHeaderSignature,
	}

	// Construct BatchList
	return batch_pb2.BatchList{
		Batches: []*batch_pb2.Batch{&batch},
	}, nil
}
