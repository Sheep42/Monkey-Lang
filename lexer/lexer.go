//lexer/lexer.go
package lexer

import "monkey/token"

type Lexer struct {
	input string		//The input
	position int 		//current position in input (points to char)
	readPosition int 	//reading pos in input (after the current position)
	ch byte				//the current char being examined
}

/** Lexer Methods **/
	//Reads line char by char and increments the Lexer position
	func (l *Lexer) readChar() {
		if l.readPosition >= len(l.input) {
			l.ch = 0
		} else {
			l.ch = l.input[l.readPosition]
		}

		l.position = l.readPosition
		l.readPosition += 1
	}

	//Returns the next character in input (without moving the current pos)
	func (l *Lexer) peekChar() byte {
		if l.readPosition >= len(l.input) {
			return 0
		} else {
			return l.input[l.readPosition]
		}
	}

	//Checks current character and returns a token for it, Returns ILLEGAL for non-mapped/non-accepted chars
	func (l *Lexer) NextToken() token.Token {
		var tok token.Token

		l.skipWhitespace()

		switch l.ch {
			case '=':
				if l.peekChar() == '=' {
					ch := l.ch

					l.readChar()

					literal := string(ch) + string(l.ch)

					tok = token.Token{
						Type: token.EQ,
						Literal: literal,
					}
				} else {
					tok = newToken(token.ASSIGN, l.ch)
				}
			case ';':
				tok = newToken(token.SEMI, l.ch)
			case '+':
				tok = newToken(token.PLUS, l.ch)
			case '-':
				tok = newToken(token.MINUS, l.ch)
			case '/':
				tok = newToken(token.SLASH, l.ch)
			case '*':
				tok = newToken(token.ASTERISK, l.ch)
			case '<':
				tok = newToken(token.LT, l.ch)
			case '>':
				tok = newToken(token.GT, l.ch)
			case '!':
				if l.peekChar() == '=' {
					ch := l.ch

					l.readChar()

					literal := string(ch) + string(l.ch)

					tok = token.Token{
						Type: token.NOT_EQ,
						Literal: literal,
					}
				} else {
					tok = newToken(token.BANG, l.ch)
				}
			case ',':
				tok = newToken(token.COMMA, l.ch)
			case '(':
				tok = newToken(token.LPAREN, l.ch)
			case ')':
				tok = newToken(token.RPAREN, l.ch)
			case '{':
				tok = newToken(token.LBRACE, l.ch)
			case '}':
				tok = newToken(token.RBRACE, l.ch)
			case 0:
				tok.Literal = ""
				tok.Type = token.EOF
			default:
				if isDigit(l.ch) { 
					//Lexes numbers
					tok.Type = token.INT
					tok.Literal = l.readNumber()

					return tok
				} else if isLetter(l.ch) {
					//Lexes keywords/user-defined identifiers
					tok.Literal = l.readIdentifier()
					tok.Type = token.LookupIdent(tok.Literal)

					return tok
				} else {
					tok = newToken(token.ILLEGAL, l.ch)
				}
		}

		l.readChar()

		return tok
	}

	//Reads an identifier and advances Lexer pos until a non-legal/whitespace character is encountered
	func (l *Lexer) readIdentifier() string {
		position := l.position

		for isLetter(l.ch) {
			l.readChar()
		}

		return l.input[position:l.position]
	}

	//Reads a number and advances Lexer pos until a non-number char is encountered
	func (l *Lexer) readNumber() string {
		position := l.position

		for isDigit(l.ch) {
			l.readChar()
		}

		return l.input[position:l.position]
	}

	//Eats whitespace
	func (l *Lexer) skipWhitespace() {
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
		}
	}
/** Utility Functions **/
	//Create a new Lexer
	func New(input string) *Lexer {
		l := &Lexer{input: input}

		l.readChar()

		return l
	}

	//Creates a new Token
	func newToken(tokenType token.TokenType, ch byte) token.Token {
		return token.Token{
			Type: tokenType, 
			Literal: string(ch),
		}
	}

	//Checks if a character is a letter
	func isLetter(ch byte) bool {
		return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || '0' <= ch && ch <= '9'
	}

	//Checks if a character is a number
	func isDigit(ch byte) bool {
		return '0' <= ch && ch <= '9'
	}