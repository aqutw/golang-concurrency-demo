#!/usr/bin/env python3
"""Build and execute every documented example, including expected failures."""
import argparse
import json
from pathlib import Path
import subprocess
import sys
import tempfile

from run import build, environment

ROOT = Path(__file__).resolve().parents[1]


def run(command, env, timeout=180):
    result = subprocess.run(command, cwd=ROOT, env=env, text=True,
                            stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                            timeout=timeout)
    return result.returncode, result.stdout


def coverage(entries):
    paths = [entry['path'] for entry in entries]
    actual = {str(p.relative_to(ROOT)) for p in (ROOT / 'examples').rglob('*.go')}
    assert len(paths) == len(set(paths)), 'Duplicate manifest paths'
    assert set(paths) == actual, 'Go files and manifest do not match'
    for entry in entries:
        assert entry['contains'], f'Missing behavior assertions: {entry["path"]}'
    return len(paths)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--repeat', type=int, default=1,
                        help='Repeat execution to check scheduling independence')
    args = parser.parse_args()
    if args.repeat < 1:
        parser.error('--repeat must be positive')
    entries = json.loads((ROOT / 'scripts/examples.json').read_text())
    count = coverage(entries)
    env = environment()
    print(f'Coverage: {count} Go examples match the manifest', flush=True)
    files = [entry['path'] for entry in entries]
    code, output = run(['gofmt', '-l', *files], env)
    assert code == 0 and not output.strip(), f'gofmt failed:\n{output}'
    code, output = run(['go', 'vet', './...'], env)
    assert code == 0, f'go vet failed:\n{output}'
    cache = ROOT / '.cache'
    cache.mkdir(exist_ok=True)
    with tempfile.TemporaryDirectory(prefix='verify-', dir=cache) as temp:
        for index, entry in enumerate(entries, 1):
            binary = str(Path(temp) / f'example-{index}')
            # Race-instrument every successful/leaking example, and the race reproducer.
            # Run the deadlock without instrumentation to retain Go's deadlock detector.
            build(entry['path'], binary, env, race=entry['kind'] != 'deadlock')
            expected = {'race': 66, 'deadlock': 2}.get(entry['kind'], 0)
            for iteration in range(args.repeat):
                code, output = run([binary], env, timeout=10)
                assert code == expected, (
                    f'{entry["path"]}: exit {code}, expected {expected}\n{output}')
                for marker in entry['contains']:
                    assert marker in output, f'{entry["path"]}: missing {marker!r}\n{output}'
                for marker, expected_count in entry['counts'].items():
                    assert output.count(marker) == expected_count, (
                        f'{entry["path"]}: {marker!r} expected {expected_count} times\n{output}')
                if entry['kind'] != 'race':
                    assert 'WARNING: DATA RACE' not in output, output
            print(f'PASS {index:02d}/{len(entries)} [{entry["kind"]}] {entry["path"]}', flush=True)
    print(f'PASS: {len(entries)} examples x {args.repeat} runs; vet, format, manifest coverage', flush=True)


if __name__ == '__main__':
    try:
        main()
    except (AssertionError, subprocess.TimeoutExpired, subprocess.CalledProcessError) as error:
        print(f'FAIL: {error}', file=sys.stderr)
        if isinstance(error, subprocess.CalledProcessError):
            print(error.stdout, file=sys.stderr)
        sys.exit(1)
