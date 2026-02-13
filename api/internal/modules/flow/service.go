package flow

import (
	"context"
	"errors"
	"fmt"
)

// Service defines the business logic interface for flows
type Service interface {
	// Flow CRUD
	CreateFlow(ctx context.Context, projectID, userID uint, req *CreateFlowRequest) (*FlowResponse, error)
	GetFlow(ctx context.Context, id uint) (*FlowDetailResponse, error)
	ListFlows(ctx context.Context, projectID uint) ([]*FlowResponse, error)
	UpdateFlow(ctx context.Context, id uint, req *UpdateFlowRequest) (*FlowResponse, error)
	DeleteFlow(ctx context.Context, id uint) error
	SaveFlow(ctx context.Context, id uint, req *SaveFlowRequest) (*FlowDetailResponse, error)

	// Step CRUD
	CreateStep(ctx context.Context, flowID uint, req *CreateStepRequest) (*StepResponse, error)
	UpdateStep(ctx context.Context, id uint, req *UpdateStepRequest) (*StepResponse, error)
	DeleteStep(ctx context.Context, id uint) error

	// Edge CRUD
	CreateEdge(ctx context.Context, flowID uint, req *CreateEdgeRequest) (*EdgeResponse, error)
	UpdateEdge(ctx context.Context, id uint, req *UpdateEdgeRequest) (*EdgeResponse, error)
	DeleteEdge(ctx context.Context, id uint) error

	// Run
	RunFlow(ctx context.Context, flowID, userID uint) (*RunResponse, error)
	ExecuteFlow(ctx context.Context, runID uint, baseURL string, events chan<- StepEvent) error
	GetRun(ctx context.Context, runID uint) (*RunResponse, error)
	ListRuns(ctx context.Context, flowID uint) ([]*RunResponse, error)
}

type service struct {
	repo Repository
}

// NewService creates a new flow service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// --- Flow CRUD ---

func (s *service) CreateFlow(ctx context.Context, projectID, userID uint, req *CreateFlowRequest) (*FlowResponse, error) {
	flow := &FlowPO{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	}
	if err := s.repo.CreateFlow(ctx, flow); err != nil {
		return nil, err
	}
	return ToFlowResponse(flow), nil
}

func (s *service) GetFlow(ctx context.Context, id uint) (*FlowDetailResponse, error) {
	flow, err := s.repo.GetFlowByID(ctx, id)
	if err != nil {
		return nil, errors.New("flow not found")
	}

	steps, err := s.repo.ListStepsByFlow(ctx, id)
	if err != nil {
		return nil, err
	}

	edges, err := s.repo.ListEdgesByFlow(ctx, id)
	if err != nil {
		return nil, err
	}

	stepResponses := make([]StepResponse, 0, len(steps))
	for _, step := range steps {
		stepResponses = append(stepResponses, *ToStepResponse(step))
	}

	edgeResponses := make([]EdgeResponse, 0, len(edges))
	for _, edge := range edges {
		edgeResponses = append(edgeResponses, *ToEdgeResponse(edge))
	}

	resp := &FlowDetailResponse{
		FlowResponse: *ToFlowResponse(flow),
		Steps:        stepResponses,
		Edges:        edgeResponses,
	}
	resp.StepCount = len(stepResponses)
	return resp, nil
}

func (s *service) ListFlows(ctx context.Context, projectID uint) ([]*FlowResponse, error) {
	flows, err := s.repo.ListFlowsByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	responses := make([]*FlowResponse, 0, len(flows))
	for _, flow := range flows {
		resp := ToFlowResponse(flow)
		// Get step count
		steps, _ := s.repo.ListStepsByFlow(ctx, flow.ID)
		resp.StepCount = len(steps)
		responses = append(responses, resp)
	}
	return responses, nil
}

func (s *service) UpdateFlow(ctx context.Context, id uint, req *UpdateFlowRequest) (*FlowResponse, error) {
	flow, err := s.repo.GetFlowByID(ctx, id)
	if err != nil {
		return nil, errors.New("flow not found")
	}

	if req.Name != nil {
		flow.Name = *req.Name
	}
	if req.Description != nil {
		flow.Description = *req.Description
	}

	if err := s.repo.UpdateFlow(ctx, flow); err != nil {
		return nil, err
	}
	return ToFlowResponse(flow), nil
}

func (s *service) DeleteFlow(ctx context.Context, id uint) error {
	_, err := s.repo.GetFlowByID(ctx, id)
	if err != nil {
		return errors.New("flow not found")
	}

	// Delete edges, steps, then flow
	if err := s.repo.DeleteEdgesByFlow(ctx, id); err != nil {
		return err
	}
	if err := s.repo.DeleteStepsByFlow(ctx, id); err != nil {
		return err
	}
	return s.repo.DeleteFlow(ctx, id)
}

