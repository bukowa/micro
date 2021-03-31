package selector

import "regexp"

type XRegexpSelector struct {
	XName string
	Regxp []*regexp.Regexp
}

func (reSel XRegexpSelector) Name() string {
	return reSel.XName
}

func (reSel XRegexpSelector) Score(value interface{}) []Score {
	var m = make([]Score, 0)
	for _, r := range reSel.Regxp {
		switch v := value.(type) {
		case string:
			i := r.FindStringIndex(v)
			if len(i) == 2 {
				m = append(m, Score{
					Selector: reSel,
					Scored:   r.String(),
					Value:    v,
				})
			}
		case []byte:
			i := r.FindIndex(v)
			if len(i) == 2 {
				m = append(m, Score{
					Selector: reSel,
					Scored:   r.String(),
					Value:    v,
				})
			}
		default:
			return m
		}
	}
	return m
}

func (reSel XRegexpSelector) New(match ...string) (Selector, error) {
	var regC = make([]*regexp.Regexp, 0)
	for _, v := range match {
		r, err := regexp.Compile(v)
		if err != nil {
			return nil, err
		}
		regC = append(regC, r)
	}
	reSel2 := XRegexpSelector{
		XName: reSel.XName,
		Regxp: regC,
	}
	return reSel2, nil
}
