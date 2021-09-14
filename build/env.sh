#!/bin/sh

set -e

if [ ! -f "build/env.sh" ]; then
    echo "$0 must be run from the root of the repository."
    exit 2
fi

# Create fake Go workspace if it doesn't exist yet.
workspace="$PWD/build/_workspace"
root="$PWD"
cardir="$workspace/src/github.com/RosettaFlow"
if [ ! -L "$cardir/Carrier-Go" ]; then
    mkdir -p "$cardir"
    cd "$cardir"
    ln -s ../../../../../. Carrier-Go
    cd "$root"
fi

# Set up the environment to use the workspace.
GOPATH="$workspace"
export GOPATH

# Run the command inside the workspace.
cd "$cardir/Carrier-Go"
PWD="$cardir/Carrier-Go"

# Launch the arguments with the configured environment.
exec "$@"
