// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package grammar

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import error
var _ = fmt.Printf
var _ = sync.Once{}
var _ = unicode.IsLetter

type WoqlLexer struct {
	*antlr.BaseLexer
	channelNames []string
	modeNames    []string
	// TODO: EOF string
}

var woqllexerLexerStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	channelNames           []string
	modeNames              []string
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func woqllexerLexerInit() {
	staticData := &woqllexerLexerStaticData
	staticData.channelNames = []string{
		"DEFAULT_TOKEN_CHANNEL", "HIDDEN",
	}
	staticData.modeNames = []string{
		"DEFAULT_MODE",
	}
	staticData.literalNames = []string{
		"", "','", "';'", "'('", "')'",
	}
	staticData.symbolicNames = []string{
		"", "COMMA", "SEMI", "LPAREN", "RPAREN", "STRPARAM", "NUMPARAM", "VARPARAM",
		"WS", "ID",
	}
	staticData.ruleNames = []string{
		"COMMA", "SEMI", "LPAREN", "RPAREN", "STRPARAM", "NUMPARAM", "VARPARAM",
		"WS", "ID", "ANYSTR", "ANYNUM", "SPACES", "NEWLINE",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 0, 9, 92, 6, -1, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2,
		4, 7, 4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2,
		10, 7, 10, 2, 11, 7, 11, 2, 12, 7, 12, 1, 0, 1, 0, 1, 1, 1, 1, 1, 2, 1,
		2, 1, 3, 1, 3, 1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1,
		6, 1, 6, 1, 7, 1, 7, 3, 7, 49, 8, 7, 1, 7, 1, 7, 1, 8, 1, 8, 4, 8, 55,
		8, 8, 11, 8, 12, 8, 56, 1, 9, 1, 9, 1, 9, 1, 9, 5, 9, 63, 8, 9, 10, 9,
		12, 9, 66, 9, 9, 1, 9, 1, 9, 1, 10, 4, 10, 71, 8, 10, 11, 10, 12, 10, 72,
		1, 10, 1, 10, 4, 10, 77, 8, 10, 11, 10, 12, 10, 78, 3, 10, 81, 8, 10, 1,
		11, 4, 11, 84, 8, 11, 11, 11, 12, 11, 85, 1, 12, 4, 12, 89, 8, 12, 11,
		12, 12, 12, 90, 0, 0, 13, 1, 1, 3, 2, 5, 3, 7, 4, 9, 5, 11, 6, 13, 7, 15,
		8, 17, 9, 19, 0, 21, 0, 23, 0, 25, 0, 1, 0, 6, 3, 0, 65, 90, 95, 95, 97,
		122, 4, 0, 48, 57, 65, 90, 95, 95, 97, 122, 4, 0, 10, 10, 12, 13, 34, 34,
		92, 92, 1, 0, 48, 57, 2, 0, 9, 9, 32, 32, 2, 0, 10, 10, 13, 13, 96, 0,
		1, 1, 0, 0, 0, 0, 3, 1, 0, 0, 0, 0, 5, 1, 0, 0, 0, 0, 7, 1, 0, 0, 0, 0,
		9, 1, 0, 0, 0, 0, 11, 1, 0, 0, 0, 0, 13, 1, 0, 0, 0, 0, 15, 1, 0, 0, 0,
		0, 17, 1, 0, 0, 0, 1, 27, 1, 0, 0, 0, 3, 29, 1, 0, 0, 0, 5, 31, 1, 0, 0,
		0, 7, 33, 1, 0, 0, 0, 9, 35, 1, 0, 0, 0, 11, 37, 1, 0, 0, 0, 13, 39, 1,
		0, 0, 0, 15, 48, 1, 0, 0, 0, 17, 52, 1, 0, 0, 0, 19, 58, 1, 0, 0, 0, 21,
		70, 1, 0, 0, 0, 23, 83, 1, 0, 0, 0, 25, 88, 1, 0, 0, 0, 27, 28, 5, 44,
		0, 0, 28, 2, 1, 0, 0, 0, 29, 30, 5, 59, 0, 0, 30, 4, 1, 0, 0, 0, 31, 32,
		5, 40, 0, 0, 32, 6, 1, 0, 0, 0, 33, 34, 5, 41, 0, 0, 34, 8, 1, 0, 0, 0,
		35, 36, 3, 19, 9, 0, 36, 10, 1, 0, 0, 0, 37, 38, 3, 21, 10, 0, 38, 12,
		1, 0, 0, 0, 39, 40, 5, 34, 0, 0, 40, 41, 5, 118, 0, 0, 41, 42, 5, 58, 0,
		0, 42, 43, 1, 0, 0, 0, 43, 44, 3, 17, 8, 0, 44, 45, 5, 34, 0, 0, 45, 14,
		1, 0, 0, 0, 46, 49, 3, 23, 11, 0, 47, 49, 3, 25, 12, 0, 48, 46, 1, 0, 0,
		0, 48, 47, 1, 0, 0, 0, 49, 50, 1, 0, 0, 0, 50, 51, 6, 7, 0, 0, 51, 16,
		1, 0, 0, 0, 52, 54, 7, 0, 0, 0, 53, 55, 7, 1, 0, 0, 54, 53, 1, 0, 0, 0,
		55, 56, 1, 0, 0, 0, 56, 54, 1, 0, 0, 0, 56, 57, 1, 0, 0, 0, 57, 18, 1,
		0, 0, 0, 58, 64, 5, 34, 0, 0, 59, 60, 5, 92, 0, 0, 60, 63, 5, 34, 0, 0,
		61, 63, 8, 2, 0, 0, 62, 59, 1, 0, 0, 0, 62, 61, 1, 0, 0, 0, 63, 66, 1,
		0, 0, 0, 64, 62, 1, 0, 0, 0, 64, 65, 1, 0, 0, 0, 65, 67, 1, 0, 0, 0, 66,
		64, 1, 0, 0, 0, 67, 68, 5, 34, 0, 0, 68, 20, 1, 0, 0, 0, 69, 71, 7, 3,
		0, 0, 70, 69, 1, 0, 0, 0, 71, 72, 1, 0, 0, 0, 72, 70, 1, 0, 0, 0, 72, 73,
		1, 0, 0, 0, 73, 80, 1, 0, 0, 0, 74, 76, 5, 46, 0, 0, 75, 77, 7, 3, 0, 0,
		76, 75, 1, 0, 0, 0, 77, 78, 1, 0, 0, 0, 78, 76, 1, 0, 0, 0, 78, 79, 1,
		0, 0, 0, 79, 81, 1, 0, 0, 0, 80, 74, 1, 0, 0, 0, 80, 81, 1, 0, 0, 0, 81,
		22, 1, 0, 0, 0, 82, 84, 7, 4, 0, 0, 83, 82, 1, 0, 0, 0, 84, 85, 1, 0, 0,
		0, 85, 83, 1, 0, 0, 0, 85, 86, 1, 0, 0, 0, 86, 24, 1, 0, 0, 0, 87, 89,
		7, 5, 0, 0, 88, 87, 1, 0, 0, 0, 89, 90, 1, 0, 0, 0, 90, 88, 1, 0, 0, 0,
		90, 91, 1, 0, 0, 0, 91, 26, 1, 0, 0, 0, 10, 0, 48, 56, 62, 64, 72, 78,
		80, 85, 90, 1, 6, 0, 0,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// WoqlLexerInit initializes any static state used to implement WoqlLexer. By default the
// static state used to implement the lexer is lazily initialized during the first call to
// NewWoqlLexer(). You can call this function if you wish to initialize the static state ahead
// of time.
func WoqlLexerInit() {
	staticData := &woqllexerLexerStaticData
	staticData.once.Do(woqllexerLexerInit)
}

// NewWoqlLexer produces a new lexer instance for the optional input antlr.CharStream.
func NewWoqlLexer(input antlr.CharStream) *WoqlLexer {
	WoqlLexerInit()
	l := new(WoqlLexer)
	l.BaseLexer = antlr.NewBaseLexer(input)
	staticData := &woqllexerLexerStaticData
	l.Interpreter = antlr.NewLexerATNSimulator(l, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	l.channelNames = staticData.channelNames
	l.modeNames = staticData.modeNames
	l.RuleNames = staticData.ruleNames
	l.LiteralNames = staticData.literalNames
	l.SymbolicNames = staticData.symbolicNames
	l.GrammarFileName = "Woql.g4"
	// TODO: l.EOF = antlr.TokenEOF

	return l
}

// WoqlLexer tokens.
const (
	WoqlLexerCOMMA    = 1
	WoqlLexerSEMI     = 2
	WoqlLexerLPAREN   = 3
	WoqlLexerRPAREN   = 4
	WoqlLexerSTRPARAM = 5
	WoqlLexerNUMPARAM = 6
	WoqlLexerVARPARAM = 7
	WoqlLexerWS       = 8
	WoqlLexerID       = 9
)
