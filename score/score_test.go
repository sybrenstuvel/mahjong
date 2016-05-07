package score

import (
    "testing"
    "encoding/json"
)

const (
    WIND_OWN = WIND_NORTH
    WIND_ROUND = WIND_WEST
)

func TestTileValid(t *testing.T) {

    assert_tile_valid := func (expected_validity bool, tile Tile) {
        validity := tile.IsValid()

        if validity != expected_validity {
            t.Errorf("Tile %q: expected validity=%v, got validity=%v\n",
                tile, expected_validity, validity)
        }
    }

    assert_tile_valid(false, BALLS_BASE)
    assert_tile_valid(false, -1)
    assert_tile_valid(false, 0xffff)
    assert_tile_valid(true, BALLS_3)
    assert_tile_valid(true, 13)
    assert_tile_valid(false, CHARS_1 - 1)
    assert_tile_valid(false, CHARS_9 + 1)
    assert_tile_valid(true, CHARS_9 - 1)
    assert_tile_valid(true, BAMBOO_1 + 1)
    assert_tile_valid(false, BAMBOO_1 - 1)
}

func TestSetValid(t *testing.T) {
    assert_set_valid := func (expected_validity, expected_chow bool, set Set) {
        validity, is_chow := IsValidSet(set)

        if validity == expected_validity && is_chow == expected_chow {
            return
        }

        as_json, err := json.Marshal(set)
        if err != nil {
            t.Fatalf("Unable to marshall set: %v", err)
        }

        t.Errorf("Set: %s\n", as_json)
        t.Errorf("Expected validity=%v is_chow=%v, got validity=%v is_chow=%v\n",
            expected_validity, expected_chow, validity, is_chow)
    }

    assert_set_valid(false, false, Set{})
    assert_set_valid(false, false, Set{Tiles: []Tile{BALLS_3, BALLS_4}, Concealed: false})
    assert_set_valid(false, false, Set{[]Tile{BALLS_BASE}, false})
    assert_set_valid(false, false, Set{[]Tile{BALLS_3}, false})
    assert_set_valid(false, false, Set{Tiles: []Tile{1, 2}})
    assert_set_valid(false, false, Set{[]Tile{DRAGON_RED, DRAGON_WHITE}, true})
    assert_set_valid(false, false, Set{[]Tile{DRAGON_RED, DRAGON_RED, DRAGON_WHITE, DRAGON_GREEN}, true})
    assert_set_valid(false, false, Set{[]Tile{DRAGON_RED, DRAGON_WHITE, DRAGON_GREEN}, true})
    assert_set_valid(false, false, Set{Tiles: []Tile{BALLS_8, BALLS_9, BALLS_9 + 1}})

    assert_set_valid(true, false, Set{Tiles: []Tile{BALLS_3, BALLS_3}})
    assert_set_valid(true, false, Set{Tiles: []Tile{BALLS_3, BALLS_3, BALLS_3}})
    assert_set_valid(true, true, Set{Tiles: []Tile{BALLS_3, BALLS_4, BALLS_5}})
    assert_set_valid(true, true, Set{Tiles: []Tile{BALLS_3, BALLS_5, BALLS_4}})
    assert_set_valid(true, false, Set{[]Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, true})
    assert_set_valid(true, false, Set{[]Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, true})

    assert_set_valid(true, false, Set{[]Tile{BALLS_9, BALLS_9}, false})
    assert_set_valid(true, true, Set{[]Tile{BALLS_2, BALLS_3, BALLS_4}, false})
    assert_set_valid(true, true, Set{[]Tile{BALLS_5, BALLS_6, BALLS_7}, false})
    assert_set_valid(true, false, Set{[]Tile{BALLS_1, BALLS_1, BALLS_1, BALLS_1}, false})
    assert_set_valid(true, false, Set{[]Tile{BALLS_8, BALLS_8, BALLS_8}, false})
}

func TestIsDragon(t *testing.T) {
    assert := func(is_dragon bool, dragon Tile) {
        if (is_dragon != dragon.IsDragon()) {
            t.Fatalf("%v is not correctly classified as dragon", dragon)
        }
    }

    assert(true, DRAGON_WHITE)
    assert(true, DRAGON_GREEN)
    assert(true, DRAGON_RED)
    assert(false, DRAGON_BASE)
    assert(false, DRAGON_WHITE + 1)
}

