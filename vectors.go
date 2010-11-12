// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The vector package implements containers for managing sequences
// of elements. Vectors grow and shrink dynamically as necessary.
package main


type PersonVector []*Person

type PlaceVector []*Place


// Initial underlying array size
const initialSize = 8


// Partial sort.Interface support

// LessInterface provides partial support of the sort.Interface.
type LessInterface interface {
	Less(y interface{}) bool
}


// Less returns a boolean denoting whether the i'th element is less than the j'th element.
func (p *PlaceVector) Less(i, j int) bool  { return (*p)[i].Id < (*p)[j].Id }
func (p *PersonVector) Less(i, j int) bool { return (*p)[i].Name < (*p)[j].Name }
