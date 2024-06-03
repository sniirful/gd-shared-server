#!/bin/bash
DIRNAME=$(dirname "$(readlink -f "$0")")
source "$DIRNAME/env.sh"

mkdir "$DIRNAME/$OUTPUT_DIRECTORY"

for i in ${!BUILD_PLATFORMS[@]}; do
    platform="${BUILD_PLATFORMS[$i]}"
    outfile="$DIRNAME/$OUTPUT_DIRECTORY/$OUTPUT_BASE_NAME.$platform"
    commandfile="$DIRNAME/$OUTPUT_DIRECTORY/command.$platform"
    if [ "$platform" == "windows" ]; then
        outfile="$outfile.exe"
        commandfile="$commandfile.bat"
    fi

    echo "Building $platform..."
    GOOS=$platform go build -o "$outfile" "$DIRNAME/main.go"
    touch "$commandfile"
done

echo "Creating necessary files and folders..."
mkdir "$DIRNAME/$OUTPUT_DIRECTORY/server"
touch "$DIRNAME/$OUTPUT_DIRECTORY/command"
touch "$DIRNAME/$OUTPUT_DIRECTORY/logfile"
# now we create the readme
echo "Last change: 2024/03/02" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "Be careful when using this product. Make sure you trust the" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "people who sent you these files. Using it improperly may" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "severely harm your device, making you vulnerable to viruses" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "or any hacking attack whatsoever." >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "Before using, you need to give this product the ability to" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "access your Google Drive account. No data will be shared" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "with third-party, but make sure you trust the people you are" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "sharing this with, as this kind of information exposes all" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "the content in your Google Drive space." >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "Go to https://console.cloud.google.com/apis/credentials and" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "create a new OAuth client ID. Once done, download the JSON" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "and copy it in this folder, renaming it into oauth.json." >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "When starting the server, the given command will be run" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "inside the server folder, which is where you will have to" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "put all your files. Only the server folder and the logfile" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "will be shared among you and your friends, so if you want" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "to change the start command, make sure to give your friends" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "an updated version of this product. If the command change is" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "important and you do not want them to run the program with" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "an outdated command, you might delete the old OAuth client" >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"
echo "and create a new one to distribute." >> "$DIRNAME/$OUTPUT_DIRECTORY/README.txt"

echo "Done."
