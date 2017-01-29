package config

// TODO add:
// 	- keyNameSanitize
type Config struct {
	PrefixStats  string      `json:"prefixStats"`
	Servers      []*Frontend `json:"servers,omitempty"`
	DumpMessages bool        `json:"dumpMessages,omitempty"`
	KeyFlush     KeyFlush    `json:"keyFlush"`
}

type Frontend struct {
	AddressIPV6 bool
	Address     string
	Port        int
}

type KeyFlush struct {
	Interval int `json:"interval"`
}
