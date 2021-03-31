package selector

type XStringSelector struct {
	XName     string
	ScoreFunc func(value string, match string) bool
	Match     []string
}

func (sel XStringSelector) New(match ...string) Selector {
	return XStringSelector{
		ScoreFunc: sel.ScoreFunc,
		XName:     sel.Name(),
		Match:     match,
	}
}

func (sel XStringSelector) Name() string {
	return sel.XName
}

func (sel XStringSelector) Score(value interface{}) []Score {
	var m = make([]Score, 0)
	switch v := value.(type) {
	case string:
		for _, match := range sel.Match {
			if ok := sel.ScoreFunc(v, match); ok {
				m = append(m, Score{
					Selector: sel,
					Scored:   match,
					Value:    v,
				})
			}
			continue
		}
	default:
		return m
	}
	return m
}
