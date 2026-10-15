# Timed repository repair workflow

For a 75–90 minute exercise:

1. **0–10 min: map the system.** Read README, go.mod, entry point and tests. Draw request -> handler -> domain -> storage, plus background work and shutdown. Identify the language/tool versions.
2. **10–20 min: establish baseline.** Run the smallest relevant test. Reproduce a reported issue with a minimal input. Record pre-existing failures separately from new changes.
3. **20–60 min: fix one cause at a time.** Add a failing regression, trace ownership and invariants, make a minimal fix, rerun the focused test. Start with crashes, incorrect status/data and data-loss bugs.
4. **60–75 min: broaden verification.** Run the relevant suite, race detector and vet. Recheck error paths, input limits, cancellation and resource cleanup. Do not claim unavailable checks passed.
5. **75–90 min: explain.** Describe symptom, root cause, changed behavior, architecture impact, tests and remaining constraints.

A good explanation fits this template:

> Input X previously caused Y because layer Z violated invariant I. The fix validates or
> updates state at boundary B. Test T reproduces the failure and now passes. Operation costs
> O(...); resource R is owned by ...; shutdown/cancellation behaves as ... . Remaining limitation: ... .

Use a separate regression_test.go for your added cases, keeping supplied acceptance tests
intact. Avoid unrelated dependency upgrades and broad rewrites. If the original task is
open-ended, make assumptions explicit before coding. External repositories may have outdated
setup instructions or require extra services; setup failure is not an application regression.
