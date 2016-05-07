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

    assert_tile_valid := func(expected_validity bool, tile Tile) {
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
    assert_set_valid := func(expected_validity, expected_chow bool, set Set) {
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
    assert_set_valid(false, false, Set{Tiles: []Tile{BALLS_BASE}, Concealed: false})
    assert_set_valid(false, false, Set{Tiles: []Tile{BALLS_3}, Concealed: false})
    assert_set_valid(false, false, Set{Tiles: []Tile{1, 2}})
    assert_set_valid(false, false, Set{Tiles: []Tile{DRAGON_RED, DRAGON_WHITE}, Concealed: true})
    assert_set_valid(false, false, Set{Tiles: []Tile{DRAGON_RED, DRAGON_RED, DRAGON_WHITE, DRAGON_GREEN}, Concealed: true})
    assert_set_valid(false, false, Set{Tiles: []Tile{DRAGON_RED, DRAGON_WHITE, DRAGON_GREEN}, Concealed: true})
    assert_set_valid(false, false, Set{Tiles: []Tile{BALLS_8, BALLS_9, BALLS_9 + 1}})

    assert_set_valid(true, false, Set{Tiles: []Tile{BALLS_3, BALLS_3}})
    assert_set_valid(true, false, Set{Tiles: []Tile{BALLS_3, BALLS_3, BALLS_3}})
    assert_set_valid(true, true, Set{Tiles: []Tile{BALLS_3, BALLS_4, BALLS_5}})
    assert_set_valid(true, true, Set{Tiles: []Tile{BALLS_3, BALLS_5, BALLS_4}})
    assert_set_valid(true, false, Set{Tiles: []Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, Concealed: true})
    assert_set_valid(true, false, Set{Tiles: []Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, Concealed: true})

    assert_set_valid(true, false, Set{Tiles: []Tile{BALLS_9, BALLS_9}})
    assert_set_valid(true, true, Set{Tiles: []Tile{BALLS_2, BALLS_3, BALLS_4}})
    assert_set_valid(true, true, Set{Tiles: []Tile{BALLS_5, BALLS_6, BALLS_7}})
    assert_set_valid(true, false, Set{Tiles: []Tile{BALLS_1, BALLS_1, BALLS_1, BALLS_1}})
    assert_set_valid(true, false, Set{Tiles: []Tile{BALLS_8, BALLS_8, BALLS_8}})
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

    assert_set_score := func(expected_score, expected_doubles int, expected_type SetType, set *Set) {
        score, doubles, _ := ScoreSet(set, WIND_OWN, WIND_ROUND)

        if score == expected_score && doubles == expected_doubles && expected_type == set.set_type {
            return
        }

        as_json, err := json.Marshal(set)
        if err != nil {
            t.Fatalf("Unable to marshall set: %v", err)
        }

        t.Errorf("Set: %s\n", as_json)
        t.Errorf("Expected score=%v doubles=%v set_type=%v, got score=%v doubles=%v set_type=%v\n",
            expected_score, expected_doubles, expected_type, score, doubles, set.set_type)
    }

    // TODO: count score & doubles for flowers & seasons.

    assert_set_score(0, 0, NO_SET, &Set{})

    // Simples: chow, pung, and kong
    assert_set_score(0, 0, PILLOW, &Set{Tiles: []Tile{BAMBOO_5, BAMBOO_5}})
    assert_set_score(0, 0, NO_SET, &Set{Tiles: []Tile{BAMBOO_5, BAMBOO_6}})
    assert_set_score(0, 0, CHOW, &Set{Tiles: []Tile{BAMBOO_5, BAMBOO_6, BAMBOO_7}})
    assert_set_score(2, 0, PUNG, &Set{Tiles: []Tile{BAMBOO_5, BAMBOO_5, BAMBOO_5}})
    assert_set_score(4, 0, PUNG, &Set{Tiles: []Tile{BAMBOO_5, BAMBOO_5, BAMBOO_5}, Concealed: true})
    assert_set_score(8, 0, KONG, &Set{Tiles: []Tile{BAMBOO_5, BAMBOO_5, BAMBOO_5, BAMBOO_5}})
    assert_set_score(16, 0, KONG, &Set{Tiles: []Tile{BAMBOO_5, BAMBOO_5, BAMBOO_5, BAMBOO_5}, Concealed: true})

    // Terminals: chow, pung, and kong
    assert_set_score(0, 0, PILLOW, &Set{Tiles: []Tile{CHARS_9, CHARS_9}})
    assert_set_score(4, 0, PUNG, &Set{Tiles: []Tile{CHARS_9, CHARS_9, CHARS_9}})
    assert_set_score(8, 0, PUNG, &Set{Tiles: []Tile{CHARS_9, CHARS_9, CHARS_9}, Concealed: true})
    assert_set_score(16, 0, KONG, &Set{Tiles: []Tile{CHARS_9, CHARS_9, CHARS_9, CHARS_9}})
    assert_set_score(32, 0, KONG, &Set{Tiles: []Tile{CHARS_9, CHARS_9, CHARS_9, CHARS_9}, Concealed: true})

    // Round winds
    assert_set_score(2, 0, PILLOW, &Set{Tiles: []Tile{WIND_WEST, WIND_WEST}})
    assert_set_score(4, 1, PUNG, &Set{Tiles: []Tile{WIND_WEST, WIND_WEST, WIND_WEST}})
    assert_set_score(8, 1, PUNG, &Set{Tiles: []Tile{WIND_WEST, WIND_WEST, WIND_WEST}, Concealed: true})
    assert_set_score(16, 1, KONG, &Set{Tiles: []Tile{WIND_WEST, WIND_WEST, WIND_WEST, WIND_WEST}})
    assert_set_score(32, 1, KONG, &Set{Tiles: []Tile{WIND_WEST, WIND_WEST, WIND_WEST, WIND_WEST}, Concealed: true})

    // Own winds
    assert_set_score(2, 0, PILLOW, &Set{Tiles: []Tile{WIND_NORTH, WIND_NORTH}})
    assert_set_score(4, 1, PUNG, &Set{Tiles: []Tile{WIND_NORTH, WIND_NORTH, WIND_NORTH}})
    assert_set_score(8, 1, PUNG, &Set{Tiles: []Tile{WIND_NORTH, WIND_NORTH, WIND_NORTH}, Concealed: true})
    assert_set_score(16, 1, KONG, &Set{Tiles: []Tile{WIND_NORTH, WIND_NORTH, WIND_NORTH, WIND_NORTH}})
    assert_set_score(32, 1, KONG, &Set{Tiles: []Tile{WIND_NORTH, WIND_NORTH, WIND_NORTH, WIND_NORTH}, Concealed: true})

    // Other winds
    assert_set_score(0, 0, PILLOW, &Set{Tiles: []Tile{WIND_EAST, WIND_EAST}})
    assert_set_score(0, 0, NO_SET, &Set{Tiles: []Tile{WIND_EAST, BAMBOO_6}})
    assert_set_score(0, 0, NO_SET, &Set{Tiles: []Tile{WIND_EAST, WIND_SOUTH, WIND_WEST}})
    assert_set_score(4, 0, PUNG, &Set{Tiles: []Tile{WIND_EAST, WIND_EAST, WIND_EAST}})
    assert_set_score(16, 0, KONG, &Set{Tiles: []Tile{WIND_EAST, WIND_EAST, WIND_EAST, WIND_EAST}})

    // Dragons: pillow, pung, and kong
    assert_set_score(0, 0, NO_SET, &Set{Tiles: []Tile{DRAGON_WHITE}})
    assert_set_score(2, 0, PILLOW, &Set{Tiles: []Tile{DRAGON_GREEN, DRAGON_GREEN}})
    assert_set_score(2, 0, PILLOW, &Set{Tiles: []Tile{DRAGON_WHITE, DRAGON_WHITE}, Concealed: true })
    assert_set_score(4, 1, PUNG, &Set{Tiles: []Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}})
    assert_set_score(8, 1, PUNG, &Set{Tiles: []Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, Concealed: true })
    assert_set_score(16, 1, KONG, &Set{Tiles: []Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}})
    assert_set_score(32, 1, KONG, &Set{Tiles: []Tile{DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE, DRAGON_WHITE}, Concealed: true })
}

