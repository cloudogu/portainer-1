package endpointgroups

import (
	"errors"
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

// @id EndpointGroupDelete
// @summary Remove an environment(endpoint) group
// @description Remove an environment(endpoint) group.
// @description **Access policy**: administrator
// @tags endpoint_groups
// @security ApiKeyAuth
// @security jwt
// @param id path int true "EndpointGroup identifier"
// @success 204 "Success"
// @failure 400 "Invalid request"
// @failure 404 "EndpointGroup not found"
// @failure 500 "Server error"
// @router /endpoint_groups/{id} [delete]
func (handler *Handler) endpointGroupDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	endpointGroupID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid environment group identifier route variable", err)
	}

	if endpointGroupID == 1 {
		return httperror.Forbidden("Unable to remove the default 'Unassigned' group", errors.New("Cannot remove the default environment group"))
	}

	endpointGroup, err := handler.DataStore.EndpointGroup().EndpointGroup(portainer.EndpointGroupID(endpointGroupID))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find an environment group with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("Unable to find an environment group with the specified identifier inside the database", err)
	}

	err = handler.DataStore.EndpointGroup().DeleteEndpointGroup(portainer.EndpointGroupID(endpointGroupID))
	if err != nil {
		return httperror.InternalServerError("Unable to remove the environment group from the database", err)
	}

	endpoints, err := handler.DataStore.Endpoint().Endpoints()
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve environment from the database", err)
	}

	for _, endpoint := range endpoints {
		if endpoint.GroupID == portainer.EndpointGroupID(endpointGroupID) {
			endpoint.GroupID = portainer.EndpointGroupID(1)
			err = handler.DataStore.Endpoint().UpdateEndpoint(endpoint.ID, &endpoint)
			if err != nil {
				return httperror.InternalServerError("Unable to update environment", err)
			}

			err = handler.updateEndpointRelations(&endpoint, nil)
			if err != nil {
				return httperror.InternalServerError("Unable to persist environment relations changes inside the database", err)
			}
		}
	}

	for _, tagID := range endpointGroup.TagIDs {
		err = handler.DataStore.Tag().UpdateTagFunc(tagID, func(tag *portainer.Tag) {
			delete(tag.EndpointGroups, endpointGroup.ID)
		})

		if handler.DataStore.IsErrObjectNotFound(err) {
			return httperror.InternalServerError("Unable to find a tag inside the database", err)
		} else if err != nil {
			return httperror.InternalServerError("Unable to persist tag changes inside the database", err)
		}
	}

	return response.Empty(w)
}
