package score

import (
	"github.com/stretchr/testify/assert"
	check "gopkg.in/check.v1"
)

type TileValueTestSuite struct{}

var _ = check.Suite(&TileValueTestSuite{})

func (s *TileValueTestSuite) TestTileValues(c *check.C) {
	// Test some values of tiles
	assert.Equal(c, 12, int(Balls2))
	assert.Equal(c, 27, int(Chars7))
	assert.Equal(c, 39, int(Bamboo9))
	assert.Equal(c, 74, int(Season4))
}
