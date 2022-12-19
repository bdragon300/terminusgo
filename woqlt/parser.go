package woqlt

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	"github.com/bdragon300/terminusgo/woql"
	"github.com/bdragon300/terminusgo/woqlt/grammar"
)

type (
	callParam  any
	callResult any
)

type call struct {
	fname  string
	params []callParam
}

type SyntaxListener struct {
	*grammar.BaseWoqlListener
	callStack []*call
	last      callResult
}

func (sl *SyntaxListener) push(c *call) {
	sl.callStack = append(sl.callStack, c)
}

func (sl *SyntaxListener) popAndInvoke() callResult {
	l := len(sl.callStack)
	if l == 0 {
		return nil
	}
	var lastCall *call
	lastCall, sl.callStack = sl.callStack[l-1], sl.callStack[:l-1]

	fn, ok := functionsVocabulary[lastCall.fname]
	if !ok {
		panic("No such function '" + lastCall.fname + "'") // FIXME: output error
	}
	result := fn(woql.NewSimpleQueryBuilder(), lastCall.params)
	return result
}

func (sl *SyntaxListener) top() *call {
	l := len(sl.callStack)
	if l == 0 {
		return nil
	}
	return sl.callStack[l-1]
}

func (sl *SyntaxListener) pushParamTop(param callParam) {
	top := sl.top()
	if top == nil {
		panic("Must not happen, parameters outside of a function")
	}
	top.params = append(top.params, param)
}

func (sl *SyntaxListener) EnterFcall(ctx *grammar.FcallContext) {
	sl.push(&call{fname: ctx.Fname().GetText()})
}

func (sl *SyntaxListener) ExitFcall(_ *grammar.FcallContext) {
	sl.last = sl.popAndInvoke()
}

func (sl *SyntaxListener) ExitParam(ctx *grammar.ParamContext) {
	node := ctx.GetChild(0)
	switch v := node.(type) {
	case antlr.TerminalNode:
		contents := v.GetText()

		switch v.GetSymbol().GetTokenType() {
		case grammar.WoqlParserSTRPARAM:
			str, _ := renderString(contents) // FIXME: handle error
			sl.pushParamTop(str)
		case grammar.WoqlParserNUMPARAM:
			parts := strings.Split(contents, ".")
			prec := uint(0)
			if len(parts) > 1 {
				prec = 53 // IEEE_754 double type mantissa bits
			}
			// FIXME: support other bases
			// FIXME: dynamic precision
			// FIXME: handle error
			num, _, _ := big.ParseFloat(contents, 10, prec, big.ToNearestEven)
			sl.pushParamTop(*num) // FIXME: handle nil pointer
		}
	case grammar.IFcallContext:
		sl.pushParamTop(sl.last)
	}
}

func ParseWOQL(query string) (*woql.QueryBuilder, error) {
	is := antlr.NewInputStream(query)
	lexer := grammar.NewWoqlLexer(is)
	tokStream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := grammar.NewWoqlParser(tokStream)
	listener := &SyntaxListener{}
	antlr.ParseTreeWalkerDefault.Walk(listener, parser.Query())

	queryObj, ok := listener.last.(*woql.QueryBuilder)
	if !ok {
		val := reflect.Indirect(reflect.ValueOf(queryObj))
		return nil, fmt.Errorf("WOQL result must be a query object, got: %v", val.Type())
	}
	return queryObj, nil
}

func renderString(stringRepr string) (string, error) {
	str, err := strconv.Unquote(stringRepr)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(str, "\\\"", "\""), nil
}
