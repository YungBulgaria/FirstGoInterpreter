package lexer

import "donkey/token"

type Lexer struct {
	input		 string
	position 	 int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch			 byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input} // Create a new Lexer called l and initialize l.input with input
	l.readChar() // Initialize l.position, l.readPosition, l.ch with the l.readChar() function 
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token // In C++ this would be like saying Token::token

	l.skipWhitespace() // Since white space is not important call a function to skip white space and start on the next valid char

	switch l.ch { // A series of conditionals where if the char of l is equal to one of this conditionals we enter that branch and set the token appropriatley 
	case '=': // Example if l.ch == '=' then we enter this branch
		if l.peekChar() == '=' { // This is a special check for '=' which we also do for '!' in order to see if the following character is also a '='
			ch := l.ch // We save the current character in ch
			l.readChar() // We move our pointer to the next char updating l in the proccess
			literal := string(ch) + string(l.ch) // We make a variable literal which concatinates the two '=' found in ch and l.ch
			tok = token.Token{Type: token.EQ, Literal: literal} // We create the token initializing the token type to EQ and the literal to '=='
		} else { // We are unable to use newToken() ^ here since literal is a string and not a char which is only one byte
			tok = newToken(token.ASSIGN, l.ch) // This does the same thing as the previous line except through the newToken function
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' { // See explanaition under case '='
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
		tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default: // If none of the cases apply the default branch is entered
		if isLetter(l.ch) { // If the current char is a letter (checked by isLetter() defined below) then we enter this branch
			tok.Literal = l.readIdentifier() // The token literal is assigned by the function readIdentifier() defined below
			tok.Type = token.LookupIdent(tok.Literal) // The token type is assigned by the token.LookupIdent() function defined in token.go -- if the token is a keyword...
			// ...then the appropriate keyword type gets assigned (e.g) FUNCTION, LET, RETURN otherwise it is just IDENT
			return tok
		} else if isDigit(l.ch) { // Same concept here but for digits
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else { // If it is neither a char or a digit it must be an illegal token so we return such
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar() // Read the next char
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position // Get current position 
	for isLetter(l.ch) {
		l.readChar() // Check if letter, if so read it
	}
	return l.input[position:l.position] // Subset the input string from the current position to the end of the letter AKA obtain current letter as identifier
}

func isLetter(ch byte) bool { // Checks if ch is a letter
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string { // Same basic functionality as readIdentifier() but for numbers
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool { // Same functionality as isLetter() but for digits
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token { // Creates a new token to return in the switch
	return token.Token{Type: tokenType, Literal:string(ch)}
}

func (l *Lexer) readChar() { // Same functionality as readIdentifier() but for chars -- used for peeking into next char.
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte { // Peeks at the char after the current one
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) skipWhitespace() { // Skips white space
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

