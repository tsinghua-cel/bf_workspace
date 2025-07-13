package library

import (
	"context"
	aiattack "github.com/tsinghua-cel/attacker-service/library/ai"
	"github.com/tsinghua-cel/attacker-service/library/confuse"
	"github.com/tsinghua-cel/attacker-service/library/exante"
	"github.com/tsinghua-cel/attacker-service/library/ext_exante"
	"github.com/tsinghua-cel/attacker-service/library/ext_sandwich"
	"github.com/tsinghua-cel/attacker-service/library/ext_staircase"
	"github.com/tsinghua-cel/attacker-service/library/ext_unrealized"
	"github.com/tsinghua-cel/attacker-service/library/ext_withholding"
	"github.com/tsinghua-cel/attacker-service/library/five"
	"github.com/tsinghua-cel/attacker-service/library/none"
	"github.com/tsinghua-cel/attacker-service/library/one"
	"github.com/tsinghua-cel/attacker-service/library/randomdelay"
	"github.com/tsinghua-cel/attacker-service/library/replay"
	"github.com/tsinghua-cel/attacker-service/library/sandwich"
	"github.com/tsinghua-cel/attacker-service/library/sandwichhonest"
	"github.com/tsinghua-cel/attacker-service/library/selfish"
	"github.com/tsinghua-cel/attacker-service/library/selfishhonest"
	"github.com/tsinghua-cel/attacker-service/library/simple"
	"github.com/tsinghua-cel/attacker-service/library/staircase"
	"github.com/tsinghua-cel/attacker-service/library/staircaseii"
	"github.com/tsinghua-cel/attacker-service/library/syncwrong"
	"github.com/tsinghua-cel/attacker-service/library/three"
	"github.com/tsinghua-cel/attacker-service/library/two"
	"github.com/tsinghua-cel/attacker-service/library/unrealized"
	"github.com/tsinghua-cel/attacker-service/library/unrealizedhonest"
	"github.com/tsinghua-cel/attacker-service/library/withholding"
	"github.com/tsinghua-cel/attacker-service/types"
	"sync"
)

type Strategy interface {
	Name() string
	Run(ctx context.Context, param types.LibraryParams, feedbacker types.FeedBacker)
	Description() string
}

var (
	allStrategies = sync.Map{}
)

func Init() {
	register(&one.Instance{})
	register(&two.Instance{})
	register(&three.Instance{})
	register(&syncwrong.Instance{})
	register(&five.Instance{})
	register(&exante.Instance{})
	register(&sandwich.Instance{})
	register(&sandwichhonest.Instance{})
	register(&withholding.Instance{})
	register(&unrealized.Instance{})
	register(&unrealizedhonest.Instance{})
	register(&staircase.Instance{})
	register(&staircaseii.Instance{})
	register(&confuse.Instance{})
	register(&randomdelay.Instance{})

	register(&ext_sandwich.Instance{})
	register(&ext_exante.Instance{})
	register(&ext_staircase.Instance{})
	register(&ext_unrealized.Instance{})
	register(&ext_withholding.Instance{})

	register(&aiattack.Instance{})
	register(&replay.Instance{})
	register(&simple.Instance{})
	register(&none.Instance{})
	register(&selfish.Instance{})
	register(&selfishhonest.Instance{})
}

func register(ins Strategy) {
	allStrategies.Store(ins.Name(), ins)
}

func GetStrategy(name string) Strategy {
	if v, ok := allStrategies.Load(name); ok {
		return v.(Strategy)
	}
	return nil
}

func GetAllStrategies() map[string]Strategy {
	strategies := make(map[string]Strategy)
	allStrategies.Range(func(k, v interface{}) bool {
		strategies[k.(string)] = v.(Strategy)
		return true
	})
	return strategies
}
func GetStrategiesList() []Strategy {
	strategies := make([]Strategy, 0)
	allStrategies.Range(func(k, v interface{}) bool {
		strategies = append(strategies, v.(Strategy))
		return true
	})
	return strategies
}
