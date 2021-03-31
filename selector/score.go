package selector

type Score struct {
	Selector `json:"selector"`
	Scored   interface{} `json:"scored"`
	Value    interface{} `json:"value"`
}
