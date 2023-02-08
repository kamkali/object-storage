package itest

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/kamkalis/object-storage/internal/domain"
	"github.com/kamkalis/object-storage/internal/domain/docker"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Discover and use Minio nodes", Ordered, func() {
	var (
		nodes   []domain.StorageNode
		timeout = 30 * time.Second
		ctx     context.Context
		cancel  context.CancelFunc
	)

	BeforeEach(func() {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	})
	AfterEach(func() {
		cancel()
	})

	When("connected cluster info specified", func() {
		d, err := docker.NewNodeDiscoverer("test-object-storage", "9000", "object-storage_test-object-storage")
		Expect(err).To(Not(HaveOccurred()))

		It("finds all minio nodes", func() {
			nodes, err = d.DiscoverNodes(ctx)
			Expect(err).To(Not(HaveOccurred()))
			Expect(nodes).To(HaveLen(3))
		})
	})

	When("nodes found", func() {
		It("can put and retrieve from each node", func() {
			for i, n := range nodes {
				Expect(n.IsAlive(ctx)).To(BeTrue())

				o := &domain.Object{
					ID:          fmt.Sprintf("someID%d", i),
					Content:     strings.NewReader(fmt.Sprintf("content%d", i)),
					ContentType: "text/plain",
					Size:        8,
				}
				Expect(n.PutObject(ctx, o)).To(Succeed())

				got, err := n.GetObject(ctx, o.ID)
				Expect(err).To(Not(HaveOccurred()))

				Expect(got.ID).To(Equal(o.ID))
				Expect(got.ContentType).To(Equal(o.ContentType))
				Expect(got.Size).To(Equal(o.Size))
				Expect(mustReadAll(got.Content)).To(Equal(mustReadAll(o.Content)))
			}
		})
	})
})

func mustReadAll(r io.Reader) []byte {
	all, err := io.ReadAll(r)
	Expect(err).To(Not(HaveOccurred()))
	return all
}

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "integration test")
}
