package router

import (
	"fmt"
)

func init() {
	fmt.Println("Loading routing data...")
}

// map: city -> list of Services (MID)

// func: MID -> Adapter

// Loaded from config file:
// map: Gateway Code -> Adapter API Call
// map: city name -> City Code (MID)
