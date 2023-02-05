package round_robin

import (
	"testing"

	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestRoundRobinLB_GetNext(t *testing.T) {
	makeMockNode := func() domain.StorageNode {
		n := mocks.NewStorageNode(t)
		n.On("IsAlive", mock.Anything).Return(true)
		return n
	}
	expectedNode := makeMockNode()
	type fields struct {
		roundRobinCount int
		servers         []domain.StorageNode
	}
	tests := []struct {
		name    string
		fields  fields
		want    domain.StorageNode
		wantErr bool
	}{
		{
			name: "gets next",
			fields: fields{
				roundRobinCount: 0,
				servers: []domain.StorageNode{
					expectedNode,
				},
			},
			want:    expectedNode,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoundRobinLB{
				roundRobinCount: tt.fields.roundRobinCount,
				servers:         tt.fields.servers,
			}
			var (
				got domain.StorageNode
				err error
			)
			for i := 0; i < 2; i++ {
				got, err = r.GetNext(context.Background())
				if (err != nil) != tt.wantErr {
					t.Errorf("GetNext() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, r.roundRobinCount, tt.fields.roundRobinCount+2)
		})
	}
}