func (s *service) SaveFlow(ctx context.Context, id uint, req *SaveFlowRequest) (*FlowDetailResponse, error) {
	flow, err := s.repo.GetFlowByID(ctx, id)
	if err != nil {
		return nil, errors.New("flow not found")
	}

	// Update flow metadata if provided
	if req.Name != nil {
		flow.Name = *req.Name
	}
	if req.Description != nil {
		flow.Description = *req.Description
	}
	if err := s.repo.UpdateFlow(ctx, flow); err != nil {
		return nil, err
	}

	// Replace all steps and edges
	if err := s.repo.DeleteEdgesByFlow(ctx, id); err != nil {
		return nil, err
	}
	if err := s.repo.DeleteStepsByFlow(ctx, id); err != nil {
		return nil, err
	}

	// Create new steps
	stepPOs := make([]*FlowStepPO, 0, len(req.Steps))
	for _, stepReq := range req.Steps {
		stepPOs = append(stepPOs, &FlowStepPO{
			FlowID:    id,
			Name:      stepReq.Name,
			SortOrder: stepReq.SortOrder,
			Method:    stepReq.Method,
			URL:       stepReq.URL,
			Headers:   stepReq.Headers,
			Body:      stepReq.Body,
			Captures:  stepReq.Captures,
			Asserts:   stepReq.Asserts,
			PositionX: stepReq.PositionX,
			PositionY: stepReq.PositionY,
		})
	}
	if err := s.repo.BatchCreateSteps(ctx, stepPOs); err != nil {
		return nil, err
	}

	// Build a mapping from sort_order to step ID for edge resolution
	stepIDMap := make(map[int]uint)
	for _, step := range stepPOs {
		stepIDMap[step.SortOrder] = step.ID
	}

	// Create new edges
	edgePOs := make([]*FlowEdgePO, 0, len(req.Edges))
	for _, edgeReq := range req.Edges {
		edgePOs = append(edgePOs, &FlowEdgePO{
			FlowID:          id,
			SourceStepID:    edgeReq.SourceStepID,
			TargetStepID:    edgeReq.TargetStepID,
			VariableMapping: edgeReq.VariableMapping,
		})
	}
	if err := s.repo.BatchCreateEdges(ctx, edgePOs); err != nil {
		return nil, err
	}

	return s.GetFlow(ctx, id)
}

// --- Step CRUD ---

func (s *service) CreateStep(ctx context.Context, flowID uint, req *CreateStepRequest) (*StepResponse, error) {
	_, err := s.repo.GetFlowByID(ctx, flowID)
	if err != nil {
		return nil, errors.New("flow not found")
	}

	step := &FlowStepPO{
		FlowID:    flowID,
		Name:      req.Name,
		SortOrder: req.SortOrder,
		Method:    req.Method,
		URL:       req.URL,
		Headers:   req.Headers,
		Body:      req.Body,
		Captures:  req.Captures,
		Asserts:   req.Asserts,
		PositionX: req.PositionX,
		PositionY: req.PositionY,
	}
	if err := s.repo.CreateStep(ctx, step); err != nil {
		return nil, err
	}
	return ToStepResponse(step), nil
}

func (s *service) UpdateStep(ctx context.Context, id uint, req *UpdateStepRequest) (*StepResponse, error) {
	step, err := s.repo.GetStepByID(ctx, id)
	if err != nil {
		return nil, errors.New("step not found")
	}

	if req.Name != nil {
		step.Name = *req.Name
	}
	if req.SortOrder != nil {
		step.SortOrder = *req.SortOrder
	}
	if req.Method != nil {
		step.Method = *req.Method
	}
	if req.URL != nil {
		step.URL = *req.URL
	}
	if req.Headers != nil {
		step.Headers = *req.Headers
	}
	if req.Body != nil {
		step.Body = *req.Body
	}
	if req.Captures != nil {
		step.Captures = *req.Captures
	}
	if req.Asserts != nil {
		step.Asserts = *req.Asserts
	}
	if req.PositionX != nil {
		step.PositionX = *req.PositionX
	}
	if req.PositionY != nil {
		step.PositionY = *req.PositionY
	}

	if err := s.repo.UpdateStep(ctx, step); err != nil {
		return nil, err
	}
	return ToStepResponse(step), nil
}

func (s *service) DeleteStep(ctx context.Context, id uint) error {
	_, err := s.repo.GetStepByID(ctx, id)
	if err != nil {
		return errors.New("step not found")
	}
	return s.repo.DeleteStep(ctx, id)
}

