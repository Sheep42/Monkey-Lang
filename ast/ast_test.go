package ast

import (
	"testing"

	"github.com/Sheep42/Monkey-Lang/token"
)

func TestString(t *testing.T) {

	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {

		t.Errorf("program.String() failed: got=%q", program.String())

	}

}
