package stats

import (
	"sync/atomic"
)

func NewProcessStats(startupTime int64) *ProcessStats {
	return &ProcessStats{
		Messages: &ProcessMessages{
			badLinesSeen:    0,
			lastMessageSeen: startupTime,
		},
	}
}

type ProcessStats struct {
	Messages *ProcessMessages
}

type ProcessMessages struct {
	badLinesSeen    int64
	lastMessageSeen int64
}

func (p *ProcessMessages) IncrBadLinesSeen() {
	atomic.AddInt64(&p.badLinesSeen, 1)
}

func (p *ProcessMessages) BadLinesSeen() int64 {
	return atomic.LoadInt64(&p.badLinesSeen)
}

func (p *ProcessMessages) SetLastMessageSeen(n int64) {
	p.lastMessageSeen = n
}

func (p *ProcessMessages) LastMessageSeen() int64 {
	return p.lastMessageSeen
}
