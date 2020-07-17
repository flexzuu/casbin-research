package sql

import "github.com/Knetic/govaluate"

/*
	Map of all valid comparators, and their string equivalents.
	Used during parsing of expressions to determine if a symbol is, in fact, a comparator.
	Also used during evaluation to determine exactly which comparator is being used.
*/
var comparatorSymbols = map[string]govaluate.OperatorSymbol{
	"==": govaluate.EQ,
	"!=": govaluate.NEQ,
	">":  govaluate.GT,
	">=": govaluate.GTE,
	"<":  govaluate.LT,
	"<=": govaluate.LTE,
	"=~": govaluate.REQ,
	"!~": govaluate.NREQ,
	"in": govaluate.IN,
}

var logicalSymbols = map[string]govaluate.OperatorSymbol{
	"&&": govaluate.AND,
	"||": govaluate.OR,
}

var bitwiseSymbols = map[string]govaluate.OperatorSymbol{
	"^": govaluate.BITWISE_XOR,
	"&": govaluate.BITWISE_AND,
	"|": govaluate.BITWISE_OR,
}

var bitwiseShiftSymbols = map[string]govaluate.OperatorSymbol{
	">>": govaluate.BITWISE_RSHIFT,
	"<<": govaluate.BITWISE_LSHIFT,
}

var additiveSymbols = map[string]govaluate.OperatorSymbol{
	"+": govaluate.PLUS,
	"-": govaluate.MINUS,
}

var multiplicativeSymbols = map[string]govaluate.OperatorSymbol{
	"*": govaluate.MULTIPLY,
	"/": govaluate.DIVIDE,
	"%": govaluate.MODULUS,
}

var exponentialSymbolsS = map[string]govaluate.OperatorSymbol{
	"**": govaluate.EXPONENT,
}

var prefixSymbols = map[string]govaluate.OperatorSymbol{
	"-": govaluate.NEGATE,
	"!": govaluate.INVERT,
	"~": govaluate.BITWISE_NOT,
}

var ternarySymbols = map[string]govaluate.OperatorSymbol{
	"?":  govaluate.TERNARY_TRUE,
	":":  govaluate.TERNARY_FALSE,
	"??": govaluate.COALESCE,
}

// this is defined separately from additiveSymbols et al because it's needed for parsing, not stage planning.
var modifierSymbols = map[string]govaluate.OperatorSymbol{
	"+":  govaluate.PLUS,
	"-":  govaluate.MINUS,
	"*":  govaluate.MULTIPLY,
	"/":  govaluate.DIVIDE,
	"%":  govaluate.MODULUS,
	"**": govaluate.EXPONENT,
	"&":  govaluate.BITWISE_AND,
	"|":  govaluate.BITWISE_OR,
	"^":  govaluate.BITWISE_XOR,
	">>": govaluate.BITWISE_RSHIFT,
	"<<": govaluate.BITWISE_LSHIFT,
}

var separatorSymbols = map[string]govaluate.OperatorSymbol{
	",": govaluate.SEPARATE,
}
