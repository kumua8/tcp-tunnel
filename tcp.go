package main

import (
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
)

type Tunnel struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func Run(tunnels []*Tunnel) {
	if len(tunnels) == 0 {
		log.Printf("[Run] no tunnel found, exited")
		return
	}
	var wg sync.WaitGroup
	for _, addr := range tunnels {
		wg.Add(1)
		go func(addr *Tunnel) {
			defer wg.Done()
			Listen(addr)
		}(addr)
	}
	wg.Wait()
}

func Listen(tunnel *Tunnel) {
	ln, err := net.Listen("tcp", tunnel.From)
	if err != nil {
		log.Printf("[Listen] listen tcp port %v failed, err: %v", tunnel.From, err)
		return
	}
	log.Printf("[Listen] listening tcp port %v", tunnel.From)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("[Listen] accept conn %v failed, err: %v", tunnel.From, err)
			continue
		}
		go func() {
			defer Recovery()
			serve(conn, tunnel)
		}()
	}
}

func serve(cliConn net.Conn, tunnel *Tunnel) {
	defer func() {
		_ = cliConn.Close()
	}()

	serverConn, err := net.Dial("tcp", tunnel.To)
	if err != nil {
		log.Printf("[serve] dial tcp conn %v failed, err: %v", tunnel.To, err)
		return
	}

	defer func() {
		_ = serverConn.Close()
	}()

	err = runTunnel(cliConn, serverConn)
	if err != nil {
		log.Printf("[server] runTunnel exec err: %v", err)
	}
}

func runTunnel(cliConn net.Conn, serverConn net.Conn) error {
	var errGroup errgroup.Group

	errGroup.Go(func() error {
		_, err := io.Copy(cliConn, serverConn)
		return err
	})
	errGroup.Go(func() error {
		_, err := io.Copy(serverConn, cliConn)
		return err
	})

	return errGroup.Wait()
}

func Recovery() {
	if r := recover(); r != nil {
		buf := make([]byte, 1<<18)
		n := runtime.Stack(buf, false)
		log.Printf("%v, STACK: %s", r, buf[0:n])
	}
}
