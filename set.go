package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

const (
	rows = 4
	cols = 3
)

var (
	red           = color.New(color.FgHiRed, color.Bold).PrintfFunc()
	green         = color.New(color.FgHiGreen, color.Bold).PrintfFunc()
	blue          = color.New(color.FgHiBlue, color.Bold).PrintfFunc()
	shapeToCapped = map[Letter][]string{
		A: {"a", "A", "𝔸"},
		B: {"b", "B", "𝔹"},
		C: {"c", "C", "ℂ"},
	}
	score = 0
)

func inc() {
	score += 1
	fmt.Printf("+1 point. You now have %d point(s).\n", score)
}

func dec() {
	score -= 1
	fmt.Printf("-1 point. You now have %d point(s).\n", score)
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	board := newBoard()

	for true {
		board.Print()
		fmt.Println(`Please enter 3 cards, or "no" if there are no sets.`)
		indices := board.readNCards(3)
		if len(indices) == 0 {
			sets := board.getSets()
			if len(sets) == 0 {
				fmt.Println("Correct, this board has no sets!")
				inc()
			} else {
				fmt.Println("Oops, there were some sets.")
				for _, set := range sets {
					for _, c := range set {
						c.Print()
						fmt.Printf(" ")
					}
					fmt.Println()
				}
				dec()
			}
			board = newBoard()
		} else {
			cards := []Card{}
			for _, index := range indices {
				cards = append(cards, board[index.r][index.c])
			}
			fmt.Printf("Selected: ")
			for _, card := range cards {
				card.Print()
				fmt.Printf(" ")
			}
			fmt.Println()
			if isSet(cards) {
				fmt.Println("These 3 cards DO form a set!")
				inc()
				for _, index := range indices {
					board[index.r][index.c] = board.newCard()
				}
			} else {
				fmt.Println("These 3 cards do not form a set.")
				dec()
			}
		}
		fmt.Println()
	}
}

func isSet(cards []Card) bool {
	if len(cards) != 3 {
		panic("please check sets with exactly 3 cards")
	}
	if (cards[0].number == cards[1].number) != (cards[0].number == cards[2].number) {
		return false
	}
	if (cards[0].color == cards[1].color) != (cards[0].color == cards[2].color) {
		return false
	}
	if (cards[0].letter == cards[1].letter) != (cards[0].letter == cards[2].letter) {
		return false
	}
	if (cards[0].caps == cards[1].caps) != (cards[0].caps == cards[2].caps) {
		return false
	}
	return true
}

func (b *Board) readNCards(n int) []Index {
	result := []Index{}
	for i := 0; i < n; i++ {
		r, c := readInput()
		if r == -1 && c == -1 {
			return nil
		}
		result = append(result, Index{r, c})
	}
	return result
}

func (b *Board) getSets() [][]Card {
	result := [][]Card{}
	cards := b.getCards()
	for i := 0; i < len(cards); i++ {
	inner:
		for j := i + 1; j < len(cards); j++ {
			card3 := getMatchingCard(cards[i], cards[j])
			for _, c := range cards {
				if card3 == c {
					result = append(result, []Card{cards[i], cards[j], card3})
					break inner
				}
			}
		}
	}
	return dedupe(result)
}

func dedupe(cs [][]Card) [][]Card {
	result := [][]Card{}
	cache := map[Card]map[Card]bool{}
	for _, cards := range cs {
		cards1, ok := cache[cards[0]]
		if !ok {
			cards1 = map[Card]bool{}
		}
		present := cards1[cards[1]]
		if !present {
			result = append(result, cards)
		}
		cache[cards[0]] = map[Card]bool{cards[1]: true, cards[2]: true}
		cache[cards[1]] = map[Card]bool{cards[0]: true, cards[2]: true}
		cache[cards[2]] = map[Card]bool{cards[0]: true, cards[1]: true}
	}
	return result
}

func (b *Board) getCards() []Card {
	result := []Card{}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			result = append(result, (*b)[r][c])
		}
	}
	return result
}

type Index struct {
	r, c int
}

func readInput() (int, int) {
	var input string
	fmt.Printf("> ")
	fmt.Scanf("%s", &input)
	if input == "no" {
		return -1, -1
	}
	row := int(input[0] - '1')
	col := int(input[1] - 'a')
	if row < 0 || row >= rows || col < 0 || col >= cols {
		panic("invalid row and column input: " + input)
	}
	return row, col
}

type Board [][]Card

func (b *Board) Print() {
	fmt.Println("   a   b   c")
	for r := 0; r < rows; r++ {
		fmt.Printf("%d ", r+1)
		for c := 0; c < cols; c++ {
			(*b)[r][c].Print()
			fmt.Printf(" ")
		}
		fmt.Println()
	}
}

type Card struct {
	number int
	color  Color
	letter Letter
	caps   Caps
}

func (c Card) Print() {
	colorFn := red
	switch c.color {
	case Red:
		colorFn = red
	case Green:
		colorFn = green
	case Blue:
		colorFn = blue
	}
	s1 := shapeToCapped[c.letter]
	symbol := s1[c.caps]
	switch c.number {
	case 1:
		colorFn(" " + symbol + " ")
	case 2:
		colorFn(" " + symbol + symbol)
	case 3:
		colorFn(symbol + symbol + symbol)
	}
}

type Color int

const (
	Red Color = iota
	Green
	Blue
)

type Letter int

const (
	A Letter = iota
	B
	C
)

type Caps int

const (
	Lower Caps = iota
	Upper
	Super
)

func newBoard() Board {
	cards := map[Card]struct{}{}
	b := make([][]Card, rows)
	for r := 0; r < rows; r++ {
		b[r] = make([]Card, cols)
		for c := 0; c < cols; c++ {
			card := randomCard()
			if _, ok := cards[card]; ok {
				c--
				continue
			}
			cards[card] = struct{}{}
			b[r][c] = card
		}
	}
	return b
}

func (b *Board) hasCard(c Card) bool {
	cards := map[Card]bool{}
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			cards[(*b)[r][c]] = true
		}
	}
	return cards[c]
}

func (b *Board) newCard() Card {
	for true {
		card := randomCard()
		if !b.hasCard(card) {
			return card
		}
	}
	return Card{}
}

func randomCard() Card {
	return Card{
		number: rand.Intn(3) + 1,
		color:  Color(rand.Intn(3)),
		letter: Letter(rand.Intn(3)),
		caps:   Caps(rand.Intn(3)),
	}
}

func getMatchingCard(c1, c2 Card) Card {
	return Card{
		number: diff3(c1.number-1, c2.number-1) + 1,
		color:  Color(diff3(int(c1.color), int(c2.color))),
		letter: Letter(diff3(int(c1.letter), int(c2.letter))),
		caps:   Caps(diff3(int(c1.caps), int(c2.caps))),
	}
}

// diff3 calculates either the number if a and b are the same, or the missing one
// out of the set of {0,1,2} if a and b are the other two members of the set.
func diff3(a, b int) int {
	if a == b {
		return a
	}
	return (2 + 1 + 0) - a - b
}
