/*
 * Detectors for hand-wide doubles.
 */

package score

// A Detector returns the number of doubles detected.
type Detector func(hand *Hand, simpleScore int) int

var detectors = map[string]Detector{
	"pure straight":         pureStraight,
	"all pungs":             allPungs,
	"full flush":            fullFlush,
	"three concealed pungs": threeConcealedPungs,
	"chow hand":             chowHand,
	"all simples":           allSimples,
	"all terminals/honours": allTerminalsHonours,
	"half-flush":            halfFlush,
	"outside hand":          outsideHand,
}

func findSetsOfType(hand *Hand, setType SetType) chan *Set {
	ch := make(chan *Set)

	go func() {
		for idx := range hand.Sets {
			set := &hand.Sets[idx]
			if set.setType&setType > 0 {
				ch <- set
			}
		}
		close(ch)
	}()

	return ch
}

func allTiles(hand *Hand) chan Tile {
	ch := make(chan Tile)

	go func() {
		for idx := range hand.Sets {
			set := &hand.Sets[idx]
			for _, tile := range set.Tiles {
				ch <- tile
			}
		}
		close(ch)
	}()

	return ch
}

func pureStraight(hand *Hand, simpleScore int) int {
	// Find the chows
	nrOfChows := 0
	suit := NoTile

	for chow := range findSetsOfType(hand, Chow) {
		switch {
		case suit == NoTile:
			suit = chow.Tiles[0].Suit()
		case chow.Tiles[0].Suit() != suit:
			return 0
			// the sets must start with 1, 4, 7
		case chow.Tiles[0].Number() != nrOfChows*3+1:
			return 0
		}

		nrOfChows++
	}

	if nrOfChows < 3 {
		return 0
	}

	return 1
}

func allPungs(hand *Hand, simpleScore int) int {
	if !hand.Winning {
		return 0
	}

	count := 0
	for _ = range findSetsOfType(hand, Pung+Kong) {
		count++
	}

	if count == 4 {
		return 1
	}

	return 0
}

func fullFlush(hand *Hand, simpleScore int) int {
	if len(hand.Sets) < 1 || len(hand.Sets[0].Tiles) < 1 {
		return 0
	}

	suit := hand.Sets[0].Tiles[0].Suit()
	if suit == NoTile {
		return 0
	}

	// Check that every tile is of the same suit
	for tile := range allTiles(hand) {
		if tile.Suit() != suit {
			return 0
		}
	}

	return 4
}

func threeConcealedPungs(hand *Hand, simpleScore int) int {
	count := 0
	for set := range findSetsOfType(hand, Pung+Kong) {
		if set.Concealed {
			count++
		}
	}

	if count >= 3 {
		return 1
	}

	return 0
}

func chowHand(hand *Hand, simpleScore int) int {
	if (hand.Winning && simpleScore > 20) || (!hand.Winning && simpleScore > 0) {
		return 0
	}

	count := 0
	for _ = range findSetsOfType(hand, Chow) {
		count++
	}

	if count == 4 {
		return 1
	}

	return 0
}

func allSimples(hand *Hand, simpleScore int) int {
	count := 0
	for tile := range allTiles(hand) {
		count++
		if !tile.IsSimple() {
			return 0
		}
	}

	// Incomplete hand.
	if count < 13 {
		return 0
	}
	return 1
}

func allTerminalsHonours(hand *Hand, simpleScore int) int {
	count := 0
	for tile := range allTiles(hand) {
		count++
		if !tile.IsTerminal() && !tile.IsHonour() {
			return 0
		}
	}

	// Incomplete hand.
	if count < 13 {
		return 0
	}
	return 1
}

func halfFlush(hand *Hand, simpleScore int) int {
	suit := NoTile
	seenHonour := false

	count := 0
	for tile := range allTiles(hand) {
		count++

		switch {
		case tile.IsHonour():
			seenHonour = true
		case suit == NoTile:
			suit = tile.Suit()
		case tile.Suit() != suit:
			return 0
		}
	}

	// That's a full flush, and is detected somewhere else.
	if !seenHonour {
		return 0
	}

	// Incomplete hand.
	if count < 13 {
		return 0
	}

	return 1
}

func outsideHand(hand *Hand, simpleScore int) int {
	if !hand.Winning {
		return 0
	}

	nrOfChows := 0
	for idx := range hand.Sets {
		set := &hand.Sets[idx]
		if set.setType == Chow {
			nrOfChows++
		}
		if !set.HasTerminalOrHonour() {
			return 0
		}
	}

	if nrOfChows == 0 {
		return 0
	}

	return 1
}
