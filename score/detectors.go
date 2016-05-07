/*
 * Detectors for hand-wide doubles.
 */

package score


// Returns the number of doubles from this detector
type Detector func (hand *Hand, simple_score int) int

var detectors = []Detector{
    pure_straight,
    all_pungs,
    full_flush,
    three_concealed_pungs,
    chow_hand,
}


func find_sets_of_type(hand *Hand, set_type SetType) chan *Set {
    ch := make(chan *Set)

    go func() {
        for idx, _ := range hand.Sets {
            set := &hand.Sets[idx]
            if set.set_type & set_type > 0 {
                ch <- set
            }
        }
        close(ch)
    }()

    return ch
}


func pure_straight(hand *Hand, simple_score int) int {
    // Find the chows
    nr_of_chows := 0
    suit := NO_TILE

    for chow := range find_sets_of_type(hand, CHOW) {
        switch {
        case suit == NO_TILE:
            suit = chow.Tiles[0].Suit()
        case chow.Tiles[0].Suit() != suit:
            return 0
        // the sets must start with 1, 4, 7
        case chow.Tiles[0].Number() != nr_of_chows * 3 + 1:
            return 0
        }

        nr_of_chows++
    }

    if nr_of_chows < 3 {
        return 0
    }

    return 1
}

func all_pungs(hand *Hand, simple_score int) int {
    if ! hand.Winning {
        return 0
    }

    count := 0
    for _ = range find_sets_of_type(hand, PUNG + KONG) {
        count++
    }

    if count == 4 {
        return 1
    }

    return 0
}

func full_flush(hand *Hand, simple_score int) int {
    if len(hand.Sets) < 1 || len(hand.Sets[0].Tiles) < 1 {
        return 0
    }

    var suit Tile = hand.Sets[0].Tiles[0].Suit()
    if suit == NO_TILE {
        return 0
    }

    // Check that every tile is of the same suit
    for _, set := range hand.Sets {
        for _, tile := range set.Tiles {
            if tile.Suit() != suit {
                return 0
            }
        }
    }

    return 4
}


func three_concealed_pungs(hand *Hand, simple_score int) int {
    count := 0
    for set := range find_sets_of_type(hand, PUNG + KONG) {
        if set.Concealed {
            count += 1
        }
    }

    if count >= 3 {
        return 1
    }

    return 0
}


func chow_hand(hand *Hand, simple_score int) int {
    if (hand.Winning && simple_score > 20) || (!hand.Winning && simple_score > 0){
        return 0
    }

    count := 0
    for _ = range find_sets_of_type(hand, CHOW) {
        count += 1
    }

    if count == 4 {
        return 1
    }

    return 0
}
