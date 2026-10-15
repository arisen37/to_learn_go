# External repository download area

No upstream repository was cloned during setup: shell Git reported
`Could not resolve host: github.com`. Research links and question guides are in the parent
catalogue. This directory's go.mod isolates downloaded third-party Go code from the local
curriculum's build/test discovery. It is not itself an interview assignment.

From the workspace root run `make clone-interviews` with working network access. The
script keeps original README questions and licenses intact, records exact commits and
measured source counts in repositories.json, skips existing matching clones without
updating them, and reports failures. It never installs dependencies or executes upstream code.
A guide belongs outside the upstream directory so it cannot be mistaken for original instructions.
