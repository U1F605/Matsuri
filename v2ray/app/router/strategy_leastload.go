package router

import (
	"context"
	"math"
	"sort"
	"time"

	core "github.com/v2fly/v2ray-core/v5"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/dice"
	"github.com/v2fly/v2ray-core/v5/features/extension"
)

// LeastLoadStrategy represents a least load balancing strategy
type LeastLoadStrategy struct {
	settings *StrategyLeastLoadConfig
	costs    *WeightManager

	observer extension.Observatory

	ctx context.Context
}

func (l *LeastLoadStrategy) GetPrincipleTarget(strings []string) []string {
	var ret []string
	nodes := l.pickOutbounds(strings)
	for _, v := range nodes {
		ret = append(ret, v.Tag)
	}
	return ret
}

// NewLeastLoadStrategy creates a new LeastLoadStrategy with settings
func NewLeastLoadStrategy(settings *StrategyLeastLoadConfig) *LeastLoadStrategy {
	return &LeastLoadStrategy{
		settings: settings,
		costs: NewWeightManager(
			settings.Costs, 1,
			func(value, cost float64) float64 {
				return value * math.Pow(cost, 0.5)
			},
		),
	}
}

// node is a minimal copy of HealthCheckResult
// we don't use HealthCheckResult directly because
// it may change by health checker during routing
type node struct {
	Tag              string
	CountAll         int
	CountFail        int
	RTTAverage       time.Duration
	RTTDeviation     time.Duration
	RTTDeviationCost time.Duration
}

func (l *LeastLoadStrategy) InjectContext(ctx context.Context) {
	l.ctx = ctx
}

func (l *LeastLoadStrategy) PickOutbound(candidates []string) string {
	selects := l.pickOutbounds(candidates)
	count := len(selects)
	if count == 0 {
		// goes to fallbackTag
		return ""
	}
	return selects[dice.Roll(count)].Tag
}

func (l *LeastLoadStrategy) pickOutbounds(candidates []string) []*node {
	qualified := l.getNodes(candidates, time.Duration(l.settings.MaxRTT))
	selects := l.selectLeastLoad(qualified)
	return selects
}

// selectLeastLoad selects nodes according to Baselines and Expected Count.
//
// The strategy always improves network response speed, not matter which mode below is configured.
// But they can still have different priorities.
//
// 1. Bandwidth priority: no Baseline + Expected Count > 0.: selects `Expected Count` of nodes.
// (one if Expected Count <= 0)
//
// 2. Bandwidth priority advanced: Baselines + Expected Count > 0.
// Select `Expected Count` amount of nodes, and also those near them according to baselines.
// In other words, it selects according to different Baselines, until one of them matches
// the Expected Count, if no Baseline matches, Expected Count applied.
//
// 3. Speed priority: Baselines + `Expected Count <= 0`.
// go through all baselines until find selects, if not, select none. Used in combination
// with 'balancer.fallbackTag', it means: selects qualified nodes or use the fallback.
func (l *LeastLoadStrategy) selectLeastLoad(nodes []*node) []*node {
	if len(nodes) == 0 {
		newError("least load: no qualified outbound").AtInfo().WriteToLog()
		return nil
	}
	expected := int(l.settings.Expected)
	availableCount := len(nodes)
	if expected > availableCount {
		return nodes
	}

	if expected <= 0 {
		expected = 1
	}
	if len(l.settings.Baselines) == 0 {
		return nodes[:expected]
	}

	count := 0
	// go through all base line until find expected selects
	for _, b := range l.settings.Baselines {
		baseline := time.Duration(b)
		for i := 0; i < availableCount; i++ {
			if nodes[i].RTTDeviationCost > baseline {
				break
			}
			count = i + 1
		}
		// don't continue if find expected selects
		if count >= expected {
			newError("applied baseline: ", baseline).AtDebug().WriteToLog()
			break
		}
	}
	if l.settings.Expected > 0 && count < expected {
		count = expected
	}
	return nodes[:count]
}

func (l *LeastLoadStrategy) getNodes(candidates []string, maxRTT time.Duration) []*node {
	if l.observer == nil {
		common.Must(core.RequireFeatures(l.ctx, func(observatory extension.Observatory) error {
			l.observer = observatory
			return nil
		}))
	}

	var ret []*node

	leastloadSort(ret)
	return ret
}

func leastloadSort(nodes []*node) {
	sort.Slice(nodes, func(i, j int) bool {
		left := nodes[i]
		right := nodes[j]
		if left.RTTDeviationCost != right.RTTDeviationCost {
			return left.RTTDeviationCost < right.RTTDeviationCost
		}
		if left.RTTAverage != right.RTTAverage {
			return left.RTTAverage < right.RTTAverage
		}
		if left.CountFail != right.CountFail {
			return left.CountFail < right.CountFail
		}
		if left.CountAll != right.CountAll {
			return left.CountAll > right.CountAll
		}
		return left.Tag < right.Tag
	})
}

func init() {
	common.Must(common.RegisterConfig((*StrategyLeastLoadConfig)(nil), nil))
}
