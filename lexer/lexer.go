package lexer

import "interp/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LEX_GE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LEX_GT, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LEX_LE, Literal: string(ch) + string(l.ch)}
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LEX_NE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LEX_LT, l.ch)
		}
	case ';':
		tok = newToken(token.LEX_SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.LEX_COMMA, l.ch)
	case '(':
		tok = newToken(token.LEX_LPAREN, l.ch)
	case ')':
		tok = newToken(token.LEX_RPAREN, l.ch)
	case '+':
		tok = newToken(token.LEX_PLUS, l.ch)
	case '-':
		tok = newToken(token.LEX_MIN, l.ch)
	case '*':
		tok = newToken(token.LEX_MULT, l.ch)
	case '/':
		tok = newToken(token.LEX_DIV, l.ch)
	case '=':
		tok = newToken(token.LEX_EQ, l.ch)
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LEX_ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LEX_COLON, l.ch)
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.LEX_EOF

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookUpIdent(tok.Literal)
			return tok

		} else if isDigit(l.ch) {
			tok.Type = token.LEX_INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.LEX_ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == 'A' || l.ch == 'a' || l.ch == 'B' || l.ch == 'b' || l.ch == 'C' ||
		l.ch == 'c' || l.ch == 'H' || l.ch == 'h' || l.ch == 'D' || l.ch == 'd' || l.ch == 'E' || l.ch == 'e' || l.ch == 'F' || l.ch == 'f' {

		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
