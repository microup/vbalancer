package rules_test

import (
	"testing"
	"vbalancer/internal/proxy/rules"
)

func TestBlacklist_CheckIpBlacklist(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		b         *rules.Blacklist
		checkedIP string
		want      bool
	}{
		{
			name: "empty blacklisted",
			b: &rules.Blacklist{
				RemoteIPs: []string{""},
			},
			checkedIP: "89.207.132.170",
			want:      false,
		},
		{
			name: "ip is blacklisted",
			b: &rules.Blacklist{
				RemoteIPs: []string{"89.207.132.170", "89.207.132.172"},
			},
			checkedIP: "89.207.132.170",
			want:      true,
		},
		{
			name: "ip is blacklisted with port",
			b: &rules.Blacklist{
				RemoteIPs: []string{"89.207.132.170", "89.207.132.172"},
			},
			checkedIP: "89.207.132.170:1234",
			want:      true,
		},		
		{
			name: "ip is not blacklisted",
			b: &rules.Blacklist{
				RemoteIPs: []string{"89.207.132.170", "89.207.132.175"},
			},
			checkedIP: "89.207.132.171",
			want:      false,
		},
		{
			name: "ip is not blacklisted with port",
			b: &rules.Blacklist{
				RemoteIPs: []string{"89.207.132.170", "89.207.132.175"},
			},
			checkedIP: "89.207.132.171:1234",
			want:      false,
		},		
	}

	for _, c := range cases {
		if got := c.b.IsIPInBlacklist(c.checkedIP); got != c.want {
			t.Errorf("name: `%s` = %v, want %v", c.name, got, c.want)
		}
	}
}
