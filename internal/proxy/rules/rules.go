package rules

import (
	"context"
	"fmt"
)

// Rules is the configuration for the rules to proxy.
type Rules struct {
	Blacklist *Blacklist `yaml:"blacklist"`
}

// Init initializes the rules.
func (r *Rules) Init(ctx context.Context) error {
	if r.Blacklist == nil {
		return nil
	}

	err := r.Blacklist.Init(ctx)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
