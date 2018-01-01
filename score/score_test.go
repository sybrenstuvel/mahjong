package score

import (
	"encoding/json"

	check "gopkg.in/check.v1"
)

const (
	windOwn   = WindNorth
	windRound = WindWest
)

type ScoreTestSuite struct{}

var _ = check.Suite(&ScoreTestSuite{})

func (s *ScoreTestSuite) TestTileValid(t *check.C) {
	assertTileValid := func(expectedValidity bool, tile Tile) {
		validity := tile.IsValid()

		if validity != expectedValidity {
			t.Errorf("Tile %q: expected validity=%v, got validity=%v\n",
				tile, expectedValidity, validity)
		}
	}

	assertTileValid(false, ballsBase)
	assertTileValid(false, -1)
	assertTileValid(false, 0xffff)
	assertTileValid(true, Balls3)
	assertTileValid(true, 13)
	assertTileValid(false, Chars1-1)
	assertTileValid(false, Chars9+1)
	assertTileValid(true, Chars9-1)
	assertTileValid(true, Bamboo1+1)
	assertTileValid(false, Bamboo1-1)
}

func (s *ScoreTestSuite) TestSetValid(t *check.C) {
	assertSetValid := func(expectedValidity, expectedChow bool, set Set) {
		validity, isChow := set.IsValid()

		if validity == expectedValidity && isChow == expectedChow {
			return
		}

		asJSON, err := json.Marshal(set)
		if err != nil {
			t.Fatalf("Unable to marshall set: %v", err)
		}

		t.Errorf("Set: %s\n", asJSON)
		t.Errorf("Expected validity=%v is_chow=%v, got validity=%v is_chow=%v\n",
			expectedValidity, expectedChow, validity, isChow)
	}

	assertSetValid(false, false, Set{})
	assertSetValid(false, false, Set{Tiles: []Tile{Balls3, Balls4}, Concealed: false})
	assertSetValid(false, false, Set{Tiles: []Tile{ballsBase}, Concealed: false})
	assertSetValid(false, false, Set{Tiles: []Tile{Balls3}, Concealed: false})
	assertSetValid(false, false, Set{Tiles: []Tile{1, 2}})
	assertSetValid(false, false, Set{Tiles: []Tile{DragonRed, DragonWhite}, Concealed: true})
	assertSetValid(false, false, Set{Tiles: []Tile{DragonRed, DragonRed, DragonWhite, DragonGreen}, Concealed: true})
	assertSetValid(false, false, Set{Tiles: []Tile{DragonRed, DragonWhite, DragonGreen}, Concealed: true})
	assertSetValid(false, false, Set{Tiles: []Tile{Balls8, Balls9, Balls9 + 1}})

	assertSetValid(true, false, Set{Tiles: []Tile{Balls3, Balls3}})
	assertSetValid(true, false, Set{Tiles: []Tile{Balls3, Balls3, Balls3}})
	assertSetValid(true, true, Set{Tiles: []Tile{Balls3, Balls4, Balls5}})
	assertSetValid(true, true, Set{Tiles: []Tile{Balls3, Balls5, Balls4}})
	assertSetValid(true, false, Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite}, Concealed: true})
	assertSetValid(true, false, Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite, DragonWhite}, Concealed: true})

	assertSetValid(true, false, Set{Tiles: []Tile{Balls9, Balls9}})
	assertSetValid(true, true, Set{Tiles: []Tile{Balls2, Balls3, Balls4}})
	assertSetValid(true, true, Set{Tiles: []Tile{Balls5, Balls6, Balls7}})
	assertSetValid(true, false, Set{Tiles: []Tile{Balls1, Balls1, Balls1, Balls1}})
	assertSetValid(true, false, Set{Tiles: []Tile{Balls8, Balls8, Balls8}})
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

	assertSetScore := func(expected_score, expected_doubles int, expected_type SetType, set *Set) {
		score, doubles, _ := set.Score(windOwn, windRound)

		if score == expected_score && doubles == expected_doubles && expected_type == set.setType {
			return
		}

		asJSON, err := json.Marshal(set)
		if err != nil {
			t.Fatalf("Unable to marshall set: %v", err)
		}

		t.Errorf("Set: %s\n", asJSON)
		t.Errorf("Expected score=%v doubles=%v set_type=%v, got score=%v doubles=%v set_type=%v\n",
			expected_score, expected_doubles, expected_type, score, doubles, set.setType)
	}

	// TODO: count score & doubles for flowers & seasons.

	assertSetScore(0, 0, NoSet, &Set{})

	// Simples: chow, pung, and kong
	assertSetScore(0, 0, Pillow, &Set{Tiles: []Tile{Bamboo5, Bamboo5}})
	assertSetScore(0, 0, NoSet, &Set{Tiles: []Tile{Bamboo5, Bamboo6}})
	assertSetScore(0, 0, Chow, &Set{Tiles: []Tile{Bamboo5, Bamboo6, Bamboo7}})
	assertSetScore(2, 0, Pung, &Set{Tiles: []Tile{Bamboo5, Bamboo5, Bamboo5}})
	assertSetScore(4, 0, Pung, &Set{Tiles: []Tile{Bamboo5, Bamboo5, Bamboo5}, Concealed: true})
	assertSetScore(8, 0, Kong, &Set{Tiles: []Tile{Bamboo5, Bamboo5, Bamboo5, Bamboo5}})
	assertSetScore(16, 0, Kong, &Set{Tiles: []Tile{Bamboo5, Bamboo5, Bamboo5, Bamboo5}, Concealed: true})

	// Terminals: chow, pung, and kong
	assertSetScore(0, 0, Pillow, &Set{Tiles: []Tile{Chars9, Chars9}})
	assertSetScore(4, 0, Pung, &Set{Tiles: []Tile{Chars9, Chars9, Chars9}})
	assertSetScore(8, 0, Pung, &Set{Tiles: []Tile{Chars9, Chars9, Chars9}, Concealed: true})
	assertSetScore(16, 0, Kong, &Set{Tiles: []Tile{Chars9, Chars9, Chars9, Chars9}})
	assertSetScore(32, 0, Kong, &Set{Tiles: []Tile{Chars9, Chars9, Chars9, Chars9}, Concealed: true})

	// Round winds
	assertSetScore(2, 0, Pillow, &Set{Tiles: []Tile{WindWest, WindWest}})
	assertSetScore(4, 1, Pung, &Set{Tiles: []Tile{WindWest, WindWest, WindWest}})
	assertSetScore(8, 1, Pung, &Set{Tiles: []Tile{WindWest, WindWest, WindWest}, Concealed: true})
	assertSetScore(16, 1, Kong, &Set{Tiles: []Tile{WindWest, WindWest, WindWest, WindWest}})
	assertSetScore(32, 1, Kong, &Set{Tiles: []Tile{WindWest, WindWest, WindWest, WindWest}, Concealed: true})

	// Own winds
	assertSetScore(2, 0, Pillow, &Set{Tiles: []Tile{WindNorth, WindNorth}})
	assertSetScore(4, 1, Pung, &Set{Tiles: []Tile{WindNorth, WindNorth, WindNorth}})
	assertSetScore(8, 1, Pung, &Set{Tiles: []Tile{WindNorth, WindNorth, WindNorth}, Concealed: true})
	assertSetScore(16, 1, Kong, &Set{Tiles: []Tile{WindNorth, WindNorth, WindNorth, WindNorth}})
	assertSetScore(32, 1, Kong, &Set{Tiles: []Tile{WindNorth, WindNorth, WindNorth, WindNorth}, Concealed: true})

	// Other winds
	assertSetScore(0, 0, Pillow, &Set{Tiles: []Tile{WindEast, WindEast}})
	assertSetScore(0, 0, NoSet, &Set{Tiles: []Tile{WindEast, Bamboo6}})
	assertSetScore(0, 0, NoSet, &Set{Tiles: []Tile{WindEast, WindSouth, WindWest}})
	assertSetScore(4, 0, Pung, &Set{Tiles: []Tile{WindEast, WindEast, WindEast}})
	assertSetScore(16, 0, Kong, &Set{Tiles: []Tile{WindEast, WindEast, WindEast, WindEast}})

	// Dragons: pillow, pung, and kong
	assertSetScore(0, 0, NoSet, &Set{Tiles: []Tile{DragonWhite}})
	assertSetScore(2, 0, Pillow, &Set{Tiles: []Tile{DragonGreen, DragonGreen}})
	assertSetScore(2, 0, Pillow, &Set{Tiles: []Tile{DragonWhite, DragonWhite}, Concealed: true})
	assertSetScore(4, 1, Pung, &Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite}})
	assertSetScore(8, 1, Pung, &Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite}, Concealed: true})
	assertSetScore(16, 1, Kong, &Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite, DragonWhite}})
	assertSetScore(32, 1, Kong, &Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite, DragonWhite}, Concealed: true})
}

