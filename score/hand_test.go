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

func (s *HandTestSuite) TestJSON(c *check.C) {
	asJSON, err := json.Marshal(Tile(Balls4))
	assert.Nil(c, err)
	assert.Equal(c, "14", string(asJSON))

	var loadedTile Tile
	err = json.Unmarshal(asJSON, &loadedTile)
	assert.Nil(c, err)
	assert.Equal(c, Balls4, loadedTile)
}
