package workflow

import "time"

// Workflow represents a BPMN workflow definition
type Workflow struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	Version       string         `json:"version"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	StartEvents   []StartEvent   `json:"startEvents"`
	EndEvents     []EndEvent     `json:"endEvents"`
	Tasks         []Task         `json:"tasks"`
	SequenceFlows []SequenceFlow `json:"sequenceFlows"`
}

// StartEvent represents a start event in the workflow
type StartEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// EndEvent represents an end event in the workflow
type EndEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Task represents a task in the workflow
type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// SequenceFlow represents a sequence flow connecting elements
type SequenceFlow struct {
	ID        string `json:"id"`
	SourceRef string `json:"sourceRef"`
	TargetRef string `json:"targetRef"`
}

// WorkflowDeployment represents a workflow deployment record
type WorkflowDeployment struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Version    string    `json:"version"`
	DeployedAt time.Time `json:"deployedAt"`
	Status     string    `json:"status"`
}

// ValidationResult represents the result of workflow validation
type ValidationResult struct {
	Valid    bool                `json:"valid"`
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Element  string `json:"element"`
	Property string `json:"property"`
	Message  string `json:"message"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Element string `json:"element"`
	Message string `json:"message"`
}
