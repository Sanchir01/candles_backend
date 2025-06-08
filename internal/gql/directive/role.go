package directive

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/Sanchir01/candles_backend/internal/feature/user"
	"github.com/Sanchir01/candles_backend/internal/gql/model"
	customMiddleware "github.com/Sanchir01/candles_backend/internal/handlers/middleware"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"log/slog"
)

type RoleDirectiveFunc = func(ctx context.Context, obj interface{}, next graphql.Resolver, role []*model.Role) (res interface{}, err error)

func RoleDirective() RoleDirectiveFunc {
	return func(
		ctx context.Context,
		obj interface{},
		next graphql.Resolver,
		role []*model.Role,
	) (res interface{}, err error) {

		ctxUserID, err := customMiddleware.GetJWTClaimsFromCtx(ctx)
		if err != nil {
			return nil, &gqlerror.Error{Message: "Unauthorized"}
		}
		if ctxUserID == nil {
			return nil, &gqlerror.Error{Message: "Unauthorized: user ID is nil"}
		}

		if role == nil {
			return nil, &gqlerror.Error{Message: "Role is nil"}
		}
		slog.Warn("this directive role", ctxUserID.Role)
		slog.Warn("User role:", ctxUserID.Role, "Required role:", *role[0])
		boolRole := hasRole(ctxUserID, role)

		if boolRole == true {
			return next(ctx)
		}

		return nil, &gqlerror.Error{Message: "Role not admin"}
	}
}
func hasRole(ctxUser *user.Claims, role []*model.Role) bool {
	for _, r := range role {
		if ctxUser.Role == *r {
			return true
		}
	}
	return false
}
