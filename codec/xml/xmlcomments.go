// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package xml

import (
	"strconv"
	"strings"

	"github.com/masagroup/soft.go/ecore"
)

// XMLComments holds the comments one document carried, each pinned to its position among its
// parent element's children. No EClass can hold a comment, so they travel beside the graph: the
// decoder collects them on the resource and the encoder writes them back under
// XML_OPTION_KEEP_COMMENTS. A comment follows its position, not the object it annotates.
type XMLComments struct {
	byAnchor map[string][]string
	count    int
}

// Add records one comment body at an anchor. Two comments at the same anchor keep the order they
// were added, which is the order they appeared in the document.
func (c *XMLComments) Add(anchor, body string) {
	if c.byAnchor == nil {
		c.byAnchor = map[string][]string{}
	}
	c.byAnchor[anchor] = append(c.byAnchor[anchor], body)
	c.count++
}

// At returns the comment bodies recorded at an anchor, in document order.
func (c *XMLComments) At(anchor string) []string {
	if c == nil {
		return nil
	}
	return c.byAnchor[anchor]
}

// Len is how many comments the document held.
func (c *XMLComments) Len() int {
	if c == nil {
		return 0
	}
	return c.count
}

// codecDataCarrier is a resource that keeps data of its document its objects cannot hold
// (ecore.EResourceImpl). A resource without it cannot keep comments, and the codec then behaves as
// if the document had none.
type codecDataCarrier interface {
	GetCodecData(key any) any
	SetCodecData(key, value any)
}

// commentsKey is the key of a document's comments in its resource's codec data.
type commentsKey struct{}

func commentsOf(resource ecore.EResource) *XMLComments {
	carrier, _ := resource.(codecDataCarrier)
	if carrier == nil {
		return nil
	}
	comments, _ := carrier.GetCodecData(commentsKey{}).(*XMLComments)
	return comments
}

func setCommentsOf(resource ecore.EResource, comments *XMLComments) {
	if carrier, _ := resource.(codecDataCarrier); carrier != nil {
		carrier.SetCodecData(commentsKey{}, comments)
	}
}

// xmlPathStep is one open element on the path from the document to the current position.
type xmlPathStep struct {
	name string
	// index is this element's position among its parent's element children, 0-based.
	index int
	// children counts the element children opened inside it so far, and so IS the position a
	// comment written here would occupy.
	children int
	// hidden marks a step that exists in the encoder's bookkeeping but writes no tag - the
	// suppressed document root, which xmlString pushes under the empty name. It takes no
	// position among its parent's children and contributes nothing to an anchor.
	hidden bool
}

// xmlPath tracks the open-element path and the sibling counts that address a comment. Decoder and
// encoder build anchors with the same code.
type xmlPath struct {
	steps []xmlPathStep
}

func newXMLPath() *xmlPath {
	// The document itself is a step, so a comment outside the root element has an anchor too.
	return &xmlPath{steps: []xmlPathStep{{name: "#document"}}}
}

// push opens an element. An empty name is a hidden step (see xmlPathStep.hidden). A name may
// arrive qualified ("xsi:type"); only the local part addresses a position, because that is all
// the decoder sees.
func (p *xmlPath) push(name string) {
	if name == "" {
		p.steps = append(p.steps, xmlPathStep{hidden: true})
		return
	}
	if i := strings.LastIndexByte(name, ':'); i >= 0 {
		name = name[i+1:]
	}
	top := &p.steps[len(p.steps)-1]
	index := top.children
	top.children++
	p.steps = append(p.steps, xmlPathStep{name: name, index: index})
}

// pop closes the innermost element. It never pops the document step, so an unbalanced stream
// leaves a usable path rather than panicking.
func (p *xmlPath) pop() {
	if len(p.steps) > 1 {
		p.steps = p.steps[:len(p.steps)-1]
	}
}

// child records an element that opens and closes in one go (a nil or a simple-content element,
// which the encoder writes without a push/pop pair). It only advances the sibling count.
func (p *xmlPath) child() {
	p.steps[len(p.steps)-1].children++
}

// here is the anchor of the current position: the path of open elements, each with its index among
// its siblings, and the number of element children already written inside the innermost one. Two
// positions in a document have the same anchor only if they are the same position.
func (p *xmlPath) here() string {
	var b strings.Builder
	for _, step := range p.steps[1:] {
		if step.hidden {
			continue
		}
		b.WriteByte('/')
		b.WriteString(step.name)
		b.WriteByte('#')
		b.WriteString(strconv.Itoa(step.index))
	}
	b.WriteByte('@')
	b.WriteString(strconv.Itoa(p.steps[len(p.steps)-1].children))
	return b.String()
}
