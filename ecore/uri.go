// *****************************************************************************
// Copyright(c) 2021 MASA Group
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
//
// *****************************************************************************

package ecore

import (
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
)

type URI struct {
	scheme   string
	username string
	password string
	host     string
	port     int
	path     string
	query    string
	fragment string
	rawURI   string
}

type URIParam interface {
	apply(*URI)
}

// funcURIParam wraps a function that modifies uri params into an
// implementation of the URIParam interface.
type funcURIParam struct {
	f func(*URI)
}

func (fup *funcURIParam) apply(uri *URI) {
	fup.f(uri)
}

func newFuncURIParam(f func(*URI)) *funcURIParam {
	return &funcURIParam{
		f: f,
	}
}

func Scheme(scheme string) URIParam {
	return newFuncURIParam(func(u *URI) {
		u.scheme = scheme
	})
}

func Username(username string) URIParam {
	return newFuncURIParam(func(u *URI) {
		u.username = username
	})
}

func Password(password string) URIParam {
	return newFuncURIParam(func(u *URI) {
		u.password = password
	})
}

func Host(host string) URIParam {
	return newFuncURIParam(func(u *URI) {
		u.host = host
	})
}

func Port(port int) URIParam {
	return newFuncURIParam(func(u *URI) {
		u.port = port
	})
}

func Path(path string) URIParam {
	return newFuncURIParam(func(u *URI) {
		u.host = path
	})
}

func Query(query string) URIParam {
	return newFuncURIParam(func(u *URI) {
		u.query = query
	})
}

func Fragment(fragment string) URIParam {
	return newFuncURIParam(func(u *URI) {
		u.fragment = fragment
	})
}

func NewURI(rawURI string) *URI {
	uri, _ := ParseURI(rawURI)
	return uri
}

func BuildURI(uriParams ...URIParam) *URI {
	uri := &URI{}
	for _, p := range uriParams {
		p.apply(uri)
	}
	return uri
}

func ParseURI(rawURI string) (*URI, error) {
	if url, err := url.Parse(rawURI); err != nil {
		return nil, err
	} else {
		uri := &URI{scheme: url.Scheme, path: url.Path, fragment: url.Fragment, query: url.RawQuery, rawURI: rawURI}
		if url.User != nil {
			uri.username = url.User.Username()
			uri.password, _ = url.User.Password()
		}
		if i := strings.IndexByte(url.Host, ':'); i >= 0 {
			uri.host = url.Host[0:i]
			uri.port, err = strconv.Atoi(url.Host[i+1:])
			if err != nil {
				return nil, err
			}
		} else {
			uri.host = url.Host
		}
		if len(url.Opaque) > 0 {
			uri.path = url.Opaque
		}
		return uri, nil
	}
}

func (uri *URI) Scheme() string {
	return uri.scheme
}

func (uri *URI) Host() string {
	return uri.host
}

func (uri *URI) Username() string {
	return uri.username
}

func (uri *URI) Password() string {
	return uri.username
}

func (uri *URI) Port() int {
	return uri.port
}

func (uri *URI) Path() string {
	return uri.path
}

func (uri *URI) Query() string {
	return uri.query
}

func (uri *URI) Fragment() string {
	return uri.query
}

func (uri *URI) IsAbsolute() bool {
	return len(uri.scheme) != 0
}

func (uri *URI) IsOpaque() bool {
	return len(uri.path) == 0
}

var emptyURI = &URI{}

func (uri *URI) IsEmpty() bool {
	return *uri == *emptyURI
}

func (uri *URI) Authority() string {
	var s strings.Builder
	if len(uri.username) > 0 || len(uri.password) > 0 {
		s.WriteString(uri.username)
		if len(uri.password) > 0 {
			s.WriteRune(':')
			s.WriteString(uri.password)
		}
		s.WriteRune('@')
	}
	s.WriteString(uri.host)
	if uri.port != 0 {
		s.WriteRune(':')
		s.WriteString(strconv.Itoa(uri.port))
	}
	return s.String()
}

