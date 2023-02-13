package rules

import "net"

type Blacklist struct {
	RemoteIPs []string `yaml:"remoteIp"`
}

// CheckIPBlacklist checks if the ip is in the blacklist.
func (b *Blacklist) IsIPInBlacklist(ip string) bool {
    host, _, err := net.SplitHostPort(ip)
    if err != nil {
        host = ip
    }

    for _, value := range b.RemoteIPs {
        if value == host {
            return true
        }
    }

    return false
}
