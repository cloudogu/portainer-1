package edgegroups

import (
	"fmt"
	"net/http"

	portainer "github.com/cloudogu/portainer-ce/api"
	"github.com/cloudogu/portainer-ce/api/dataservices"
	"github.com/cloudogu/portainer-ce/api/internal/slices"
	httperror "github.com/portainer/libhttp/error"
	"github.com/portainer/libhttp/response"
)

type decoratedEdgeGroup struct {
	portainer.EdgeGroup
	HasEdgeStack  bool `json:"HasEdgeStack"`
	HasEdgeGroup  bool `json:"HasEdgeGroup"`
	EndpointTypes []portainer.EndpointType
}

// @id EdgeGroupList
// @summary list EdgeGroups
// @description **Access policy**: administrator
// @tags edge_groups
// @security ApiKeyAuth
// @security jwt
// @produce json
// @success 200 {array} decoratedEdgeGroup "EdgeGroups"
// @failure 500
// @failure 503 "Edge compute features are disabled"
// @router /edge_groups [get]
func (handler *Handler) edgeGroupList(w http.ResponseWriter, r *http.Request) *httperror.HandlerError {
	edgeGroups, err := handler.DataStore.EdgeGroup().EdgeGroups()
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve Edge groups from the database", err)
	}

	edgeStacks, err := handler.DataStore.EdgeStack().EdgeStacks()
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve Edge stacks from the database", err)
	}

	usedEdgeGroups := make(map[portainer.EdgeGroupID]bool)

	for _, stack := range edgeStacks {
		for _, groupID := range stack.EdgeGroups {
			usedEdgeGroups[groupID] = true
		}
	}

	edgeJobs, err := handler.DataStore.EdgeJob().EdgeJobs()
	if err != nil {
		return httperror.InternalServerError("Unable to retrieve Edge jobs from the database", err)
	}

	decoratedEdgeGroups := []decoratedEdgeGroup{}
	for _, orgEdgeGroup := range edgeGroups {
		usedByEdgeJob := false
		for _, edgeJob := range edgeJobs {
			if slices.Contains(edgeJob.EdgeGroups, portainer.EdgeGroupID(orgEdgeGroup.ID)) {
				usedByEdgeJob = true
				break
			}
		}

		edgeGroup := decoratedEdgeGroup{
			EdgeGroup:     orgEdgeGroup,
			EndpointTypes: []portainer.EndpointType{},
		}
		if edgeGroup.Dynamic {
			endpointIDs, err := handler.getEndpointsByTags(edgeGroup.TagIDs, edgeGroup.PartialMatch)
			if err != nil {
				return httperror.InternalServerError("Unable to retrieve environments and environment groups for Edge group", err)
			}

			edgeGroup.Endpoints = endpointIDs
		}

		endpointTypes, err := getEndpointTypes(handler.DataStore.Endpoint(), edgeGroup.Endpoints)
		if err != nil {
			return httperror.InternalServerError("Unable to retrieve environment types for Edge group", err)
		}

		edgeGroup.EndpointTypes = endpointTypes

		edgeGroup.HasEdgeStack = usedEdgeGroups[edgeGroup.ID]

		edgeGroup.HasEdgeGroup = usedByEdgeJob

		decoratedEdgeGroups = append(decoratedEdgeGroups, edgeGroup)
	}

	return response.JSON(w, decoratedEdgeGroups)
}

func getEndpointTypes(endpointService dataservices.EndpointService, endpointIds []portainer.EndpointID) ([]portainer.EndpointType, error) {
	typeSet := map[portainer.EndpointType]bool{}
	for _, endpointID := range endpointIds {
		endpoint, err := endpointService.Endpoint(endpointID)
		if err != nil {
			return nil, fmt.Errorf("failed fetching environment: %w", err)
		}

		typeSet[endpoint.Type] = true
	}

	endpointTypes := make([]portainer.EndpointType, 0, len(typeSet))
	for endpointType := range typeSet {
		endpointTypes = append(endpointTypes, endpointType)
	}

	return endpointTypes, nil
}
