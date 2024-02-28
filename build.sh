#!/bin/bash
DIRNAME=$(dirname "$(readlink -f "$0")")
source "$DIRNAME/env.sh"

mkdir "$DIRNAME/$OUTPUT_DIRECTORY"

for i in ${!BUILD_PLATFORMS[@]}; do
    platform="${BUILD_PLATFORMS[$i]}"
    outfile="$DIRNAME/$OUTPUT_DIRECTORY/$OUTPUT_BASE_NAME-$platform"
    if [ "$platform" == "windows" ]; then
        outfile="$outfile.exe"
    fi

    echo "Building $platform..."
    GOOS=$platform go build -o "$outfile" "$DIRNAME/main.go"
done

echo "Creating necessary files and folders..."
mkdir "$DIRNAME/$OUTPUT_DIRECTORY/server"
touch "$DIRNAME/$OUTPUT_DIRECTORY/start.command"
touch "$DIRNAME/$OUTPUT_DIRECTORY/logfile"
# now we create the readme
echo "TODO" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"

echo "Done."
