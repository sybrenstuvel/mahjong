package score

import "sort"

func IsValidTile(tile Tile) bool {
    if BALLS_1   <= tile && tile <= BALLS_9 { return true }
    if CHARS_1   <= tile && tile <= CHARS_9 { return true }
    if BAMBOO_1  <= tile && tile <= BAMBOO_9 { return true }
    if WIND_EAST <= tile && tile <= SEASON_4 { return true }

    return false;
}


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
        if ! IsValidTile(tile) {
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

// Returns (score, doubles)
func ScoreSet(set *Set, wind_own, wind_round Tile) (int, int) {
    // TODO: count score & doubles for flowers & seasons.
    if len(set.Tiles) < 2 { return 0, 0 }

    is_valid, is_chow := IsValidSet(*set)
    if !is_valid || is_chow { return 0, 0 }

    // If we're here, we know it's a pillow/pung/kong, so the
    // length and first tile determine the score.
    tile := set.Tiles[0]
    scoring_wind := tile == wind_own || tile == wind_round

    var multiplier int
    if set.Concealed {
        multiplier = 2
    } else {
        multiplier = 1
    }
    // A kong scores 4 times as much as a pung.
    if len(set.Tiles) == 4 {
        multiplier *= 4
    }

    switch len(set.Tiles) {
    case 2:
        switch {
        case scoring_wind: return 2,0
        case tile.IsDragon(): return 2, 0
        default: return 0, 0
        }
    case 3, 4:
        switch {
        case tile.IsTerminal(): return 4 * multiplier, 0
        case tile.IsWind():
            var double int
            if scoring_wind {
                double = 1
            } else {
                double = 0
            }
            return 4 * multiplier, double
        case tile.IsDragon(): return 4 * multiplier, 1
        default: return 2 * multiplier, 0
        }
    }

    return 0, 0
}

// Returns the score for the given hand.
func Score(hand Hand) int {
    return 0
}
