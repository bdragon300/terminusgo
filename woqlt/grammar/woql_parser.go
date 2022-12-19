// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package grammar // Woql
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type WoqlParser struct {
	*antlr.BaseParser
}

var woqlParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func woqlParserInit() {
	staticData := &woqlParserStaticData
	staticData.literalNames = []string{
		"", "','", "';'", "'('", "')'",
	}
	staticData.symbolicNames = []string{
		"", "COMMA", "SEMI", "LPAREN", "RPAREN", "STRPARAM", "NUMPARAM", "VARPARAM",
		"WS", "ID",
	}
	staticData.ruleNames = []string{
		"query", "fcall", "fname", "paramlist", "param",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 9, 35, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7, 4,
		1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 3, 1, 3,
		1, 3, 5, 3, 24, 8, 3, 10, 3, 12, 3, 27, 9, 3, 1, 4, 1, 4, 1, 4, 1, 4, 3,
		4, 33, 8, 4, 1, 4, 0, 0, 5, 0, 2, 4, 6, 8, 0, 0, 33, 0, 10, 1, 0, 0, 0,
		2, 13, 1, 0, 0, 0, 4, 18, 1, 0, 0, 0, 6, 20, 1, 0, 0, 0, 8, 32, 1, 0, 0,
		0, 10, 11, 3, 2, 1, 0, 11, 12, 5, 0, 0, 1, 12, 1, 1, 0, 0, 0, 13, 14, 3,
		4, 2, 0, 14, 15, 5, 3, 0, 0, 15, 16, 3, 6, 3, 0, 16, 17, 5, 4, 0, 0, 17,
		3, 1, 0, 0, 0, 18, 19, 5, 9, 0, 0, 19, 5, 1, 0, 0, 0, 20, 25, 3, 8, 4,
		0, 21, 22, 5, 1, 0, 0, 22, 24, 3, 8, 4, 0, 23, 21, 1, 0, 0, 0, 24, 27,
		1, 0, 0, 0, 25, 23, 1, 0, 0, 0, 25, 26, 1, 0, 0, 0, 26, 7, 1, 0, 0, 0,
		27, 25, 1, 0, 0, 0, 28, 33, 5, 6, 0, 0, 29, 33, 5, 7, 0, 0, 30, 33, 5,
		5, 0, 0, 31, 33, 3, 2, 1, 0, 32, 28, 1, 0, 0, 0, 32, 29, 1, 0, 0, 0, 32,
		30, 1, 0, 0, 0, 32, 31, 1, 0, 0, 0, 33, 9, 1, 0, 0, 0, 2, 25, 32,
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

// WoqlParserInit initializes any static state used to implement WoqlParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewWoqlParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func WoqlParserInit() {
	staticData := &woqlParserStaticData
	staticData.once.Do(woqlParserInit)
}

// NewWoqlParser produces a new parser instance for the optional input antlr.TokenStream.
func NewWoqlParser(input antlr.TokenStream) *WoqlParser {
	WoqlParserInit()
	this := new(WoqlParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &woqlParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "java-escape"

	return this
}

// WoqlParser tokens.
const (
	WoqlParserEOF      = antlr.TokenEOF
	WoqlParserCOMMA    = 1
	WoqlParserSEMI     = 2
	WoqlParserLPAREN   = 3
	WoqlParserRPAREN   = 4
	WoqlParserSTRPARAM = 5
	WoqlParserNUMPARAM = 6
	WoqlParserVARPARAM = 7
	WoqlParserWS       = 8
	WoqlParserID       = 9
)

// WoqlParser rules.
const (
	WoqlParserRULE_query     = 0
	WoqlParserRULE_fcall     = 1
	WoqlParserRULE_fname     = 2
	WoqlParserRULE_paramlist = 3
	WoqlParserRULE_param     = 4
)

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsQueryContext differentiates from other interfaces.
	IsQueryContext()
}

type QueryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryContext() *QueryContext {
	var p = new(QueryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WoqlParserRULE_query
	return p
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WoqlParserRULE_query

	return p
}

func (s *QueryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryContext) Fcall() IFcallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFcallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFcallContext)
}

func (s *QueryContext) EOF() antlr.TerminalNode {
	return s.GetToken(WoqlParserEOF, 0)
}

func (s *QueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.ExitQuery(s)
	}
}

func (p *WoqlParser) Query() (localctx IQueryContext) {
	this := p
	_ = this

	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, WoqlParserRULE_query)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(10)
		p.Fcall()
	}
	{
		p.SetState(11)
		p.Match(WoqlParserEOF)
	}

	return localctx
}

// IFcallContext is an interface to support dynamic dispatch.
type IFcallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFcallContext differentiates from other interfaces.
	IsFcallContext()
}

type FcallContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFcallContext() *FcallContext {
	var p = new(FcallContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WoqlParserRULE_fcall
	return p
}

func (*FcallContext) IsFcallContext() {}

func NewFcallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FcallContext {
	var p = new(FcallContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WoqlParserRULE_fcall

	return p
}

func (s *FcallContext) GetParser() antlr.Parser { return s.parser }

func (s *FcallContext) Fname() IFnameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFnameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFnameContext)
}

