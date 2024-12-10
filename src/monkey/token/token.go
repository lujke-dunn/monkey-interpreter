package token

// Make a struct of a token 
type TokenType string

type Token struct {
	Type TokenType // Allows us to determine the type of token.
	Literal string
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	} 
	return IDENT
}

var keywords = map[string]TokenType{ // a collection of all valid keyword types
	"fn": FUNCTION,
	"let": LET,
}

// Types of identifiers that our token will recognise 
const (
	// Special tokens 
	ILLEGAL = "ILLEGAL"
	EOF = "EOF"

	// Identifiers and literal tokens 
    IDENT = "IDENT" // Variable names x, y, z 
	INT = "INT" // Numbers: 12345

	// Operators 
	ASSIGN = "="
	PLUS = "+" 
	MINUS = "-"
	SLASH = "/"
	BANG = "!"
	ASTERISK = "*"

	LT = "<"
	GT = ">"

	// Delimiters
	COMMA = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"

	LBRACE = "{"
	RBRACE = "}"

	// Keywords 
	FUNCTION = "FUNCTION"
	LET = "LET"
)