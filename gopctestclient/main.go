// Copyright 2018-2020 opcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	//	"reflect"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"

	//"io"
	"io/ioutil"
	"strconv"
	"strings"
)

var endpoint = "opc.tcp://localhost:4840"

type neededstruct struct {
	name     string
	min      float64
	max      float64
	value    interface{}
	nodenumb string
	panic    bool
}

func (v *neededstruct) update() {
	v.value = read(v.nodenumb)
}

func (v *neededstruct) alarm() {
	v.panic = true
}

func (v *neededstruct) unalarm() {
	v.panic = false
}

func (v *neededstruct) check() bool {
	if v.value.(float64) > v.max || v.value.(float64) < v.min {
		v.alarm()
		return true
	} else {
		return false
	}

}

func read(nodeID string) interface{} {

	ctx := context.Background()

	c := opcua.NewClient(endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err := c.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		log.Fatalf("invalid node id: %v", err)
	}

	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead: []*ua.ReadValueID{
			&ua.ReadValueID{NodeID: id},
		},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	resp, err := c.Read(req)
	if err != nil {
		log.Fatalf("Read failed: %s", err)
	}
	if resp.Results[0].Status != ua.StatusOK {
		log.Fatalf("Status not OK: %v", resp.Results[0].Status)
	}
	//log.Printf("%#v", resp.Results[0].Value.Value())
	//fmt.Println(reflect.TypeOf(resp.Results[0].Value.Value()))
	return resp.Results[0].Value.Value()

}

func main() {

	values := []neededstruct{}
	dat, err := ioutil.ReadFile("data.txt")
	if err != nil {
		panic(err)
	}
	ps := strings.Split(string(dat), "\n")
	tnodes := browseNode()
	for _, s := range ps {

		if len(s) == 0 {
			continue
		}
		ss := strings.Split(s, ":")
		minfval, _ := strconv.ParseFloat(ss[1], 64)
		maxfval, _ := strconv.ParseFloat(ss[2], 64)
		numid := ""
		for _, v := range tnodes {

			if v.Records()[0] == ss[0] {
				numid = v.Records()[2]
			}
		}

		values = append(values, neededstruct{ss[0], minfval, maxfval, read(numid), numid, false})
	}

	totg := make(chan []neededstruct)
	fromtg := make(chan string)
	go tgbot(totg, fromtg)
	for true {

		badvalues := []neededstruct{}
		for i := 0; i < len(values); i++ {
			values[i].update()

			if values[i].panic == false {
				if values[i].check() {
					fmt.Println("ALARM")
					badvalues = append(badvalues, values[i])
				}
			}
			//fmt.Println(values[i])
		}

		if len(badvalues) > 0 {
			totg <- badvalues
		}
		select {
		case msg := <-fromtg:
			if msg == "reset" {

				for i := 0; i < len(values); i++ {
					values[i].unalarm()
				}
			}
		default:

		}
		time.Sleep(5 * time.Second)
	}
}
