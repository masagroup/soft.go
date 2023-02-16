package ecore

import (
	"io"
	"strings"
)

type xmlStringSegment struct {
	strings.Builder
	lineWidth int
}

type xmlString struct {
	currentSegment     *xmlStringSegment
	firstElementMark   *xmlStringSegment
	indentation        string
	segments           []*xmlStringSegment
	indents            []string
	elementNames       []string
	lineWidth          int
	depth              int
	lastElementIsStart bool
}

const MaxInt = int(^uint(0) >> 1)

func newXmlString() *xmlString {
	segment := &xmlStringSegment{}
	s := &xmlString{
		segments:           []*xmlStringSegment{segment},
		currentSegment:     segment,
		lineWidth:          MaxInt,
		depth:              0,
		indentation:        "    ",
		indents:            []string{""},
		lastElementIsStart: false,
	}
	return s
}

func (s *xmlString) write(w io.Writer) error {
	for _, segment := range s.segments {
		if _, err := w.Write([]byte(segment.String())); err != nil {
			return err
		}
	}
	return nil
}

func (s *xmlString) add(newString string) {
	if s.lineWidth != MaxInt {
		s.currentSegment.lineWidth += len(newString)
	}
	s.currentSegment.WriteString(newString)
}

func (s *xmlString) addLine() {
	s.add("\n")
	s.currentSegment.lineWidth = 0
}

func (s *xmlString) startElement(name string) {
	if s.lastElementIsStart {
		s.closeStartElement()
	}
	s.elementNames = append(s.elementNames, name)
	if len(name) > 0 {
		s.depth++
		s.add(s.getElementIndent())
		s.add("<")
		s.add(name)
		s.lastElementIsStart = true
		if s.firstElementMark == nil {
			s.firstElementMark = s.mark()
		}
	}
}

func (s *xmlString) closeStartElement() {
	s.add(">")
	s.addLine()
	s.lastElementIsStart = false
}

func (s *xmlString) endElement() {
	if s.lastElementIsStart {
		s.endEmptyElement()
	} else {
		name := s.removeLast()
		if name != "" {
			s.add(s.getElementIndentWithExtra(1))
			s.add("</")
			s.add(name)
			s.add(">")
			s.addLine()
		}
	}
}

func (s *xmlString) endEmptyElement() {
	s.removeLast()
	s.add("/>")
	s.addLine()
	s.lastElementIsStart = false
}

func (s *xmlString) removeLast() string {
	end := len(s.elementNames) - 1
	result := s.elementNames[end]
	s.elementNames = s.elementNames[:end]
	if result != "" {
		s.depth--
	}
	return result
}

func (s *xmlString) addAttribute(name string, value string) {
	s.startAttribute(name)
	s.addAttributeContent(value)
	s.endAttribute()
}

func (s *xmlString) startAttribute(name string) {
	if s.currentSegment.lineWidth > s.lineWidth {
		s.addLine()
		s.add(s.getAttributeIndent())
	} else {
		s.add(" ")
	}
	s.add(name)
	s.add("=\"")
}

func (s *xmlString) addAttributeContent(content string) {
	s.add(content)
}

func (s *xmlString) endAttribute() {
	s.add("\"")
}

func (s *xmlString) addNil(name string) {
	if s.lastElementIsStart {
		s.closeStartElement()
	}

	s.depth++
	s.add(s.getElementIndent())
	s.add("<")
	s.add(name)
	if s.currentSegment.lineWidth > s.lineWidth {
		s.addLine()
		s.add(s.getAttributeIndent())
	} else {
		s.add(" ")
	}
	s.add("xsi:nil=\"true\"/>")
	s.depth--
	s.addLine()
	s.lastElementIsStart = false
}

func (s *xmlString) addContent(name string, content string) {
	if s.lastElementIsStart {
		s.closeStartElement()
	}
	s.depth++
	s.add(s.getElementIndent())
	s.add("<")
	s.add(name)
	s.add(">")
	s.add(content)
	s.add("</")
	s.depth--
	s.add(name)
	s.add(">")
	s.addLine()
	s.lastElementIsStart = false
}

func (s *xmlString) getElementIndent() string {
	return s.getElementIndentWithExtra(0)
}

func (s *xmlString) getElementIndentWithExtra(extra int) string {
	nesting := s.depth + extra - 1
	for i := len(s.indents) - 1; i < nesting; i++ {
		s.indents = append(s.indents, s.indents[i]+"  ")
	}
	return s.indents[nesting]
}

func (s *xmlString) getAttributeIndent() string {
	nesting := s.depth + 1
	for i := len(s.indents) - 1; i < nesting; i++ {
		s.indents = append(s.indents, s.indents[i]+"  ")
	}
	return s.indents[nesting]
}

func (s *xmlString) mark() *xmlStringSegment {
	r := s.currentSegment
	s.currentSegment = &xmlStringSegment{}
	s.segments = append(s.segments, s.currentSegment)
	return r
}

func (s *xmlString) resetToFirstElementMark() {
	s.resetToMark(s.firstElementMark)
}

func (s *xmlString) resetToMark(segment *xmlStringSegment) {
	if segment != nil {
		s.currentSegment = segment
	}
}
