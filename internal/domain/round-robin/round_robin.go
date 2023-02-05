package round_robin

import (
	"math"
	"sync"

	"github.com/kamkalis/object-storage/internal/domain"
	"golang.org/x/net/context"
)

type RoundRobinLB struct {
	mu              sync.Mutex
	roundRobinCount int
	servers         []domain.StorageNode
}

func (r *RoundRobinLB) GetNext(ctx context.Context) (domain.StorageNode, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	node := r.servers[r.roundRobinCount%len(r.servers)]
	for !node.IsAlive(ctx) {
		r.roundRobinCount++
		node = r.servers[r.roundRobinCount%len(r.servers)]
	}

	if r.roundRobinCount == math.MaxInt {
		r.roundRobinCount = -1
	}
	r.roundRobinCount++

	return node, nil
}

func (r *RoundRobinLB) ReBalance(ctx context.Context, nodes []domain.StorageNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.servers = nodes
	r.roundRobinCount = 0
	return nil
}
