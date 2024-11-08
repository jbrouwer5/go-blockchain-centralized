/*
 *
 * Copyright 2015 gRPC authors.
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
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strconv"
	"time"

	pb "blockchain/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	nVersion = int32(1)
)

var (
	knownPeers = []string{}
	myAddr = "start"
)

// SERVER STUFF

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedDistributionServer
}

// func (s *server) Registrar(ctx context.Context, r *pb.Registration) (*pb.IPReply, error) {
// 	return &pb.IPReply{Ip: r.AddrMe}, nil
// }

func (s *server) Handshake(ctx context.Context, h *pb.Handshake) (*pb.KnownPeers, error) {
	log.Print(h.AddrMe)
	// knownPeers = MakeHandshake(h.AddrMe, myAddr, knownPeers)
	return &pb.KnownPeers{Ips: knownPeers}, nil
}
// END SERVER STUFF

func MakeConnection(addr string, myAddr string) (ip string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDistributionClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	r, err := c.Registrar(ctx, &pb.Registration{NVersion: nVersion, NTime: time.Now().Unix(), AddrMe: myAddr})
	if err != nil {
			log.Fatalf("could not register: %v", err)
	}
	return r.Ip
}

func MakeHandshake (addr string, addrMe string, peers []string) (ips []string) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	} else {
		log.Print("pass1")
	}
	defer conn.Close()
	c := pb.NewDistributionClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	r, err := c.Handshake(ctx, &pb.Handshake{NVersion: nVersion, NTime: time.Now().Unix(), AddrMe: addrMe, BestHeight: int32(1)})
	if err != nil {
			log.Fatalf("could not handshake: %v", err)
	} else {
		log.Print("pass2")
	}
	
	thisLevel := []string{}
	for i:=0;i<len(r.Ips);i++ {
		if !slices.Contains(peers, r.Ips[i]) && !slices.Contains(thisLevel, r.Ips[i]) &&  r.Ips[i] != addrMe{
			thisLevel = append(thisLevel, r.Ips[i])
		}
	}

	nextLevel := []string{}
	for len(thisLevel) > 0{
		for i:=0;i<len(thisLevel);i++ {
			peers = append(peers, thisLevel[i])
			newIps := MakeHandshake(thisLevel[i], addrMe, peers)
			nextLevel = append(nextLevel, newIps...)
		}

		thisLevel = []string{}
		for i:=0;i<len(nextLevel);i++ {
			if !slices.Contains(peers, nextLevel[i]) && !slices.Contains(thisLevel, nextLevel[i]){
				thisLevel = append(thisLevel, nextLevel[i])
			}
		}
		nextLevel = []string{}
	}
	
	return peers
}



func main() {
	port := os.Getenv("PORT")
	// START SERVER
	flag.Parse()
	intPort, _ := strconv.Atoi(port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", intPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDistributionServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	
	// SERVER STARTED
	// JCoin := NewJCoin()

	flag.Parse()

	// Register with the DNS 
	
	addr := "dns-seed:58333"
	myAddr = "127.0.0.1:" + fmt.Sprint(port)
	log.Print(myAddr)
	ip := MakeConnection(addr, myAddr)
	if ip == "" || ip == addr{
		log.Print("No known nodes")
	} else {
		log.Print("Known node ", ip)
		// Handshake with newly found node 
		knownPeers = append(knownPeers, ip)
		knownPeers = MakeHandshake(ip, myAddr, knownPeers)
		log.Print(knownPeers)
	}

	transactionPool := NewTxnMemoryPool()
	go SimulateTransactions(transactionPool)

	// SERVING SERVER
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

