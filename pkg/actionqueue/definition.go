package actionqueue

import (
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/proto"
)

// DriverConstructorFunc creates a DriverFactory from shared dependencies.
// Queue driver packages register one of these to be called during service startup.
type DriverConstructorFunc func(deps any) any

// Definition describes an action queue backed by a Temporal workflow.
type Definition struct {
	// Type-level fields (set by generated code in init())
	Name               string
	Signal             string
	ActionProto        proto.Message
	WorkflowIDFunc     func(proto.Message) string
	TenantIDFunc       func(proto.Message) string
	EntityIDsFunc      func(proto.Message) []string
	WorkflowIDFromArgs func(tenantID string, entityIDs []string) string

	// Driver constructor (set by queue driver packages via SetDriverConstructor)
	DriverConstructor DriverConstructorFunc

	// Runtime fields (set when driver is bound)
	ActivityMain    any
	WorkflowFunc    any
	ActivityOptions func() workflow.ActivityOptions
}
