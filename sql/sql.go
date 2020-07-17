package sql

import (
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"regexp"
	"strings"
	"time"
)

type EvaluableExpression struct {
	*govaluate.EvaluableExpression
	Sub string
}

/*
	Returns a string representing this expression as if it were written in SQL.
	This function assumes that all parameters exist within the same table, and that the table essentially represents
	a serialized object of some sort (e.g., hibernate).
	If your data model is more normalized, you may need to consider iterating through each actual token given by `Tokens()`
	to create your query.

	Boolean values are considered to be "1" for true, "0" for false.

	Times are formatted according to this.QueryDateFormat.
*/
func (e EvaluableExpression) ToSQLQuery() (string, error) {

	var stream *tokenStream
	var transactions *expressionOutputStream
	var transaction string
	var err error

	stream = newTokenStream(e.Tokens())
	transactions = new(expressionOutputStream)

	for stream.hasNext() {

		transaction, err = e.findNextSQLString(stream, transactions)
		if err != nil {
			return "", err
		}

		transactions.add(transaction)
	}

	return transactions.createString(" "), nil
}

func (e EvaluableExpression) findNextSQLString(stream *tokenStream, transactions *expressionOutputStream) (string, error) {

	var token govaluate.ExpressionToken
	var ret string

	token = stream.next()

	switch token.Kind {

	case govaluate.STRING:
		ret = fmt.Sprintf("'%v'", token.Value)
	case govaluate.PATTERN:
		ret = fmt.Sprintf("'%s'", token.Value.(*regexp.Regexp).String())
	case govaluate.TIME:
		ret = fmt.Sprintf("'%s'", token.Value.(time.Time).Format(e.QueryDateFormat))

	case govaluate.LOGICALOP:
		switch logicalSymbols[token.Value.(string)] {

		case govaluate.AND:
			ret = "AND"
		case govaluate.OR:
			ret = "OR"
		}

	case govaluate.BOOLEAN:
		if token.Value.(bool) {
			ret = "1"
		} else {
			ret = "0"
		}

	case govaluate.VARIABLE:
		if(token.Value.(string) == "r_sub") {
			ret = fmt.Sprintf("'%s'", e.Sub)
			return ret, nil
		}

		ret = fmt.Sprintf("[%s]", token.Value.(string))


	case govaluate.NUMERIC:
		ret = fmt.Sprintf("%g", token.Value.(float64))

	case govaluate.COMPARATOR:
		switch comparatorSymbols[token.Value.(string)] {

		case govaluate.EQ:
			ret = "="
		case govaluate.NEQ:
			ret = "<>"
		case govaluate.REQ:
			ret = "RLIKE"
		case govaluate.NREQ:
			ret = "NOT RLIKE"
		default:
			ret = fmt.Sprintf("%s", token.Value.(string))
		}

	case govaluate.TERNARY:

		switch ternarySymbols[token.Value.(string)] {

		case govaluate.COALESCE:

			left := transactions.rollback()
			right, err := e.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("COALESCE(%v, %v)", left, right)
		case govaluate.TERNARY_TRUE:
			fallthrough
		case govaluate.TERNARY_FALSE:
			return "", errors.New("Ternary operators are unsupported in SQL output")
		}
	case govaluate.PREFIX:
		switch prefixSymbols[token.Value.(string)] {

		case govaluate.INVERT:
			ret = fmt.Sprintf("NOT")
		default:

			right, err := e.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("%s%s", token.Value.(string), right)
		}
	case govaluate.MODIFIER:

		switch modifierSymbols[token.Value.(string)] {

		case govaluate.EXPONENT:

			left := transactions.rollback()
			right, err := e.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("POW(%s, %s)", left, right)
		case govaluate.MODULUS:

			left := transactions.rollback()
			right, err := e.findNextSQLString(stream, transactions)
			if err != nil {
				return "", err
			}

			ret = fmt.Sprintf("MOD(%s, %s)", left, right)
		default:
			ret = fmt.Sprintf("%s", token.Value.(string))
		}
	case govaluate.CLAUSE:
		ret = "("
	case govaluate.CLAUSE_CLOSE:
		ret = ")"
	case govaluate.SEPARATOR:
		ret = ","
	case govaluate.ACCESSOR:

		switch v := token.Value.(type) {
		case []string: {
			if v[0] == "r_ctx" {
				ret = strings.ToLower(v[1])
			}
		}
		}


	default:
		errorMsg := fmt.Sprintf("Unrecognized query token '%s' of kind '%s'", token.Value, token.Kind)
		return "", errors.New(errorMsg)
	}

	return ret, nil
}
