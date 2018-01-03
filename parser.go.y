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

%type<expr>      program expression functionDefine call binaryOperator unaryOperator number
%type<ident>     identifier
%type<expList>   callArguments
%type<expList>   expressionList
%type<identList> defineArguments

%token<token> NUMBER IDENTIFIER NEWLINE DEFINE

%right ';'
%right '=' DEFINE

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
	| identifier
	{ $$ = $1 }
	| call
	{ $$ = $1 }
	| functionDefine
	{ $$ = $1 }

number
	: NUMBER
	{
		num, _ := strconv.ParseInt($1.Literal, 10, 64)
		$$ = Number(num)
	}

identifier
	: IDENTIFIER
	{
		$$ = Identifier($1.Literal)
	}

call
	: binaryOperator
	| unaryOperator
	| expression '(' callArguments ')'
	{
		$$ = FunctionCall {
			Function: $1,
			Arguments: $3,
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
			Function: Identifier("-"),
			Arguments: []Expression{$2},
		}
	}
	| '!' expression
	{
		$$ = FunctionCall {
			Function: Identifier("!"),
			Arguments: []Expression{$2},
		}
	}

binaryOperator
	: expression '+' expression
	{
		$$ = FunctionCall {
			Function: Identifier("+"),
			Arguments: []Expression{$1, $3},
		}
	}
	| expression '-' expression
	{
		$$ = FunctionCall {
			Function: Identifier("-"),
			Arguments: []Expression{$1, $3},
		}
	}
	| expression '*' expression
	{
		$$ = FunctionCall {
			Function: Identifier("*"),
			Arguments: []Expression{$1, $3},
		}
	}
	| expression '/' expression
	{
		$$ = FunctionCall {
			Function: Identifier("/"),
			Arguments: []Expression{$1, $3},
		}
	}
	| identifier '=' expression
	{
		$$ = FunctionCall {
			Function: Identifier("="),
			Arguments: []Expression{$1, $3},
		}
	}
	| identifier DEFINE expression
	{
		$$ = FunctionCall {
			Function: Identifier(":="),
			Arguments: []Expression{$1, $3},
		}
	}

functionDefine
	: '(' defineArguments ')' '{' expressionList '}'
	{
		$$ = FunctionDefine {
			Arguments: $2,
			Expression: $5,
		}
	}

defineArguments
	:
	{
		$$ = []Identifier{}
	}
	| identifier
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

%%
