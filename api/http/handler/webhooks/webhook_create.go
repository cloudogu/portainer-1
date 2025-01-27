package webhooks

import (
	"errors"
	"net/http"

	"github.com/cloudogu/portainer-ce/api/http/security"
	"github.com/cloudogu/portainer-ce/api/internal/registryutils/access"

	portainer "github.com/cloudogu/portainer-ce/api"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"

	"github.com/asaskevich/govalidator"
	"github.com/gofrs/uuid"
)

type webhookCreatePayload struct {
	ResourceID  string
	EndpointID  int
	RegistryID  portainer.RegistryID
	WebhookType int
}

func (payload *webhookCreatePayload) Validate(r *http.Request) error {
	if govalidator.IsNull(payload.ResourceID) {
		return errors.New("Invalid ResourceID")
	}
	if payload.EndpointID == 0 {
		return errors.New("Invalid EndpointID")
	}
	if payload.WebhookType != 1 {
		return errors.New("Invalid WebhookType")
	}
	return nil
}

// @summary Create a webhook
// @description **Access policy**: authenticated
// @security ApiKeyAuth
// @security jwt
// @tags webhooks
// @accept json
// @produce json
// @param body body webhookCreatePayload true "Webhook data"
// @success 200 {object} portainer.Webhook
// @failure 400
// @failure 409
// @failure 500
// @router /webhooks [post]
func (handler *Handler) webhookCreate(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	var payload webhookCreatePayload
	err := request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return httperror.BadRequest("Invalid request payload", err)
	}

	webhook, err := handler.DataStore.Webhook().WebhookByResourceID(payload.ResourceID)
	if err != nil && !handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.InternalServerError("An error occurred retrieving webhooks from the database", err)
	}
	if webhook != nil {
		return &httperror.HandlerError{StatusCode: http.StatusConflict, Message: "A webhook for this resource already exists", Err: errors.New("A webhook for this resource already exists")}
	}

	endpointID := portainer.EndpointID(payload.EndpointID)

	securityContext, err := security.RetrieveRestrictedRequestContext(r)
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve user info from request context", err)
	}

	if !securityContext.IsAdmin {
		return httperror.Forbidden("Not authorized to create a webhook", errors.New("not authorized to create a webhook"))
	}

	if payload.RegistryID != 0 {
		tokenData, err := security.RetrieveTokenData(r)
		if err != nil {
			return httperror.InternalServerError("Unable to retrieve user authentication token", err)
		}

		_, err = access.GetAccessibleRegistry(handler.DataStore, tokenData.ID, endpointID, payload.RegistryID)
		if err != nil {
			return httperror.Forbidden("Permission deny to access registry", err)
		}
	}

	token, err := uuid.NewV4()
	if err != nil {
		return httperror.InternalServerError("Error creating unique token", err)
	}

	webhook = &portainer.Webhook{
		Token:       token.String(),
		ResourceID:  payload.ResourceID,
		EndpointID:  endpointID,
		RegistryID:  payload.RegistryID,
		WebhookType: portainer.WebhookType(payload.WebhookType),
	}

	err = handler.DataStore.Webhook().Create(webhook)
	if err != nil {
		return httperror.InternalServerError("Unable to persist the webhook inside the database", err)
	}

	return response.JSON(w, webhook)
}
