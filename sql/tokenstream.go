package sql

import "github.com/Knetic/govaluate"

type tokenStream struct {
tokens      []govaluate.ExpressionToken
index       int
tokenLength int
}

func newTokenStream(tokens []govaluate.ExpressionToken) *tokenStream {

var ret *tokenStream

ret = new(tokenStream)
ret.tokens = tokens
ret.tokenLength = len(tokens)
return ret
}

func (this *tokenStream) rewind() {
this.index -= 1
}

func (this *tokenStream) next() govaluate.ExpressionToken {

var token govaluate.ExpressionToken

token = this.tokens[this.index]

this.index += 1
return token
}

func (this tokenStream) hasNext() bool {

return this.index < this.tokenLength
}
