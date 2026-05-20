package run

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kest-labs/kest/api/internal/contracts"
	"github.com/kest-labs/kest/api/internal/modules/request"
	"github.com/kest-labs/kest/api/internal/modules/workspace"
	"github.com/kest-labs/kest/api/pkg/handler"
	"github.com/kest-labs/kest/api/pkg/response"
)

type Handler struct {
	contracts.BaseModule
	requestService   request.Service
	service          Service
	workspaceService workspace.Service
}

func NewHandler(requestService request.Service, service Service, workspaceService workspace.Service) *Handler {
	return &Handler{
		requestService:   requestService,
		service:          service,
		workspaceService: workspaceService,
	}
}

func (h *Handler) Name() string {
	return "run"
}

func (h *Handler) Run(c *gin.Context) {
	workspaceID, ok := h.authorizeWorkspace(c, workspace.RoleWrite)
	if !ok {
		return
	}

	collectionID, ok := handler.ParseID(c, "cid")
	if !ok {
		return
	}

	requestID, ok := handler.ParseID(c, "rid")
	if !ok {
		return
	}

	var req RunRequest
	if !handler.BindJSON(c, &req) {
		return
	}

	reqModel, err := h.requestService.GetByID(c.Request.Context(), requestID, collectionID, workspaceID)
	if err != nil {
		if err == request.ErrRequestNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == request.ErrInvalidCollection {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error(), err)
		return
	}

	userID, ok := handler.GetUserID(c)
	if !ok {
		return
	}

	resp, err := h.service.RunRequest(c.Request.Context(), workspaceID, reqModel, userID, &req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *Handler) Create(c *gin.Context) {
	workspaceID, ok := h.authorizeWorkspace(c, workspace.RoleWrite)
	if !ok {
		return
	}

	userID, ok := handler.GetUserID(c)
	if !ok {
		return
	}

	var req CreateRunRequest
	if !handler.BindJSON(c, &req) {
		return
	}

	run, err := h.service.RecordRun(c.Request.Context(), workspaceID, userID, &req)
	if err != nil {
		if errors.Is(err, ErrInvalidRunInput) {
			response.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		response.InternalServerError(c, err.Error(), err)
		return
	}

	response.Created(c, toRunResponse(run))
}

func (h *Handler) List(c *gin.Context) {
	workspaceID, ok := h.authorizeWorkspace(c, workspace.RoleRead)
	if !ok {
		return
	}

	sourceType := c.Query("source_type")
	sourceID := c.Query("source_id")
	page := handler.QueryInt(c, "page", 1)
	perPage := handler.QueryInt(c, "per_page", 20)

	runs, total, err := h.service.ListRuns(c.Request.Context(), workspaceID, sourceType, sourceID, page, perPage)
	if err != nil {
		response.InternalServerError(c, err.Error(), err)
		return
	}

	response.Success(c, gin.H{
		"items": toRunResponseList(runs),
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"per_page": perPage,
			"pages":    (total + int64(perPage) - 1) / int64(perPage),
		},
	})
}

func (h *Handler) Get(c *gin.Context) {
	workspaceID, ok := h.authorizeWorkspace(c, workspace.RoleRead)
	if !ok {
		return
	}

	runID, ok := handler.ParseID(c, "runid")
	if !ok {
		return
	}

	run, err := h.service.GetRunByID(c.Request.Context(), runID)
	if err != nil {
		if errors.Is(err, ErrRunNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error(), err)
		return
	}
	if run.WorkspaceID != workspaceID {
		response.NotFound(c, ErrRunNotFound.Error())
		return
	}

	response.Success(c, toRunResponse(run))
}

func (h *Handler) authorizeWorkspace(c *gin.Context, requiredRole string) (string, bool) {
	workspaceID, ok := handler.ParseID(c, "id")
	if !ok {
		return "", false
	}

	userID, ok := handler.GetUserID(c)
	if !ok {
		return "", false
	}

	allowed, err := h.workspaceService.HasPermission(workspaceID, userID, requiredRole, false)
	if err != nil || !allowed {
		response.Error(c, http.StatusForbidden, "workspace not found or access denied")
		return "", false
	}

	return workspaceID, true
}
