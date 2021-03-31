package selector

type Registry interface {
	Register(key interface{}, sel Selector) error
	Scored(v interface{}) map[interface{}][]Score
}

type BaseRegistry map[interface{}][]Selector

func (b BaseRegistry) Register(key interface{}, sel ...Selector) error {
	if v, ok := b[key]; ok {
		b[key] = append(v, sel...)
		return nil
	}
	b[key] = sel
	return nil
}

func (b BaseRegistry) Scored(v interface{}) map[interface{}][]Score {
	var m = make(map[interface{}][]Score)
	for k, selectors := range b {
		for _, sel := range selectors {
			if score := sel.Score(v); score != nil {
				if v2, ok := m[k]; ok {
					m[k] = append(v2, score...)
					continue
				}
				m[k] = score
			}
		}
	}
	return m
}
