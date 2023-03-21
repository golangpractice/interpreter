package lexer

import "monkey/token"

// A Lexer reads the input string and parses it, keeping track of the current char and the next char positions
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// New will initialize and return a new Lexer
func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	lexer.readChar()
	return lexer
}

// NextToken will go through the input to the lexer and return the next token
func (lexer *Lexer) NextToken() token.Token {
	var tok token.Token

	lexer.skipWhitespace()

	switch lexer.ch {
	case '=':
		if lexer.peekChar() == '=' {
			ch := lexer.ch
			lexer.readChar()
			literal := string(ch) + string(lexer.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, lexer.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, lexer.ch)
	case '(':
		tok = newToken(token.LPAREN, lexer.ch)
	case ')':
		tok = newToken(token.RPAREN, lexer.ch)
	case ',':
		tok = newToken(token.COMMA, lexer.ch)
	case '+':
		tok = newToken(token.PLUS, lexer.ch)
	case '-':
		tok = newToken(token.MINUS, lexer.ch)
	case '!':
		if lexer.peekChar() == '=' {
			ch := lexer.ch
			lexer.readChar()
			literal := string(ch) + string(lexer.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, lexer.ch)
		}
	case '/':
		tok = newToken(token.SLASH, lexer.ch)
	case '*':
		tok = newToken(token.ASTERISK, lexer.ch)
	case '<':
		tok = newToken(token.LT, lexer.ch)
	case '>':
		tok = newToken(token.GT, lexer.ch)
	case '{':
		tok = newToken(token.LBRACE, lexer.ch)
	case '}':
		tok = newToken(token.RBRACE, lexer.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(lexer.ch) {
			tok.Literal = lexer.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lexer.ch) {
			tok.Type = token.INT
			tok.Literal = lexer.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, lexer.ch)
		}
	}

	lexer.readChar()
	return tok
}

// readChar will give us the next character and advance our position in the input string
//
// First, we check whether we have reached the end of input. If so, set lexer.ch = 0 which is the ASCII code for the
// "NUL" character. This signifies that either nothing has been read yet, or we arrived at the "end of file"
func (lexer *Lexer) readChar() {
	if lexer.readPosition >= len(lexer.input) {
		lexer.ch = 0
	} else {
		lexer.ch = lexer.input[lexer.readPosition]
	}
	lexer.position = lexer.readPosition
	lexer.readPosition += 1
}

// peekChar returns the next char in the input string without incrementing position or readPosition
func (lexer *Lexer) peekChar() byte {
	if lexer.readPosition >= len(lexer.input) {
		return 0
	} else {
		return lexer.input[lexer.readPosition]
	}
}

// readIdentifier returns the current identifier string being read from the input
//
// This is done by storing the initial input position and then repeatedly calling isLetter. As long as
// isLetter returns TRUE, then readChar is called to advance the position in the input string. Once isLetter returns
// FALSE, we will return the string of characters between the initial input position and the current one.
func (lexer *Lexer) readIdentifier() string {
	position := lexer.position
	for isLetter(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

// readNumber returns the current digit string being read from the input
//
// This is done by storing the initial input position and then repeatedly calling isDigit. As long as
// isDigit returns TRUE, then readChar is called to advance the position in the input string. Once isDigit returns
// FALSE, we will return the string of characters between the initial input position and the current one.
func (lexer *Lexer) readNumber() string {
	position := lexer.position
	for isDigit(lexer.ch) {
		lexer.readChar()
	}
	return lexer.input[position:lexer.position]
}

// skipWhitespace will skip over whitespace characters in the given input to the Lexer.
//
// If a given input string contains whitespace, then readChar will be called repeatedly to advance our position
// in the input string until we have passed over all the whitespace.
func (lexer *Lexer) skipWhitespace() {
	for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
		lexer.readChar()
	}
}

// isLetter returns TRUE if the given ch param is a letter or and underscore "_", else FALSE
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit returns TRUE if the given ch param is numeric, else FALSE
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// newToken will return a new Token given a TokenType and literal as params
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
