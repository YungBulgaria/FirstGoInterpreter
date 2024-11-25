package parser

import (
	"donkey/ast"
	"donkey/lexer"
	"donkey/token"
	"fmt"
)

type Parser struct {
	l *lexer.Lexer
	errors []string
	curToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser { // Create a new parser object initializing l with the lexer object we created 
	p := &Parser{
		l: 		l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() { // Sets curToken to peekToken and peekToken to p.l.NextToken()
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // Creates var program which is an ast::Program object
	program.Statements = []ast.Statement{} // Initialize the program statements var with an empty arr

	for p.curToken.Type != token.EOF { // If the current token is not an EOF then enter this block
		stmt := p.parseStatement() // Set the current statement to the result of parseStatement()
		if stmt != nil { // If p.parseStatement() was not nill then we...
			program.Statements = append(program.Statements, stmt) // ...append the statement to the program.Statements var
		}
		p.nextToken() // Call the nextToken() function to update p.curToken and p.PeekToken()
	}
	return program // Return program which contains an arr of statements in the general order (as of now) LET, IDENT, ... , SEMICOLON
}

func (p *Parser) parseStatement() ast.Statement { 
	switch p.curToken.Type { // Sets curToken.Type based on a switch
	case token.LET: // If it is a LET statement then we return p.parseLetStatement() which is defined below
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token :p.curToken} // Sets the stmt var to the ast::LetStatement.Token which is initialized with p.curToken

	if !p.expectPeek(token.IDENT) { // If the token directly after the LetStatement is not an IDENT token then return nil 
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal} // Set the current statement's name while initializing the Identifier object with... 
																			  // ...Token = p.curToken (same as before) and Value = p.curToken.Literal (self explanatory)

	if !p.expectPeek(token.ASSIGN) { // If the token after the identifier is not an ASSIGN token then we return nil
		return nil
	}

	// TODO: We're skipping the expressions until we 
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) { // Skips expressions and goes to EOL if it is not a semicolon we fetch the next token and recursively repeat this whole process
		p.nextToken()
	}

	return stmt // If the above line evals to true then we just return the statement and resume operations above
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token :p.curToken}

	p.nextToken()
	
	// TODO: We're skipping the expressions until we
	// encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}


 // Helper functions for executing the parsing above 
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
// Helper functions for exectuing the parsing above

func (p *Parser) Errors() []string { // Returns the erros in parser
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) { // Simple function that checks the peeked token to see if it doesn't match up with its expected type
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", 
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg) // Appends error message to the error var in the parser struct (see further details in parser_test.go)
}