package deck

import "testing"

func TestShuffle(t *testing.T) {
	testDeck := New()
	testDeck.Shuffle()
}

func TestRemoveAndShuffle(t *testing.T) {
	testDeck := New()
	testDeck.Remove(Card{Value: 1, Stave: "triangle"})
	testDeck.Remove(Card{Value: 10, Stave: "circle"})
	testDeck.Remove(Card{Value: 8, Stave: "square"})
	testDeck.Remove(Card{Value: 7, Stave: "triangle"})
	testDeck.Shuffle()
}

func TestDealAndShuffle(t *testing.T) {
	testDeck := New()
	_ = testDeck.Deal(2)
	_ = testDeck.Deal(2)
	testDeck.Shuffle()
}
