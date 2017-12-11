package score

import (
	"encoding/json"

	check "gopkg.in/check.v1"
)

const (
	WIND_OWN   = WindNorth
	WIND_ROUND = WindWest
)

type ScoreTestSuite struct{}

var _ = check.Suite(&ScoreTestSuite{})

func (s *ScoreTestSuite) TestTileValid(t *check.C) {
	assert_tile_valid := func(expected_validity bool, tile Tile) {
		validity := tile.IsValid()

		if validity != expected_validity {
			t.Errorf("Tile %q: expected validity=%v, got validity=%v\n",
				tile, expected_validity, validity)
		}
	}

	assert_tile_valid(false, ballsBase)
	assert_tile_valid(false, -1)
	assert_tile_valid(false, 0xffff)
	assert_tile_valid(true, Balls3)
	assert_tile_valid(true, 13)
	assert_tile_valid(false, Chars1-1)
	assert_tile_valid(false, Chars9+1)
	assert_tile_valid(true, Chars9-1)
	assert_tile_valid(true, Bamboo1+1)
	assert_tile_valid(false, Bamboo1-1)
}

func (s *ScoreTestSuite) TestSetValid(t *check.C) {
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
	assert_set_valid(false, false, Set{Tiles: []Tile{Balls3, Balls4}, Concealed: false})
	assert_set_valid(false, false, Set{Tiles: []Tile{ballsBase}, Concealed: false})
	assert_set_valid(false, false, Set{Tiles: []Tile{Balls3}, Concealed: false})
	assert_set_valid(false, false, Set{Tiles: []Tile{1, 2}})
	assert_set_valid(false, false, Set{Tiles: []Tile{DragonRed, DragonWhite}, Concealed: true})
	assert_set_valid(false, false, Set{Tiles: []Tile{DragonRed, DragonRed, DragonWhite, DragonGreen}, Concealed: true})
	assert_set_valid(false, false, Set{Tiles: []Tile{DragonRed, DragonWhite, DragonGreen}, Concealed: true})
	assert_set_valid(false, false, Set{Tiles: []Tile{Balls8, Balls9, Balls9 + 1}})

	assert_set_valid(true, false, Set{Tiles: []Tile{Balls3, Balls3}})
	assert_set_valid(true, false, Set{Tiles: []Tile{Balls3, Balls3, Balls3}})
	assert_set_valid(true, true, Set{Tiles: []Tile{Balls3, Balls4, Balls5}})
	assert_set_valid(true, true, Set{Tiles: []Tile{Balls3, Balls5, Balls4}})
	assert_set_valid(true, false, Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite}, Concealed: true})
	assert_set_valid(true, false, Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite, DragonWhite}, Concealed: true})

	assert_set_valid(true, false, Set{Tiles: []Tile{Balls9, Balls9}})
	assert_set_valid(true, true, Set{Tiles: []Tile{Balls2, Balls3, Balls4}})
	assert_set_valid(true, true, Set{Tiles: []Tile{Balls5, Balls6, Balls7}})
	assert_set_valid(true, false, Set{Tiles: []Tile{Balls1, Balls1, Balls1, Balls1}})
	assert_set_valid(true, false, Set{Tiles: []Tile{Balls8, Balls8, Balls8}})
}

func (s *ScoreTestSuite) TestIsDragon(t *check.C) {
	assert := func(is_dragon bool, dragon Tile) {
		if is_dragon != dragon.IsDragon() {
			t.Fatalf("%v is not correctly classified as dragon", dragon)
		}
	}

	assert(true, DragonWhite)
	assert(true, DragonGreen)
	assert(true, DragonRed)
	assert(false, dragonBase)
	assert(false, DragonWhite+1)
}

func (s *ScoreTestSuite) TestIsWind(t *check.C) {
	assert := func(is_wind bool, wind Tile) {
		if is_wind != wind.IsWind() {
			t.Fatalf("%v is not correctly classified as wind", wind)
		}
	}

	assert(true, WindEast)
	assert(true, WindSouth)
	assert(true, WindWest)
	assert(true, WindNorth)
	assert(false, windBase)
	assert(false, WindNorth+1)
}

