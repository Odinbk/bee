package core

import (
	"github.com/PuerkitoBio/goquery"
)

type Sensor interface {
}

type HTMLSensor struct {
}

func (s HTMLSensor) FindItemByCssSelector(doc goquery.Document, selector string) (count int, isFound bool) {
	isFound = false
	selection := doc.Find(selector)
	if selection != nil {
		isFound = true
	}

	return count, isFound
}
