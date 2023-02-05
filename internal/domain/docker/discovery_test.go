package docker

import (
	"github.com/kamkalis/object-storage/internal/domain"
)

var _ domain.NodeDiscoverer = (*NodeDiscoverer)(nil)
