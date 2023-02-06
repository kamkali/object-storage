package consistent_hash

import (
	"github.com/kamkalis/object-storage/internal/domain"
)

type RingNode struct {
	domain.StorageNode
	HashID uint32
}

type RingNodes []*RingNode

func (n RingNodes) Len() int           { return len(n) }
func (n RingNodes) Less(i, j int) bool { return n[i].HashID < n[j].HashID }
func (n RingNodes) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
