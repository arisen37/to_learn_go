#!/usr/bin/env python3
"""Check local curriculum integrity without running learner tests or upstream code."""
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[1]
errors = []
assignments = sorted((ROOT / 'assignments').glob('*/*/exercise.go'))
counts = {difficulty: 0 for difficulty in ('easy', 'medium', 'hard')}
for source in assignments:
    relative = source.parent.relative_to(ROOT / 'assignments')
    counts[relative.parts[0]] += 1
    solution = ROOT / 'solutions' / relative
    for required in [source.parent/'README.md', source.parent/'exercise_test.go',
                     solution/'solution.go', solution/'README.md', solution/'exercise_test.go']:
        if not required.is_file():
            errors.append(f'Missing {required.relative_to(ROOT)}')
    a, b = source.parent/'exercise_test.go', solution/'exercise_test.go'
    if a.exists() and b.exists() and a.read_bytes() != b.read_bytes():
        errors.append(f'Supplied test copies differ: {relative}; put learner regressions in regression_test.go')
if counts != {'easy':7, 'medium':7, 'hard':7}:
    errors.append(f'Expected 7 assignments per difficulty, got {counts}')
left = ROOT/'interviews/repair-shop/repair_test.go'
right = ROOT/'solutions/interviews/repair-shop/repair_test.go'
if left.read_bytes() != right.read_bytes():
    errors.append('Repair acceptance test copies differ')
production = [p for p in (ROOT/'interviews/repair-shop').rglob('*.go') if not p.name.endswith('_test.go')]
lines = sum(len(p.read_text().splitlines()) for p in production)
if not 500 <= lines <= 600:
    errors.append(f'Repair starter must have 500–600 physical Go lines, got {lines}')
# Check local Markdown links; external source links are deliberately not fetched here.
markdown = [ROOT/'README.md', ROOT/'Roadmap.md']
for folder in ['lessons','docs','examples','assignments','solutions','interviews']:
    markdown.extend(p for p in (ROOT/folder).rglob('*.md') if 'upstream' not in p.relative_to(ROOT).parts)
for page in markdown:
    for match in re.finditer(r'(?<!!)\[[^\]]+\]\(([^)]+)\)', page.read_text()):
        target = match.group(1).split('#', 1)[0]
        if not target or '://' in target or target.startswith('mailto:'):
            continue
        target = target.strip('<>')
        if not (page.parent/target).exists():
            errors.append(f'{page.relative_to(ROOT)}: missing link {target}')
print(f'Assignments: {counts}; exact test copies checked: {len(assignments)+1}')
print(f'Repair starter: {len(production)} implementation files, {lines} physical Go lines (tests excluded)')
print(f'Local Markdown links checked in {len(markdown)} documents')
if errors:
    print('\n'.join(errors), file=sys.stderr)
    raise SystemExit(1)
print('Curriculum integrity passed')