// --- Edge CRUD ---

func (s *service) CreateEdge(ctx context.Context, flowID uint, req *CreateEdgeRequest) (*EdgeResponse, error) {
	_, err := s.repo.GetFlowByID(ctx, flowID)
	if err != nil {
		return nil, errors.New("flow not found")
	}

	edge := &FlowEdgePO{
		FlowID:          flowID,
		SourceStepID:    req.SourceStepID,
		TargetStepID:    req.TargetStepID,
		VariableMapping: req.VariableMapping,
	}
	if err := s.repo.CreateEdge(ctx, edge); err != nil {
		return nil, err
	}
	return ToEdgeResponse(edge), nil
}

func (s *service) UpdateEdge(ctx context.Context, id uint, req *UpdateEdgeRequest) (*EdgeResponse, error) {
	edge, err := s.repo.GetEdgeByID(ctx, id)
	if err != nil {
		return nil, errors.New("edge not found")
	}

	if req.SourceStepID != nil {
		edge.SourceStepID = *req.SourceStepID
	}
	if req.TargetStepID != nil {
		edge.TargetStepID = *req.TargetStepID
	}
	if req.VariableMapping != nil {
		edge.VariableMapping = *req.VariableMapping
	}

	if err := s.repo.UpdateEdge(ctx, edge); err != nil {
		return nil, err
	}
	return ToEdgeResponse(edge), nil
}

func (s *service) DeleteEdge(ctx context.Context, id uint) error {
	_, err := s.repo.GetEdgeByID(ctx, id)
	if err != nil {
		return errors.New("edge not found")
	}
	return s.repo.DeleteEdge(ctx, id)
}

// --- Run ---

func (s *service) ExecuteFlow(ctx context.Context, runID uint, baseURL string, events chan<- StepEvent) error {
	run, err := s.repo.GetRunByID(ctx, runID)
	if err != nil {
		close(events)
		return errors.New("run not found")
	}

	// Get flow steps and edges
	steps, err := s.repo.ListStepsByFlow(ctx, run.FlowID)
	if err != nil {
		close(events)
		return err
	}

	edges, err := s.repo.ListEdgesByFlow(ctx, run.FlowID)
	if err != nil {
		close(events)
		return err
	}

	// Convert to value slices
	stepValues := make([]FlowStepPO, 0, len(steps))
	for _, step := range steps {
		stepValues = append(stepValues, *step)
	}
	edgeValues := make([]FlowEdgePO, 0, len(edges))
	for _, edge := range edges {
		edgeValues = append(edgeValues, *edge)
	}

	runner := NewRunner(s.repo, baseURL)
	return runner.Execute(ctx, run, stepValues, edgeValues, events)
}

func (s *service) RunFlow(ctx context.Context, flowID, userID uint) (*RunResponse, error) {
	flowDetail, err := s.GetFlow(ctx, flowID)
	if err != nil {
		return nil, err
	}

	if len(flowDetail.Steps) == 0 {
		return nil, errors.New("flow has no steps")
	}

	// Create run record
	run := &FlowRunPO{
		FlowID:      flowID,
		Status:      RunStatusPending,
		TriggeredBy: userID,
	}
	if err := s.repo.CreateRun(ctx, run); err != nil {
		return nil, err
	}

	// Create step result placeholders
	for _, step := range flowDetail.Steps {
		result := &FlowStepResultPO{
			RunID:  run.ID,
			StepID: step.ID,
			Status: RunStatusPending,
		}
		if err := s.repo.CreateStepResult(ctx, result); err != nil {
			return nil, fmt.Errorf("failed to create step result: %w", err)
		}
	}

	return ToRunResponse(run), nil
}

func (s *service) GetRun(ctx context.Context, runID uint) (*RunResponse, error) {
	run, err := s.repo.GetRunByID(ctx, runID)
	if err != nil {
		return nil, errors.New("run not found")
	}

	resp := ToRunResponse(run)

	results, err := s.repo.ListStepResultsByRun(ctx, runID)
	if err != nil {
		return nil, err
	}

	stepResults := make([]StepResultResponse, 0, len(results))
	for _, result := range results {
		stepResults = append(stepResults, *ToStepResultResponse(result))
	}
	resp.StepResults = stepResults

	return resp, nil
}

func (s *service) ListRuns(ctx context.Context, flowID uint) ([]*RunResponse, error) {
	runs, err := s.repo.ListRunsByFlow(ctx, flowID)
	if err != nil {
		return nil, err
	}

	responses := make([]*RunResponse, 0, len(runs))
	for _, run := range runs {
		responses = append(responses, ToRunResponse(run))
	}
	return responses, nil
}
