// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package grammar // Woql
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// BaseWoqlListener is a complete listener for a parse tree produced by WoqlParser.
type BaseWoqlListener struct{}

var _ WoqlListener = &BaseWoqlListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseWoqlListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseWoqlListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseWoqlListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseWoqlListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterQuery is called when production query is entered.
func (s *BaseWoqlListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BaseWoqlListener) ExitQuery(ctx *QueryContext) {}

// EnterFcall is called when production fcall is entered.
func (s *BaseWoqlListener) EnterFcall(ctx *FcallContext) {}

// ExitFcall is called when production fcall is exited.
func (s *BaseWoqlListener) ExitFcall(ctx *FcallContext) {}

// EnterFname is called when production fname is entered.
func (s *BaseWoqlListener) EnterFname(ctx *FnameContext) {}

// ExitFname is called when production fname is exited.
func (s *BaseWoqlListener) ExitFname(ctx *FnameContext) {}

// EnterParamlist is called when production paramlist is entered.
func (s *BaseWoqlListener) EnterParamlist(ctx *ParamlistContext) {}

// ExitParamlist is called when production paramlist is exited.
func (s *BaseWoqlListener) ExitParamlist(ctx *ParamlistContext) {}

// EnterParam is called when production param is entered.
func (s *BaseWoqlListener) EnterParam(ctx *ParamContext) {}

// ExitParam is called when production param is exited.
func (s *BaseWoqlListener) ExitParam(ctx *ParamContext) {}
