#!/bin/bash
DIRNAME=$(dirname "$(readlink -f "$0")")
source "$DIRNAME/env.sh"

rm -rf "$DIRNAME/$OUTPUT_DIRECTORY/"