func assertScore(t *check.C, expectedScore int, hand *Hand) {
	hand.WindRound = windRound
	hand.WindOwn = windOwn
	score := Score(hand)

	if score == expectedScore {
		return
	}

	asJSON, err := json.Marshal(hand)
	if err != nil {
		t.Fatalf("Unable to marshall hand: %v", err)
	}

	t.Errorf("Hand: %s\n", asJSON)
	t.Errorf("Expected score=%v, got score=%v\n", expectedScore, score)
}

func (s *ScoreTestSuite) TestScoreEmptyHand(t *check.C) {
	hand := &Hand{}
	assertScore(t, 0, hand)
	if hand.Winning {
		t.Error("Hand should have been recognised as non-winning.")
	}
}

func (s *ScoreTestSuite) TestScoreTwoPairsPung(t *check.C) {
	// Two pairs of dragons and a pung of balls 1
	hand := &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonRed, DragonRed}},
		Set{Tiles: []Tile{DragonGreen, DragonGreen}},
		Set{Tiles: []Tile{Balls1, Balls1, Balls1}},
	}}
	assertScore(t, 8, hand)
	if hand.Winning {
		t.Error("Hand should have been recognised as non-winning.")
	}
}

func (s *ScoreTestSuite) TestScoreTwoChowsPungSimples(t *check.C) {
	// Two chows and a pung of simples
	hand := &Hand{Sets: []Set{
		Set{Tiles: []Tile{Bamboo4, Bamboo5, Bamboo6}},
		Set{Tiles: []Tile{Balls1, Balls3, Balls2}},
		Set{Tiles: []Tile{Chars4, Chars4, Chars4}},
	}}
	assertScore(t, 2, hand)
	if hand.Winning {
		t.Error("Hand should have been recognised as non-winning.")
	}
}

