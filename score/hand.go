package score

import "encoding/json"

type Tile int

const (
	NO_TILE Tile = 0

	// Balls 11-19
	BALLS_BASE Tile = 10
	BALLS_1    Tile = BALLS_BASE + 1
	BALLS_2    Tile = BALLS_BASE + 2
	BALLS_3    Tile = BALLS_BASE + 3
	BALLS_4    Tile = BALLS_BASE + 4
	BALLS_5    Tile = BALLS_BASE + 5
	BALLS_6    Tile = BALLS_BASE + 6
	BALLS_7    Tile = BALLS_BASE + 7
	BALLS_8    Tile = BALLS_BASE + 8
	BALLS_9    Tile = BALLS_BASE + 9

	// Characters 21-29
	CHARS_BASE Tile = 20
	CHARS_1    Tile = CHARS_BASE + 1
	CHARS_2    Tile = CHARS_BASE + 2
	CHARS_3    Tile = CHARS_BASE + 3
	CHARS_4    Tile = CHARS_BASE + 4
	CHARS_5    Tile = CHARS_BASE + 5
	CHARS_6    Tile = CHARS_BASE + 6
	CHARS_7    Tile = CHARS_BASE + 7
	CHARS_8    Tile = CHARS_BASE + 8
	CHARS_9    Tile = CHARS_BASE + 9

	// Bamboo 31-39
	BAMBOO_BASE Tile = 30
	BAMBOO_1    Tile = BAMBOO_BASE + 1
	BAMBOO_2    Tile = BAMBOO_BASE + 2
	BAMBOO_3    Tile = BAMBOO_BASE + 3
	BAMBOO_4    Tile = BAMBOO_BASE + 4
	BAMBOO_5    Tile = BAMBOO_BASE + 5
	BAMBOO_6    Tile = BAMBOO_BASE + 6
	BAMBOO_7    Tile = BAMBOO_BASE + 7
	BAMBOO_8    Tile = BAMBOO_BASE + 8
	BAMBOO_9    Tile = BAMBOO_BASE + 9

	MAY_CHOW_BELOW = BAMBOO_9 + 1

	// Winds 41-44
	WIND_BASE  Tile = 40
	WIND_EAST  Tile = WIND_BASE + 1
	WIND_SOUTH Tile = WIND_BASE + 2
	WIND_WEST  Tile = WIND_BASE + 3
	WIND_NORTH Tile = WIND_BASE + 4

	// Dragons 51-54
	DRAGON_BASE  Tile = 50
	DRAGON_RED   Tile = DRAGON_BASE + 1
	DRAGON_GREEN Tile = DRAGON_BASE + 2
	DRAGON_WHITE Tile = DRAGON_BASE + 3

	// Flowers 61-64
	FLOWER_BASE Tile = 60
	FLOWER_1    Tile = FLOWER_BASE + 1
	FLOWER_2    Tile = FLOWER_BASE + 2
	FLOWER_3    Tile = FLOWER_BASE + 3
	FLOWER_4    Tile = FLOWER_BASE + 4

	// Seasons 71-74
	SEASON_BASE Tile = 70
	SEASON_1    Tile = SEASON_BASE + 1
	SEASON_2    Tile = SEASON_BASE + 2
	SEASON_3    Tile = SEASON_BASE + 3
	SEASON_4    Tile = SEASON_BASE + 4
)

func (tile *Tile) IsValid() bool {
	if *tile < MAY_CHOW_BELOW {
		modulo := int(*tile) % 10
		return 1 <= modulo && modulo <= 9
	}

	return tile.IsHonour() || tile.IsFlower() || tile.IsSeason()
}

func (tile *Tile) MarshalJSON() ([]byte, error) {
	return json.Marshal(tile.String())
}

func (tile Tile) IsWind() bool {
	return tile > WIND_BASE && tile-WIND_BASE <= 4
}

func (tile Tile) IsDragon() bool {
	return tile > DRAGON_BASE && tile-DRAGON_BASE <= 3
}

func (tile Tile) IsFlower() bool {
	return tile > FLOWER_BASE && tile-FLOWER_BASE <= 4
}

func (tile Tile) IsSeason() bool {
	return tile > SEASON_BASE && tile-SEASON_BASE <= 4
}

func (tile Tile) IsHonour() bool {
	return tile.IsDragon() || tile.IsWind()
}

func (tile Tile) IsTerminal() bool {
	number := int(tile) % 10
	return tile < MAY_CHOW_BELOW && (number == 1 || number == 9)
}

func (tile Tile) IsSimple() bool {
	number := int(tile) % 10
	return tile < MAY_CHOW_BELOW && number > 1 && number < 9
}

func (tile Tile) Suit() Tile {
	if tile >= MAY_CHOW_BELOW {
		return NO_TILE
	}

	number := int(tile) % 10
	return Tile(int(tile) - number)
}

func (tile Tile) Number() int {
	if tile >= MAY_CHOW_BELOW {
		return 0
	}
	return int(tile) % 10
}

// implements sort.Interface for []Tile based on tile order.
type ByTileOrder []Tile

func (a ByTileOrder) Len() int           { return len(a) }
func (a ByTileOrder) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTileOrder) Less(i, j int) bool { return a[i] < a[j] }

// implements sort.Interface for []Set based on tile order.
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

type SetType int

const (
	NO_SET = 0
	PILLOW = 1
	CHOW   = 2
	PUNG   = 4
	KONG   = 8
)

// A set consists of one to four tiles.
type Set struct {
	Tiles     []Tile
	Concealed bool
	set_type  SetType
}

func (set *Set) HasTerminalOrHonour() bool {
	for _, tile := range set.Tiles {
		if tile.IsTerminal() || tile.IsHonour() {
			return true
		}
	}
	return false
}

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
