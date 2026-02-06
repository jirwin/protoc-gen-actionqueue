package actionqueue

import (
	"go.temporal.io/sdk/workflow"
	"google.golang.org/protobuf/proto"
)

// Definition describes an action queue backed by a Temporal workflow.
type Definition struct {
	Name            string
	Signal          string
	ActionProto     proto.Message
	WorkflowIDFunc  func(proto.Message) string
	TenantIDFunc    func(proto.Message) string
	ActivityMain    any
	WorkflowFunc    any
	ActivityOptions func() workflow.ActivityOptions
}
