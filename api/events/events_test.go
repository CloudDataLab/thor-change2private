// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package events_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/vechain/thor/api/events"
	"github.com/vechain/thor/block"
	"github.com/vechain/thor/logdb"
	"github.com/vechain/thor/thor"
	"github.com/vechain/thor/tx"
)

var contractAddr = thor.BytesToAddress([]byte("contract"))
var ts *httptest.Server

func TestEvents(t *testing.T) {
	initEventServer(t)
	defer ts.Close()
	getEvents(t)
}

func getEvents(t *testing.T) {
	t0 := thor.BytesToBytes32([]byte("topic0"))
	t1 := thor.BytesToBytes32([]byte("topic1"))
	limit := 5
	filter := &events.Filter{
		Range: &logdb.Range{
			Unit: "",
			From: 0,
			To:   10,
		},
		Options: &logdb.Options{
			Offset: 0,
			Limit:  uint64(limit),
		},
		Order:   "",
		Address: &contractAddr,
		TopicSets: []*events.TopicSet{
			&events.TopicSet{
				Topic0: &t0,
			},
			&events.TopicSet{
				Topic1: &t1,
			},
		},
	}
	f, err := json.Marshal(&filter)
	if err != nil {
		t.Fatal(err)
	}
	res := httpPost(t, ts.URL+"/events?address="+contractAddr.String(), f)
	var logs []*events.FilteredEvent
	if err := json.Unmarshal(res, &logs); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, limit, len(logs), "should be `limit` logs")
}

func initEventServer(t *testing.T) {
	db, err := logdb.NewMem()
	if err != nil {
		t.Fatal(err)
	}
	txEv := &tx.Event{
		Address: contractAddr,
		Topics:  []thor.Bytes32{thor.BytesToBytes32([]byte("topic0")), thor.BytesToBytes32([]byte("topic1"))},
		Data:    []byte("data"),
	}

	header := new(block.Builder).Build().Header()
	for i := 0; i < 100; i++ {
		if err := db.Prepare(header).ForTransaction(thor.BytesToBytes32([]byte("txID")), thor.BytesToAddress([]byte("txOrigin"))).
			Insert(tx.Events{txEv}, nil).Commit(); err != nil {
			if err != nil {
				t.Fatal(err)
			}
		}
		header = new(block.Builder).ParentID(header.ID()).Build().Header()
	}

	router := mux.NewRouter()
	events.New(db).Mount(router, "/events")
	ts = httptest.NewServer(router)
}

func httpPost(t *testing.T, url string, data []byte) []byte {
	res, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	r, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	return r
}
