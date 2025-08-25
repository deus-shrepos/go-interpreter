package interpreter

import (
	"testing"

	"github.com/go-interpreter/internal/interpreter"
	"github.com/go-interpreter/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment_NewEnvironment(t *testing.T) {
	tests := []struct {
		name      string
		enclosing *interpreter.Environment
		wantNil   bool
	}{
		{
			name:      "global scope",
			enclosing: nil,
			wantNil:   true,
		},
		{
			name:      "inner scope",
			enclosing: interpreter.NewEnvironment(nil),
			wantNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := interpreter.NewEnvironment(tt.enclosing)
			assert.NotNil(t, env)
			assert.NotNil(t, env.Values)
			if tt.wantNil {
				assert.Nil(t, env.Enclosing)
			} else {
				assert.NotNil(t, env.Enclosing)
			}
		})
	}
}

func TestEnvironment_Define(t *testing.T) {
	tests := []struct {
		name     string
		varName  string
		value    any
		wantType any
	}{
		{
			name:     "string value",
			varName:  "str",
			value:    "test",
			wantType: "",
		},
		{
			name:     "integer value",
			varName:  "num",
			value:    42,
			wantType: 0,
		},
		{
			name:     "nil value",
			varName:  "empty",
			value:    nil,
			wantType: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := interpreter.NewEnvironment(nil)
			env.Define(tt.varName, tt.value)

			got, exists := env.Values[tt.varName]
			assert.True(t, exists)
			assert.IsType(t, tt.wantType, got)
			assert.Equal(t, tt.value, got)
		})
	}
}

func TestEnvironment_Get(t *testing.T) {
	tests := []struct {
		name      string
		setupEnv  func() *interpreter.Environment
		token     token.Token
		wantValue any
		wantErr   bool
	}{
		{
			name: "existing variable in current scope",
			setupEnv: func() *interpreter.Environment {
				env := interpreter.NewEnvironment(nil)
				env.Define("x", 42)
				return env
			},
			token: token.Token{
				Lexeme: "x",
				Line:   1,
				Char:   1,
			},
			wantValue: 42,
			wantErr:   false,
		},
		{
			name: "existing variable in outer scope",
			setupEnv: func() *interpreter.Environment {
				outer := interpreter.NewEnvironment(nil)
				outer.Define("x", 42)
				return interpreter.NewEnvironment(outer)
			},
			token: token.Token{
				Lexeme: "x",
				Line:   1,
				Char:   1,
			},
			wantValue: 42,
			wantErr:   false,
		},
		{
			name: "undefined variable",
			setupEnv: func() *interpreter.Environment {
				return interpreter.NewEnvironment(nil)
			},
			token: token.Token{
				Lexeme: "undefined",
				Line:   1,
				Char:   1,
			},
			wantValue: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := tt.setupEnv()
			got, err := env.Get(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, got)
			}
		})
	}
}

func TestEnvironment_Assign(t *testing.T) {
	tests := []struct {
		name     string
		setupEnv func() *interpreter.Environment
		token    token.Token
		value    any
		wantErr  bool
	}{
		{
			name: "assign to existing variable",
			setupEnv: func() *interpreter.Environment {
				env := interpreter.NewEnvironment(nil)
				env.Define("x", 42)
				return env
			},
			token: token.Token{
				Lexeme: "x",
				Line:   1,
				Char:   1,
			},
			value:   100,
			wantErr: false,
		},
		{
			name: "assign to variable in outer scope",
			setupEnv: func() *interpreter.Environment {
				outer := interpreter.NewEnvironment(nil)
				outer.Define("x", 42)
				return interpreter.NewEnvironment(outer)
			},
			token: token.Token{
				Lexeme: "x",
				Line:   1,
				Char:   1,
			},
			value:   100,
			wantErr: false,
		},
		{
			name: "assign to undefined variable",
			setupEnv: func() *interpreter.Environment {
				return interpreter.NewEnvironment(nil)
			},
			token: token.Token{
				Lexeme: "undefined",
				Line:   1,
				Char:   1,
			},
			value:   100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := tt.setupEnv()
			err := env.Assign(tt.token, tt.value)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				got, err := env.Get(tt.token)
				assert.NoError(t, err)
				assert.Equal(t, tt.value, got)
			}
		})
	}
}
