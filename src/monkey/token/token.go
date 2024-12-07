package token

// Make a struct of a token 
type TokenType string

type Token struct {
	Type TokenType // Allows us to determine the type of token.
	Literal string
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