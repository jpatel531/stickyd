package config

type Config struct {
	PrefixStats  string      `json:"prefixStats"`
	Servers      []*Frontend `json:"servers,omitempty"`
	DumpMessages bool        `json:"dumpMessages,omitempty"`
}

type Frontend struct {
	AddressIPV6 bool
	Address     string
	Port        int
}
