package stat

import (
	"sync"
	"sync/atomic"
	"time"
)

type QueryStat struct {
	ClusterTable string
	Count        int64
	Sum          int64
}

type StatReport struct {
	mu   sync.RWMutex
	runs map[string]*QueryStat
}

func NewStat() *StatReport {
	return &StatReport{
		runs: make(map[string]*QueryStat),
	}
}

func (m *StatReport) getItem(key string) *QueryStat {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.runs[key]
}

// copy 复制, 复制完后，初始化原值
func (m *StatReport) copyItem(item *QueryStat) *QueryStat {
	return &QueryStat{
		Count: atomic.SwapInt64(&item.Count, 0),
		Sum:   atomic.SwapInt64(&item.Sum, 0),
	}

}

func (m *StatReport) StatInfo() []*QueryStat {
	m.mu.RLock()
	defer m.mu.RUnlock()

	items := make([]*QueryStat, 0)
	for key, item := range m.runs {
		n := m.copyItem(item)
		n.ClusterTable = key
		items = append(items, n)
	}

	return items

}

func (m *StatReport) addItem(key string) *QueryStat {
	m.mu.Lock()
	defer m.mu.Unlock()
	// recheck again
	if m.runs[key] == nil {
		item := &QueryStat{}
		m.copyItem(item)

		m.runs[key] = item
	}

	return m.runs[key]
}

func (m *StatReport) IncQuery(cluster, table string, elapse time.Duration) {
	key := cluster + "." + table

	item := m.getItem(key)
	if item == nil {
		item = m.addItem(key)
	}

	micro := elapse.Nanoseconds() / 1000

	atomic.AddInt64(&item.Count, 1)
	atomic.AddInt64(&item.Sum, micro)
}
