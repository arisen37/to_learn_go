# Repository repair and company interview practice

Start with [the included Go repair lab](repair-shop/README.md): a **514-line starter**
across four implementation files including its command (after gofmt; tests excluded).
It is an original teaching exercise with eight intended defects, acceptance tests and
an [explained reference solution](../solutions/interviews/repair-shop/README.md).
It is **not** a company question. Source line counts include comments and blank lines;
`scripts/audit-curriculum.py` measures the actual implementation, not repository size.

## Public company questions found on GitHub

Researched on **2026-10-15** from the repositories' own README pages. These are published
practice sources, not leaked/private questions or proof that a company still uses them.
Only the Dr.Fit source closely matches the Go repository-repair format. Others are useful
existing-code extension exercises and are labelled accordingly. No external source is
claimed to be exactly 500–600 lines; sizes remain unknown until downloaded.

| Source | Language / task type | Local question and architecture guide |
| --- | --- | --- |
| [Dr.Fit author-published hiring task](https://github.com/YordanBahov-bg/dr-fit-hiring-test) | Go + React Native; reported bugs and small features; personal author account, employer usage not independently verified | [Guide](company-guides/dr-fit.md) |
| [Noyo company repository](https://github.com/noyo-technologies/backend-coding-challenge) | Python; existing application extensions and incomplete tests; explicitly not a planted-bug hunt | [Guide](company-guides/noyo.md) |
| [Barstool Sports company repository](https://github.com/BarstoolSports/backend-challenge) | Node.js; modify API behavior, add notes and authorization/pagination tests | [Guide](company-guides/barstool.md) |
| [DejaBlue company repository](https://github.com/deja-blue/software-interview) | Go + TypeScript; pricing, optimization, GraphQL streaming and UI | [Guide](company-guides/deja-blue.md) |

## Clone status and original questions

Web research succeeded; Git clones failed because the shell could not resolve github.com.
Therefore `upstream/` does **not** contain these external source repositories yet, and no
upstream tests were run. [repositories.json](repositories.json) records that state without
invented commits, sizes or success claims. To fetch all four with Git/network available:

```sh
make clone-interviews
# or, from the root:
python3 scripts/clone-interviews.py
```

The script clones to upstream/<name>, preserves each original README/question and license,
records the exact commit and production-source line count, and reports every failure.
Clones are local study copies; the script does not fork, submit a PR, contact a company,
or install/run dependencies. Original README files remain the authoritative questions.
Our guides summarize scope and architecture with source links; they are not official answer keys.

## Practice order

1. Original Go repair lab: reproduce a symptom, write a regression, fix, explain.
2. Dr.Fit: backend issues first, then optional frontend integration.
3. Barstool or Noyo: transfer the same tracing/testing method to another language.
4. DejaBlue: advanced business rules, algorithms and concurrency in an existing stack.

Use [the timed workflow](WORKFLOW.md). Keep solutions private when the upstream requests it;
this curriculum does not publish or submit them. For external work, write your own tests
and reasoning under a separate local working copy. Their dependencies and environments
are not part of the dependency-free core curriculum.
