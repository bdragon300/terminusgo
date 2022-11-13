// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package grammar // Woql
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// WoqlListener is a complete listener for a parse tree produced by WoqlParser.
type WoqlListener interface {
	antlr.ParseTreeListener

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterFcall is called when entering the fcall production.
	EnterFcall(c *FcallContext)

	// EnterFname is called when entering the fname production.
	EnterFname(c *FnameContext)

	// EnterParamlist is called when entering the paramlist production.
	EnterParamlist(c *ParamlistContext)

	// EnterParam is called when entering the param production.
	EnterParam(c *ParamContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitFcall is called when exiting the fcall production.
	ExitFcall(c *FcallContext)

	// ExitFname is called when exiting the fname production.
	ExitFname(c *FnameContext)

	// ExitParamlist is called when exiting the paramlist production.
	ExitParamlist(c *ParamlistContext)

	// ExitParam is called when exiting the param production.
	ExitParam(c *ParamContext)
}