func (s *FcallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(WoqlParserLPAREN, 0)
}

func (s *FcallContext) Paramlist() IParamlistContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamlistContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamlistContext)
}

func (s *FcallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(WoqlParserRPAREN, 0)
}

func (s *FcallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FcallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FcallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.EnterFcall(s)
	}
}

func (s *FcallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.ExitFcall(s)
	}
}

func (p *WoqlParser) Fcall() (localctx IFcallContext) {
	this := p
	_ = this

	localctx = NewFcallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, WoqlParserRULE_fcall)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(13)
		p.Fname()
	}
	{
		p.SetState(14)
		p.Match(WoqlParserLPAREN)
	}
	{
		p.SetState(15)
		p.Paramlist()
	}
	{
		p.SetState(16)
		p.Match(WoqlParserRPAREN)
	}

	return localctx
}

// IFnameContext is an interface to support dynamic dispatch.
type IFnameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFnameContext differentiates from other interfaces.
	IsFnameContext()
}

type FnameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFnameContext() *FnameContext {
	var p = new(FnameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WoqlParserRULE_fname
	return p
}

func (*FnameContext) IsFnameContext() {}

func NewFnameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FnameContext {
	var p = new(FnameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WoqlParserRULE_fname

	return p
}

func (s *FnameContext) GetParser() antlr.Parser { return s.parser }

func (s *FnameContext) ID() antlr.TerminalNode {
	return s.GetToken(WoqlParserID, 0)
}

func (s *FnameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FnameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FnameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.EnterFname(s)
	}
}

func (s *FnameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.ExitFname(s)
	}
}

func (p *WoqlParser) Fname() (localctx IFnameContext) {
	this := p
	_ = this

	localctx = NewFnameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, WoqlParserRULE_fname)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(18)
		p.Match(WoqlParserID)
	}

	return localctx
}

// IParamlistContext is an interface to support dynamic dispatch.
type IParamlistContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParamlistContext differentiates from other interfaces.
	IsParamlistContext()
}

type ParamlistContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParamlistContext() *ParamlistContext {
	var p = new(ParamlistContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WoqlParserRULE_paramlist
	return p
}

func (*ParamlistContext) IsParamlistContext() {}

func NewParamlistContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamlistContext {
	var p = new(ParamlistContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WoqlParserRULE_paramlist

	return p
}

func (s *ParamlistContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamlistContext) AllParam() []IParamContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IParamContext); ok {
			len++
		}
	}

	tst := make([]IParamContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IParamContext); ok {
			tst[i] = t.(IParamContext)
			i++
		}
	}

	return tst
}

func (s *ParamlistContext) Param(i int) IParamContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamContext)
}

func (s *ParamlistContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(WoqlParserCOMMA)
}

func (s *ParamlistContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(WoqlParserCOMMA, i)
}

func (s *ParamlistContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamlistContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamlistContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.EnterParamlist(s)
	}
}

func (s *ParamlistContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.ExitParamlist(s)
	}
}

func (p *WoqlParser) Paramlist() (localctx IParamlistContext) {
	this := p
	_ = this

	localctx = NewParamlistContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, WoqlParserRULE_paramlist)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(20)
		p.Param()
	}
	p.SetState(25)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == WoqlParserCOMMA {
		{
			p.SetState(21)
			p.Match(WoqlParserCOMMA)
		}
		{
			p.SetState(22)
			p.Param()
		}

		p.SetState(27)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IParamContext is an interface to support dynamic dispatch.
type IParamContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsParamContext differentiates from other interfaces.
	IsParamContext()
}

type ParamContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParamContext() *ParamContext {
	var p = new(ParamContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = WoqlParserRULE_param
	return p
}

func (*ParamContext) IsParamContext() {}

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext {
	var p = new(ParamContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = WoqlParserRULE_param

	return p
}

func (s *ParamContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamContext) NUMPARAM() antlr.TerminalNode {
	return s.GetToken(WoqlParserNUMPARAM, 0)
}

func (s *ParamContext) VARPARAM() antlr.TerminalNode {
	return s.GetToken(WoqlParserVARPARAM, 0)
}

func (s *ParamContext) STRPARAM() antlr.TerminalNode {
	return s.GetToken(WoqlParserSTRPARAM, 0)
}

func (s *ParamContext) Fcall() IFcallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFcallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFcallContext)
}

func (s *ParamContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.EnterParam(s)
	}
}

func (s *ParamContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(WoqlListener); ok {
		listenerT.ExitParam(s)
	}
}

func (p *WoqlParser) Param() (localctx IParamContext) {
	this := p
	_ = this

	localctx = NewParamContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, WoqlParserRULE_param)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(32)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case WoqlParserNUMPARAM:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(28)
			p.Match(WoqlParserNUMPARAM)
		}

	case WoqlParserVARPARAM:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(29)
			p.Match(WoqlParserVARPARAM)
		}

	case WoqlParserSTRPARAM:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(30)
			p.Match(WoqlParserSTRPARAM)
		}

	case WoqlParserID:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(31)
			p.Fcall()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}
