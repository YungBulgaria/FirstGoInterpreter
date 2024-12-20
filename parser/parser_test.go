package parser

import (
	"testing"
	"donkey/ast"
	"donkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x  = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input) // Initialize a lexer with the input above
	p := New(l) // Initialize the parser with the lexer

	program := p.ParseProgram() // Initialize program with p.ParseProgram() defined in parser.go:35
	checkParserErrors(t, p) // Call the check error function defined on line 73.

	if program == nil { // If the program is nil return an error
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 { // If the length of program.Statements does not match up with the amount of statements in the input you provided (3 in this case) return an error
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct { // Define tests to confirm proper execution
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests { // Iterate through each expectedIdentifier in tests 
		stmt := program.Statements[i] // Set a var called stmt = program.Statements[i] which will return the ith statement generated by running p.ParseProgram()

		if !testLetStatement(t, stmt, tt.expectedIdentifier) { // Calls testLetStatement() defined below
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool { // Takes a test, statement, and a string (in our case the expectedIdentifier from tests)
	if s.TokenLiteral() != "let" { // If the token is not "let" which should be the first token of any of our statements
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral()) // Return an error
		return false // Fail the test
	}

	letStmt, ok := s.(*ast.LetStatement) // Check to see if ast LetStatement var is properly initialized? Not 100% sure on ,ok syntax...
	// Okay I think I understand what's happening up here.
	// We assign letStmt to s.(*ast.LetStatement) which takes s which we know is program.Statements[i]
	// program.Statements[i] returns a statement object from the statements array in our program struct
	// so we say s.(*ast.LetStatement) which effecively means statment.(ast::LetStatment)...
	// ...or just statement.LetStatment (99% sure)...
	// ...SOOOO this casts? letStmt to be a LetStatement object as defined by the LetStatement struct in ast.go:23
	// This lets us properly check all the LetStatment vars like Token, Name, and Value
	// And yeah I think ok just means that it was correctly intialized as a LetStatement and not some other statement.
	if !ok {							 // ...but I think it is just checking that the variable was properly initialized.
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name { // Check to see if letStmt name == name from tests
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name { // Check to see if the TokenLiteral == name from tests
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}

	return true // Return true meaning passed all test cases
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", 
				returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", 
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", 
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func checkParserErrors(t *testing.T, p *Parser) { // Checks for any errors during parsing WITHOUT halting execution upon hitting an error 
	errors := p.Errors() // Initializes error var with p.Errors() described in parser.go:102
	if len(errors) == 0 { // No errors so return nothing
		return
	}

	t.Errorf("parser has %d errors", len(errors)) // Return the amount of errors based on the length of the error array stored in the parser struct
	for _, msg := range errors { // For each error and its message, print error and its message
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow() // Fail the test
}