func (s *ScoreTestSuite) TestIsValid(t *check.C) {
	assert := func(is_valid bool, valid Tile) {
		if is_valid != valid.IsValid() {
			t.Fatalf("%v is not correctly classified as valid", valid)
		}
	}

	assert(true, WindEast)
	assert(true, Balls1)
	assert(true, Balls9)
	assert(true, Season4)
	assert(false, ballsBase)
	assert(false, WindNorth+1)
	assert(false, Balls9+1)
	assert(false, Season4+1)
}

func (s *ScoreTestSuite) TestSetScores(t *check.C) {

	assert_set_score := func(expected_score, expected_doubles int, expected_type SetType, set *Set) {
		score, doubles, _ := ScoreSet(set, WIND_OWN, WIND_ROUND)

		if score == expected_score && doubles == expected_doubles && expected_type == set.setType {
			return
		}

		as_json, err := json.Marshal(set)
		if err != nil {
			t.Fatalf("Unable to marshall set: %v", err)
		}

		t.Errorf("Set: %s\n", as_json)
		t.Errorf("Expected score=%v doubles=%v set_type=%v, got score=%v doubles=%v set_type=%v\n",
			expected_score, expected_doubles, expected_type, score, doubles, set.setType)
	}

	// TODO: count score & doubles for flowers & seasons.

	assert_set_score(0, 0, NoSet, &Set{})

	// Simples: chow, pung, and kong
	assert_set_score(0, 0, Pillow, &Set{Tiles: []Tile{Bamboo5, Bamboo5}})
	assert_set_score(0, 0, NoSet, &Set{Tiles: []Tile{Bamboo5, Bamboo6}})
	assert_set_score(0, 0, Chow, &Set{Tiles: []Tile{Bamboo5, Bamboo6, Bamboo7}})
	assert_set_score(2, 0, Pung, &Set{Tiles: []Tile{Bamboo5, Bamboo5, Bamboo5}})
	assert_set_score(4, 0, Pung, &Set{Tiles: []Tile{Bamboo5, Bamboo5, Bamboo5}, Concealed: true})
	assert_set_score(8, 0, Kong, &Set{Tiles: []Tile{Bamboo5, Bamboo5, Bamboo5, Bamboo5}})
	assert_set_score(16, 0, Kong, &Set{Tiles: []Tile{Bamboo5, Bamboo5, Bamboo5, Bamboo5}, Concealed: true})

	// Terminals: chow, pung, and kong
	assert_set_score(0, 0, Pillow, &Set{Tiles: []Tile{Chars9, Chars9}})
	assert_set_score(4, 0, Pung, &Set{Tiles: []Tile{Chars9, Chars9, Chars9}})
	assert_set_score(8, 0, Pung, &Set{Tiles: []Tile{Chars9, Chars9, Chars9}, Concealed: true})
	assert_set_score(16, 0, Kong, &Set{Tiles: []Tile{Chars9, Chars9, Chars9, Chars9}})
	assert_set_score(32, 0, Kong, &Set{Tiles: []Tile{Chars9, Chars9, Chars9, Chars9}, Concealed: true})

	// Round winds
	assert_set_score(2, 0, Pillow, &Set{Tiles: []Tile{WindWest, WindWest}})
	assert_set_score(4, 1, Pung, &Set{Tiles: []Tile{WindWest, WindWest, WindWest}})
	assert_set_score(8, 1, Pung, &Set{Tiles: []Tile{WindWest, WindWest, WindWest}, Concealed: true})
	assert_set_score(16, 1, Kong, &Set{Tiles: []Tile{WindWest, WindWest, WindWest, WindWest}})
	assert_set_score(32, 1, Kong, &Set{Tiles: []Tile{WindWest, WindWest, WindWest, WindWest}, Concealed: true})

	// Own winds
	assert_set_score(2, 0, Pillow, &Set{Tiles: []Tile{WindNorth, WindNorth}})
	assert_set_score(4, 1, Pung, &Set{Tiles: []Tile{WindNorth, WindNorth, WindNorth}})
	assert_set_score(8, 1, Pung, &Set{Tiles: []Tile{WindNorth, WindNorth, WindNorth}, Concealed: true})
	assert_set_score(16, 1, Kong, &Set{Tiles: []Tile{WindNorth, WindNorth, WindNorth, WindNorth}})
	assert_set_score(32, 1, Kong, &Set{Tiles: []Tile{WindNorth, WindNorth, WindNorth, WindNorth}, Concealed: true})

	// Other winds
	assert_set_score(0, 0, Pillow, &Set{Tiles: []Tile{WindEast, WindEast}})
	assert_set_score(0, 0, NoSet, &Set{Tiles: []Tile{WindEast, Bamboo6}})
	assert_set_score(0, 0, NoSet, &Set{Tiles: []Tile{WindEast, WindSouth, WindWest}})
	assert_set_score(4, 0, Pung, &Set{Tiles: []Tile{WindEast, WindEast, WindEast}})
	assert_set_score(16, 0, Kong, &Set{Tiles: []Tile{WindEast, WindEast, WindEast, WindEast}})

	// Dragons: pillow, pung, and kong
	assert_set_score(0, 0, NoSet, &Set{Tiles: []Tile{DragonWhite}})
	assert_set_score(2, 0, Pillow, &Set{Tiles: []Tile{DragonGreen, DragonGreen}})
	assert_set_score(2, 0, Pillow, &Set{Tiles: []Tile{DragonWhite, DragonWhite}, Concealed: true})
	assert_set_score(4, 1, Pung, &Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite}})
	assert_set_score(8, 1, Pung, &Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite}, Concealed: true})
	assert_set_score(16, 1, Kong, &Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite, DragonWhite}})
	assert_set_score(32, 1, Kong, &Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite, DragonWhite}, Concealed: true})
}