func (uri *URI) String() string {
	if len(uri.rawURI) == 0 {
		var s strings.Builder
		if len(uri.scheme) > 0 {
			s.WriteString(uri.scheme)
			s.WriteRune(':')
		}
		if autority := uri.Authority(); len(autority) > 0 {
			s.WriteString("//")
			s.WriteString(autority)
		}
		s.WriteString(uri.path)
		if len(uri.query) > 0 {
			s.WriteRune('?')
			s.WriteString(uri.query)
		}
		if len(uri.fragment) > 0 {
			s.WriteRune('#')
			s.WriteString(uri.fragment)
		}
		uri.rawURI = s.String()
	}
	return uri.rawURI
}

func (uri *URI) Normalize() *URI {
	if uri.IsOpaque() {
		return uri.Copy()
	}

	np := normalize(uri.path)
	if np == uri.path {
		return uri.Copy()
	}

	return &URI{
		scheme:   uri.scheme,
		username: uri.username,
		password: uri.password,
		host:     uri.host,
		port:     uri.port,
		path:     np,
		query:    uri.query,
		fragment: uri.fragment,
	}
}

func (uri *URI) Resolve(ref *URI) *URI {
	return resolve(uri, ref)
}

func (uri *URI) Relativize(ref *URI) *URI {
	return relativize(uri, ref)
}

func (uri *URI) Copy() *URI {
	return &URI{
		scheme:   uri.scheme,
		username: uri.username,
		password: uri.password,
		host:     uri.host,
		port:     uri.port,
		path:     uri.path,
		query:    uri.query,
		fragment: uri.fragment,
		rawURI:   uri.rawURI,
	}
}

func (uri *URI) TrimFragment() *URI {
	return &URI{
		scheme:   uri.scheme,
		username: uri.username,
		password: uri.password,
		host:     uri.host,
		port:     uri.port,
		path:     uri.path,
		query:    uri.query,
	}
}

func (uri *URI) TrimQuery() *URI {
	return &URI{
		scheme:   uri.scheme,
		username: uri.username,
		password: uri.password,
		host:     uri.host,
		port:     uri.port,
		path:     uri.path,
		fragment: uri.fragment,
	}
}

func (uri *URI) ReplacePrefix(oldPrefix *URI, newPrefix *URI) *URI {
	if uri.scheme != oldPrefix.scheme ||
		uri.username != oldPrefix.username ||
		uri.password != oldPrefix.password ||
		uri.host != oldPrefix.host ||
		uri.port != oldPrefix.port {
		return nil
	}
	if oldLen := len(oldPrefix.path); len(uri.path) >= oldLen && uri.path[0:oldLen] == oldPrefix.path {
		return &URI{scheme: newPrefix.scheme,
			username: uri.username,
			password: uri.password,
			host:     uri.host,
			port:     uri.port,
			path:     newPrefix.path + uri.path[oldLen:],
			query:    uri.query,
			fragment: uri.fragment}
	}
	return nil

}

func normalize(path string) string {
	runes := []rune(path)
	// Does this path need normalization?
	ns := needsNormalization(runes) // Number of segments
	if ns < 0 {
		// Nope -- just return it
		return path
	}

	segs := make([]int, ns)
	split(runes, segs)

	// Remove dots
	removeDots(runes, segs)

	// Prevent scheme-name confusion
	maybeAddLeadingDot(runes, segs)

	// Join the remaining segments and return the result
	newSize := join(runes, segs)
	runes = runes[0:newSize]
	return string(runes)
}

// The following algorithm for path normalization avoids the creation of a
// string object for each segment, as well as the use of a string buffer to
// compute the final result, by using a single char array and editing it in
// place.  The array is first split into segments, replacing each slash
// with '\0' and creating a segment-index array, each element of which is
// the index of the first char in the corresponding segment.  We then walk
// through both arrays, removing ".", "..", and other segments as necessary
// by setting their entries in the index array to -1.  Finally, the two
// arrays are used to rejoin the segments and compute the final result.
//
// This code is based upon src/solaris/native/java/io/canonicalize_md.c

