package user

import (
	"github.com/bookpanda/cas-sso/backend/internal/dto"
	"github.com/bookpanda/cas-sso/backend/internal/model"
)

func ModelToDto(in *model.User) *dto.User {
	return &dto.User{
		ID:        in.ID.String(),
		Email:     in.Email,
		Firstname: in.Firstname,
		Lastname:  in.Lastname,
	}
}

func UpdateRequestToModel(in *dto.UpdateUserProfileRequest) (*model.User, error) {
	user := &model.User{
		Firstname: in.Firstname,
		Lastname:  in.Lastname,
	}

	return user, nil
}

func UpdateUserBodyToRequest(id string, body *dto.UpdateUserProfileBody) *dto.UpdateUserProfileRequest {
	return &dto.UpdateUserProfileRequest{
		ID:        id,
		Firstname: body.Firstname,
		Lastname:  body.Lastname,
	}
}
