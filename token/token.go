package token

const (
	LEX_ILLEGAL = "ILLEGAL"
	LEX_EOF     = "EOF"

	LEX_IDENT = "IDENT" // add, foobar, x, y, ...
	LEX_INT   = "INT"   // 1343456B, 1343456b, 1343456C, 1343456c, 1343456D, 1343456H
	LEX_FLOAT = "FLOAT" // 3.14159, 1.5E+10, .25E-5

	// Operators
	LEX_ASSIGN = ":="
	LEX_PLUS   = "+"
	LEX_MIN    = "-"
	LEX_MULT   = "*"
	LEX_DIV    = "/"
	LEX_GT     = ">"
	LEX_LT     = "<"
	LEX_EQ     = "="
	LEX_NE     = "<>"
	LEX_LE     = "<="
	LEX_GE     = ">="

	// Delimiters
	LEX_COMMA     = ","
	LEX_SEMICOLON = ";"
	LEX_COLON     = ":"

	LEX_LPAREN = "("
	LEX_RPAREN = ")"

	// Keywords
	KW_GOTO    = "GOTO"
	KW_INTEGER = "INTEGER"
	KW_REAL    = "REAL"
	KW_READ    = "READ"
	KW_WRITE   = "WRITE"
	KW_IF      = "IF"
	KW_ELSE    = "ELSE"
	KW_THEN    = "THEN"
	KW_END     = "END"
	KW_LOOP    = "LOOP"
	KW_BEGIN   = "BEGIN"
	KW_SKIP    = "SKIP"
	KW_SPACE   = "SPACE"
	KW_TAB     = "TAB"
	KW_MOD     = "MOD"
	KW_OF      = "of"
)

var keywords = map[string]TokenType{
	"integer": KW_INTEGER,
	"real":    KW_REAL,
	"read":    KW_READ,
	"goto":    KW_GOTO,
	"if":      KW_IF,
	"else":    KW_ELSE,
	"write":   KW_WRITE,
	"then":    KW_THEN,
	"end":     KW_END,
	"loop":    KW_LOOP,
	"begin":   KW_BEGIN,
	"skip":    KW_SKIP,
	"tab":     KW_TAB,
	"space":   KW_SPACE,
	"mod":     KW_MOD,
	"of":      KW_OF,
}

func LookUpIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return LEX_IDENT
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
