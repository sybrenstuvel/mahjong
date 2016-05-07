package score

import (
    "sort"
)


// Returns (is valid, is chow)
func IsValidSet(set Set) (bool, bool) {
    if len(set.Tiles) < 2 || len(set.Tiles) > 4 {
        return false, false
    }

    // Sorting is needed to determine chows.
    sort.Sort(ByTileOrder(set.Tiles))

    first_tile := set.Tiles[0]
    var may_chow bool = first_tile < MAY_CHOW_BELOW

    var is_same bool
    var is_chow bool

    for idx, tile := range set.Tiles {
        if ! tile.IsValid() {
            return false, false
        }

        // Either all tiles should be the same, or sequential.
        is_same = tile == first_tile
        is_chow = may_chow && int(tile) == int(first_tile) + idx

        if !is_same && !is_chow { return false, false }
    }

    if is_chow {
        valid := len(set.Tiles) == 3
        return valid, valid
    }

    return true, false
}

// Returns (score, doubles, valid)
func ScoreSet(set *Set, wind_own, wind_round Tile) (int, int, bool) {
    // TODO: count score & doubles for flowers & seasons.
    if len(set.Tiles) < 2 {
        set.set_type = NO_SET
        return 0, 0, false
    }

    is_valid, is_chow := IsValidSet(*set)
    if !is_valid {
        set.set_type = NO_SET
        return 0, 0, is_valid
    }
    if is_chow {
        set.set_type = CHOW
        return 0, 0, is_valid
    }

    // If we're here, we know it's a pillow/pung/kong, so the
    // length and first tile determine the score.
    tile := set.Tiles[0]
    scoring_wind := tile == wind_own || tile == wind_round

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
        set.set_type = KONG
    } else {
        set.set_type = PUNG
    }

    switch len(set.Tiles) {
    case 2:
        set.set_type = PILLOW
        switch {
        case scoring_wind: return 2,0, true
        case tile.IsDragon(): return 2, 0, true
        default: return 0, 0, true
        }
    case 3, 4:
        switch {
        case tile.IsTerminal(): return 4 * multiplier, 0, true
        case tile.IsWind():
            var double int
            if scoring_wind {
                double = 1
            } else {
                double = 0
            }
            return 4 * multiplier, double, true
        case tile.IsDragon(): return 4 * multiplier, 1, true
        default: return 2 * multiplier, 0, true
        }
    }

    panic("Impossible situation turned out to be possible after all.")
}


// Returns the score for the given hand.
func Score(hand *Hand) int {
    total_score := 0
    total_doubles := 0
    nr_of_pungs := 0
    nr_of_pillows := 0

    // Sorting the sets makes it easier to detect pure straights, nine gates and others.
    sort.Sort(SortSetsByTileOrder(hand.Sets))

    // Start by summing up the tile set scores.
    for idx, _ := range hand.Sets {
        set := &hand.Sets[idx]
        set_score, set_doubles, is_valid := ScoreSet(set, hand.WindOwn, hand.WindRound)

        if !is_valid { continue; }

        switch len(set.Tiles) {
        case 2:
            nr_of_pillows += 1
        case 3, 4:
            nr_of_pungs += 1
        }

        total_score += set_score
        total_doubles += set_doubles
    }

    // Detect winning hand
    if nr_of_pungs == 4 && nr_of_pillows == 1 {
        total_score += 20
        hand.Winning = true
    } else {
        hand.Winning = false
    }

    // Count doubles
    for _, detector := range detectors {
        total_doubles += detector(hand, total_score)
    }

    return total_score * 1 << uint(total_doubles)
}
