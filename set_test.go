package main

import (
	"testing"
)

func TestGetMatchingCard1(t *testing.T) {
	c1 := Card{
		number: 1,
		color:  Green,
		letter: A,
		caps:   Lower,
	}
	c2 := Card{
		number: 3,
		color:  Red,
		letter: A,
		caps:   Super,
	}
	want := Card{
		number: 2,
		color:  Blue,
		letter: A,
		caps:   Upper,
	}

	if got := getMatchingCard(c1, c2); got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestGetMatchingCard2(t *testing.T) {
	c1 := Card{
		number: 2,
		color:  Blue,
		letter: A,
		caps:   Upper,
	}
	c2 := Card{
		number: 2,
		color:  Green,
		letter: B,
		caps:   Upper,
	}
	want := Card{
		number: 2,
		color:  Red,
		letter: C,
		caps:   Upper,
	}

	if got := getMatchingCard(c1, c2); got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
