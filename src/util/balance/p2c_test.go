package p2c

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mathx"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

func init() {
	logx.Disable()
}

func TestP2cPicker_PickNil(t *testing.T) {
	builder := new(p2cPickerBuilder)
	picker := builder.Build(base.PickerBuildInfo{})
	_, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/",
		Ctx:            context.Background(),
	})
	assert.NotNil(t, err)
}

func TestP2cPicker_Pick(t *testing.T) {
	tests := []struct {
		name       string
		candidates int
		threshold  float64
	}{
		{
			name:       "single",
			candidates: 1,
			threshold:  0.9,
		},
		{
			name:       "two",
			candidates: 2,
			threshold:  0.5,
		},
		{
			name:       "multiple",
			candidates: 100,
			threshold:  0.95,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			const total = 10000
			builder := new(p2cPickerBuilder)
			ready := make(map[balancer.SubConn]base.SubConnInfo)
			for i := 0; i < test.candidates; i++ {
				ready[new(mockClientConn)] = base.SubConnInfo{
					Address: resolver.Address{
						Addr:       strconv.Itoa(i),
						Attributes: attributes.New(tagName, "mbp"),
					},
				}
			}

			picker := builder.Build(base.PickerBuildInfo{
				ReadySCs: ready,
			})
			var wg sync.WaitGroup
			wg.Add(total)
			for i := 0; i < total; i++ {
				nctx := context.Background()
				if i%20 == 0 {
					nctx = context.WithValue(nctx, tagName, "mbp")
				}
				result, err := picker.Pick(balancer.PickInfo{
					FullMethodName: "/",
					Ctx:            nctx,
				})
				assert.Nil(t, err)
				if i%100 == 0 {
					err = status.Error(codes.DeadlineExceeded, "deadline")
				}
				go func() {
					runtime.Gosched()
					result.Done(balancer.DoneInfo{
						Err: err,
					})
					wg.Done()
				}()
			}

			wg.Wait()
			dist := make(map[interface{}]int)
			conns := picker.(*p2cPicker).conns
			for _, conn := range conns {
				dist[conn.addr.Addr] = int(conn.requests)
			}

			entropy := mathx.CalcEntropy(dist)
			assert.True(t, entropy > test.threshold, fmt.Sprintf("entropy is %f, less than %f",
				entropy, test.threshold))
		})
	}
}

type mockClientConn struct{}

func (m mockClientConn) UpdateAddresses(addresses []resolver.Address) {
}

func (m mockClientConn) Connect() {
}
