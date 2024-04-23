package rules_test

import (
	"context"
	"testing"
	"time"
	"vbalancer/internal/proxy/rules"

	cache "github.com/microup/vcache"
	"github.com/stretchr/testify/assert"
)

const CachedDurationToEvict = 5 * time.Second

//nolint:funlen
func TestBlacklist_CheckIpBlacklist(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		b         *rules.Blacklist
		checkedIP string
		want      bool
	}{
		{
			name: "empty blacklisted",
			b: &rules.Blacklist{
				CacheDurationToEvict: CachedDurationToEvict,
				Cache:                cache.New(time.Second, CachedDurationToEvict),
				RemoteIP:             []string{""},
			},
			checkedIP: "89.207.132.170",
			want:      false,
		},
		{
			name: "ip is blacklisted",
			b: &rules.Blacklist{
				CacheDurationToEvict: CachedDurationToEvict,
				Cache:                cache.New(time.Second, CachedDurationToEvict),
				RemoteIP:             []string{"89.207.132.170", "89.207.132.172"},
			},
			checkedIP: "89.207.132.170",
			want:      true,
		},
		{
			name: "ip is blacklisted with port",

			b: &rules.Blacklist{
				CacheDurationToEvict: CachedDurationToEvict,
				Cache:                cache.New(time.Second, CachedDurationToEvict),
				RemoteIP:             []string{"89.207.132.170", "89.207.132.172"},
			},
			checkedIP: "89.207.132.170:1234",
			want:      true,
		},
		{
			name: "ip is not blacklisted",

			b: &rules.Blacklist{
				CacheDurationToEvict: CachedDurationToEvict,
				Cache:                cache.New(time.Second, CachedDurationToEvict),
				RemoteIP:             []string{"89.207.132.170", "89.207.132.175"},
			},
			checkedIP: "89.207.132.171",
			want:      false,
		},
		{
			name: "ip is not blacklisted with port",
			b: &rules.Blacklist{
				CacheDurationToEvict: CachedDurationToEvict,
				Cache:                cache.New(time.Second, CachedDurationToEvict),
				RemoteIP:             []string{"89.207.132.170", "89.207.132.175"},
			},
			checkedIP: "89.207.132.171:1234",
			want:      false,
		},
	}

	ctx := context.Background()

	for _, test := range testCases {
		err := test.b.Init(ctx)

		assert.NoError(t, err, "name: `%s`", test.name)

		assert.Equalf(t, test.b.IsBlacklistIP(test.checkedIP), test.want, "name: `%s`", test.name)
	}
}
