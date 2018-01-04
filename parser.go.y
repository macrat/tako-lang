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
	call      FunctionCall
	expList   ExpressionList
	identList []Identifier
	object    Object
}

%type<expr>      program expression functionDefine number condition condFunction
%type<call>      call binaryOperator unaryOperator takeMember
%type<ident>     identifier
%type<expList>   callArguments expressionList
%type<identList> defineArguments
%type<object>    object objectList

%token<token> NUMBER IDENTIFIER NEWLINE DEFINE_OPERATOR CALCULATE_DEFINE_OPERATOR COMPARE_OPERATOR IF ELSE FUNCTION_SEP

%right ';'
%right DEFINE_OPERATOR
%right COMPARE_OPERATOR

%left  '+' '-'
%left  '*' '/'

%right '!'

%right '.'

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
	| '(' expressionList ')'
	{ $$ = $2 }

expression
	: number
	| identifier
	{ $$ = $1 }
	| call
	{ $$ = $1 }
	| functionDefine
	| condition
	| '(' expression ')'
	{ $$ = $2 }
	| object
	{ $$ = $1 }

object
	: '[' objectList ']'
	{ $$ = $2 }

objectList
	:
	{
		$$ = Object { Named: make(map[string]Expression) }
	}
	| expression
	{
		$$ = Object { Indexed: []Expression{$1}, Named: make(map[string]Expression) }
	}
	| objectList ',' expression
	{
		$$ = $1
		$$.Indexed = append($1.Indexed, $3)
	}
	| objectList ',' NEWLINE expression
	{
		$$ = $1
		$$.Indexed = append($1.Indexed, $4)
	}
	| identifier ':' expression
	{
		$$ = Object { Named: map[string]Expression{$1.Key: $3} }
	}
	| objectList ',' identifier ':' expression
	{
		$$ = $1
		$$.Named[$3.Key] = $5
	}
	| objectList ',' NEWLINE identifier ':' expression
	{
		$$ = $1
		$$.Named[$4.Key] = $6
	}
	| objectList ','
	| objectList ',' NEWLINE

number
	: NUMBER
	{
		num, _ := strconv.ParseInt($1.Literal, 10, 64)
		$$ = Number(num)
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
	| takeMember '(' callArguments ')'
	{
		$$ = FunctionCall {
			Function: $1,
			Arguments: append([]Expression{$1.Arguments[0]}, $3...),
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
			Function: NewIdentifier("-:"),
			Arguments: []Expression{$2},
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| '!' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier("!:"),
			Arguments: []Expression{$2},
			Pos: yylex.(*Lexer).lastPosition,
		}
	}

binaryOperator
	: expression '+' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier(":+:"),
			Arguments: []Expression{$1, $3},
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| expression '-' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier(":-:"),
			Arguments: []Expression{$1, $3},
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| expression '*' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier(":*:"),
			Arguments: []Expression{$1, $3},
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| expression '/' expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier(":/:"),
			Arguments: []Expression{$1, $3},
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| expression COMPARE_OPERATOR expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier(":" + $2.Literal + ":"),
			Arguments: []Expression{$1, $3},
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| identifier DEFINE_OPERATOR expression
	{
		$$ = FunctionCall {
			Function: NewIdentifier(":" + $2.Literal + ":"),
			Arguments: []Expression{$1, $3},
			Pos: $1.Position(),
		}
	}
	| takeMember DEFINE_OPERATOR expression
	{
		funcName := $1.Function.(Identifier).Key
		$$ = FunctionCall {
			Function: NewIdentifier(funcName[:len(funcName)-1] + "=:"),
			Arguments: append($1.Arguments, $3),
			Pos: $1.Position(),
		}
	}
	| takeMember CALCULATE_DEFINE_OPERATOR expression
	{
		funcName := $1.Function.(Identifier).Key
		$$ = FunctionCall {
			Function: NewIdentifier(funcName[:len(funcName)-1] + "=:"),
			Arguments: append(
				$1.Arguments,
				FunctionCall {
					Function: NewIdentifier(":" + $2.Literal + ":"),
					Arguments: []Expression{$1, $3},
					Pos: $1.Position(),
				},
			),
			Pos: $1.Position(),
		}
	}
	| takeMember

takeMember
	: expression '.' identifier
	{
		$$ = FunctionCall {
			Function: NewIdentifier(":.:"),
			Arguments: []Expression{$1, $3},
			Pos: yylex.(*Lexer).lastPosition,
		}
	}
	| expression '[' expression ']'
	{
		$$ = FunctionCall {
			Function: NewIdentifier(":[]:"),
			Arguments: []Expression{$1, $3},
			Pos: yylex.(*Lexer).lastPosition,
		}
	}

functionDefine
	: '(' defineArguments FUNCTION_SEP expressionList '}'
	{
		$$ = FunctionDefine {
			Arguments: $2,
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
	: condition
	| '{' expression '}'
	{
		$$ = FunctionDefine {
			Arguments: []Identifier{},
			Expression: $2,
			Pos: yylex.(*Lexer).lastPosition,
		}
	}

%%
