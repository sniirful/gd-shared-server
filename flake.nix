{
  description = "Nix Flake-based FHS development environment for gd-shared-server";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.05";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }: flake-utils.lib.eachSystem flake-utils.lib.allSystems (system:
    let
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShell = pkgs.mkShell rec {
        fhsName = "gd-shared-server-fhs-env";
        fhsScriptName = "gd-shared-server-fhs-env-script";
        fhsPackages = with pkgs; [
          bash
          git
          go_1_21
        ];

        packages = [
          (pkgs.writeShellScriptBin fhsScriptName ''
            #!/bin/bash
            go version
            bash
          '')
          (pkgs.buildFHSEnv {
            name = fhsName;
            targetPkgs = pkgs: fhsPackages;

            # this whole shenanigan was made because putting
            # the raw commands in the "runScript" just wouldn't
            # work
            runScript = fhsScriptName;
          })
        ];
        # TODO: all this mess creates a level-4 shell (`echo $SHLVL`)
        shellHook = ''
          ${fhsName}
          # the "exit 0" is necessary for not having to type in
          # "exit" another time after quitting the nix develop
          exit 0
        '';
      };
    }
  );
}
