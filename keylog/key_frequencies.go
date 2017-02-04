package keylog

type keyFrequencies []keyFrequency

type keyFrequency struct {
	key       string
	frequency int
}

func (k keyFrequencies) Len() int {
	return len(k)
}

func (k keyFrequencies) Less(i, j int) bool {
	return k[i].frequency > k[j].frequency
}

func (k keyFrequencies) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}
