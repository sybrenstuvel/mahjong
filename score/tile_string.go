// Code generated by "stringer -type=Tile"; DO NOT EDIT

package score

import "fmt"

const _Tile_name = "BALLS_BASEBALLS_1BALLS_2BALLS_3BALLS_4BALLS_5BALLS_6BALLS_7BALLS_8BALLS_9CHARS_BASECHARS_1CHARS_2CHARS_3CHARS_4CHARS_5CHARS_6CHARS_7CHARS_8CHARS_9BAMBOO_BASEBAMBOO_1BAMBOO_2BAMBOO_3BAMBOO_4BAMBOO_5BAMBOO_6BAMBOO_7BAMBOO_8BAMBOO_9WIND_BASEWIND_EASTWIND_SOUTHWIND_WESTWIND_NORTHDRAGON_BASEDRAGON_REDDRAGON_GREENDRAGON_WHITEFLOWER_BASEFLOWER_1FLOWER_2FLOWER_3FLOWER_4SEASON_BASESEASON_1SEASON_2SEASON_3SEASON_4"

var _Tile_index = [...]uint16{0, 10, 17, 24, 31, 38, 45, 52, 59, 66, 73, 83, 90, 97, 104, 111, 118, 125, 132, 139, 146, 157, 165, 173, 181, 189, 197, 205, 213, 221, 229, 238, 247, 257, 266, 276, 287, 297, 309, 321, 332, 340, 348, 356, 364, 375, 383, 391, 399, 407}

func (i Tile) String() string {
	i -= 10
	if i < 0 || i >= Tile(len(_Tile_index)-1) {
		return fmt.Sprintf("Tile(%d)", i+10)
	}
	return _Tile_name[_Tile_index[i]:_Tile_index[i+1]]
}
