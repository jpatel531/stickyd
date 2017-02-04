package config

// TODO add:
// 	- keyNameSanitize
type Config struct {
	PrefixStats  string      `json:"prefixStats"`
	Servers      []*Frontend `json:"servers,omitempty"`
	DumpMessages bool        `json:"dumpMessages,omitempty"`
	KeyFlush     KeyFlush    `json:"keyFlush"`
	MgmtAddress  string      `json:"mgmtAddress,omitempty"`
	MgmtPort     int         `json:"mgmtPort,omitempty"`
}

type Frontend struct {
	AddressIPV6 bool   `json:"addressIPV6,omitempty"`
	Address     string `json:"address,omitempty"`
	Port        int    `json:"port,omitempty"`
}

type KeyFlush struct {
	Interval int    `json:"interval,omitempty"`
	Percent  int    `json:"percent,omitempty"`
	Log      string `json:"log,omitempty"`
}
