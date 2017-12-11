package score

import (
	"encoding/json"

	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

type HandTestSuite struct{}

var _ = check.Suite(&HandTestSuite{})

func (s *HandTestSuite) TestTileValues(c *check.C) {
	// Test some values of tiles
	assert.Equal(c, 12, int(Balls2))
	assert.Equal(c, 27, int(Chars7))
	assert.Equal(c, 39, int(Bamboo9))
	assert.Equal(c, 74, int(Season4))
}

func (s *HandTestSuite) TestTileJSON(c *check.C) {
	asJSON, err := json.Marshal(Tile(Balls4))
	assert.Nil(c, err)
	assert.Equal(c, "14", string(asJSON))

	var loadedTile Tile
	err = json.Unmarshal(asJSON, &loadedTile)
	assert.Nil(c, err)
	assert.Equal(c, Balls4, loadedTile)
}

func (s *HandTestSuite) TestSetJSON(c *check.C) {
	assertTwoWay := func(set Set) {
		asJSON, err := json.Marshal(set)
		assert.Nil(c, err)

		var loadedSet Set
		err = json.Unmarshal(asJSON, &loadedSet)
		assert.Nil(c, err)
		assert.Equal(c, set, loadedSet)
	}

	assertTwoWay(Set{})
	assertTwoWay(Set{Tiles: []Tile{Balls3, Balls4}, Concealed: false})
	assertTwoWay(Set{Tiles: []Tile{Balls3}, Concealed: false})
	assertTwoWay(Set{Tiles: []Tile{DragonRed, DragonWhite}, Concealed: true})
	assertTwoWay(Set{Tiles: []Tile{DragonRed, DragonRed, DragonWhite, DragonGreen}, Concealed: true})
	assertTwoWay(Set{Tiles: []Tile{DragonRed, DragonWhite, DragonGreen}, Concealed: true})

	assertTwoWay(Set{Tiles: []Tile{Balls3, Balls3}})
	assertTwoWay(Set{Tiles: []Tile{Balls3, Balls3, Balls3}})
	assertTwoWay(Set{Tiles: []Tile{Balls3, Balls4, Balls5}})
	assertTwoWay(Set{Tiles: []Tile{Balls3, Balls5, Balls4}})
	assertTwoWay(Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite}, Concealed: true})
	assertTwoWay(Set{Tiles: []Tile{DragonWhite, DragonWhite, DragonWhite, DragonWhite}, Concealed: true})

	assertTwoWay(Set{Tiles: []Tile{Balls9, Balls9}})
	assertTwoWay(Set{Tiles: []Tile{Balls2, Balls3, Balls4}})
	assertTwoWay(Set{Tiles: []Tile{Balls5, Balls6, Balls7}})
	assertTwoWay(Set{Tiles: []Tile{Balls1, Balls1, Balls1, Balls1}})
	assertTwoWay(Set{Tiles: []Tile{Balls8, Balls8, Balls8}})

	assertMarshalTileNotValid := func(set Set) {
		_, err := json.Marshal(set)
		if err == nil {
			c.Fatal("expected an error, did not find one")
		}
		assert.Contains(c, err.Error(), ErrTileNotValid.Error())
	}
	assertMarshalTileNotValid(Set{Tiles: []Tile{1, 2}})
	assertMarshalTileNotValid(Set{Tiles: []Tile{ballsBase}, Concealed: false})
	assertMarshalTileNotValid(Set{Tiles: []Tile{Balls8, Balls9, Balls9 + 1}})

	var loadedSet Set
	err := json.Unmarshal([]byte("{\"tiles\":[1,2],\"concealed\":false}"), &loadedSet)
	assert.Equal(c, ErrTileNotValid, err)
}
