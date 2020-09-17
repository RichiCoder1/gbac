package providers

import (
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
	r "github.com/richicoder1/gbac/pkg/rego/providers/redis"
)

type Provider interface {
	Can(user, action, resource, owner string) (bool, error)
}

func GetProvider(ctx *rego.BuiltinContext) (Provider, error) {
	config := ctx.Runtime.Get(ast.StringTerm("config")).Get(ast.StringTerm("gbac"))
	if config == nil {
		return nil, errors.Errorf("Missing config for plugin gbac")
	}
	switch {
	case config.Get(ast.StringTerm("redis")) != nil:
		redisTerm := config.Get(ast.StringTerm("redis"))
		rawConfig, err := redisTerm.MarshalJSON()
		if err != nil {
			return nil, err
		}
		return r.New(rawConfig)
	}
	return nil, errors.Errorf("Unable to find provider for given config.")
}
