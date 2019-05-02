// Copyright (c) 2019 srfrog - https://srfrog.me
// Use of this source code is governed by the license in the LICENSE file.

// Package dict is a Go implementation of Python dict, which are hashable object maps [1].
// Dictionaries complement Go map and slice types to provide a simple interface to
// store and access key-value data with relatively fast performance at the cost of extra
// memory. This is done by using the features of both maps and slices.
//
// A dict object is made of a key and value parts. A keys field is a slice that holds the order of
// values entered, which is the insertion order by default. Each key value contains the hash,
// index, or key ID, and the name of the key. The key value is used to find the value matching
// the key name in the values map, using the hash index. The values field in a dict object
// holds the values for a given key name, indexed by the hash index value.
//
// The key names must be a supported hashable types. The hashable types are int, uint, float,
// string, and types that implement fmt.Stringer. The key ID is made using string values.
// The values stored in a dict can be any Go type, including other dict objects.
//
// The func New() creates a new dict. It can take values to initialize the object. These can
// be slices, maps, channels and scalar values to create the dict. When using maps, the map
// keys must be hashable types that will be used as dict key IDs.
//
// 1- https://docs.python.org/3.7/library/stdtypes.html#dict
package dict

// Version is the version of this package.
const Version = "0.0.2"
