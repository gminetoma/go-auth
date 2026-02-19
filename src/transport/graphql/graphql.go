package graphql

import (
	"net/http"

	authApplication "github.com/gminetoma/go-auth/src/auth/application"

	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

type SetupGraphQLParams struct {
	AuthService         authApplication.AuthService
	EnablePlayground    bool
	EnableIntrospection bool
}

func SetupGraphQL(mux *http.ServeMux, params SetupGraphQLParams) {
	handler := gqlHandler.New(NewExecutableSchema(Config{
		Resolvers: &Resolver{
			AuthService: params.AuthService,
		},
	}))

	handler.AddTransport(transport.Options{})
	handler.AddTransport(transport.POST{})

	if params.EnableIntrospection {
		handler.Use(extension.Introspection{})
	}

	mux.Handle("/graphql", handler)
	mux.Handle("/playground", playground.Handler("GraphQL", "/graphql"))
}
