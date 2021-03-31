package selector

import "testing"

func compareScore(t *testing.T, s1, s2 Score) {
	if s1.Value != s2.Value {
		t.Error()
	}
	if s1.Scored != s2.Scored {
		t.Error(s1.Scored, s2.Scored)
	}
	if s1.Name() != s2.Name() {
		t.Error()
	}
}

func TestSelectors(t *testing.T) {
	var selectors = map[string]Selector{
		"prefix":     StringPrefix.New("y", "ye", "y"),
		"suffix":     StringSuffix.New("s", "es", "x"),
		"equal":      StringEqual.New("s", "es", "yes"),
		"equal_fold": StringEqualFold.New("YeS", "no", "x", "yes"),
		"contains":   StringContains.New("e", "s", "y", "x", "y"),
	}

	r, err := RegexpMatch.New("[yes]", "[y]")
	if err != nil {
		t.Error(err)
	}
	selectors["regexp"] = r

	type want struct {
		scored map[interface{}][]Score
	}
	var tests = map[string]want{
		"prefix": {
			scored: map[interface{}][]Score{
				"prefix": {
					{
						Selector: selectors["prefix"],
						Scored:   "y",
						Value:    "yes",
					},
					{
						Selector: selectors["prefix"],
						Scored:   "ye",
						Value:    "yes",
					},
				},
			},
		},
		"suffix": {
			scored: map[interface{}][]Score{
				"suffix": {
					{
						Selector: selectors["suffix"],
						Scored:   "s",
						Value:    "yes",
					},
					{
						Selector: selectors["suffix"],
						Scored:   "es",
						Value:    "yes",
					},
				},
			},
		},
		"equal": {
			scored: map[interface{}][]Score{
				"equal": {
					{
						Selector: selectors["equal"],
						Scored:   "yes",
						Value:    "yes",
					},
				},
			},
		},
		"equal_fold": {
			scored: map[interface{}][]Score{
				"equal_fold": {
					{
						Selector: selectors["equal_fold"],
						Scored:   "YeS",
						Value:    "yes",
					},
					{
						Selector: selectors["equal_fold"],
						Scored:   "yes",
						Value:    "yes",
					},
				},
			},
		},
		"contains": {
			scored: map[interface{}][]Score{
				"contains": {
					{
						Selector: selectors["contains"],
						Scored:   "e",
						Value:    "yes",
					},
					{
						Selector: selectors["contains"],
						Scored:   "s",
						Value:    "yes",
					},
					{
						Selector: selectors["contains"],
						Scored:   "y",
						Value:    "yes",
					},
				},
			},
		},
		"regexp": {
			scored: map[interface{}][]Score{
				"regexp": {
					{
						Selector: selectors["regexp"],
						Scored:   "[yes]",
						Value:    "yes",
					},
					{
						Selector: selectors["regexp"],
						Scored:   "[y]",
						Value:    "yes",
					},
				},
			},
		},
	}

	var registry = make(BaseRegistry, len(selectors))
	for k, v := range selectors {
		registry.Register(k, v)
	}
	scored := registry.Scored("yes")

	for k, _ := range selectors {
		for _, v := range tests[k].scored {
			for i, s := range v {
				compareScore(t, scored[k][i], s)
			}
		}
	}
}
