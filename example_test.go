// Copyright (c) 2019 srfrog - https://srfrog.me
// Use of this source code is governed by the license in the LICENSE file.

package dict_test

import (
	"fmt"
	"time"

	"github.com/srfrog/dict"
)

type User struct {
	email  string
	name   string
	bday   time.Time
	social *dict.Dict // embedded dict
}

// Example shows a simple example of creating a new Dict and inserting values.
func Example() {
	users := []*User{
		{
			email: "madison.horton90@example.com",
			name:  "Madison Horton",
			bday:  time.Date(1971, 7, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			email: "amanda.wallace56@example.com",
			name:  "Amanda Wallace",
			bday:  time.Date(1978, 12, 6, 0, 0, 0, 0, time.UTC),
			social: dict.New(map[string]string{
				"twitter":   "@amandawall",
				"instagram": "@amanda_chill",
				"riot":      "xXxAmAnAcExXx",
			}),
		},
		{
			email: "morris.ryan98@example.com",
			name:  "Morris Ryan",
			bday:  time.Date(1984, 12, 5, 0, 0, 0, 0, time.UTC),
			social: dict.New(map[string]string{
				"linkedin": "ryan.p.morris",
			}),
		},
		{
			email: "riley.lawson20@example.com",
			name:  "Riley Lawson",
			bday:  time.Date(1984, 6, 7, 0, 0, 0, 0, time.UTC),
		},
		{
			email: "angel.perry56@example.com",
			name:  "Angel Perry",
			bday:  time.Date(1985, 8, 4, 0, 0, 0, 0, time.UTC),
		},
	}

	d := dict.New()

	// Add all users to dict. Use email as their map key.
	for _, user := range users {
		d.Set(user.email, user)
	}

	// Sanity: check that we in fact added users.
	if d.IsEmpty() {
		fmt.Println("The users dict is empty!")
		return
	}

	// Get user Amanda by email and print her info if found.
	// Get() returns an interface that will be nil if nothing is found, so we need
	// to make a type check to prevent type-assertion panic.
	user, ok := d.Get("amanda.wallace56@example.com").(*User)
	if !ok {
		fmt.Println("User was not found")
		return
	}

	// User was found, print the info.
	fmt.Println("Name:", user.name)
	fmt.Println("Birth year:", user.bday.Year())
	fmt.Println("Twitter:", user.social.Get("twitter"))

	// Amanda doesn't have WhatsApp listed, so this Get will return a nil value.
	fmt.Println("WhatsApp:", user.social.Get("whatsapp"))

	// We are done, clean up.
	if !d.Clear() {
		fmt.Println("Failed to clear")
	}

	// Output:
	// Name: Amanda Wallace
	// Birth year: 1978
	// Twitter: @amandawall
	// WhatsApp: <nil>
}

type Car struct {
	Model, BrandID string
	Recalls        int
	History        *dict.Dict
}

// ExampleDict_Update shows a dict that is updated with another using Update().
func ExampleDict_Update() {
	// Map of cars, indexed by VIN.
	// Data source: NHTSA.gov
	vins := map[string]*Car{
		"2C3KA43R08H129584": &Car{
			Model:   "2008 CHRYSLER 300",
			Recalls: 3,
			BrandID: "ACB9976A-DB5F-4D57-B9A8-9F5C53D87C7C",
		},
		"1N6AD07U78C416152": &Car{
			Model:   "2008 NISSAN FRONTIER SE-V6 RWD",
			Recalls: 0,
			BrandID: "003096EE-C8FC-4C2F-ADEF-406F86C1F70B",
		},
		"WDDGF8AB8EA940372": &Car{
			Model:   "2014 Mercedes-Benz C300W4",
			Recalls: 2,
			BrandID: "57B7B707-4357-4306-9FD6-1EDCA43CF77B",
		},
	}

	// Create new dict and initialize with vins map.
	d := dict.New(vins)

	// Add a couple more VINs.
	d.Set("1N4AL2AP4BN404580", &Car{
		Model:   "2011 NISSAN ALTIMA 2.5 S CVT",
		Recalls: 0,
		BrandID: "003096EE-C8FC-4C2F-ADEF-406F86C1F70B",
	})
	d.Set("4T1BE46K48U762452", &Car{
		Model:   "2008 TOYOTA Camry",
		Recalls: 3,
		BrandID: "C5764FE4-F1E8-46BE-AFC6-A2FC90110387",
		History: dict.New().
			Set("2008/10/21", "Vehicle sold").
			Set("2010/08/03", "Car was washed without soap"),
	})

	// Check total
	fmt.Println("Total VIN Count:", d.Len())

	// We store brands in their own dict, maybe sourced from another DB.
	brandIDs := dict.New().
		Set("ACB9976A-DB5F-4D57-B9A8-9F5C53D87C7C", "Chrysler").
		Set("003096EE-C8FC-4C2F-ADEF-406F86C1F70B", "Nissan").
		Set("57B7B707-4357-4306-9FD6-1EDCA43CF77B", "Mercedes-Benz").
		Set("C5764FE4-F1E8-46BE-AFC6-A2FC90110387", "Toyota")

	// // Keep VINs and BrandIDs in the same dict.
	d.Update(brandIDs)

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
		fmt.Println("Brand:", d.Get(car.BrandID))
		fmt.Println("Model:", car.Model)
		fmt.Println("Recalls:", car.Recalls)
		fmt.Println("Logs:", car.History.Keys())
	}

	// Output:
	// Total VIN Count: 5
	// ---
	// VIN: 2C3KA43R08H129584
	// Brand: Chrysler
	// Model: 2008 CHRYSLER 300
	// Recalls: 3
	// Logs: []
	// ---
	// VIN: 4T1BE46K48U762452
	// Brand: Toyota
	// Model: 2008 TOYOTA Camry
	// Recalls: 3
	// Logs: [2008/10/21 2010/08/03]
}
