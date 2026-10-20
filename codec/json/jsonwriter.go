// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package json

import (
	"bufio"
	"io"
	"unicode/utf8"
)

// jsonWriter writes compact JSON: it separates the elements of the open objects and arrays, and
// escapes strings as encoding/json does without HTML escaping. The first write error is kept in
// err and the following writes are ignored.
type jsonWriter struct {
	w        *bufio.Writer
	buffer   []byte
	first    []bool // per open object or array: no element written yet
	afterKey bool   // a key was written, its value follows without separator
	err      error
}

func newJSONWriter(w io.Writer) *jsonWriter {
	return &jsonWriter{w: bufio.NewWriter(w)}
}

func (w *jsonWriter) write(b []byte) {
	if w.err == nil {
		_, w.err = w.w.Write(b)
	}
}

func (w *jsonWriter) writeByte(b byte) {
	if w.err == nil {
		w.err = w.w.WriteByte(b)
	}
}

// separate writes the comma before an element that is not the first of its object or array.
func (w *jsonWriter) separate() {
	if w.afterKey {
		w.afterKey = false
		return
	}
	if n := len(w.first); n > 0 {
		if w.first[n-1] {
			w.first[n-1] = false
		} else {
			w.writeByte(',')
		}
	}
}

func (w *jsonWriter) beginObject() {
	w.separate()
	w.writeByte('{')
	w.first = append(w.first, true)
}

func (w *jsonWriter) endObject() {
	w.first = w.first[:len(w.first)-1]
	w.writeByte('}')
}

func (w *jsonWriter) beginArray() {
	w.separate()
	w.writeByte('[')
	w.first = append(w.first, true)
}

func (w *jsonWriter) endArray() {
	w.first = w.first[:len(w.first)-1]
	w.writeByte(']')
}

func (w *jsonWriter) key(name string) {
	w.separate()
	w.writeString(name)
	w.writeByte(':')
	w.afterKey = true
}

func (w *jsonWriter) stringValue(s string) {
	w.separate()
	w.writeString(s)
}

func (w *jsonWriter) writeString(s string) {
	w.buffer = appendString(w.buffer[:0], s)
	w.write(w.buffer)
}

// flush writes the buffered data and returns the first write error.
func (w *jsonWriter) flush() error {
	if w.err == nil {
		w.err = w.w.Flush()
	}
	return w.err
}

const hex = "0123456789abcdef"

// appendString appends s as a JSON string, escaped as encoding/json does with SetEscapeHTML(false):
// quote, backslash and control characters escaped, invalid UTF-8 replaced with U+FFFD, U+2028 and
// U+2029 escaped.
func appendString(dst []byte, s string) []byte {
	dst = append(dst, '"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if b >= 0x20 && b != '"' && b != '\\' {
				i++
				continue
			}
			dst = append(dst, s[start:i]...)
			switch b {
			case '\\', '"':
				dst = append(dst, '\\', b)
			case '\b':
				dst = append(dst, '\\', 'b')
			case '\f':
				dst = append(dst, '\\', 'f')
			case '\n':
				dst = append(dst, '\\', 'n')
			case '\r':
				dst = append(dst, '\\', 'r')
			case '\t':
				dst = append(dst, '\\', 't')
			default:
				dst = append(dst, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			dst = append(dst, s[start:i]...)
			dst = append(dst, `�`...)
			i += size
			start = i
			continue
		}
		if c == ' ' || c == ' ' {
			dst = append(dst, s[start:i]...)
			dst = append(dst, '\\', 'u', '2', '0', '2', hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	dst = append(dst, s[start:]...)
	return append(dst, '"')
}
