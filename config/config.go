package config

// TODO add:
// 	- keyNameSanitize
type Config struct {
	PrefixStats      string      `json:"prefixStats,omitempty"`
	Servers          []*Frontend `json:"servers,omitempty"`
	DumpMessages     bool        `json:"dumpMessages,omitempty"`
	KeyFlush         KeyFlush    `json:"keyFlush,omitempty"`
	MgmtAddress      string      `json:"mgmtAddress,omitempty"`
	MgmtPort         int         `json:"mgmtPort,omitempty"`
	Backends         []string    `json:"backends,omitempty"`
	FlushInterval    int         `json:"flushInterval,omitempty"`
	PercentThreshold []int       `json:"percentThreshold,omitempty"`
	GraphitePort     int         `json:"graphitePort,omitempty"`
	GraphiteHost     string      `json:"graphiteHost,omitempty"`
	Graphite         Graphite    `json:"graphite,omitempty"`
}

type Frontend struct {
	Server      string `json:"server,omitempty"`
	AddressIPV6 bool   `json:"addressIPV6,omitempty"`
	Address     string `json:"address,omitempty"`
	Port        int    `json:"port,omitempty"`
}

type KeyFlush struct {
	Interval int    `json:"interval,omitempty"`
	Percent  int    `json:"percent,omitempty"`
	Log      string `json:"log,omitempty"`
}

type Graphite struct {
	LegacyNamespace *bool   `json:"legacyNamespace,omitempty"`
	GlobalPrefix    *string `json:"globalPrefix,omitempty"`
	GlobalSuffix    *string `json:"globalSuffix,omitempty"`
	PrefixCounter   string  `json:"prefixCounter,omitempty"`
	PrefixTimer     string  `json:"prefixTimer,omitempty"`
	PrefixGauge     string  `json:"prefixGauge,omitempty"`
	PrefixSet       string  `json:"prefixSet,omitempty"`
}
