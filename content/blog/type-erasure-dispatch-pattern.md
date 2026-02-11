---
title: "Type-Erasure Dispatch Pattern in Go"
slug: type-erasure-dispatch-pattern-go
date: 2026-02-11
description: "Use type-erasure dispatch to run one orchestration workflow across many concrete types in Go without hard-coded type switches."
tags:
  - golang
  - architecture
  - design-pattern
  - polymorphism
---

Coordinators tend to rot when they need to manage many concrete types. You start with a small workflow, then add one
type switch, then another, and eventually every new backend or model change touches orchestration code.

The **type-erasure dispatch pattern** is how I got out of this in a multi-database sync system.

## The Scaling Failure That Forced the Refactor

Originally I had 6 different database variations to support behind one sync workflow. The design tied capabilities directly to
typed interfaces and backend implementations.

Every time I added a new structure:

- I had to update the central interface
- then update all 6 databases to match

At that point I still had 3 more databases planned and about 5 more capabilities to add. The growth path was not
maintainable. The surface area was exploding as a cross-product of databases and capabilities, and orchestration kept
getting harder to reason about.

## Core Roles

1. `Capability ID`  
A stable string/enum key (for example: `page`, `template`, `product`).

2. `Item Contract`  
A small type-erased interface for all syncable records (for example: `Key()`, `Checksum()`, `ModifiedAt()`).

3. `Typed Repository`  
Strongly-typed data access for one item type (`Repo[T]`).

4. `Capability Repository`  
A runtime, type-erased repository contract used by generic workflows.

5. `Adapter Backend`  
A concrete backend that advertises supported capabilities and implements typed operations.

6. `Capability Wrapper/Bridge`  
Adapts typed operations into the type-erased capability repository.

7. `Adapter Facade`  
Holds `map[capability]CapabilityRepo` and resolves repos by capability.

8. `Orchestrator`  
Runs generic operations (`push`, `pull`, `sync`) by:
- negotiating shared capabilities
- dispatching each capability to resolved repos
- applying conflict/dirtiness rules consistently

## Reference Shape (Pseudocode)

```go
type RepoItem interface {
    Key() string
    Checksum() string
    StoredChecksum() string
    ModifiedAt() time.Time
}

type Repo[T RepoItem] interface {
    List(ctx context.Context) ([]T, error)
    Read(ctx context.Context, key string) (*T, error)
    Upsert(ctx context.Context, item *T) error
}

type CapabilityRepo interface {
    Capability() string
    List(ctx context.Context) ([]RepoItem, error)
    Read(ctx context.Context, key string) (RepoItem, error)
    Upsert(ctx context.Context, item RepoItem) error
    DirtyKeys(ctx context.Context) (map[string]struct{}, error)
    MarkClean(ctx context.Context, items []RepoItem) error
}

type Backend interface {
    Capabilities() []string
    GetCapability(capability string) CapabilityRepo
    // typed operations...
}
```

## Runtime Flow

1. A backend declares supported capabilities.
2. An adapter facade builds a capability-to-repo map.
3. The orchestrator intersects capability sets between source and target.
4. For each shared capability:
- resolve local and remote capability repos
- execute a generic operation (`list/read/upsert`)
- apply common policy (dirty-only push, conflict resolution, checksum state updates)
5. Aggregate per-item/per-capability errors into one operation result.

## Why This Works

- Generic logic stays generic.
- Backend differences are isolated behind capability resolution.
- Feature availability is explicit and discoverable at runtime.
- Adding a new capability is incremental instead of invasive.

## Extension Workflow

1. Add a new capability identifier.
2. Define the item type to satisfy the base item contract.
3. Implement a typed repo for each backend that should support it.
4. Register/wrap each typed repo into the capability map.
5. Optionally add orchestration order, labels, and mode-specific policies.
6. Add tests for:
- unsupported capability behavior
- dispatch correctness
- sync behavior and conflict policy

## Error Model

Recommended error categories:

- `ErrNoCapability`: requested capability is not supported
- `ErrReadonly`: write attempted on read-only target
- `ErrUnsupportedMode`: operation mode not available for this capability

The key win is that I no longer edit orchestration when adding every new concrete type. I add capabilities and backend
support incrementally, and the dispatcher handles the rest.
