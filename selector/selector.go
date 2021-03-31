package selector

type Selector interface {
	Name() string
	Score(value interface{}) []Score
}

func NewBaseSelector(name string, f func(v interface{}) []Score) Selector {
	return &BaseSelector{
		name: name,
		f:    f,
	}
}

type BaseSelector struct {
	name string
	f    func(interface{}) []Score
}

func (b BaseSelector) Name() string {
	return b.name
}

func (b BaseSelector) Score(v interface{}) []Score {
	return b.f(v)
}
