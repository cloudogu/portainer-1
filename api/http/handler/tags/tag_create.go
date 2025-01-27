package tags

import (
	"errors"
	"net/http"

	"github.com/asaskevich/govalidator"
	portainer "github.com/cloudogu/portainer-ce/api"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

type tagCreatePayload struct {
	// Name
	Name string `validate:"required" example:"org/acme"`
}

func (payload *tagCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.Name) {
		return errors.New("Invalid tag name")
	}
	return nil
}

// @id TagCreate
// @summary Create a new tag
// @description Create a new tag.
// @description **Access policy**: administrator
// @tags tags
// @security ApiKeyAuth
// @security jwt
// @accept json
// @produce json
// @param body body tagCreatePayload true "Tag details"
// @success 200 {object} portainer.Tag "Success"
// @failure 409 "Tag name exists"
// @failure 500 "Server error"
// @router /tags [post]
func (handler *Handler) tagCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload tagCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return httperror.BadRequest("Invalid request payload", err)
	}

	tags, err := handler.DataStore.Tag().Tags()
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve tags from the database", err)
	}

	for _, tag := range tags {
		if tag.Name == payload.Name {
			return &httperror.HandlerError{StatusCode: http.StatusConflict, Message: "This name is already associated to a tag", Err: errors.New("A tag already exists with this name")}
		}
	}

	tag := &portainer.Tag{
		Name:           payload.Name,
		EndpointGroups: map[portainer.EndpointGroupID]bool{},
		Endpoints:      map[portainer.EndpointID]bool{},
	}

	err = handler.DataStore.Tag().Create(tag)
	if err != nil {
		return httperror.InternalServerError("Unable to persist the tag inside the database", err)
	}

	return response.JSON(w, tag)
}
