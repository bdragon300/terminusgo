grammar Woql;

query: fcall EOF;
fcall: fname LPAREN paramlist RPAREN;
fname: ID;
paramlist: param (',' param)*;
param: NUMPARAM | VARPARAM | STRPARAM | fcall;


COMMA : ',' ;
SEMI : ';' ;
LPAREN : '(' ;
RPAREN : ')' ;

STRPARAM : ANYSTR;
NUMPARAM : ANYNUM;
VARPARAM : '"v:' ID '"';

WS: (SPACES | NEWLINE) -> skip;
ID: [a-zA-Z_][a-zA-Z_0-9]+ ;

fragment ANYSTR : '"' ( '\\"' | ~[\\\r\n\f"] )* '"';
fragment ANYNUM : [0-9]+('.' [0-9]+)?;
fragment SPACES : [ \t]+;
fragment NEWLINE: [\n\r]+;
