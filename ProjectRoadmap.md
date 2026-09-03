# Geareo Save Editor — Project Roadmap

## Status: Foundation complete ✅

- Layered architecture: `core` → `modules/world` (or `model`/`save`) → `app` (`Service`/`Controller`) → `ui`
- `SaveEditor`: loads/parses save files, holds `Path` + objectified data + passthrough
- Full model layer: structs for `Circuit`, `CircuitEntity`, `CircuitEntityData` variants (with correct marshal/unmarshal, value vs pointer receivers sorted out)
- `UI` interface + working CLI implementation (menus, prompts, lists)
- File picking via native dialogs (`utils`/`dialogs`, zenity)
- Backup-on-write system, with restore-from-backup flow
- Import/export for circuits, including **multi-file import with two-pass ID remapping** (so sub-circuit references survive collisions)
- Validation system (`ValidateSavefile` / `validateCircuit`) — duplicate ID detection, broken sub-circuit reference detection

This is the hard part. Everything below builds on top of a working foundation.

---

## Phase 1 — Round out core editing (CRUD)

You have Get / Import / Export. The missing symmetric operations:

- [ ] `DeleteCircuit(id)` — with a validation check first (does anything reference this circuit?)
- [ ] `DuplicateCircuit(id)` — in-save copy, reusing the same ID-remap logic you already built for import
- [ ] `RenameCircuit(id, newName)`
- [ ] `MoveCircuit` / transform position, rotation
- [ ] Per-part editing: add/remove/move individual `CircuitEntity` parts within a circuit
- [ ] Bulk operations: use `GetCircuitPartsBetweenPositions` (already built) to select a region, then delete/move/copy the whole selection at once

**Why now:** cheapest phase — reuses patterns (two-pass remap, `SaveEditor` mechanics / `Service` rules split) you've already proven work.

---

## Phase 2 — Tighten validation & safety

- [ ] Run `ValidateSavefile()` automatically before every `Save()`/`SaveAs()`, warn via `ui.Confirm` if issues exist (not just a manual menu item)
- [ ] Run validation automatically after every import batch — safety net for the remap logic
- [ ] Add a "malformed data" issue type for failed type assertions (flagged last message, worth actually adding)
- [ ] Dry-run / preview mode for destructive ops — show what would be deleted/changed before committing
- [ ] Undo/redo — since `SaveEditor` is a plain struct tree with no external resources, snapshotting is cheap: deep-copy before each mutation, push onto a stack, pop to undo

**Why now:** you're about to add real destructive operations (Phase 1) — this phase is what makes using them safe to experiment with.

---

## Phase 3 — Autobuild / templates

This was your original scoped feature — now's the time, since CRUD + validation give you safe primitives to build on.

- [ ] Step/Wizard runner (sketched earlier — sequential prompts with validation/retry per step)
- [ ] One real build type end-to-end (e.g. a straight wire run, or a small logic gate grid)
- [ ] Template library: save a selected region of parts as a named, reusable template
- [ ] Apply a template at an arbitrary origin + rotation
- [ ] Stretch: parameterized templates ("N-bit adder" takes a width param)

---

## Phase 4 — UX polish

- [ ] Search: find circuits/parts by name substring, not just exact ID/index
- [ ] Recent files list (reuse the platform-specific directory logic from `utils`/`dialogs`)
- [ ] Consistent `ClickAny()`/pause behavior across every handler (some currently skip it — small cleanup pass)
- [ ] Config file for preferences (default save dir override, backup retention count/policy)
- [ ] Logging to file for your own debugging + future bug reports (sketched earlier, not yet wired in — confirm if you did this)

---

## Phase 5 — Distribution

- [ ] Single-binary builds for Windows/Mac/Linux (GoReleaser or plain `go build` cross-compilation)
- [ ] README with setup + usage instructions
- [ ] Stretch: swap CLI for a GUI backend (Fyne/Wails) — your `UI` interface was built specifically so this is a new implementation, not a rewrite

---

## Suggested order

1. **Phase 1** (CRUD) — highest value, cheapest, reuses existing patterns
2. **Phase 2** (safety) — do this *while* building Phase 1, not strictly after — add undo/redo before you're relying on Delete/Duplicate day-to-day
3. **Phase 3** (autobuild) — the "fun" feature, most satisfying once safety net exists
4. **Phase 4 & 5** — whenever the tool is good enough that friction/distribution starts to matter more than new features

## Architecture rules to keep following (so you don't relitigate them)

- `SaveEditor` = data mechanics (get/add/remove, no decisions)
- `Service` = business rules (validation, ID collision policy, composing multiple `SaveEditor` calls)
- `Controller` = UI orchestration only (prompts, menu dispatch, display) — never touches raw data or JSON
- New "get X within Y" functions take the already-resolved `Y` as a pointer, don't re-look-it-up by name
- Range over slices by index (`for i := range`) whenever you need a pointer into the real data, never `for _, v := range` + `&v`
- Match type assertions to how `UnmarshalJSON` actually assigns the field (value vs pointer) — check, don't assume
