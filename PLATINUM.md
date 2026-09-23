# Platinum Logic maintenance

This is the public, source-compatible Platinum Logic fork of
[ogen-go/ogen](https://github.com/ogen-go/ogen), based on upstream v1.24.0 at
`0d865e7e568f1b36e5e6788e39aa5cd14e02999f`. Upstream's Apache-2.0 license and
copyright notices remain intact. The Go module path stays `github.com/ogen-go/ogen`.

## Patch scope

- Use one complete, pairwise-disjoint JSON field-type discriminator when
  overlapping enum/constant fields otherwise prevent union inference. Null and
  omitted-field ambiguity still fail closed; generation order is deterministic.
- Enforce JSON constant values during decoding, including optional constants
  when present, using the existing semantic JSON comparison. Ordinary enum
  validation remains the separate generated validation boundary.
- Resolve external encoding interfaces through their underlying type, including
  Go 1.27 interface aliases, and return an error instead of a nil dereference.

The integration fixtures are generic, public reproductions. Private application
schemas and generated clients are deliberately absent.

## Verification and releases

The `Platinum ogen` workflow runs fixture-generation drift, full-module race
tests and Go vet independently on Go 1.27.1. It accepts only trusted maintenance
branch/tag pushes and explicit dispatches. This public dependency fork uses
isolated GitHub-hosted runners, without access to internal organization runners.
No application credentials or deployment steps belong here. Keep inherited
upstream workflows disabled.

Local equivalents, from this module with `GOWORK=off`:

```sh
go generate ./...
git diff --exit-code
go test -race -count=1 -timeout=15m ./...
go vet ./...
```

Upstream's root fixture suite includes skipped `Full` placeholders. Those are
not executed integration tests; the separate `internal/integration` package
executes the generated fixtures. The race-enabled full example corpus needs a
larger timeout than upstream's five-minute shell wrapper on current toolchains.

Keep changes on `platinum-v1.24`. Publish immutable annotated tags with the
`v1.24.0-platinum.N` convention only after verification. Consumers keep the
upstream requirement and use a remote Go module replacement pinned to that exact
fork tag (or immutable pseudo-version), never a moving branch, local path, or
edited module cache. Verify downstream regeneration and client validation before
accepting a consumer update.

For an upstream update, compare the generic patches with the new release, retain
the regression fixtures, run every fork check and the affected downstream gates,
then publish a new immutable tag. Remove a patch only when upstream provides the
same tested behavior. Upstream issue/PR submission is a separate reviewed action.

Repository settings and initial/default maintenance branch are owned by the
`platinum-logic-infra-pulumi` project's `dev` stack. Pulumi must never reset a
maintenance branch tip or republish a release tag.
