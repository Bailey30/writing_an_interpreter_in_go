package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// the ORDER and RELATION to each other of these matter regarding operator precedence
const (
	_ int = iota // gives the following constants incrementing numbers as values
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	// we use these to check if the appropriate map (infix or prefix) has a parsing function associated with curToken.Type
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// functions that parse expressions and return expression nodes for the AST
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression // the argument is "left side" of infix operator thats being parsed
)

// adds expression parsing functions to list (map) so they can be extracted when needed
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	// associated tokens with parsing functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier) // registers the function parseIdentifier to token.IDENT - parseIdentifier just returns an ast.Identifier node with the current token information at the time the function is called. Its just the function that is associated right now

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

// t is the expcted token used when calling expectPeek()
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// move along each token creating an AST node when a statement or expression token is encountered
	// nextToken is then called within this methods until the statement has been parsed
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		// since let and return statements are the only two real statement types in Monkey
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// create a let statement node with token.LET token
	stmt := &ast.LetStatement{Token: p.curToken}

	// next token should be an indentifier e.g let x
	// also goes to next token
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// create an identifier node with the now current token.IDENT
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// expects an equals sign next
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// next token until a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: parse the expression

	// keep going until semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// checking if there is a prefix function associated with this type of token
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	// call the function if there is one
	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// enforces the correctness of the order of tokens by checking the type of the next token
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