func TestIsWind(t *testing.T) {
    assert := func(is_wind bool, wind Tile) {
        if (is_wind != wind.IsWind()) {
            t.Fatalf("%v is not correctly classified as wind", wind)
        }
    }

    assert(true, WIND_EAST)
    assert(true, WIND_SOUTH)
    assert(true, WIND_WEST)
    assert(true, WIND_NORTH)
    assert(false, WIND_BASE)
    assert(false, WIND_NORTH + 1)
}

func TestIsValid(t *testing.T) {
    assert := func(is_valid bool, valid Tile) {
        if (is_valid != valid.IsValid()) {
            t.Fatalf("%v is not correctly classified as valid", valid)
        }
    }

    assert(true, WIND_EAST)
    assert(true, BALLS_1)
    assert(true, BALLS_9)
    assert(true, SEASON_4)
    assert(false, BALLS_BASE)
    assert(false, WIND_NORTH + 1)
    assert(false, BALLS_9 + 1)
    assert(false, SEASON_4 + 1)
}

func TestSetScores(t *testing.T) {

    assert_set_score := func(expected_score, expected_doubles int, set *Set) {
        score, doubles, _ := ScoreSet(set, WIND_OWN, WIND_ROUND)

        if score == expected_score && doubles == expected_doubles {
            return
        }

        as_json, err := json.Marshal(set)
        if err != nil {
            t.Fatalf("Unable to marshall set: %v", err)
        }

        t.Errorf("Set: %s\n", as_json)
        t.Errorf("Expected score=%v with %v doubles, got score=%v with %v doubles\n",
            expected_score, expected_doubles, score, doubles)
    }

    // TODO: count score & doubles for flowers & seasons.

    assert_set_score(0, 0, &Set{})

    // Simples: chow, pung, and kong
    assert_set_score(0, 0, &Set{[]Tile{BAMBOO_5, BAMBOO_5}, false})
    assert_set_score(0, 0, &Set{[]Tile{BAMBOO_5, BAMBOO_6}, false})
    assert_set_score(0, 0, &Set{[]Tile{BAMBOO_5, BAMBOO_6, BAMBOO_7}, false})
    assert_set_score(2, 0, &Set{[]Tile{BAMBOO_5, BAMBOO_5, BAMBOO_5}, false})
    assert_set_score(4, 0, &Set{[]Tile{BAMBOO_5, BAMBOO_5, BAMBOO_5}, true})
    assert_set_score(8, 0, &Set{[]Tile{BAMBOO_5, BAMBOO_5, BAMBOO_5, BAMBOO_5}, false})
    assert_set_score(16, 0, &Set{[]Tile{BAMBOO_5, BAMBOO_5, BAMBOO_5, BAMBOO_5}, true})

    // Terminals: chow, pung, and kong
    assert_set_score(0, 0, &Set{[]Tile{CHARS_9, CHARS_9}, false})
    assert_set_score(4, 0, &Set{[]Tile{CHARS_9, CHARS_9, CHARS_9}, false})
    assert_set_score(8, 0, &Set{[]Tile{CHARS_9, CHARS_9, CHARS_9}, true})
    assert_set_score(16, 0, &Set{[]Tile{CHARS_9, CHARS_9, CHARS_9, CHARS_9}, false})
    assert_set_score(32, 0, &Set{[]Tile{CHARS_9, CHARS_9, CHARS_9, CHARS_9}, true})

    // Round winds
    assert_set_score(2, 0, &Set{[]Tile{WIND_ROUND, WIND_ROUND}, false})
    assert_set_score(4, 1, &Set{[]Tile{WIND_ROUND, WIND_ROUND, WIND_ROUND}, false})
    assert_set_score(8, 1, &Set{[]Tile{WIND_ROUND, WIND_ROUND, WIND_ROUND}, true})
    assert_set_score(16, 1, &Set{[]Tile{WIND_ROUND, WIND_ROUND, WIND_ROUND, WIND_ROUND}, false})
    assert_set_score(32, 1, &Set{[]Tile{WIND_ROUND, WIND_ROUND, WIND_ROUND, WIND_ROUND}, true})

    // Own winds
    assert_set_score(2, 0, &Set{[]Tile{WIND_OWN, WIND_OWN}, false})
    assert_set_score(4, 1, &Set{[]Tile{WIND_OWN, WIND_OWN, WIND_OWN}, false})
    assert_set_score(8, 1, &Set{[]Tile{WIND_OWN, WIND_OWN, WIND_OWN}, true})
    assert_set_score(16, 1, &Set{[]Tile{WIND_OWN, WIND_OWN, WIND_OWN, WIND_OWN}, false})
    assert_set_score(32, 1, &Set{[]Tile{WIND_OWN, WIND_OWN, WIND_OWN, WIND_OWN}, true})

    // Other winds
    assert_set_score(0, 0, &Set{[]Tile{WIND_EAST, WIND_EAST}, false})
    assert_set_score(0, 0, &Set{[]Tile{WIND_EAST, BAMBOO_6}, false})
    assert_set_score(0, 0, &Set{[]Tile{WIND_EAST, WIND_SOUTH, WIND_WEST}, false})
    assert_set_score(4, 0, &Set{[]Tile{WIND_EAST, WIND_EAST, WIND_EAST}, false})
    assert_set_score(8, 0, &Set{[]Tile{WIND_EAST, WIND_EAST, WIND_EAST}, true})
    assert_set_score(16, 0, &Set{[]Tile{WIND_EAST, WIND_EAST, WIND_EAST, WIND_EAST}, false})
    assert_set_score(32, 0, &Set{[]Tile{WIND_EAST, WIND_EAST, WIND_EAST, WIND_EAST}, true})

    // Dragons: pillow, pung, and kong
    assert_set_score(0, 0, &Set{[]Tile{DRAGON_WHITE}, false})
    assert_set_score(2, 0, &Set{[]Tile{DRAGON_WHITE, DRAGON_WHITE}, false})
    assert_set_score(2, 0, &Set{[]Tile{DRAGON_WHITE, DRAGON_WHITE}, true})
    assert_set_score(4, 1, &Set{[]Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, false})
    assert_set_score(8, 1, &Set{[]Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, true})
    assert_set_score(16, 1, &Set{[]Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, false})
    assert_set_score(32, 1, &Set{[]Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, true})
}

