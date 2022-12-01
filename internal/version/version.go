package version

// DON'T TOUCH THIS VALUE | AUTO INCREASE WHEN BUILT AND DEPLOYED TO REMOTE SERVER.
var version = 1

func Get() uint32 {
	return uint32(version)
}
