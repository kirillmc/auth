package access

import (
	"context"

	"github.com/kirillmc/auth/internal/model"
	"github.com/pkg/errors"
)

var accessibleRoles map[string]model.Role

func (s *serv) accessibleRoles(ctx context.Context) (map[string]model.Role, error) {
	if accessibleRoles == nil {
		accessibleRolesTemp, err := s.userRepository.GetAccessibleRoles(ctx)
		if err != nil {
			accessibleRoles = make(map[string]model.Role)
			return accessibleRoles, errors.New("failed to get endpoints-roles data")
		}

		accessibleRoles = accessibleRolesTemp
	}
	return accessibleRoles, nil
}
