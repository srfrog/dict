# Dict
[![GoDoc](https://godoc.org/github.com/srfrog/dict?status.svg)](https://godoc.org/github.com/srfrog/dict)
[![Go Report Card](https://goreportcard.com/badge/github.com/srfrog/dict?svg=1)](https://goreportcard.com/report/github.com/srfrog/dict)
[![Coverage Status](https://coveralls.io/repos/github/srfrog/dict/badge.svg?branch=master)](https://coveralls.io/github/srfrog/dict?branch=master)
[![Build Status](https://travis-ci.com/srfrog/dict.svg?branch=master)](https://travis-ci.com/srfrog/dict)

*Python dictionary data type (dict) in Go*

Package dict is a Go implementation of Python [dict][1], which are hashable object maps.
Dictionaries complement Go map and slice types to provide a simple interface to
store and access key-value data with relatively fast performance at the cost of extra
memory. This is done by using the features of both maps and slices.

## Quick Start

Install using "go get":

	go get github.com/srfrog/dict

Then import from your source:

	import "github.com/srfrog/dict"

View [example_test.go][2] for an extended example of basic usage and features.

## Features

- [x] Initialize a new dict with scalars, slices, maps, channels and other dictionaries.
- [x] Go types int, uint, float, string and fmt.Stringer are hashable for dict keys.
- [x] Go map keys are used for dict keys if they are hashable.
- [x] Dict items are sorted in their insertion order, unlike Go maps.
- [ ] Go routine safe with minimal mutex locking (WIP)
- [x] Builtin JSON support for marshalling and unmarshalling
- [ ] sql.Scanner support via optional sub-package (WIP)
- [ ] Plenty of tests and examples to get you started quickly (WIP)

## Documentation

The full code documentation is located at GoDoc:

[http://godoc.org/github.com/srfrog/dict](http://godoc.org/github.com/srfrog/dict)

The source code is thoroughly commented, have a look.

## Usage

Minimal example showing basic usage:

```go
package main

import (
   "github.com/srfrog/dict"
)

type Car struct {
   Model, BrandID string
}

func main() {
   // Map of car models, indexed by VIN.
   // Data source: NHTSA.gov
   vins := map[string]*Car{
      "2C3KA43R08H129584": &Car{
         Model:   "2008 CHRYSLER 300",
         BrandID: "ACB9976A-DB5F-4D57-B9A8-9F5C53D87C7C",
      },
      "1N6AD07U78C416152": &Car{
         Model:   "2008 NISSAN FRONTIER SE-V6 RWD",
         BrandID: "003096EE-C8FC-4C2F-ADEF-406F86C1F70B",
      },
      "WDDGF8AB8EA940372": &Car{
         Model:   "2014 Mercedes-Benz C300W4",
         BrandID: "57B7B707-4357-4306-9FD6-1EDCA43CF77B",
      },
   }

   // Create new dict and initialize with vins map.
   d := dict.New(vins)

   // Add a couple more VINs.
   d.Set("1N4AL2AP4BN404580", &Car{
      Model:   "2011 NISSAN ALTIMA 2.5 S CVT",
      BrandID: "003096EE-C8FC-4C2F-ADEF-406F86C1F70B",
   })
   d.Set("4T1BE46K48U762452", &Car{
      Model:   "2008 TOYOTA Camry",
      BrandID: "C5764FE4-F1E8-46BE-AFC6-A2FC90110387",
   })

   // Check current total
   fmt.Println("Total VIN Count:", d.Len())

   // Print VINs that have 3 or more recalls
   for item := range d.Items() {
      car, ok := item.Value.(*Car)
      if !ok {
         continue // Not a Car
      }
      if car.Recalls < 3 {
         continue // Not enough recalls
      }
      fmt.Println("---")
      fmt.Println("VIN:", item.Key)
   }
}
```

[1]: https://docs.python.org/3.7/library/stdtypes.html#dict
[2]: https://github.com/srfrog/dict/blob/master/example_test.go
