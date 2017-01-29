package stats

type Stats struct {
	Counters *Counters
}

func New(prefix string) *Stats {
	return &Stats{
		Counters: newCounters(prefix),
	}
}
