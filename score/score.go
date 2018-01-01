package score

import (
	"sort"

	log "github.com/sirupsen/logrus"
)

// IsValid returns (is valid, is chow) for this set.
func (set *Set) IsValid() (bool, bool) {
	if len(set.Tiles) < 2 || len(set.Tiles) > 4 {
		return false, false
	}

	// Sorting is needed to determine chows.
	sort.Sort(ByTileOrder(set.Tiles))

	firstTile := set.Tiles[0]
	mayChow := firstTile < mayChowBelow

	var isSame bool
	var isChow bool

	for idx, tile := range set.Tiles {
		if !tile.IsValid() {
			return false, false
		}

		// Either all tiles should be the same, or sequential.
		isSame = tile == firstTile
		isChow = mayChow && int(tile) == int(firstTile)+idx

		if !isSame && !isChow {
			return false, false
		}
	}

	if isChow {
		valid := len(set.Tiles) == 3
		return valid, valid
	}

	return true, false
}

// Score returns (basic score, doubles, valid) for the given set.
// The returned 'basic' score is not yet multiplied by the doubles.
func (set *Set) Score(windOwn, windRound Tile) (int, int, bool) {
	// TODO: count score & doubles for flowers & seasons.
	if len(set.Tiles) < 2 {
		set.setType = NoSet
		return 0, 0, false
	}

	isValid, isChow := set.IsValid()
	if !isValid {
		set.setType = NoSet
		return 0, 0, false
	}
	if isChow {
		set.setType = Chow
		return 0, 0, true
	}

	// If we're here, we know it's a pillow/pung/kong, so the
	// length and first tile determine the score.
	tile := set.Tiles[0]
	scoringWind := tile == windOwn || tile == windRound

	// A concealed pung/kong scores double
	var multiplier int
	if set.Concealed {
		multiplier = 2
	} else {
		multiplier = 1
	}

	// A kong scores 4 times as much as a pung.
	if len(set.Tiles) == 4 {
		multiplier *= 4
		set.setType = Kong
	} else {
		set.setType = Pung
	}

	switch len(set.Tiles) {
	case 2:
		set.setType = Pillow
		switch {
		case scoringWind:
			return 2, 0, true
		case tile.IsDragon():
			return 2, 0, true
		default:
			return 0, 0, true
		}
	case 3, 4:
		switch {
		case tile.IsTerminal():
			return 4 * multiplier, 0, true
		case tile.IsWind():
			var double int
			if scoringWind {
				double = 1
			} else {
				double = 0
			}
			return 4 * multiplier, double, true
		case tile.IsDragon():
			return 4 * multiplier, 1, true
		default:
			return 2 * multiplier, 0, true
		}
	}

	panic("Impossible situation turned out to be possible after all.")
}

// Score calculates the score for the given hand.
func Score(hand *Hand) int {
	totalScore := 0
	totalDoubles := 0
	nrOfPungs := 0
	nrOfPillows := 0

	log.WithField("hand", hand).Debug("calculating hand score")

	// Sorting the sets makes it easier to detect pure straights, nine gates and others.
	sort.Sort(SortSetsByTileOrder(hand.Sets))

	// Start by summing up the tile set scores.
	for idx := range hand.Sets {
		set := &hand.Sets[idx]
		setScore, setDoubles, isValid := set.Score(hand.WindOwn, hand.WindRound)
		log.WithFields(log.Fields{
			"set-idx": idx,
			"score":   setScore,
			"doubles": setDoubles,
			"valid":   isValid,
		}).Debug("set score calculated")
		if !isValid {
			continue
		}

		switch len(set.Tiles) {
		case 2:
			nrOfPillows++
		case 3, 4:
			nrOfPungs++
		}

		totalScore += setScore
		totalDoubles += setDoubles
	}

	// Detect winning hand
	if nrOfPungs == 4 && nrOfPillows == 1 {
		totalScore += 20
		hand.Winning = true
	} else {
		hand.Winning = false
	}

	// Count doubles
	for label, detector := range detectors {
		doubles := detector(hand, totalScore)
		log.WithFields(log.Fields{
			"detector": label,
			"doubles":  doubles,
		}).Debug("ran detector")
		totalDoubles += doubles
	}

	finalScore := totalScore * 1 << uint(totalDoubles)
	log.WithFields(log.Fields{
		"tile-score": totalScore,
		"doubles":    totalDoubles,
		"score":      finalScore,
	}).Debug("hand score calculated")

	return finalScore
}