// Check the given path to see if it might need normalization.  A path
// might need normalization if it contains duplicate slashes, a "."
// segment, or a ".." segment.  Return -1 if no further normalization is
// possible, otherwise return the number of segments found.
//
// This method takes a string argument rather than a char array so that
// this test can be performed without invoking path.toCharArray().
//
func needsNormalization(path []rune) int {
	normal := true
	ns := 0              // Number of segments
	end := len(path) - 1 // Index of last char in path
	p := 0               // Index of next char in path

	// Skip initial slashes
	for p <= end {
		if path[p] != '/' {
			break
		}
		p++
	}

	if p > 1 {
		normal = false
	}

	// Scan segments
	for p <= end {
		// Looking at "." or ".." ?
		if (path[p] == '.') && ((p == end) || ((path[p+1] == '/') || ((path[p+1] == '.') && ((p+1 == end) || (path[p+2] == '/'))))) {
			normal = false
		}
		ns++

		// Find beginning of next segment
		for p <= end {
			c := path[p]
			p++
			if c != '/' {
				continue
			}

			// Skip redundant slashes
			for p <= end {
				if path[p] != '/' {
					break
				}
				normal = false
				p++
			}
			break
		}
	}
	if normal {
		return -1
	} else {
		return ns
	}
}

// Split the given path into segments, replacing slashes with nulls and
// filling in the given segment-index array.
//
// Preconditions:
//   segs.length == Number of segments in path
//

// Postconditions:
//   All slashes in path replaced by '\0'
//   segs[i] == Index of first char in segment i (0 <= i < segs.length)
//
func split(path []rune, segs []int) {
	end := len(path) - 1 // Index of last char in path
	p := 0               // Index of next char in path
	i := 0               // Index of current segment

	// Skip initial slashes
	for p <= end {
		if path[p] != '/' {
			break
		}
		path[p] = rune(0)
		p++
	}

	for p <= end {
		// Note start of segment
		segs[i] = p
		p++
		i++
		// Find beginning of next segment
		for p <= end {
			c := path[p]
			p++
			if c != '/' {
				continue
			}

			path[p-1] = rune(0)

			// Skip redundant slashes
			for p <= end {
				if path[p] != '/' {
					break
				}
				path[p] = rune(0)
				p++
			}
			break
		}
	}
}

// Remove "." segments from the given path, and remove segment pairs
// consisting of a non-".." segment followed by a ".." segment.
//
func removeDots(path []rune, segs []int) {
	ns := len(segs)
	end := len(path) - 1
	for i := 0; i < ns; i++ {
		dots := 0 // Number of dots found (0, 1, or 2)

		// Find next occurrence of "." or ".."
		for ok := true; ok; ok = i < ns {
			p := segs[i]
			if path[p] == '.' {
				if p == end {
					dots = 1
					break
				} else if path[p+1] == rune(0) {
					dots = 1
					break
				} else if (path[p+1] == '.') && ((p+1 == end) || (path[p+2] == rune(0))) {
					dots = 2
					break
				}
			}
			i++
		}

		if (i > ns) || (dots == 0) {
			break
		}

		if dots == 1 {
			// Remove this occurrence of "."
			segs[i] = -1
		} else {
			// If there is a preceding non-".." segment, remove both that
			// segment and this occurrence of ".."; otherwise, leave this
			// ".." segment as-is.
			var j int
			for j = i - 1; j >= 0; j-- {
				if segs[j] != -1 {
					break
				}
			}

			if j >= 0 {
				q := segs[j]
				if !((path[q] == '.') && (path[q+1] == '.') && (path[q+2] == rune(0))) {
					segs[i] = -1
					segs[j] = -1
				}
			}
		}
	}
}

// DEVIATION: If the normalized path is relative, and if the first
// segment could be parsed as a scheme name, then prepend a "." segment
//
func maybeAddLeadingDot(path []rune, segs []int) {
	if path[0] == rune(0) {
		// The path is absolute
		return
	}

	ns := len(segs)
	f := 0 // Index of first segment
	for f < ns {
		if segs[f] >= 0 {
			break
		}
		f++
	}

	if f >= ns || f == 0 {
		// The path is empty, or else the original first segment survived,
		// in which case we already know that no leading "." is needed
		return
	}
	p := segs[f]
	for p < len(path) && path[p] != ':' && path[p] != rune(0) {
		p++
	}

	if p >= len(path) || path[p] == rune(0) {
		// No colon in first segment, so no "." needed
		return
	}

	// At this point we know that the first segment is unused,
	// hence we can insert a "." segment at that position
	path[0] = '.'
	path[1] = rune(0)
	segs[0] = 0
}

