#!/usr/bin/env bash
# Run from any directory after GitHub CLI login, with normal Git/network access.
set -euo pipefail
cd "$(dirname "$0")/.."
repo_name=golang-concurrency-demo
owner=$(gh api user --jq .login)
if [[ "$owner" != aqutw ]]; then
  echo "This project's module path is github.com/aqutw/$repo_name; sign in as aqutw before publishing." >&2
  exit 1
fi
repo="$owner/$repo_name"
python3 scripts/verify.py
git_root=$(git rev-parse --show-toplevel 2>/dev/null || true)
if [[ -z "$git_root" ]]; then
  git init -b main
elif [[ "$git_root" != "$PWD" ]]; then
  echo 'This directory is inside another repository; publication stopped.' >&2
  exit 1
fi
publish_branch=$(git branch --show-current)
if [[ -z "$publish_branch" ]]; then
  echo 'Check out the branch to publish before running this script.' >&2
  exit 1
fi
if git remote get-url origin >/dev/null 2>&1; then
  remote=$(git remote get-url origin)
  case "$remote" in
    "https://github.com/$repo"|"https://github.com/$repo.git"|"git@github.com:$repo.git") ;;
    *) echo "Unexpected origin: $remote" >&2; exit 1 ;;
  esac
fi
# Stage only the artifacts intended for this public project.
git add .gitignore .github README.md THIRD_PARTY.md go.mod go.sum \
  examples scripts vendor
if ! git diff --cached --quiet; then
  git commit -m 'Add reproducible Go concurrency examples and verification'
fi
# List authenticated owned repositories so network/auth errors cannot be mistaken for 404.
visibility=$(gh repo list "$owner" --limit 1000 --json name,visibility \
  --jq '.[] | select(.name == "golang-concurrency-demo") | .visibility')
if [[ -z "$visibility" ]]; then
  gh repo create "$repo" --public --description 'Reproducible Go concurrency examples with race, deadlock, and leak verification'
elif [[ "$visibility" != PUBLIC ]]; then
  echo "$repo already exists and is private; publication stopped before changing its visibility." >&2
  exit 1
fi
if ! git remote get-url origin >/dev/null 2>&1; then
  git remote add origin "https://github.com/$repo.git"
fi
git push -u origin "$publish_branch"
remote_head=$(git ls-remote origin "refs/heads/$publish_branch" | cut -f1)
[[ "$remote_head" == "$(git rev-parse HEAD)" ]]
gh repo view "$repo" --json url,visibility --jq '.url + " (" + .visibility + ")'
