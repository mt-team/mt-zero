package p2c

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/syncx"
	"github.com/tal-tech/go-zero/core/timex"
	"github.com/tal-tech/go-zero/zrpc/internal/codes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const (
	// Name is the name of p2c balancer.
	Name = "p2c_ewma"

	decayTime       = int64(time.Second * 10) // default value from finagle
	forcePick       = int64(time.Second)
	initSuccess     = 1000
	throttleSuccess = initSuccess / 2
	penalty         = int64(math.MaxInt32)
	pickTimes       = 3
	logInterval     = time.Minute
	tagName         = "xTag"
)

var emptyPickResult balancer.PickResult

func init() {
	balancer.Register(newBuilder())
}

type p2cPickerBuilder struct{}

func (b *p2cPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	readySCs := info.ReadySCs
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var conns []*subConn
	for conn, connInfo := range readySCs {
		var xTag string
		if connInfo.Address.Attributes != nil {
			if v := connInfo.Address.Attributes.Value(tagName); v != nil {
				xTag = v.(string)
			}

		}
		conns = append(conns, &subConn{
			addr:    connInfo.Address,
			conn:    conn,
			success: initSuccess,
			tag:     xTag,
		})
	}

	return &p2cPicker{
		conns: conns,
		r:     rand.New(rand.NewSource(time.Now().UnixNano())),
		stamp: syncx.NewAtomicDuration(),
	}
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, new(p2cPickerBuilder), base.Config{HealthCheck: true})
}

type p2cPicker struct {
	conns []*subConn
	r     *rand.Rand
	stamp *syncx.AtomicDuration
	lock  sync.Mutex
}

func (p *p2cPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	var tag string
	if v := info.Ctx.Value(tagName); v != nil {
		tag = v.(string)
	}

	var chosen *subConn
	switch len(p.conns) {
	case 0:
		return emptyPickResult, balancer.ErrNoSubConnAvailable
	case 1:
		chosen = p.choose(p.conns[0], nil)
	case 2:
		conns := p.matchTag(tag, p.conns)
		if len(conns) == 1 {
			chosen = conns[0]
			break
		}
		chosen = p.choose(p.conns[0], p.conns[1])
	default:
		conns := p.matchTag(tag, p.conns)
		switch len(conns) {
		case 0:
			chosen = p.moreChoose(p.conns)
		case 1:
			chosen = conns[0]
		default:
			chosen = p.moreChoose(conns)
		}
	}

	atomic.AddInt64(&chosen.inflight, 1)
	atomic.AddInt64(&chosen.requests, 1)

	return balancer.PickResult{
		SubConn: chosen.conn,
		Done:    p.buildDoneFunc(chosen),
	}, nil
}

func (p *p2cPicker) moreChoose(conns []*subConn) *subConn {
	var node1, node2 *subConn
	for i := 0; i < pickTimes; i++ {
		a := p.r.Intn(len(conns))
		b := p.r.Intn(len(conns) - 1)
		if b >= a {
			b++
		}
		node1 = conns[a]
		node2 = conns[b]
		if node1.healthy() && node2.healthy() {
			break
		}
	}
	return p.choose(node1, node2)
}

func (p *p2cPicker) matchTag(connsTag string, pickTag []*subConn) []*subConn {
	if connsTag == "" {
		return pickTag
	}
	pTag := make([]*subConn, 0, len(pickTag))
	for _, v := range pickTag {
		if v.tag == connsTag {
			pTag = append(pTag, v)
		}
	}
	return pTag
}

func (p *p2cPicker) buildDoneFunc(c *subConn) func(info balancer.DoneInfo) {
	start := int64(timex.Now())
	return func(info balancer.DoneInfo) {
		atomic.AddInt64(&c.inflight, -1)
		now := timex.Now()
		last := atomic.SwapInt64(&c.last, int64(now))
		td := int64(now) - last
		if td < 0 {
			td = 0
		}
		w := math.Exp(float64(-td) / float64(decayTime))
		lag := int64(now) - start
		if lag < 0 {
			lag = 0
		}
		olag := atomic.LoadUint64(&c.lag)
		if olag == 0 {
			w = 0
		}
		atomic.StoreUint64(&c.lag, uint64(float64(olag)*w+float64(lag)*(1-w)))
		success := initSuccess
		if info.Err != nil && !codes.Acceptable(info.Err) {
			success = 0
		}
		osucc := atomic.LoadUint64(&c.success)
		atomic.StoreUint64(&c.success, uint64(float64(osucc)*w+float64(success)*(1-w)))

		stamp := p.stamp.Load()
		if now-stamp >= logInterval {
			if p.stamp.CompareAndSwap(stamp, now) {
				p.logStats()
			}
		}
	}
}

func (p *p2cPicker) choose(c1, c2 *subConn) *subConn {
	start := int64(timex.Now())
	if c2 == nil {
		atomic.StoreInt64(&c1.pick, start)
		return c1
	}

	if c1.load() > c2.load() {
		c1, c2 = c2, c1
	}

	pick := atomic.LoadInt64(&c2.pick)
	if start-pick > forcePick && atomic.CompareAndSwapInt64(&c2.pick, pick, start) {
		return c2
	}

	atomic.StoreInt64(&c1.pick, start)
	return c1
}

func (p *p2cPicker) logStats() {
	var stats []string

	p.lock.Lock()
	defer p.lock.Unlock()

	for _, conn := range p.conns {
		stats = append(stats, fmt.Sprintf("conn: %s, load: %d, reqs: %d",
			conn.addr.Addr, conn.load(), atomic.SwapInt64(&conn.requests, 0)))
	}

	logx.Statf("p2c - %s", strings.Join(stats, "; "))
}

type subConn struct {
	addr     resolver.Address
	conn     balancer.SubConn
	lag      uint64
	inflight int64
	success  uint64
	requests int64
	last     int64
	pick     int64
	tag      string
}

func (c *subConn) healthy() bool {
	return atomic.LoadUint64(&c.success) > throttleSuccess
}

func (c *subConn) load() int64 {
	// plus one to avoid multiply zero
	lag := int64(math.Sqrt(float64(atomic.LoadUint64(&c.lag) + 1)))
	load := lag * (atomic.LoadInt64(&c.inflight) + 1)
	if load == 0 {
		return penalty
	}

	return load
}
