#!/bin/bash


cp ./scripts/git-hooks/commit-msg ./scripts/git-hooks/pre-commit .git/hooks/
chmod +x .git/hooks/commit-msg
chmod +x .git/hooks/pre-commit
