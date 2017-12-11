package score

import (
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

type TileValueTestSuite struct {
}

var _ = check.Suite(&TileValueTestSuite{})

func (s *TileValueTestSuite) TestTileValues(c *check.C) {
	// Test some values of tiles
	assert.Equal(c, 12, int(BALLS_2))
	assert.Equal(c, 27, int(CHARS_7))
	assert.Equal(c, 39, int(BAMBOO_9))
	assert.Equal(c, 74, int(SEASON_4))
}
