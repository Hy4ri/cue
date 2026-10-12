{
  lib,
  stdenvNoCC,
  fetchurl,
}:

let
  version = "1.3.2";
  releases = {
    x86_64-linux = {
      arch = "x86_64";
      hash = "sha256-PFS2jlSH/1YUTorxH55idsQEbhC3YmyUfoYHfcpoz3A=";
    };
    aarch64-linux = {
      arch = "arm64";
      hash = "sha256-judalHfGoEHGYECYD9KBhBtH02u9H8CJjY+K3HHu0FM=";
    };
  };
  release = releases.${stdenvNoCC.hostPlatform.system};
in
stdenvNoCC.mkDerivation {
  pname = "cue";
  inherit version;

  src = fetchurl {
    url = "https://github.com/SuperCoolPencil/cue/releases/download/v${version}/cue_${version}_Linux_${release.arch}.tar.gz";
    inherit (release) hash;
  };
  sourceRoot = ".";

  dontConfigure = true;
  dontBuild = true;
  # The release binaries are static; keep them unchanged.
  dontPatchELF = true;
  dontStrip = true;

  installPhase = ''
    runHook preInstall

    install -Dm755 cue $out/bin/cue

    runHook postInstall
  '';

  meta = {
    description = "Terminal media client for Plex and Jellyfin";
    homepage = "https://github.com/SuperCoolPencil/cue";
    license = lib.licenses.mit;
    mainProgram = "cue";
    platforms = builtins.attrNames releases;
  };
}