func TestScore(t *testing.T) {

    assert_score := func(expected_score int, hand *Hand) {
        hand.WindRound = WIND_ROUND
        hand.WindOwn = WIND_OWN
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
    if hand.Winning {
        t.Error("Hand should have been recognised as non-winning.")
    }

    // Two pairs of dragons and a pung of balls 1
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{DRAGON_RED, DRAGON_RED}},
        Set{Tiles: []Tile{DRAGON_GREEN, DRAGON_GREEN}},
        Set{Tiles: []Tile{BALLS_1, BALLS_1, BALLS_1}},
    }}
    assert_score(8, hand)
    if hand.Winning {
        t.Error("Hand should have been recognised as non-winning.")
    }

    // Two chows and a pung of simples
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{BAMBOO_4, BAMBOO_5, BAMBOO_6}},
        Set{Tiles: []Tile{BALLS_1, BALLS_3, BALLS_2}},
        Set{Tiles: []Tile{CHARS_4, CHARS_4, CHARS_4}},
    }}
    assert_score(2, hand)
    if hand.Winning {
        t.Error("Hand should have been recognised as non-winning.")
    }

    // Kong of dragons and a concealed pung of simples
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{DRAGON_GREEN, DRAGON_GREEN, DRAGON_GREEN, DRAGON_GREEN}},
        Set{Tiles: []Tile{CHARS_4, CHARS_4, CHARS_4}, Concealed: true},
    }}
    assert_score(40, hand)
    if hand.Winning {
        t.Error("Hand should have been recognised as non-winning.")
    }

    // Winning hand with no doubles.
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{BALLS_9, BALLS_9}},
        Set{Tiles: []Tile{BAMBOO_2, BAMBOO_3, BAMBOO_4}},
        Set{Tiles: []Tile{BALLS_5, BALLS_6, BALLS_7}},
        Set{Tiles: []Tile{BALLS_1, BALLS_1, BALLS_1, BALLS_1}},
        Set{Tiles: []Tile{BALLS_8, BALLS_8, BALLS_8}},
    }}
    assert_score(16 + 2 + 20, hand)
    if ! hand.Winning {
        t.Error("Hand should have been recognised as winning.")
    }
    if full_flush(hand, 38) != 0 {
        t.Error("Hand should not be detected as full flush")
    }

    // Winning hand with full flush
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{BALLS_9, BALLS_9}},
        Set{Tiles: []Tile{BALLS_2, BALLS_3, BALLS_4}},
        Set{Tiles: []Tile{BALLS_5, BALLS_6, BALLS_7}},
        Set{Tiles: []Tile{BALLS_1, BALLS_1, BALLS_1, BALLS_1}},
        Set{Tiles: []Tile{BALLS_8, BALLS_8, BALLS_8}},
    }}
    if full_flush(hand, 38) == 0 {
        t.Error("Hand should be detected as full flush")
    }
    assert_score((16 + 2 + 20) * (1 << 4), hand) // 608
    if ! hand.Winning {
        t.Error("Hand should have been recognised as winning.")
    }

    // Winning hand with pure straight
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{DRAGON_GREEN, DRAGON_GREEN}},
        Set{Tiles: []Tile{BALLS_7, BALLS_8, BALLS_9}},
        Set{Tiles: []Tile{BALLS_1, BALLS_3, BALLS_2}},
        Set{Tiles: []Tile{BALLS_5, BALLS_6, BALLS_4}},
        Set{Tiles: []Tile{BAMBOO_8, BAMBOO_8, BAMBOO_8}},
    }}
    assert_score((2 + 2 + 20) * 2, hand)  // 48
    if ! hand.Winning {
        t.Error("Hand should have been recognised as winning.")
    }

    // Non-winning hand with pure straight
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{DRAGON_GREEN}},
        Set{Tiles: []Tile{BALLS_7, BALLS_8, BALLS_9}},
        Set{Tiles: []Tile{BALLS_1, BALLS_3, BALLS_2}},
        Set{Tiles: []Tile{BALLS_5, BALLS_6, BALLS_4}},
        Set{Tiles: []Tile{BAMBOO_8, BAMBOO_8, BAMBOO_8}},
    }}
    assert_score(4, hand)
    if hand.Winning {
        t.Error("Hand should have been recognised as non-winning.")
    }

    // Winning hand with all pungs
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{DRAGON_GREEN, DRAGON_GREEN}},
        Set{Tiles: []Tile{BALLS_7, BALLS_7, BALLS_7}},
        Set{Tiles: []Tile{CHARS_2, CHARS_2, CHARS_2}},
        Set{Tiles: []Tile{BALLS_5, BALLS_5, BALLS_5}},
        Set{Tiles: []Tile{BAMBOO_8, BAMBOO_8, BAMBOO_8}},
    }}
    assert_score((2 + 8 + 20) * 2, hand) // 60
    if all_pungs(hand, 60) == 0 {
        t.Error("Hand should have been recognised as all pungs.")
    }
    if ! hand.Winning {
        t.Error("Hand should have been recognised as winning.")
    }

    // Non-winning hand with three concealed pungs
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{DRAGON_GREEN, DRAGON_GREEN}},
        Set{Tiles: []Tile{BALLS_7, BALLS_7, BALLS_7}, Concealed: true},
        Set{Tiles: []Tile{CHARS_2, CHARS_2, CHARS_2}, Concealed: true},
        Set{Tiles: []Tile{BALLS_5, BALLS_5, BALLS_5, BALLS_5}, Concealed: true},
        Set{Tiles: []Tile{BAMBOO_1, BAMBOO_2, BAMBOO_3}},
    }}
    assert_score((2 + 4 + 4 + 16 + 20) * 2, hand) // 92
    if three_concealed_pungs(hand, 92) == 0 {
        t.Error("Hand should have been recognised as three concealed pungs.")
    }
    if ! hand.Winning {
        t.Error("Hand should have been recognised as winning.")
    }

    // Winning chow hand
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{BALLS_2, BALLS_2}},
        Set{Tiles: []Tile{BALLS_1, BALLS_2, BALLS_3}},
        Set{Tiles: []Tile{CHARS_1, CHARS_2, CHARS_3}},
        Set{Tiles: []Tile{BALLS_5, BALLS_6, BALLS_7}},
        Set{Tiles: []Tile{BAMBOO_1, BAMBOO_2, BAMBOO_3}},
    }}
    assert_score(40, hand)
    if chow_hand(hand, 20) == 0 {
        t.Error("Hand should have been recognised as chow hand.")
    }
    if ! hand.Winning {
        t.Error("Hand should have been recognised as winning.")
    }

    // Winning all-simples hand
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{BALLS_2, BALLS_2}},
        Set{Tiles: []Tile{BALLS_2, BALLS_3, BALLS_4}},
        Set{Tiles: []Tile{CHARS_2, CHARS_2, CHARS_2}},
        Set{Tiles: []Tile{BALLS_5, BALLS_6, BALLS_7}},
        Set{Tiles: []Tile{BAMBOO_2, BAMBOO_3, BAMBOO_4}},
    }}
    assert_score(22 * 2, hand)
    if all_simples(hand, 44) == 0 {
        t.Error("Hand should have been recognised as all simples.")
    }

    // None-winning all-simples hand
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{CHARS_6}},
        Set{Tiles: []Tile{BALLS_2}},
        Set{Tiles: []Tile{BALLS_2, BALLS_3, BALLS_4}},
        Set{Tiles: []Tile{CHARS_2, CHARS_2, CHARS_2}},
        Set{Tiles: []Tile{BALLS_5, BALLS_6, BALLS_7}},
        Set{Tiles: []Tile{BAMBOO_2, BAMBOO_3, BAMBOO_4}},
    }}
    assert_score(4, hand)
    if all_simples(hand, 4) == 0 {
        t.Error("Hand should have been recognised as all simples.")
    }

    // None-winning terminals & honours hand
    hand = &Hand{Sets: []Set{
        Set{Tiles: []Tile{DRAGON_GREEN}},
        Set{Tiles: []Tile{WIND_EAST}},
        Set{Tiles: []Tile{WIND_WEST, WIND_WEST, WIND_WEST}}, // 4 + 1d
        Set{Tiles: []Tile{CHARS_1, CHARS_1, CHARS_1}}, // 4
        Set{Tiles: []Tile{WIND_NORTH, WIND_NORTH, WIND_NORTH}}, // 4 + 1d
        Set{Tiles: []Tile{BAMBOO_9, BAMBOO_9, BAMBOO_9}}, // 4
    }}
    assert_score(16 * 8, hand)
    if all_terminals_honours(hand, 128) == 0 {
        t.Error("Hand should have been recognised as all terminals & honours.")
    }
}

func TestTile_Suit(t *testing.T) {
    assert := func(expect_suit, tile Tile) {
        if tile.Suit() != expect_suit {
            t.Errorf("Tile %q doesn't have expected suit %q but %q",
                tile, expect_suit, tile.Suit())
        }
    }

    assert(BAMBOO_BASE, BAMBOO_1)
    assert(BAMBOO_BASE, BAMBOO_9)
    assert(BALLS_BASE, BALLS_1)
    assert(BALLS_BASE, BALLS_9)
    assert(CHARS_BASE, CHARS_1)
    assert(CHARS_BASE, CHARS_9)
    assert(NO_TILE, BAMBOO_9 + 1)
    assert(NO_TILE, SEASON_3)
}
