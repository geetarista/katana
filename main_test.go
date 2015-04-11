package main

import "testing"

func TestRoll_Bounds(t *testing.T) {
	var err error
	var n int

	for i := 0; i < 10000; i++ {
		n, err = roll()

		if n < 1 || n > 6 {
			t.Fatal("Roll outside bounds: %d", roll)
		}

		if err != nil {
			t.Fatal("Error rolling: %s", err)
		}
	}
}

func TestChuck_Bounds(t *testing.T) {
	var err error
	var n int

	for i := 0; i < 10000; i++ {
		n, err = chuck()
		if err != nil {
			t.Fatal("Error chucking: %s", err)
		}

		if _, ok := wordMap[n]; !ok {
			t.Fatal("Index outside bounds: %d", n)
		}
	}
}
