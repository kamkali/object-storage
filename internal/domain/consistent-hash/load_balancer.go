package consistent_hash

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"

	"github.com/google/uuid"
	"github.com/kamkalis/object-storage/internal/domain"
	"golang.org/x/net/context"
)

// RingLoadBalancer implements LoadBalancer based on consistent-hashing algorithm
type RingLoadBalancer struct {
	mu sync.Mutex

	nodes RingNodes
}

func NewRingLoadBalancer() *RingLoadBalancer {
	return &RingLoadBalancer{
		nodes: RingNodes{},
	}
}

func (r *RingLoadBalancer) GetNode(ctx context.Context, key uuid.UUID) (domain.StorageNode, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.nodes) == 0 {
		return nil, fmt.Errorf("no nodes in the ring")
	}

	i := sort.Search(r.nodes.Len(), func(i int) bool {
		return r.nodes[i].HashID >= crc32.ChecksumIEEE([]byte(key.String()))
	})
	if i >= r.nodes.Len() {
		i = 0
	}
	return r.nodes[i], nil
}

func (r *RingLoadBalancer) ReBalance(ctx context.Context, nodes []domain.StorageNode) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.nodes = RingNodes{}
	for _, server := range nodes {
		r.nodes = append(r.nodes, &RingNode{
			StorageNode: server,
			HashID:      crc32.ChecksumIEEE([]byte(server.ID())),
		})
	}
	sort.Sort(r.nodes)

	return nil
}
