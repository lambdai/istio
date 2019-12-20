// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"istio.io/istio/tools/istio-iptables/pkg/config"
)

var istioLocalIPv6 = net.IP {0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,6}

type Validator struct {
	Config *Config
}

type Config struct {
	ServerListenAddress string
	ServerOriginalPort  uint16
	ServerOriginalIP    net.IP
}
type Service struct {
	Config *Config
}
type Client struct {
	Config *Config
}

func (validator *Validator) Run() error {
	s := Service{
		validator.Config,
	}
	sError := make(chan error, 1)
	sTimer := time.NewTimer(300 * time.Second)
	defer sTimer.Stop()
	go func() {
		sError <- s.Run()
	}()

	// infinite loop
	go func() {
		c := Client{Config: validator.Config}
		for {
			_ = c.Run()
			// Avoid spamming the request to the validation server.
			// Since the TIMEWAIT socket is cleaned up in 60 second,
			// it's maintaining 60 TIMEWAIT sockets. Not big deal.
			time.Sleep(time.Second)
		}
	}()
	select {
	case <-sTimer.C:
		return errors.New("validation timeout")
	case err := <-sError:
		if err == nil {
			fmt.Println("validation passed!")
		} else {
			fmt.Println("validation failed:" + err.Error())
		}
		return err
	}
}

func NewValidator(config *config.Config, hostIP net.IP) *Validator {
	fmt.Println("in new validator: " + hostIP.String())
	// It's tricky here:
	// Connect to 127.0.0.6 will redirect to 127.0.0.1
	// Connect to ::6       will redirect to ::6
	isIpv6 := hostIP.To4() == nil
	listenIP := net.IPv4(127, 0,0,1)
	serverIP := net.IPv4(127, 0,0,6)
	formatString := "%s:%s"
	if isIpv6 {
		listenIP = istioLocalIPv6
		serverIP = istioLocalIPv6
		formatString = "[%s]:%s"
	}
	fmt.Printf("%s:%s", listenIP.String(), config.ProxyPort)
	return &Validator{
		Config: &Config{
			ServerListenAddress: fmt.Sprintf(formatString, listenIP.String(), config.ProxyPort),
			ServerOriginalPort:  config.IptablesProbePort,
			ServerOriginalIP:    serverIP,
		},
	}
}

// Write human readable response
func echo(conn io.WriteCloser, echo []byte) {
	_, _ = conn.Write(echo)
	_ = conn.Close()
}

func (s *Service) Run() error {
	l, err := net.Listen("tcp", s.Config.ServerListenAddress)
	if err != nil {
		fmt.Println("Error on listening:", err.Error())
		return err
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + s.Config.ServerListenAddress)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			continue
		}
		_, port, err := GetOriginalDestination(conn)
		if err != nil {
			fmt.Println("Error getting original dst: " + err.Error())
			conn.Close()
			continue
		}

		// echo original port for debugging.
		// Since the write amount is small it should fit in sock buffer and never blocks.
		echo(conn, []byte(strconv.Itoa(int(port))))
		// Handle connections
		// Since the write amount is small it should fit in sock buffer and never blocks.
		if port != s.Config.ServerOriginalPort {
			// This could be probe request from no where
			continue
		}
		// Server recovers the magical original port
		return nil
	}
}

func (c *Client) Run() error {
	serverOriginalAddress := fmt.Sprintf("%s:%d", c.Config.ServerOriginalIP, c.Config.ServerOriginalPort)
	conn, err := net.Dial("tcp", serverOriginalAddress)
	if err != nil {
		fmt.Printf("Error connecting to %s: %s\n", serverOriginalAddress, err.Error())
		return err
	}
	conn.Close()
	return nil
}