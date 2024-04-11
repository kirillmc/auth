package user

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/kirillmc/auth/internal/model"
	"github.com/kirillmc/platform_common/pkg/db"
)

const (
	accessRolesBd  = "roles_to_endpoints"
	endpointColumn = "endpoint"
)

func (r *repo) GetAccessibleRoles(ctx context.Context) (map[string]model.Role, error) {
	builder := sq.Select(endpointColumn, roleColumn).
		PlaceholderFormat(sq.Dollar).
		From(accessRolesBd)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Print("Error in repo layer1")
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetAccessibleRoles",
		QueryRaw: query,
	}

	accessibleRoles := make(map[string]model.Role)
	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		log.Print("Error in repo layer1")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var endpoint string
		var role int

		err = rows.Scan(&endpoint, &role)
		if err != nil {
			log.Print("Error in repo layer1")
			return nil, err
		}
		accessibleRoles[endpoint] = model.Role(role)
	}

	return accessibleRoles, nil
}
