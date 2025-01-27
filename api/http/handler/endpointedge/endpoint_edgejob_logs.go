package endpointedge

import (
	"net/http"
	"strconv"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/http/middlewares"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/request"
	"github.com/portainer/libhttp/response"
)

type logsPayload struct {
	FileContent string
}

func (payload *logsPayload) Validate(r *http.Request) error {
	return nil
}

// endpointEdgeJobsLogs
// @summary Inspect an EdgeJob Log
// @description **Access policy**: public
// @tags edge, endpoints
// @accept json
// @produce json
// @param id path string true "environment(endpoint) Id"
// @param jobID path string true "Job Id"
// @success 200
// @failure 500
// @failure 400
// @router /endpoints/{id}/edge/jobs/{jobID}/logs [post]
func (handler *Handler) endpointEdgeJobsLogs(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	endpoint, err := middlewares.FetchEndpoint(r)
	if err != nil {
		return httperror.BadRequest("Unable to find an environment on request context", err)
	}

	err = handler.requestBouncer.AuthorizedEdgeEndpointOperation(r, endpoint)
	if err != nil {
		return httperror.Forbidden("Permission denied to access environment", err)
	}

	edgeJobID, err := request.RetrieveNumericRouteVariableValue(r, "jobID")
	if err != nil {
		return httperror.BadRequest("Invalid edge job identifier route variable", err)
	}

	var payload logsPayload
	err = request.DecodeAndValidateJSONPayload(r, &payload)
	if err != nil {
		return httperror.BadRequest("Invalid request payload", err)
	}

	edgeJob, err := handler.DataStore.EdgeJob().EdgeJob(portainer.EdgeJobID(edgeJobID))
	if handler.DataStore.IsErrObjectNotFound(err) {
		return httperror.NotFound("Unable to find an edge job with the specified identifier inside the database", err)
	} else if err != nil {
		return httperror.InternalServerError("Unable to find an edge job with the specified identifier inside the database", err)
	}

	err = handler.FileService.StoreEdgeJobTaskLogFileFromBytes(strconv.Itoa(edgeJobID), strconv.Itoa(int(endpoint.ID)), []byte(payload.FileContent))
	if err != nil {
		return httperror.InternalServerError("Unable to save task log to the filesystem", err)
	}

	meta := portainer.EdgeJobEndpointMeta{CollectLogs: false, LogsStatus: portainer.EdgeJobLogsStatusCollected}
	if _, ok := edgeJob.GroupLogsCollection[endpoint.ID]; ok {
		edgeJob.GroupLogsCollection[endpoint.ID] = meta
	} else {
		edgeJob.Endpoints[endpoint.ID] = meta
	}

	err = handler.DataStore.EdgeJob().UpdateEdgeJob(edgeJob.ID, edgeJob)

	handler.ReverseTunnelService.AddEdgeJob(endpoint.ID, edgeJob)

	if err != nil {
		return httperror.InternalServerError("Unable to persist edge job changes to the database", err)
	}

	return response.JSON(w, nil)
}
