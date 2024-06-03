# Overview
“gd-shared-server” stands for “Google Drive Shared Server”. It allows you and your friends to host any server, such as Minecraft, while always keeping the most recent version on a Google Drive account and thus automating the task of having to manually share and download the latest files once the server is closed.

This can be used when you have no hosting server available and want to use the direct players’ PCs to host the server they will play in.

<!-- TODO: add pictures or a GIF -->

## Security implications
“gd-shared-server” clearly warns its users about the security issues of this hosting method and therefore implies that it will be used with very trusted people. The two main implications are:
- this project makes use of full access control to a Google Drive account. Even if this project per se will not perform any malicious actions on your files, the authentication files required for it to work can easily be used by a bad actor to manage all the files in the account. **It is strongly recommended to use an empty or a new Drive because of this**;
- this project directly executes commands and applications that are blindly downloaded from Google Drive and uploaded from anyone who has access to it. This means that **a bad actor can easily upload, either manually or via this project, any malicious scripts or executables, thus potentially infecting your PC and possibly leading to data theft, data encryption or even data loss**.

The user is considered fully aware of all of this and, as such, none of the maintainers are to be considered responsible for any malicious activity caused by the use of this project.

## How to use
At the time of writing this, no GitHub Releases are provided. As such, you will have to build the project yourself. See [Building](#building).

The content that will always be shared between anyone who has access to the account is everything that you put inside the `server` folder. To run the actual server, you need to:
1. read the README.txt and acknowledge what’s written inside;
2. put the server files inside the `server` folder;
3. edit the `command` file (or `command.platform` specific file) to set your server entry point (e.g. `java -jar minecraft-server.jar`). The command file is assumed to be Batch on Windows and Bash on all other platforms. The command will be launched inside the `server` folder;
4. start the application by executing the file `server.platform`;
5. follow the instructions given by the application.

## Building
At the time of writing this, to build this project it is required that you use Bash and Go >= 1.21.5.

If you’re using NixOS and running a Flake-based system, you can simply run the following command in this folder:
```bash
nix develop
```
and you will have all the necessary dependencies set up for you.

Once the dependencies are satisfied, you have to:
1. edit the `env.sh` file and edit the `BUILD_PLATFORMS` variable to remove the ones you do not need and add the ones you do need. For a list of platforms - as well as architectures - you can check [here](https://pkg.go.dev/internal/platform);
2. run the `build.sh` script.

By default, the output folder is `dist`. After building, you can simply move everything from `dist` to anywhere you want and start using it like explained in [How to use](#how-to-use).
