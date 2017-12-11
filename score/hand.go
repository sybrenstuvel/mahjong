package score

import (
	"encoding/json"
	"errors"
)

// Tile represents a single MJ tile
type Tile int

// This block contains all possible tiles in the MJ game.
// They are grouped by suit.
const (
	NoTile Tile = 0

	// Balls 11-19
	ballsBase Tile = 10
	Balls1    Tile = ballsBase + 1
	Balls2    Tile = ballsBase + 2
	Balls3    Tile = ballsBase + 3
	Balls4    Tile = ballsBase + 4
	Balls5    Tile = ballsBase + 5
	Balls6    Tile = ballsBase + 6
	Balls7    Tile = ballsBase + 7
	Balls8    Tile = ballsBase + 8
	Balls9    Tile = ballsBase + 9

	// Characters 21-29
	charsBase Tile = 20
	Chars1    Tile = charsBase + 1
	Chars2    Tile = charsBase + 2
	Chars3    Tile = charsBase + 3
	Chars4    Tile = charsBase + 4
	Chars5    Tile = charsBase + 5
	Chars6    Tile = charsBase + 6
	Chars7    Tile = charsBase + 7
	Chars8    Tile = charsBase + 8
	Chars9    Tile = charsBase + 9

	// Bamboo 31-39
	bambooBase Tile = 30
	Bamboo1    Tile = bambooBase + 1
	Bamboo2    Tile = bambooBase + 2
	Bamboo3    Tile = bambooBase + 3
	Bamboo4    Tile = bambooBase + 4
	Bamboo5    Tile = bambooBase + 5
	Bamboo6    Tile = bambooBase + 6
	Bamboo7    Tile = bambooBase + 7
	Bamboo8    Tile = bambooBase + 8
	Bamboo9    Tile = bambooBase + 9

	mayChowBelow = Bamboo9 + 1

	// Winds 41-44
	windBase  Tile = 40
	WindEast  Tile = windBase + 1
	WindSouth Tile = windBase + 2
	WindWest  Tile = windBase + 3
	WindNorth Tile = windBase + 4

	// Dragons 51-54
	dragonBase  Tile = 50
	DragonRed   Tile = dragonBase + 1
	DragonGreen Tile = dragonBase + 2
	DragonWhite Tile = dragonBase + 3

	// Flowers 61-64
	flowerBase Tile = 60
	Flower1    Tile = flowerBase + 1
	Flower2    Tile = flowerBase + 2
	Flower3    Tile = flowerBase + 3
	Flower4    Tile = flowerBase + 4

	// Seasons 71-74
	seasonBase Tile = 70
	Season1    Tile = seasonBase + 1
	Season2    Tile = seasonBase + 2
	Season3    Tile = seasonBase + 3
	Season4    Tile = seasonBase + 4
)

// Standard errors.
var (
	ErrTileNotValid = errors.New("tile not valid")
)

// IsValid returns true if the number of the tile represents a valid tile.
func (tile Tile) IsValid() bool {
	if tile < mayChowBelow {
		modulo := int(tile) % 10
		return 1 <= modulo && modulo <= 9
	}

	return tile.IsHonour() || tile.IsFlower() || tile.IsSeason()
}

// MarshalJSON converts a tile to JSON.
func (tile Tile) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(tile))
}

// UnmarshalJSON reads a tile from JSON
func (tile *Tile) UnmarshalJSON(data []byte) error {
	tilenr := 0
	if err := json.Unmarshal(data, &tilenr); err != nil {
		return err
	}
	*tile = Tile(tilenr)
	if !tile.IsValid() {
		return ErrTileNotValid
	}
	return nil
}

// IsWind returns true for wind tiles.
func (tile Tile) IsWind() bool {
	return tile > windBase && tile-windBase <= 4
}

// IsDragon returns true for dragon tiles.
func (tile Tile) IsDragon() bool {
	return tile > dragonBase && tile-dragonBase <= 3
}

// IsFlower returns true for flower tiles.
func (tile Tile) IsFlower() bool {
	return tile > flowerBase && tile-flowerBase <= 4
}

// IsSeason returns true for season tiles.
func (tile Tile) IsSeason() bool {
	return tile > seasonBase && tile-seasonBase <= 4
}

// IsHonour returns true for honour tiles.
func (tile Tile) IsHonour() bool {
	return tile.IsDragon() || tile.IsWind()
}

// IsTerminal returns true for terminal tiles.
func (tile Tile) IsTerminal() bool {
	number := int(tile) % 10
	return tile < mayChowBelow && (number == 1 || number == 9)
}

// IsSimple returns true for simple tiles.
func (tile Tile) IsSimple() bool {
	number := int(tile) % 10
	return tile < mayChowBelow && number > 1 && number < 9
}

// Suit returns the suit base number of the tile.
// By itself it's rather worthless, but it does allow comparison of suits between tiles.
func (tile Tile) Suit() Tile {
	if tile >= mayChowBelow {
		return NoTile
	}

	number := int(tile) % 10
	return Tile(int(tile) - number)
}

// Number returns the number of the tile (like 9 for a bamboo-9),
// or 0 if the tile has no number (like a red dragon or flower).
// A chow consists of sequential tiles for which this function
// returns a non-zero number.
func (tile Tile) Number() int {
	if tile >= mayChowBelow {
		return 0
	}
	return int(tile) % 10
}

// ByTileOrder implements sort.Interface for []Tile based on tile order.
type ByTileOrder []Tile

func (a ByTileOrder) Len() int           { return len(a) }
func (a ByTileOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTileOrder) Less(i, j int) bool { return a[i] < a[j] }

// SortSetsByTileOrder implements sort.Interface for []Set based on tile order.
type SortSetsByTileOrder []Set

func (a SortSetsByTileOrder) Len() int      { return len(a) }
func (a SortSetsByTileOrder) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SortSetsByTileOrder) Less(i, j int) bool {
	// TODO: sort properly instead of hard-coded looking at idx 0 and 1.
	if len(a[j].Tiles) == 0 {
		return true
	}
	if len(a[i].Tiles) == 0 {
		return false
	}

	if a[i].Tiles[0] == a[j].Tiles[0] {
		if len(a[j].Tiles) == 1 {
			return true
		}
		if len(a[i].Tiles) == 1 {
			return false
		}
		return a[i].Tiles[1] < a[j].Tiles[1]
	}

	return a[i].Tiles[0] < a[j].Tiles[0]
}

// SetType represents a type of set, see the constants defined below.
type SetType int

// Constants indicating the type of set.
const (
	NoSet  SetType = 0
	Pillow SetType = 1
	Chow   SetType = 2
	Pung   SetType = 4
	Kong   SetType = 8
)

// Set consists of one to four tiles.
type Set struct {
	Tiles     []Tile
	Concealed bool
	setType   SetType
}

// HasTerminalOrHonour returns true if the set contains a terminal or an honour.
func (set *Set) HasTerminalOrHonour() bool {
	for _, tile := range set.Tiles {
		if tile.IsTerminal() || tile.IsHonour() {
			return true
		}
	}
	return false
}

// Hand represents a hand (which may be non-winning) and consistst of sets and win conditions.
type Hand struct {
	Sets                 []Set
	WindOwn              Tile
	WindRound            Tile
	LastChance           bool
	WinSelfDrawn         bool
	WinOnReplacementTile bool
	LastTileOfWall       bool
	RobbedTheKong        bool
	OutInDraw            bool
	Winning              bool
}
