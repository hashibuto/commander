package commander

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer_LeadingTrailing(t *testing.T) {
	line := "  hello   there  macaroni    "
	tokenGroups := Tokenize(line)
	assert.Len(t, tokenGroups, 1)
	tokens := tokenGroups[0].Tokens

	assert.Len(t, tokens, 3)
	assert.Equal(t, tokens[0], "hello")
	assert.Equal(t, tokens[1], "there")
	assert.Equal(t, tokens[2], "macaroni")
}
