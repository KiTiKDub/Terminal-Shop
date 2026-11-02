package main

import (
	"strings"

	figure "github.com/common-nighthawk/go-figure"
)

func (m *model) setHeaderString(headerText string) {
	headerFigure := figure.NewFigure(headerText, "", true)
	m.header = headerFigure.String()
}

func (m *model) getHeaderLineCount() int {
	return strings.Count(m.header, "\n") + 1
}