// Join the segments in the given path according to the given segment-index
// array, ignoring those segments whose index entries have been set to -1,
// and inserting slashes as needed.  Return the length of the resulting
// path.
//
// Preconditions:
//   segs[i] == -1 implies segment i is to be ignored
//   path computed by split, as above, with '\0' having replaced '/'
//
// Postconditions:
//   path[0] .. path[return value] == Resulting path
//
func join(path []rune, segs []int) int {
	ns := len(segs)      // Number of segments
	end := len(path) - 1 // Index of last char in path
	p := 0               // Index of next path char to write
	if path[p] == rune(0) {
		// Restore initial slash for absolute paths
		path[p] = '/'
		p++
	}
	for i := 0; i < ns; i++ {
		q := segs[i] // Current segment
		if q == -1 {
			// Ignore this segment
			continue
		}

		if p == q {
			// We're already at this segment, so just skip to its end
			for p <= end && path[p] != rune(0) {
				p++
			}

			if p <= end {
				// Preserve trailing slash
				path[p] = '/'
				p++
			}
		} else if p < q {
			// Copy q down to p
			for q <= end && path[q] != rune(0) {
				path[p] = path[q]
				p++
				q++
			}
			if q <= end {
				// Preserve trailing slash
				path[p] = '/'
				p++
			}
		}
	}
	return p
}

// If both URIs are hierarchical, their scheme and authority components are
// identical, and the base path is a prefix of the child's path, then
// return a relative URI that, when resolved against the base, yields the
// child; otherwise, return the child.
//
func relativize(base *URI, child *URI) *URI {
	// check if child if opaque
	if child.IsOpaque() || base.IsOpaque() {
		return child.Copy()
	}
	if base.scheme != child.scheme || base.Authority() != child.Authority() {
		return child.Copy()
	}

	bp := normalize(base.path)
	cp := normalize(child.path)
	if bp != cp {
		if !strings.HasSuffix(bp, "/") {
			i := strings.LastIndexByte(bp, '/')
			bp = bp[0 : i+1]
		}

		if !strings.HasPrefix(cp, bp) {
			return child.Copy()
		}
	}
	return &URI{
		path:     cp[len(bp):],
		query:    child.query,
		fragment: child.fragment,
	}
}

func resolve(base *URI, child *URI) *URI {
	// check if child if opaque first
	if child.IsOpaque() || base.IsOpaque() {
		return child
	}

	// Reference to current document (lone fragment)
	childAuthority := child.Authority()
	if len(child.scheme) == 0 && len(childAuthority) == 0 && len(child.path) == 0 && len(child.fragment) != 0 && len(child.query) == 0 {
		if len(base.fragment) == 0 && child.fragment == base.fragment {
			return base.Copy()
		}

		return &URI{
			scheme:   base.scheme,
			username: base.username,
			password: base.password,
			host:     base.host,
			port:     base.port,
			path:     base.path,
			query:    base.query,
			fragment: child.fragment,
		}
	}

	// Child is absolute
	if len(child.scheme) != 0 {
		return child.Copy()
	}

	ru := &URI{
		scheme:   base.scheme,
		query:    child.query,
		fragment: child.fragment,
	}

	// Authority
	if len(childAuthority) == 0 {
		ru.host = base.host
		ru.username = base.username
		ru.password = base.password
		ru.port = base.port
		if len(child.path) > 0 && child.path[0] == '/' {
			ru.path = child.path
		} else {
			ru.path = resolvePath(base.path, child.path, base.IsAbsolute())
		}
	} else {
		ru.host = child.host
		ru.username = child.username
		ru.password = child.password
		ru.port = child.port
		ru.path = child.path
	}
	// Recombine (nothing to do here)
	return ru
}

func resolvePath(base string, child string, isAbsolute bool) string {
	i := strings.LastIndexByte(base, '/')
	cn := len(child)
	var path string
	if cn == 0 {
		if i >= 0 {
			path = base[0 : i+1]
		}
	} else {
		if i >= 0 {
			path = path + base[0:i+1]
		}
		path = path + child
	}
	return normalize(path)
}

func CreateFileURI(path string) *URI {
	p := filepath.ToSlash(path)
	if len(p) == 0 {
		return &URI{}
	} else {
		if isAbsolute := filepath.IsAbs(p); isAbsolute {
			return &URI{scheme: "file", path: p}
		} else {
			return &URI{path: p}
		}
	}
}

func CreateMemoryURI(path string) *URI {
	p := filepath.ToSlash(path)
	if len(p) == 0 {
		return nil
	} else {
		return &URI{scheme: "memory", path: path}
	}
}
