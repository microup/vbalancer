package rules

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/microup/vcache"
)

// Blacklist defines the blacklist configuration.
type Blacklist struct {
	Cache *cache.VCache 
	CacheDurationToEvict time.Duration `yaml:"cacheDurationToEvict"`
	RemoteIP []string `yaml:"remoteIp"`
}

// Init initializes the blacklist.
func (b *Blacklist) Init(ctx context.Context) error {
	b.Cache = cache.New(time.Second, b.CacheDurationToEvict)
	
	err := b.Cache.StartEvict(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// IsBlacklistIP checks if the ip is in the blacklist.
func (b *Blacklist) IsBlacklistIP(checkIP string) bool {
	val, found := b.Cache.Get(checkIP)
	if found {
		b, ok := val.(bool)
		if ok {
			return b
		}
    } 

	host, _, err := net.SplitHostPort(checkIP)
	if err != nil {
		host = checkIP
	}

	for _, value := range b.RemoteIP {
		if value == host {
			_ = b.Cache.Add(checkIP, true)

			return true
		}
	}

	_ = b.Cache.Add(checkIP, false)

	return false
}
