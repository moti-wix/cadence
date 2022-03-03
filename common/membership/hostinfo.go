// The MIT License (MIT)

// Copyright (c) 2017-2020 Uber Technologies Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package membership

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	PortTchannel = "tchannel"
	PortGRPC     = "grpc"
)

// PortMap is a map of port names to port numbers.
type PortMap map[string]uint16

// HostInfo is a type that contains the info about a cadence host
type HostInfo struct {
	addr     string // ip:port returned by peer provider
	ip       string // @todo should we set this to net.IP ?
	identity string
	portMap  PortMap // ports host is listening to
}

// NewHostInfo creates a new HostInfo instance
func NewHostInfo(addr string) HostInfo {
	ip, _, _ := net.SplitHostPort(addr)
	return HostInfo{
		addr: addr,
		ip:   ip,
	}
}

// String formats a PortMap into a string of name:port pairs
func (m PortMap) String() string {
	res := make([]string, 0, len(m))
	for name, port := range m {
		res = append(res, fmt.Sprintf("%s:%d", name, port))
	}
	return strings.Join(res, ", ")
}

// NewDetailedHostInfo creates a new HostInfo instance with identity and portmap information
func NewDetailedHostInfo(addr string, identity string, portMap PortMap) HostInfo {
	ip, _, _ := net.SplitHostPort(addr)
	return HostInfo{
		addr:     addr,
		ip:       ip,
		identity: identity,
		portMap:  portMap,
	}
}

// GetAddress returns the ip:port address
func (hi HostInfo) GetAddress() string {
	return hi.addr
}

// GetNamedAddress returns the ip:port address
func (hi HostInfo) GetNamedAddress(port string) (string, error) {
	if port, set := hi.portMap[port]; set {
		return net.JoinHostPort(hi.ip, strconv.Itoa(int(port))), nil
	}

	return "", fmt.Errorf("port %q is not set for %+v", port, hi)
}

// Belongs tells if ip:port is assigned to this member
func (hi HostInfo) Belongs(address string) (bool, error) {

	if hi.addr == address {
		return true, nil
	}

	ip, port, err := net.SplitHostPort(address)
	if err != nil {
		return false, err
	}

	if ip != hi.ip {
		return false, nil
	}

	for _, number := range hi.portMap {
		if port == strconv.Itoa(int(number)) {
			return true, nil
		}
	}
	return false, nil
}

// Identity implements ringpop's Membership interface
func (hi HostInfo) Identity() string {
	// if identity is not set, return address
	if hi.identity == "" {
		return hi.addr
	}

	return hi.identity
}

// Label is a noop function to conform to ringpop hashring member interface
func (hi HostInfo) Label(key string) (value string, has bool) {
	return "", false
}

// SetLabel is a noop function to conform to ringpop hashring member interface
func (hi HostInfo) SetLabel(key string, value string) {
}

// String will return a human-readable host details
func (hi HostInfo) String() string {
	return fmt.Sprintf("addr: %s, identity: %s, portMap: %s", hi.addr, hi.identity, hi.portMap)
}
