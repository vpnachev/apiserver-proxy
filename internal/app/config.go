// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"net"
	"time"

	utiliptables "github.com/gardener/apiserver-proxy/internal/iptables"
	"github.com/gardener/apiserver-proxy/internal/netif"
)

// ConfigParams lists the configuration options that can be provided to node-cache
type ConfigParams struct {
	LocalPort     string        // port to listen for dns requests
	Interface     string        // Name of the interface to be created
	Interval      time.Duration // specifies how often to run iptables rules check
	SetupIptables bool          // enable iptables setup
	Cleanup       bool          // clean the created interface and iptables
	IPAddress     string        // IP address on which the proxy is listening
}

// CacheApp contains all the config required to run node-cache.
type CacheApp struct {
	iptables      utiliptables.Interface
	iptablesRules []iptablesRule
	params        *ConfigParams
	netManager    netif.Manager
	localIP       net.IP
	localIPStr    string
	exitChan      <-chan struct{}
}

type iptablesRule struct {
	table utiliptables.Table
	chain utiliptables.Chain
	args  []string
}