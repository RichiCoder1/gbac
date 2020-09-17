package rego

import (
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/richicoder1/gbac/pkg/rego/providers"
)

// RegisterBuiltins defines the gbac builtins
func RegisterBuiltins() {
	rego.RegisterBuiltin4(
		&rego.Function{
			Name:    "gbac.can",
			Decl:    types.NewFunction(types.Args(types.S, types.S, types.S, types.S), types.B),
			Memoize: true,
		},
		func(bctx rego.BuiltinContext, userTerm, actionTerm, resourceTerm, parentTerm *ast.Term) (*ast.Term, error) {
			var user, action, resource, parent string

			if err := ast.As(userTerm.Value, &user); err != nil {
				return nil, err
			}
			if err := ast.As(actionTerm.Value, &action); err != nil {
				return nil, err
			}
			if err := ast.As(resourceTerm.Value, &resource); err != nil {
				return nil, err
			}
			if err := ast.As(parentTerm.Value, &parent); err != nil {
				return nil, err
			}

			provider, err := providers.GetProvider(&bctx)
			if err != nil {
				return nil, err
			}

			allow, err := provider.Can(user, action, resource, parent)
			if err != nil {
				return nil, err
			}
			return ast.BooleanTerm(allow), nil
		},
	)
}
