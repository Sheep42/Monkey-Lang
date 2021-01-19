//token/token.go

package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

//Define our token types
const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//Identifiers + Literals
	IDENT  = "IDENT"  //add, foobar, x, y ...
	INT    = "INT"    //Integer literal
	STRING = "STRING" // String literal

	//Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	//Delimiters
	COMMA = ","
	SEMI  = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

//Define language keywords/map them to their token type
var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

/** Utility Functions **/
//Look up identifiers in the keywords map
func LookupIdent(ident string) TokenType {
	//If keyword exists return matching token type
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	//Otherwise return the IDENT token type
	return IDENT
}
