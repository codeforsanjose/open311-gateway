package geo

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestLookup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestLookup <<<<<<<<<<<<<<<<<<<<<<<<<<")
	req := &Request{
		Address:  "New York City, NY",
		Provider: GOOGLE,
	}
	resp, err := req.Lookup(nil)
	if err != nil {
		t.Fatalf("Lookup error: %v", err)
	}
	fmt.Printf("Response: %s\n", spew.Sdump(resp))
	if s := resp.Status; s != "OK" {
		t.Fatalf(`Status == %q, want "OK"`, s)
	}
	if l := len(resp.Results); l != 1 {
		t.Fatalf("len(Results) == %d, want 1", l)
	}
	addr := "New York, NY, USA"
	if a := resp.Found; a != addr {
		t.Errorf("Address == %q, want %q", a, addr)
	}
}

func TestLookup2(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestLookup2 <<<<<<<<<<<<<<<<<<<<<<<<<<")
	req := &Request{
		Address:  "17200 Quail Ct, Morgan Hill, CA",
		Provider: GOOGLE,
	}
	resp, err := req.Lookup(nil)
	if err != nil {
		t.Fatalf("Lookup error: %v", err)
	}
	fmt.Printf("Response: %s\n", spew.Sdump(resp))
	if s := resp.Status; s != "OK" {
		t.Fatalf(`Status == %q, want "OK"`, s)
	}
	if l := len(resp.Results); l != 1 {
		t.Fatalf("len(Results) == %d, want 1", l)
	}
	addr := "17200 Quail Ct, Morgan Hill, CA 95037, USA"
	if a := resp.Found; a != addr {
		t.Errorf("Address == %q, want %q", a, addr)
	}
}

func TestReverseLookup(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestReverseLookup <<<<<<<<<<<<<<<<<<<<<<<<<<")
	loc := Point{37.338208, -121.886329}
	req := &Request{
		Location: &loc,
		Provider: GOOGLE,
	}
	resp, err := req.Lookup(nil)
	if err != nil {
		t.Fatalf("Lookup error: %v", err)
	}
	fmt.Printf("Found: %s\n", resp.Found)
	fmt.Printf("Response: %s\n", spew.Sdump(resp))
	if s := resp.Status; s != "OK" {
		t.Fatalf(`Status == %q, want "OK"`, s)
	}
}

func TestLookupWithBounds(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestLookupWithBounds <<<<<<<<<<<<<<<<<<<<<<<<<<")
	req := &Request{
		Address:  "Winnetka",
		Provider: GOOGLE,
	}
	bounds := &Bounds{Point{34.172684, -118.604794},
		Point{34.236144, -118.500938}}
	req.Bounds = bounds
	resp, err := req.Lookup(nil)
	if err != nil {
		t.Fatalf("Lookup error: %v", err)
	}
	fmt.Printf("Response: %s\n", spew.Sdump(resp))
	if s := resp.Status; s != "OK" {
		t.Fatalf(`Status == %q, want "OK"`, s)
	}
	if l := len(resp.Results); l != 1 {
		t.Fatalf("len(Results) == %d, want 1", l)
	}
	addr := "Winnetka, Los Angeles, CA, USA"
	if a := resp.Found; a != addr {
		t.Errorf("Address == %q, want %q", a, addr)
	}
}

func TestLookupWithLanguage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestLookupWithLanguage <<<<<<<<<<<<<<<<<<<<<<<<<<")
	req := &Request{
		Address:  "札幌市",
		Provider: GOOGLE,
	}
	req.Language = "ja"
	resp, err := req.Lookup(nil)
	if err != nil {
		t.Fatalf("Lookup error: %v", err)
	}
	if s := resp.Status; s != "OK" {
		t.Fatalf(`Status == %q, want "OK"`, s)
	}
	if l := len(resp.Results); l != 1 {
		t.Fatalf("len(Results) == %d, want 1", l)
	}
	addr := "日本, 北海道札幌市"
	if a := resp.Found; a != addr {
		t.Errorf("Address == %q, want %q", a, addr)
	}
}

func TestLookupWithRegion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestLookupWithRegion <<<<<<<<<<<<<<<<<<<<<<<<<<")
	req := &Request{
		Address:  "Toledo",
		Provider: GOOGLE,
	}
	req.Region = "es"
	resp, err := req.Lookup(nil)
	if err != nil {
		t.Fatalf("Lookup error: %v", err)
	}
	if s := resp.Status; s != "OK" {
		t.Fatalf(`Status == %q, want "OK"`, s)
	}
	if l := len(resp.Results); l != 1 {
		t.Fatalf("len(Results) == %d, want 1", l)
	}
	addr := "Toledo, Toledo, Spain"
	if a := resp.Found; a != addr {
		t.Errorf("Address == %q, want %q", a, addr)
	}
}

// func TestGetBounds(t *testing.T) {
// 	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestGetBounds <<<<<<<<<<<<<<<<<<<<<<<<<<")
// 	center := Point{37.138698, -121.615391}
// 	var radius float64 = 500

// 	ne, sw, _ := GetBounds(&center, radius)
// 	fmt.Printf("NE: %v | %v   SW: %v | %v\n", ne.Lat, ne.Lng, sw.Lat, sw.Lng)
// }

func TestGetLatLng(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestLookupWithRegion <<<<<<<<<<<<<<<<<<<<<<<<<<")
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestGetLatLng <<<<<<<<<<<<<<<<<<<<<<<<<<")
	addr := "17200 Quail Ct., Morgan Hill, CA"
	lat, lng, err := GetLatLng(addr)
	if err != nil {
		t.Fatalf("Lookup error: %v", err)
	}
	fmt.Printf("Coordinates: %v | %v\n", lat, lng)
}

func TestGetAddress(t *testing.T) {
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestGetAddress <<<<<<<<<<<<<<<<<<<<<<<<<<")
	var (
		lat = 37.151181
		lng = -121.602626
	)

	addr, err := GetAddress(lat, lng)
	if err != nil {
		t.Fatalf("Lookup error: %v", err)
	}
	expected := "17200 Quail Ct, Morgan Hill, CA 95037, USA"
	if addr != expected {
		t.Errorf("Address == %q, want %q", addr, expected)
	}
	fmt.Printf("Address: %s\n", addr)
}

func TestGetCity(t *testing.T) {
	fmt.Println("\n>>>>>>>>>>>>>>>>>>> TestGetCity <<<<<<<<<<<<<<<<<<<<<<<<<<")
	var (
		lat = 37.151181
		lng = -121.602626
	)

	city, err := GetCity(lat, lng)
	if err != nil {
		t.Fatalf("Lookup error: %v", err)
	}
	expected := "Morgan Hill"
	if city != expected {
		t.Errorf("Address == %q, want %q", city, expected)
	}
	fmt.Printf("Address: %s\n", city)
}
