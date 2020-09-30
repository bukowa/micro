package html

/*
Copyright © 2020 Mateusz Kurowski

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

// GoQueryCollectFunc is a function that gathers items from given goquery.Document.
type GoQueryCollectFunc = func(doc *goquery.Document) []string

// GoQueryCollector is a basic Collector implementation based on 'github.com/PuerkitoBio/goquery'.
type GoQueryCollector struct {
	Collector
	CollectFunc []GoQueryCollectFunc
}

// NewGoQueryCollector returns new GoQueryCollector.
func NewGoQueryCollector(collect ...GoQueryCollectFunc) Collector {
	return &GoQueryCollector{
		CollectFunc: collect,
	}
}

func (c *GoQueryCollector) Collect(r io.Reader, w io.Writer) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	for _, f := range c.CollectFunc {
		for _, v := range f(doc) {
			if _, err := w.Write([]byte(v)); err != nil {
				return err
			}
		}
	}
	return nil
}

// CollectSelectorAttribute gathers attributes for given selectors.
// Example:
//   var selectors = map[string][]string{"a": {"href"}}
var CollectSelectorAttributes = func(selectors map[string][]string) GoQueryCollectFunc {
	return func(doc *goquery.Document) (values []string) {
		for selector, attributes := range selectors {
			doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
				for _, attribute := range attributes {
					if value, exists := selection.Attr(attribute); exists {
						values = append(values, value)
					}
				}
			})
		}
		return values
	}
}
