package bookManager

import (
	"context"
	"sync"
	"time"
)

type BookManager struct {
	mu    sync.RWMutex
	books map[string]*OrderBook
}

func NewBookManager() *BookManager {
	return &BookManager{
		books: make(map[string]*OrderBook, 16),
	}
}

func (m *BookManager) GetOrCreate(productID string) *OrderBook {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ob, ok := m.books[productID]; ok {
		return ob
	}
	ob := NewOrderBook(productID)
	m.books[productID] = ob
	return ob
}

// RebuildAllHeaps 重建Heaps
func (m *BookManager) RebuildAllHeaps() {
	m.mu.RLock()
	books := make([]*OrderBook, 0, len(m.books))
	for _, ob := range m.books {
		books = append(books, ob)
	}
	m.mu.RUnlock()

	for _, ob := range books {
		ob.RebuildHeaps()
	}
}

type TopNSnapshot struct {
	ProductID string  `json:"product_id"`
	Ts        int64   `json:"ts"` // epoch ms
	Bids      []Level `json:"bids"`
	Asks      []Level `json:"asks"`
}

func (m *BookManager) SnapshotTopNAll(n int) []TopNSnapshot {
	m.mu.RLock()
	books := make([]*OrderBook, 0, len(m.books))
	for _, ob := range m.books {
		books = append(books, ob)
	}
	m.mu.RUnlock()

	nowMs := time.Now().UnixMilli()
	out := make([]TopNSnapshot, 0, len(books))

	for _, ob := range books {
		bids, asks := ob.SnapshotTopN(n)
		out = append(out, TopNSnapshot{
			ProductID: ob.productID,
			Ts:        nowMs,
			Bids:      bids,
			Asks:      asks,
		})
	}
	return out
}

func (m *BookManager) StartSnapshotTimerAsync(ctx context.Context, interval time.Duration, n int, callback func([]TopNSnapshot)) {

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				snapshots := m.SnapshotTopNAll(n)
				if len(snapshots) == 0 {
					continue
				}
				// 执行回调（建议考虑是否需要 go callback(response) 异步处理）
				go callback(snapshots)

			}
		}
	}()
}
