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
	expList   []Expression
	identList []Identifier
}

%type<expr>      program expression functionDefine call binaryOperator unaryOperator number
%type<ident>     identifier
%type<expList>   callArguments
%type<identList> defineArguments

%token<token> NUMBER IDENTIFIER

%right ';'
%right '='

%left  '+' '-'
%left  '*' '/'

%right '!'

%%

program
	: expression
	{
		$$ = $1
		yylex.(*Lexer).result = $$
	}

expression
	: expression ';'
	| number
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
	| callArguments ','

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
	: expression ';' expression
	{
		$$ = FunctionCall {
			Function: Identifier(";"),
			Arguments: []Expression{$1, $3},
		}
	}
	| expression '+' expression
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

functionDefine
	: '(' defineArguments ')' '{' expression '}'
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
	| defineArguments ','

%%