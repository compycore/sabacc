package deck

import "testing"

func TestShuffle(t *testing.T) {
	testDeck := New()
	if len(testDeck) != 62 {
		t.Errorf("Deck size not correct; want: %d, got: %d", 62, len(testDeck))
	}

	// There's no real way to test Shuffle but I want to run it and make sure it doesn't throw an out of bounds error
	testDeck = Shuffle(testDeck)
}

func TestRemove(t *testing.T) {
	testDeck := New()

	testDeck = Remove(testDeck, Card{Value: 1, Stave: "triangle"})
	testDeck = Remove(testDeck, Card{Value: 10, Stave: "circle"})
	testDeck = Remove(testDeck, Card{Value: 8, Stave: "square"})
	testDeck = Remove(testDeck, Card{Value: 7, Stave: "triangle"})

	if len(testDeck) != 58 {
		t.Errorf("Deck size not correct; want: %d, got: %d", 58, len(testDeck))
	}
}

func TestDeal(t *testing.T) {
	testDeck := New()

	testDeck, hand := Deal(testDeck, 2)
	if len(testDeck) != 60 {
		t.Errorf("Deck size not correct; want: %d, got: %d", 60, len(testDeck))
	}

	if len(hand) != 2 {
		t.Errorf("Hand size not correct; want: %d, got: %d", 2, len(hand))
	}
}