func (s *ScoreTestSuite) TestKongDragonConcealedPung(t *check.C) {
	// Kong of dragons and a concealed pung of simples
	hand := &Hand{Sets: []Set{
		Set{Tiles: []Tile{DragonGreen, DragonGreen, DragonGreen, DragonGreen}},
		Set{Tiles: []Tile{Chars4, Chars4, Chars4}, Concealed: true},
	}}
	assertScore(t, 40, hand)
	if hand.Winning {
		t.Error("Hand should have been recognised as non-winning.")
	}
}

func (s *ScoreTestSuite) TestScoreChowConcealedPung(t *check.C) {
	// Chow and concealed pung of dragons
	hand := &Hand{Sets: []Set{
		Set{Tiles: []Tile{Balls1, Balls2, Balls3}, Concealed: false},
		Set{Tiles: []Tile{DragonGreen, DragonGreen, DragonGreen, DragonGreen}, Concealed: true},
	}}
	assertScore(t, 64, hand)
	if hand.Winning {
		t.Error("Hand should have been recognised as non-winning.")
	}
}

func (s *ScoreTestSuite) TestScoreWinningHandNoDoubles(t *check.C) {
	// Winning hand with no doubles.
	hand := &Hand{Sets: []Set{
		Set{Tiles: []Tile{Balls9, Balls9}},
		Set{Tiles: []Tile{Bamboo2, Bamboo3, Bamboo4}},
		Set{Tiles: []Tile{Balls5, Balls6, Balls7}},
		Set{Tiles: []Tile{Balls1, Balls1, Balls1, Balls1}},
		Set{Tiles: []Tile{Balls8, Balls8, Balls8}},
	}}
	assertScore(t, 16+2+20, hand)
	if !hand.Winning {
		t.Error("Hand should have been recognised as winning.")
	}
	if fullFlush(hand, 38) != 0 {
		t.Error("Hand should not be detected as full flush")
	}
}

func (s *ScoreTestSuite) TestScoreTODOSplit(t *check.C) {
	var hand *Hand
	// Winning hand with full flush
	hand = &Hand{Sets: []Set{
		Set{Tiles: []Tile{Balls9, Balls9}},
		Set{Tiles: []Tile{Balls2, Balls3, Balls4}},
		Set{Tiles: []Tile{Balls5, Balls6, Balls7}},
		Set{Tiles: []Tile{Balls1, Balls1, Balls1, Balls1}},
		Set{Tiles: []Tile{Balls8, Balls8, Balls8}},
	}}
	if fullFlush(hand, 38) == 0 {
		t.Error("Hand should be detected as full flush")
	}
	assertScore(t, (16+2+20)*(1<<4), hand) // 608
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
	assertScore(t, (2+2+20)*2, hand) // 48
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
	assertScore(t, 4, hand)
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
	assertScore(t, (2+8+20)*2, hand) // 60
	if allPungs(hand, 60) == 0 {
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
	assertScore(t, (2+4+4+16+20)*2, hand) // 92
	if threeConcealedPungs(hand, 92) == 0 {
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
	assertScore(t, 40, hand)
	if chowHand(hand, 20) == 0 {
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
	assertScore(t, 22*2, hand)
	if allSimples(hand, 44) == 0 {
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
	assertScore(t, 4, hand)
	if allSimples(hand, 4) == 0 {
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
	assertScore(t, 16*8, hand)
	if allTerminalsHonours(hand, 128) == 0 {
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
	assertScore(t, (20+2+4+2)*4, hand) // 56
	if halfFlush(hand, 56) == 0 {
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
	assertScore(t, 24, hand)
	if halfFlush(hand, 24) == 0 {
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
	assertScore(t, (20+4+4+4)*4, hand) // 128
	if outsideHand(hand, 128) == 0 {
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
