%{
package main

import (
	"strconv"
)

%}

%union {
	expr      Expression
	token     Token
	ident     Identifier
	expList   ExpressionList
	identList []Identifier
}

%type<expr>      program expression functionDefine functionDefineWithoutArgument call binaryOperator unaryOperator number boolean null condition condFunction
%type<ident>     identifier
%type<expList>   callArguments
%type<expList>   expressionList
%type<identList> defineArguments

%token<token> NUMBER BOOLEAN NULL IDENTIFIER NEWLINE DEFINE_OPERATOR COMPARE_OPERATOR IF ELSE

%right ';'
%right DEFINE_OPERATOR
%right COMPARE_OPERATOR

%left  '+' '-'
%left  '*' '/'

%right '!'

%%

program
	: expressionList
	{
		$$ = $1
		yylex.(*Lexer).result = $$
	}

expressionList
	:
	{
		$$ = ExpressionList{}
	}
	| expression
	{
		$$ = ExpressionList{$1}
	}
	| expressionList ';' expression
	{
		$$ = append($1, $3)
	}
	| expressionList NEWLINE expression
	{
		$$ = append($1, $3)
	}
	| expressionList ';'
	| expressionList NEWLINE

expression
	: number
	{ $$ = $1 }
	| boolean
	{ $$ = $1 }
	| null
	{ $$ = $1 }
	| identifier
	{ $$ = $1 }
	| call
	{ $$ = $1 }
	| functionDefine
	{ $$ = $1 }
	| condition
	{ $$ = $1 }

number
	: NUMBER
	{
		num, _ := strconv.ParseInt($1.Literal, 10, 64)
		$$ = Number(num)
	}

boolean
	: BOOLEAN
	{
		$$ = Boolean($1.Literal == "true")
	}

null
	: NULL
	{
		$$ = Null{}
	}

identifier
	: IDENTIFIER
	{
		$$ = Identifier{
			Key: $1.Literal,
			Pos: $1.Pos,
		}
	}

call
	: binaryOperator
	| unaryOperator
	| expression '(' callArguments ')'
	{
		$$ = FunctionCall {
			Function: $1,
			Arguments: $3,
			Pos: yylex.(*Lexer).lastPosition,
		}
	}

callArguments
	:
	{
		$$ = []Expression{}
	}
	| expression
	{
		$$ = []Expression{$1}
	}
	| callArguments ',' expression
	{
		$$ = append($1, $3)
	}
	| callArguments ',' NEWLINE expression
	{
		$$ = append($1, $4)
	}
	| callArguments ',' NEWLINE

unaryOperator
	: '-' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier("-"),
			Arguments: []Expression{$2},
		}
	}
	| '!' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier("!"),
			Arguments: []Expression{$2},
		}
	}

binaryOperator
	: expression '+' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier("+"),
			Arguments: []Expression{$1, $3},
		}
	}
	| expression '-' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier("-"),
			Arguments: []Expression{$1, $3},
		}
	}
	| expression '*' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier("*"),
			Arguments: []Expression{$1, $3},
		}
	}
	| expression '/' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier("/"),
			Arguments: []Expression{$1, $3},
		}
	}
	| expression COMPARE_OPERATOR expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier($2.Literal),
			Arguments: []Expression{$1, $3},
		}
	}
	| identifier DEFINE_OPERATOR expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier($2.Literal),
			Arguments: []Expression{$1, $3},
		}
	}

functionDefine
	: '(' defineArguments ')' '{' expressionList '}'
	{
		$$ = FunctionDefine {
			Arguments: $2,
			Expression: $5,
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| functionDefineWithoutArgument

functionDefineWithoutArgument
	: '(' ')' '{' expressionList '}'
	{
		$$ = FunctionDefine {
			Arguments: []Identifier{},
			Expression: $4,
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| '{' expressionList '}'
	{
		$$ = FunctionDefine {
			Arguments: []Identifier{},
			Expression: $2,
			Pos: yylex.(*Lexer).lastPosition,
		}
	}

defineArguments
	: identifier
	{
		$$ = []Identifier{$1}
	}
	| defineArguments ',' identifier
	{
		$$ = append($1, $3)
	}
	| defineArguments ',' NEWLINE identifier
	{
		$$ = append($1, $4)
	}
	| defineArguments ',' NEWLINE

condition
	: IF expression condFunction
	{
		$$ = Condition{
			Condition: $2,
			Then: $3,
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| IF expression condFunction ELSE condFunction
	{
		$$ = Condition{
			Condition: $2,
			Then: $3,
			Else: $5,
			Pos: yylex.(*Lexer).lastPosition,
		}
	}

condFunction
	: expression
	{
		$$ = FunctionDefine {
			Arguments: []Identifier{},
			Expression: $1,
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| '{' expression '}'
	{
		$$ = FunctionDefine {
			Arguments: []Identifier{},
			Expression: $2,
			Pos: yylex.(*Lexer).lastPosition,
		}
	}

%%
