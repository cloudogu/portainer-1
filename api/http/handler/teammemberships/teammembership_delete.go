package teammemberships

import (
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/http/errors"
	"github.com/cloudogu/portainer-ce/api/http/security"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

// @id TeamMembershipDelete
// @summary Remove a team membership
// @description Remove a team membership. Access is only available to administrators leaders of the associated team.
// @description **Access policy**: administrator
// @tags team_memberships
// @security ApiKeyAuth
// @security jwt
// @param id path int true "TeamMembership identifier"
// @success 204 "Success"
// @failure 400 "Invalid request"
// @failure 403 "Permission denied"
// @failure 404 "TeamMembership not found"
// @failure 500 "Server error"
// @router /team_memberships/{id} [delete]
func (handler *Handler) teamMembershipDelete(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	membershipID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid membership identifier route variable", err)
	}

	membership, err := handler.DataStore.TeamMembership().TeamMembership(portainer.TeamMembershipID(membershipID))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find a team membership with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("Unable to find a team membership with the specified identifier inside the database", err)
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve info from request context", err)
	}

	if !security.AuthorizedTeamManagement(membership.TeamID, securityContext) {
		return httperror.Forbidden("Permission denied to delete the membership", errors.ErrResourceAccessDenied)
	}

	err = handler.DataStore.TeamMembership().DeleteTeamMembership(portainer.TeamMembershipID(membershipID))
	if err != nil {
		return httperror.InternalServerError("Unable to remove the team membership from the database", err)
	}

	return response.Empty(w)
}
