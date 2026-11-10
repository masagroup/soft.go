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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// lookupEncoding returns the encoding of an IANA name or alias, or else of a WHATWG label, as
// GB2312: nil when neither knows it.
func lookupEncoding(name string) encoding.Encoding {
	if e, err := ianaindex.IANA.Encoding(name); err == nil && e != nil {
		return e
	}
	e, _ := charset.Lookup(name)
	return e
}

// newDocumentReader returns r in UTF-8, the only encoding encoding/xml reads, and "UTF-16" when the
// first bytes of r are UTF-16 (XML 1.0 appendix F): a byte order mark, or a declaration without
// one. Another encoding is decoded as its declaration names it (charsetReader).
func newDocumentReader(r io.Reader) (io.Reader, string) {
	br := bufio.NewReader(r)
	head, _ := br.Peek(4)
	var e encoding.Encoding
	switch {
	case bytes.HasPrefix(head, []byte{0xFE, 0xFF}), bytes.Equal(head, []byte{0, '<', 0, '?'}):
		e = unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
	case bytes.HasPrefix(head, []byte{0xFF, 0xFE}), bytes.Equal(head, []byte{'<', 0, '?', 0}):
		e = unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	default:
		return br, ""
	}
	return transform.NewReader(br, e.NewDecoder()), "UTF-16"
}

// charsetReader is the CharsetReader of encoding/xml for a document in detected, the encoding
// newDocumentReader found: a document already decoded from UTF-16 must declare UTF-16.
func charsetReader(detected string) func(label string, input io.Reader) (io.Reader, error) {
	return func(label string, input io.Reader) (io.Reader, error) {
		if detected != "" {
			if strings.HasPrefix(strings.ToUpper(label), "UTF-16") {
				return input, nil
			}
			return nil, fmt.Errorf("encoding %q of a document in %s", label, detected)
		}
		if e := lookupEncoding(label); e != nil {
			return e.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unsupported encoding %q", label)
	}
}

// writeEncoded writes the segments of str to w in the encoding name, a character the encoding
// cannot represent as a character reference, &#8364;, as EMF does.
func writeEncoded(w io.Writer, name string, str *xmlString) error {
	e := lookupEncoding(name)
	if e == nil {
		return fmt.Errorf("unsupported encoding %q", name)
	}
	if e == unicode.UTF8 {
		return str.write(w)
	}
	encoded := transform.NewWriter(w, encoding.HTMLEscapeUnsupported(e.NewEncoder()))
	if err := str.write(encoded); err != nil {
		return err
	}
	return encoded.Close()
}

// contentEscaper and attributeEscaper escape only what XML requires, as EMF does: '&', '<' and '>',
// a '"' in an attribute, and the whitespaces a parser would normalize.
var (
	contentEscaper   = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\r", "&#13;")
	attributeEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;",
		"\r", "&#13;", "\n", "&#xA;", "\t", "&#x9;")
)

// escapeAttribute escapes s as the value of an attribute
func escapeAttribute(s string) string {
	return attributeEscaper.Replace(s)
}

// escapeContent escapes s as the content of an element
func escapeContent(s string) string {
	return contentEscaper.Replace(s)
}
