package edgestacks

import (
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/http/middlewares"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

// @id EdgeStackStatusDelete
// @summary Delete an EdgeStack status
// @description Authorized only if the request is done by an Edge Environment(Endpoint)
// @tags edge_stacks
// @produce json
// @param id path string true "EdgeStack Id"
// @success 200 {object} portainer.EdgeStack
// @failure 500
// @failure 400
// @failure 404
// @failure 403
// @router /edge_stacks/{id}/status/{endpoint_id} [delete]
func (handler *Handler) edgeStackStatusDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	stackID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid stack identifier route variable", err)
	}

	endpoint, err := middlewares.FetchEndpoint(r)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve a valid endpoint from the handler context", err)
	}

	err = handler.requestBouncer.AuthorizedEdgeEndpointOperation(r, endpoint)
	if err != nil {
		return httperror.Forbidden("Permission denied to access environment", err)
	}

	var edgeStack *portainer.EdgeStack
	err = handler.DataStore.EdgeStack().UpdateEdgeStackFunc(portainer.EdgeStackID(stackID), func(stack *portainer.EdgeStack) {
		delete(stack.Status, endpoint.ID)
		edgeStack = stack
	})
	if err != nil {
		return handler.handlerDBErr(err, "Unable to persist the stack changes inside the database")
	}

	return response.JSON(w, edgeStack)
}
