package rpc

import (
	"strings"
	"sync"
	"time"
)

// Discovery for target app address, with load balancing and fail over.
//
//    1. The discovery will refresh the service list every 10 seconds by
//       default. The selection strategy is roundrobin, for keeping the
//       connections balanced.
//    2. Once a connection to a specified host failed, discard the host and
//       try the next. And loop the process until find a connectable host
//       or the retry time reaches the max time. If all the hosts are
//       discarded, refresh the host list.
//    3. The discovery also accepts an extra address pair, for fallback
//       when the discovery agent is down. If no backup address found, zone
//       will try to pick one from the discarded addresses.
//
type Discovery struct {
	// target means which service should to be found.
	// eg. comment for Comment service.
	target string

	// a fallback choice if discovery agent is down.
	address *Address

	// ttl specified the interval for the discovery to refresh the service list.
	ttl time.Duration

	// the real woker to fetch the service list.
	// diplomat *diplomat.Diplomat

	// record the last time the address is acessed,
	// used for refreshing the service list after time expired.
	lastCleanupTime *time.Time

	// record the unnormal address to prevent failure.
	discardedAddrs map[string]struct{}

	mu sync.Mutex
}

type Address struct {
	IP   string
	Port string
}

// newAddressFromString accept a string in "ip:port" format.
func newAddressFromString(addrString string) *Address {
	splitedAddr := strings.Split(addrString, ":")

	return &Address{
		IP:   splitedAddr[0],
		Port: splitedAddr[1],
	}
}

func (addr *Address) String() string {
	return addr.IP + ":" + addr.Port
}

func (addr *Address) Valid() bool {
	if addr.IP != "" && addr.Port != "" {
		return true
	}

	return false
}

func NewDiscovery(targetName string) *Discovery {
	return &Discovery{
		target:         targetName,
		ttl:            10 * time.Second,
		//diplomat:       diplomat.Discover(),
		discardedAddrs: map[string]struct{}{},
	}
}

func NewDiscoveryWithFallback(targetName string, address *Address) *Discovery {
	return &Discovery{
		target:         targetName,
		address:        address,
		ttl:            10 * time.Second,
		//diplomat:       diplomat.Discover(),
		discardedAddrs: map[string]struct{}{},
	}
}

// GetAddress try to get a usable address from consul.
func (d *Discovery) GetAddress() (*Address, error) {


	consulServiceMap := map[string]*Address {
		"zvideo-service": &Address{
			"0.0.0.0",
			"8000",
		},
	}

	return consulServiceMap[d.target],nil

	// 服务发现，得到IP:PORT;
	// return  &Address{IP: "0.0.0.0", Port: "9000"}, nil
}

func (d *Discovery) DiscardAddress(address *Address) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if address == nil {
		return
	}
	d.discardedAddrs[address.String()] = struct{}{}
}

func (d *Discovery) cleanup() {
	// Don't lock
	now := time.Now()
	if d.lastCleanupTime != nil {
		if d.lastCleanupTime.Add(d.ttl).Before(now) {
			d.discardedAddrs = map[string]struct{}{}
			d.lastCleanupTime = &now
		}
	} else {
		d.lastCleanupTime = &now
	}
}
