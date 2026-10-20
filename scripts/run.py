#!/usr/bin/env python3
"""Build and run one example, with support for the macOS system linker."""
import argparse
import os
from pathlib import Path
import subprocess
import sys
import tempfile

ROOT = Path(__file__).resolve().parents[1]


def environment():
    env = os.environ.copy()
    env.update(GOCACHE=str(ROOT / '.cache/build'),
               GOMODCACHE=str(ROOT / '.cache/mod'), GOFLAGS='-mod=vendor',
               GOPROXY='off', GORACE='exitcode=66 atexit_sleep_ms=0')
    (ROOT / '.cache').mkdir(exist_ok=True)
    return env


def build(path, binary, env, race=False):
    command = ['go', 'build']
    if race:
        command.append('-race')
    use_system_linker = sys.platform == 'darwin' and Path(path).parent.name == 'pprof'
    if use_system_linker:
        # Only the net/http example needs this compatibility path. External
        # linking introduces cgo, which would suppress runtime deadlock detection.
        # Go 1.22's internal linker can omit LC_UUID on newer macOS releases.
        command.append('-ldflags=-linkmode=external')
    command += ['-o', str(binary), './' + str(Path(path).parent)]
    subprocess.run(command, cwd=ROOT, env=env, check=True,
                   stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True, timeout=180)
    if use_system_linker:
        # Old Go + new ld can also leave an invalid ad-hoc signature.
        subprocess.run(['codesign', '--force', '--sign', '-', str(binary)],
                       check=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT,
                       text=True, timeout=30)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('example', help='For example: basic/04_atomic')
    parser.add_argument('--race', action='store_true')
    args = parser.parse_args()
    example = (ROOT / 'examples' / args.example / 'main.go').resolve()
    if not example.is_relative_to(ROOT / 'examples') or not example.is_file():
        parser.error('Select a directory inside examples containing main.go')
    env = environment()
    with tempfile.TemporaryDirectory(prefix='run-', dir=ROOT / '.cache') as temp:
        binary = Path(temp) / 'example'
        build(str(example.relative_to(ROOT)), binary, env, race=args.race)
        return subprocess.run([str(binary)], cwd=ROOT, env=env, timeout=10).returncode


if __name__ == '__main__':
    try:
        sys.exit(main())
    except subprocess.CalledProcessError as error:
        print(error.stdout, file=sys.stderr)
        sys.exit(error.returncode)
    except subprocess.TimeoutExpired as error:
        print(error, file=sys.stderr)
        sys.exit(1)
