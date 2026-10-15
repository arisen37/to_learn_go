#!/usr/bin/env python3
"""Clone researched public questions without executing third-party source code."""
from pathlib import Path
import json
import subprocess
import sys
from datetime import datetime, timezone

ROOT = Path(__file__).resolve().parents[1]
MANIFEST = ROOT / 'interviews/repositories.json'
EXCLUDED = {'node_modules', 'vendor', '.git', 'dist', 'build', 'generated', '__pycache__'}
EXTENSIONS = {'.go', '.py', '.js', '.jsx', '.ts', '.tsx'}


def git(*args, cwd=None):
    return subprocess.run(['git', *args], cwd=cwd, text=True, capture_output=True,
                          check=True, timeout=60).stdout.strip()


def inventory(destination):
    # Physical source lines, excluding recognizable test/generated/dependency files.
    total = 0
    files = 0
    for relative in git('ls-files', cwd=destination).splitlines():
        source = Path(relative)
        name = source.name.lower()
        if (source.suffix not in EXTENSIONS or EXCLUDED.intersection(source.parts)
                or any(part.lower() in {'test', 'tests', '__tests__'} for part in source.parts)
                or name.endswith('_test.go') or name.startswith('test_')
                or '.test.' in name or '.spec.' in name or name.endswith('.d.ts')):
            continue
        data = (destination / source).read_text(errors='replace')
        if 'Code generated' in data[:2048]:
            continue
        total += len(data.splitlines())
        files += 1
    licenses = [str(p.relative_to(destination)) for p in destination.glob('LICEN*')]
    return total, files, licenses


def main():
    records = json.loads(MANIFEST.read_text())
    failures = []
    for record in records:
        destination = ROOT / 'interviews/upstream' / record['name']
        remote = record['url'] + '.git'
        try:
            if destination.exists():
                if not (destination / '.git').is_dir():
                    raise RuntimeError(f'{destination} exists but is not a Git clone; left unchanged')
                actual = git('remote', 'get-url', 'origin', cwd=destination)
                if actual.removesuffix('.git') != remote.removesuffix('.git'):
                    raise RuntimeError(f'{destination}: unexpected origin; left unchanged')
                print(f"Using existing clone (no pull/reset): {record['name']}", flush=True)
            else:
                print(f"Cloning {record['url']}", flush=True)
                git('clone', '--depth', '1', remote, str(destination))
            commit = git('rev-parse', 'HEAD', cwd=destination)
            lines, files, licenses = inventory(destination)
            readmes = [str(p.relative_to(destination)) for p in destination.glob('README*')]
            record.update(clone_status='cloned', commit=commit, production_lines=lines,
                          production_source_files=files, licenses=licenses, question_files=readmes,
                          checked_at=datetime.now(timezone.utc).isoformat(), tests_run=False)
            print(f"  {commit}: {files} source files, {lines} physical source lines; questions: {readmes}")
        except (subprocess.SubprocessError, OSError, RuntimeError) as exc:
            detail = getattr(exc, 'stderr', None) or str(exc)
            record.update(clone_status='failed: ' + detail.strip(), tests_run=False)
            failures.append(record['name'])
            print(f"  FAILED {record['name']}: {detail.strip()}", file=sys.stderr)
        # Persist progress after each repository; a failed clone never looks successful.
        MANIFEST.write_text(json.dumps(records, indent=2) + '\n')
    if failures:
        print('Not downloaded: ' + ', '.join(failures), file=sys.stderr)
        return 1
    print('All sources cloned. Read each original README before installing or running it.')
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
