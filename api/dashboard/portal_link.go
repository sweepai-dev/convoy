package dashboard

import (
	"fmt"
	"net/http"

	"github.com/frain-dev/convoy/pkg/log"

	"github.com/frain-dev/convoy/api/models"
	"github.com/frain-dev/convoy/database/postgres"
	"github.com/frain-dev/convoy/datastore"
	"github.com/frain-dev/convoy/services"
	"github.com/frain-dev/convoy/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	m "github.com/frain-dev/convoy/internal/pkg/middleware"
)

func (a *DashboardHandler) CreatePortalLink(w http.ResponseWriter, r *http.Request) {
	var newPortalLink models.PortalLink
	if err := util.ReadJSON(r, &newPortalLink); err != nil {
		_ = render.Render(w, r, util.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	project, err := a.retrieveProject(r)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	if err = a.A.Authz.Authorize(r.Context(), "project.manage", project); err != nil {
		_ = render.Render(w, r, util.NewErrorResponse("Unauthorized", http.StatusForbidden))
		return
	}

	cp := services.CreatePortalLinkService{
		PortalLinkRepo: postgres.NewPortalLinkRepo(a.A.DB),
		EndpointRepo:   postgres.NewEndpointRepo(a.A.DB),
		Portal:         &newPortalLink,
		Project:        project,
	}

	portalLink, err := cp.Run(r.Context())
	if err != nil {
		_ = render.Render(w, r, util.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	baseUrl, err := a.retrieveHost()
	if err != nil {
		_ = render.Render(w, r, util.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	pl := portalLinkResponse(portalLink, baseUrl)
	_ = render.Render(w, r, util.NewServerResponse("Portal link created successfully", pl, http.StatusCreated))
}

func (a *DashboardHandler) GetPortalLinkByID(w http.ResponseWriter, r *http.Request) {
	project, err := a.retrieveProject(r)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	portalLink, err := postgres.NewPortalLinkRepo(a.A.DB).FindPortalLinkByID(r.Context(), project.UID, chi.URLParam(r, "portalLinkID"))
	if err != nil {
		if err == datastore.ErrPortalLinkNotFound {
			_ = render.Render(w, r, util.NewServerResponse(err.Error(), nil, http.StatusNotFound))
			return
		}

		_ = render.Render(w, r, util.NewServerResponse("error retrieving portal link", nil, http.StatusBadRequest))
		return
	}

	baseUrl, err := a.retrieveHost()
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}
	pl := portalLinkResponse(portalLink, baseUrl)

	_ = render.Render(w, r, util.NewServerResponse("Portal link fetched successfully", pl, http.StatusOK))
}

func (a *DashboardHandler) UpdatePortalLink(w http.ResponseWriter, r *http.Request) {
	var updatePortalLink models.PortalLink
	err := util.ReadJSON(r, &updatePortalLink)
	if err != nil {
		_ = render.Render(w, r, util.NewErrorResponse(err.Error(), http.StatusBadRequest))
		return
	}

	project, err := a.retrieveProject(r)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	if err = a.A.Authz.Authorize(r.Context(), "project.manage", project); err != nil {
		_ = render.Render(w, r, util.NewErrorResponse("Unauthorized", http.StatusForbidden))
		return
	}

	portalLink, err := postgres.NewPortalLinkRepo(a.A.DB).FindPortalLinkByID(r.Context(), project.UID, chi.URLParam(r, "portalLinkID"))
	if err != nil {
		if err == datastore.ErrPortalLinkNotFound {
			_ = render.Render(w, r, util.NewServerResponse(err.Error(), nil, http.StatusNotFound))
			return
		}

		_ = render.Render(w, r, util.NewServerResponse("error retrieving portal link", nil, http.StatusBadRequest))
		return
	}

	upl := services.UpdatePortalLinkService{
		PortalLinkRepo: postgres.NewPortalLinkRepo(a.A.DB),
		EndpointRepo:   postgres.NewEndpointRepo(a.A.DB),
		Project:        project,
		Update:         &updatePortalLink,
		PortalLink:     portalLink,
	}

	portalLink, err = upl.Run(r.Context())
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	baseUrl, err := a.retrieveHost()
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}
	pl := portalLinkResponse(portalLink, baseUrl)

	_ = render.Render(w, r, util.NewServerResponse("Portal link updated successfully", pl, http.StatusAccepted))
}

func (a *DashboardHandler) RevokePortalLink(w http.ResponseWriter, r *http.Request) {
	project, err := a.retrieveProject(r)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	if err = a.A.Authz.Authorize(r.Context(), "project.manage", project); err != nil {
		_ = render.Render(w, r, util.NewErrorResponse("Unauthorized", http.StatusForbidden))
		return
	}

	portalLink, err := postgres.NewPortalLinkRepo(a.A.DB).FindPortalLinkByID(r.Context(), project.UID, chi.URLParam(r, "portalLinkID"))
	if err != nil {
		if err == datastore.ErrPortalLinkNotFound {
			_ = render.Render(w, r, util.NewServerResponse(err.Error(), nil, http.StatusNotFound))
			return
		}

		_ = render.Render(w, r, util.NewServerResponse("error retrieving portal link", nil, http.StatusBadRequest))
		return
	}

	err = postgres.NewPortalLinkRepo(a.A.DB).RevokePortalLink(r.Context(), project.UID, portalLink.UID)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	_ = render.Render(w, r, util.NewServerResponse("Portal link revoked successfully", nil, http.StatusOK))
}

func (a *DashboardHandler) LoadPortalLinksPaged(w http.ResponseWriter, r *http.Request) {
	pageable := m.GetPageableFromContext(r.Context())
	project, err := a.retrieveProject(r)
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	var endpointID string
	endpointIDs := getEndpointIDs(r)

	if len(endpointIDs) > 0 {
		endpointID = endpointIDs[0]
	}

	filter := &datastore.FilterBy{EndpointID: endpointID}

	portalLinks, paginationData, err := postgres.NewPortalLinkRepo(a.A.DB).LoadPortalLinksPaged(r.Context(), project.UID, filter, pageable)
	if err != nil {
		log.WithError(err).Println("an error occurred while fetching portal links")
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	plResponse := []*models.PortalLinkResponse{}
	baseUrl, err := a.retrieveHost()
	if err != nil {
		_ = render.Render(w, r, util.NewServiceErrResponse(err))
		return
	}

	for _, portalLink := range portalLinks {
		pl := portalLinkResponse(&portalLink, baseUrl)
		plResponse = append(plResponse, pl)
	}

	_ = render.Render(w, r, util.NewServerResponse("Portal links fetched successfully", pagedResponse{Content: plResponse, Pagination: &paginationData}, http.StatusOK))
}

func portalLinkResponse(pl *datastore.PortalLink, baseUrl string) *models.PortalLinkResponse {
	return &models.PortalLinkResponse{
		UID:               pl.UID,
		ProjectID:         pl.ProjectID,
		Name:              pl.Name,
		URL:               fmt.Sprintf("%s/portal?token=%s", baseUrl, pl.Token),
		Token:             pl.Token,
		OwnerID:           pl.OwnerID,
		Endpoints:         pl.Endpoints,
		EndpointCount:     len(pl.EndpointsMetadata),
		CanManageEndpoint: pl.CanManageEndpoint,
		EndpointsMetadata: pl.EndpointsMetadata,
		CreatedAt:         pl.CreatedAt,
		UpdatedAt:         pl.UpdatedAt,
	}
}