func (s *ScoreTestSuite) TestScore(t *check.C) {

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
		Set{Tiles: []Tile{DragonRed, DragonRed}},
		Set{Tiles: []Tile{DragonGreen, DragonGreen}},
		Set{Tiles: []Tile{Balls1, Balls1, Balls1}},
	}}
	assert_score(8, hand)
	if hand.Winning {
		t.Error("Hand should have been recognised as non-winning.")
	}

	// Two chows and a pung of simples
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{Bamboo4, Bamboo5, Bamboo6}},
		Set{Tiles: []Tile{Balls1, Balls3, Balls2}},
		Set{Tiles: []Tile{Chars4, Chars4, Chars4}},
	}}
	assert_score(2, hand)
	if hand.Winning {
		t.Error("Hand should have been recognised as non-winning.")
	}

	// Kong of dragons and a concealed pung of simples
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonGreen, DragonGreen, DragonGreen, DragonGreen}},
		Set{Tiles: []Tile{Chars4, Chars4, Chars4}, Concealed: true},
	}}
	assert_score(40, hand)
	if hand.Winning {
		t.Error("Hand should have been recognised as non-winning.")
	}

	// Winning hand with no doubles.
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{Balls9, Balls9}},
		Set{Tiles: []Tile{Bamboo2, Bamboo3, Bamboo4}},
		Set{Tiles: []Tile{Balls5, Balls6, Balls7}},
		Set{Tiles: []Tile{Balls1, Balls1, Balls1, Balls1}},
		Set{Tiles: []Tile{Balls8, Balls8, Balls8}},
	}}
	assert_score(16+2+20, hand)
	if !hand.Winning {
		t.Error("Hand should have been recognised as winning.")
	}
	if full_flush(hand, 38) != 0 {
		t.Error("Hand should not be detected as full flush")
	}

	// Winning hand with full flush
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{Balls9, Balls9}},
		Set{Tiles: []Tile{Balls2, Balls3, Balls4}},
		Set{Tiles: []Tile{Balls5, Balls6, Balls7}},
		Set{Tiles: []Tile{Balls1, Balls1, Balls1, Balls1}},
		Set{Tiles: []Tile{Balls8, Balls8, Balls8}},
	}}
	if full_flush(hand, 38) == 0 {
		t.Error("Hand should be detected as full flush")
	}
	assert_score((16+2+20)*(1<<4), hand) // 608
	if !hand.Winning {
		t.Error("Hand should have been recognised as winning.")
	}

	// Winning hand with pure straight
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonGreen, DragonGreen}},
		Set{Tiles: []Tile{Balls7, Balls8, Balls9}},
		Set{Tiles: []Tile{Balls1, Balls3, Balls2}},
		Set{Tiles: []Tile{Balls5, Balls6, Balls4}},
		Set{Tiles: []Tile{Bamboo8, Bamboo8, Bamboo8}},
	}}
	assert_score((2+2+20)*2, hand) // 48
	if !hand.Winning {
		t.Error("Hand should have been recognised as winning.")
	}

	// Non-winning hand with pure straight
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonGreen}},
		Set{Tiles: []Tile{Balls7, Balls8, Balls9}},
		Set{Tiles: []Tile{Balls1, Balls3, Balls2}},
		Set{Tiles: []Tile{Balls5, Balls6, Balls4}},
		Set{Tiles: []Tile{Bamboo8, Bamboo8, Bamboo8}},
	}}
	assert_score(4, hand)
	if hand.Winning {
		t.Error("Hand should have been recognised as non-winning.")
	}

	// Winning hand with all pungs
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonGreen, DragonGreen}},
		Set{Tiles: []Tile{Balls7, Balls7, Balls7}},
		Set{Tiles: []Tile{Chars2, Chars2, Chars2}},
		Set{Tiles: []Tile{Balls5, Balls5, Balls5}},
		Set{Tiles: []Tile{Bamboo8, Bamboo8, Bamboo8}},
	}}
	assert_score((2+8+20)*2, hand) // 60
	if all_pungs(hand, 60) == 0 {
		t.Error("Hand should have been recognised as all pungs.")
	}
	if !hand.Winning {
		t.Error("Hand should have been recognised as winning.")
	}

	// Non-winning hand with three concealed pungs
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonGreen, DragonGreen}},
		Set{Tiles: []Tile{Balls7, Balls7, Balls7}, Concealed: true},
		Set{Tiles: []Tile{Chars2, Chars2, Chars2}, Concealed: true},
		Set{Tiles: []Tile{Balls5, Balls5, Balls5, Balls5}, Concealed: true},
		Set{Tiles: []Tile{Bamboo1, Bamboo2, Bamboo3}},
	}}
	assert_score((2+4+4+16+20)*2, hand) // 92
	if three_concealed_pungs(hand, 92) == 0 {
		t.Error("Hand should have been recognised as three concealed pungs.")
	}
	if !hand.Winning {
		t.Error("Hand should have been recognised as winning.")
	}

	// Winning chow hand
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{Balls2, Balls2}},
		Set{Tiles: []Tile{Balls1, Balls2, Balls3}},
		Set{Tiles: []Tile{Chars1, Chars2, Chars3}},
		Set{Tiles: []Tile{Balls5, Balls6, Balls7}},
		Set{Tiles: []Tile{Bamboo1, Bamboo2, Bamboo3}},
	}}
	assert_score(40, hand)
	if chow_hand(hand, 20) == 0 {
		t.Error("Hand should have been recognised as chow hand.")
	}
	if !hand.Winning {
		t.Error("Hand should have been recognised as winning.")
	}

	// Winning all-simples hand
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{Balls2, Balls2}},
		Set{Tiles: []Tile{Balls2, Balls3, Balls4}},
		Set{Tiles: []Tile{Chars2, Chars2, Chars2}},
		Set{Tiles: []Tile{Balls5, Balls6, Balls7}},
		Set{Tiles: []Tile{Bamboo2, Bamboo3, Bamboo4}},
	}}
	assert_score(22*2, hand)
	if all_simples(hand, 44) == 0 {
		t.Error("Hand should have been recognised as all simples.")
	}

	// None-winning all-simples hand
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{Chars6}},
		Set{Tiles: []Tile{Balls2}},
		Set{Tiles: []Tile{Balls2, Balls3, Balls4}},
		Set{Tiles: []Tile{Chars2, Chars2, Chars2}},
		Set{Tiles: []Tile{Balls5, Balls6, Balls7}},
		Set{Tiles: []Tile{Bamboo2, Bamboo3, Bamboo4}},
	}}
	assert_score(4, hand)
	if all_simples(hand, 4) == 0 {
		t.Error("Hand should have been recognised as all simples.")
	}

	// None-winning terminals & honours hand
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonGreen}},
		Set{Tiles: []Tile{WindEast}},
		Set{Tiles: []Tile{WindWest, WindWest, WindWest}},    // 4 + 1d
		Set{Tiles: []Tile{Chars1, Chars1, Chars1}},          // 4
		Set{Tiles: []Tile{WindNorth, WindNorth, WindNorth}}, // 4 + 1d
		Set{Tiles: []Tile{Bamboo9, Bamboo9, Bamboo9}},       // 4
	}}
	assert_score(16*8, hand)
	if all_terminals_honours(hand, 128) == 0 {
		t.Error("Hand should have been recognised as all terminals & honours.")
	}

	// Winning half-flush
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonGreen, DragonGreen}},     // 2
		Set{Tiles: []Tile{WindWest, WindWest, WindWest}}, // 4 + 1d
		Set{Tiles: []Tile{Chars1, Chars2, Chars3}},
		Set{Tiles: []Tile{Chars4, Chars5, Chars6}},
		Set{Tiles: []Tile{Chars2, Chars2, Chars2}}, // 2
	}}
	assert_score((20+2+4+2)*4, hand) // 56
	if half_flush(hand, 56) == 0 {
		t.Error("Hand should have been recognised as half-flush.")
	}

	// Non-winning half-flush
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonGreen}},
		Set{Tiles: []Tile{WindEast}},
		Set{Tiles: []Tile{WindWest, WindWest, WindWest}}, // 4 + 1d
		Set{Tiles: []Tile{Chars1, Chars2, Chars3}},
		Set{Tiles: []Tile{Chars4, Chars5, Chars6}},
		Set{Tiles: []Tile{Chars2, Chars2, Chars2}}, // 2
	}}
	assert_score(24, hand)
	if half_flush(hand, 24) == 0 {
		t.Error("Hand should have been recognised as half-flush.")
	}

	// Outside hand
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{WindEast, WindEast}},
		Set{Tiles: []Tile{WindWest, WindWest, WindWest}}, // 4 + 1d
		Set{Tiles: []Tile{Chars3, Chars2, Chars1}},
		Set{Tiles: []Tile{Bamboo9, Bamboo9, Bamboo9}}, // 4
		Set{Tiles: []Tile{Balls1, Balls1, Balls1}},    // 4
	}}
	assert_score((20+4+4+4)*4, hand) // 128
	if outside_hand(hand, 128) == 0 {
		t.Error("Hand should have been recognised as outside hand.")
	}
}

func (s *ScoreTestSuite) TestTile_Suit(t *check.C) {
	assert := func(expect_suit, tile Tile) {
		if tile.Suit() != expect_suit {
			t.Errorf("Tile %q doesn't have expected suit %q but %q",
				tile, expect_suit, tile.Suit())
		}
	}

	assert(bambooBase, Bamboo1)
	assert(bambooBase, Bamboo9)
	assert(ballsBase, Balls1)
	assert(ballsBase, Balls9)
	assert(charsBase, Chars1)
	assert(charsBase, Chars9)
	assert(NoTile, Bamboo9+1)
	assert(NoTile, Season3)
}