func TestSimpleScores(t *testing.T) {

    assert_score := func(expected_score int, hand *Hand) {
        score := Score(hand)

        if score == expected_score {
            return
        }

        as_json, err := json.Marshal(hand)
        if err != nil {
            t.Fatalf("Unable to marshall hand: %v", err)
        }

        t.Errorf("Hand: %s\n", as_json)
        t.Errorf("Expected score=%v, got score=%v\n", expected_score, score)
    }

    var hand *Hand

    // Empty hand
    hand = &Hand{}
    assert_score(0, hand)

    // Two pairs of dragons and a pung of balls 1
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{DRAGON_RED, DRAGON_RED}},
        Set{Tiles: []Tile{DRAGON_GREEN, DRAGON_GREEN}},
        Set{Tiles: []Tile{BALLS_1, BALLS_1, BALLS_1}},
    }}
    assert_score(8, hand)

    // Two chows and a pung of simples
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{BAMBOO_4, BAMBOO_5, BAMBOO_6}},
        Set{Tiles: []Tile{BALLS_2, BALLS_3, BALLS_4}},
        Set{Tiles: []Tile{CHARS_4, CHARS_4, CHARS_4}},
    }}
    assert_score(2, hand)

    // Kong of dragons and a concealed pung of simples
    hand = &Hand{Sets: []Set{
        Set{[]Tile{DRAGON_GREEN, DRAGON_GREEN, DRAGON_GREEN, DRAGON_GREEN}, false},
        Set{[]Tile{CHARS_4, CHARS_4, CHARS_4}, true},
    }}
    assert_score(40, hand)

    // Winning hand
    hand = &Hand{Sets: []Set{
        Set{[]Tile{BALLS_9, BALLS_9}, false},
        Set{[]Tile{BALLS_2, BALLS_3, BALLS_4}, false},
        Set{[]Tile{BALLS_5, BALLS_6, BALLS_7}, false},
        Set{[]Tile{BALLS_1, BALLS_1, BALLS_1, BALLS_1}, false},
        Set{[]Tile{BALLS_8, BALLS_8, BALLS_8}, false},
    }}
    // Wrong score, as we don't take most doubles into account.
    assert_score(16+2+20, hand)
    if ! hand.Winning {
        t.Error("Hand should have been recognised as winning.")
    }
}
