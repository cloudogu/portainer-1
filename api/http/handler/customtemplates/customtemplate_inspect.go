package customtemplates

import (
	"net/http"
	"strconv"

	portainer "github.com/cloudogu/portainer-ce/api"
	httperrors "github.com/cloudogu/portainer-ce/api/http/errors"
	"github.com/cloudogu/portainer-ce/api/http/security"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

// @id CustomTemplateInspect
// @summary Inspect a custom template
// @description Retrieve details about a template.
// @description **Access policy**: authenticated
// @tags custom_templates
// @security ApiKeyAuth
// @security jwt
// @produce json
// @param id path int true "Template identifier"
// @success 200 {object} portainer.CustomTemplate "Success"
// @failure 400 "Invalid request"
// @failure 404 "Template not found"
// @failure 500 "Server error"
// @router /custom_templates/{id} [get]
func (handler *Handler) customTemplateInspect(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	customTemplateID, err := request.RetrieveNumericRouteVariableValue(r, "id")
	if err != nil {
		return httperror.BadRequest("Invalid Custom template identifier route variable", err)
	}

	customTemplate, err := handler.DataStore.CustomTemplate().CustomTemplate(portainer.CustomTemplateID(customTemplateID))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find a custom template with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("Unable to find a custom template with the specified identifier inside the database", err)
	}

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve user info from request context", err)
	}

	resourceControl, err := handler.DataStore.ResourceControl().ResourceControlByResourceIDAndType(strconv.Itoa(customTemplateID), portainer.CustomTemplateResourceControl)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve a resource control associated to the custom template", err)
	}

	access := userCanEditTemplate(customTemplate, securityContext)
	if !access {
		return httperror.Forbidden("Access denied to resource", httperrors.ErrResourceAccessDenied)
	}

	if resourceControl != nil {
		customTemplate.ResourceControl = resourceControl
	}

	return response.JSON(w, customTemplate)
}
