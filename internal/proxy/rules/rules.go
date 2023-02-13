package rules

// Rules is the configuration for the rules to proxy.
type Rules struct {
	Blacklist *Blacklist `yaml:"blacklist"`
}
