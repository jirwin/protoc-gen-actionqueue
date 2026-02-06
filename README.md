# protoc-gen-actionqueue

A protoc plugin that generates action queue registration code from proto annotations, reducing boilerplate when creating new action queues.

## Installation

```bash
go install github.com/jirwin/protoc-gen-actionqueue@latest
```

## Usage

### Proto Annotation

Add the `actionqueue.v1.action` option to your action message:

```protobuf
syntax = "proto3";

package c1.models.myentity.v1;

import "actionqueue/v1/options.proto";
import "google/protobuf/timestamp.proto";

message MyEntityAction {
  option (actionqueue.v1.action) = {
    name: "my-entity-queue"
    signal: "signal:start:my-entity-queue"
    workflow_id_fields: ["tenant_id", "entity_id"]
    // tenant_id_field: "tenant_id"  // optional, defaults to "tenant_id"
  };

  string tenant_id = 1;
  string entity_id = 2;
  string id = 3;
  MyEntityEvent event = 4;
  google.protobuf.Timestamp deadline = 5;
}

enum MyEntityEvent {
  MY_ENTITY_EVENT_UNSPECIFIED = 0;
  MY_ENTITY_EVENT_PROCESS = 1;
}
```

### Generated Code

The plugin generates a `*.pb.actionqueue.go` file with:

```go
// Constants
const (
    MyEntityQueueName   = "my-entity-queue"
    MyEntityQueueSignal = "signal:start:my-entity-queue"
)

// WorkflowID function
func MyEntityActionWorkflowID(action proto.Message) string {
    a := action.(*MyEntityAction)
    return fmt.Sprintf("my-entity-queue:%s:%s", a.GetTenantId(), a.GetEntityId())
}

// TenantID function
func MyEntityActionTenantID(action proto.Message) string {
    return action.(*MyEntityAction).GetTenantId()
}

// Definition factory
func MyEntityQueueDefinition(
    activityMain any,
    workflowFunc any,
    activityOptions func() workflow.ActivityOptions,
) *actionqueue.Definition {
    return &actionqueue.Definition{
        Name:            MyEntityQueueName,
        Signal:          MyEntityQueueSignal,
        ActionProto:     &MyEntityAction{},
        WorkflowIDFunc:  MyEntityActionWorkflowID,
        TenantIDFunc:    MyEntityActionTenantID,
        ActivityMain:    activityMain,
        WorkflowFunc:    workflowFunc,
        ActivityOptions: activityOptions,
    }
}
```

### buf.gen.yaml Configuration

```yaml
version: v2
plugins:
  - local: protoc-gen-actionqueue
    out: pkg/pb
    opt:
      - paths=source_relative
```

## Options

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique workflow name (e.g., "my-entity-queue") |
| `signal` | string | Yes | Signal name to wake the workflow |
| `workflow_id_fields` | repeated string | Yes | Field names to compose the workflow ID |
| `tenant_id_field` | string | No | Field name for tenant extraction (default: "tenant_id") |

## Validation

The generator validates:
- `name` is non-empty
- `signal` is non-empty
- All `workflow_id_fields` exist as scalar fields in the message
- `tenant_id_field` exists (or default "tenant_id" exists)

## Development

```bash
make build     # Build the plugin
make generate  # Generate options.pb.go
make test      # Run tests
make clean     # Clean build artifacts
```

## License

Apache 2.0
