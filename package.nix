{
  lib,
  buildGoModule,
  installShellFiles,
  makeWrapper,
  mpv,
  ffmpeg-full,
  mkvtoolnix-cli,
  xdg-utils,
}:

buildGoModule {
  pname = "cue";
  version = "unstable";

  src = lib.cleanSource ./.;
  vendorHash = "sha256-VJEbnosZZYYLWy3+Mscd1v0rYWBS71C1lTk3zxeSXsQ=";

  subPackages = [ "." ];
  env.CGO_ENABLED = "0";
  ldflags = [
    "-s"
    "-w"
    "-X main.Version=unstable"
  ];

  nativeBuildInputs = [
    installShellFiles
    makeWrapper
  ];
  nativeCheckInputs = [
    mpv
    ffmpeg-full
    mkvtoolnix-cli
  ];

  checkPhase = ''
    runHook preCheck
    go test ./...
    runHook postCheck
  '';

  postInstall = ''
    installShellCompletion --cmd cue \
      --bash <($out/bin/cue completion bash) \
      --zsh <($out/bin/cue completion zsh) \
      --fish <($out/bin/cue completion fish)

    wrapProgram $out/bin/cue \
      --suffix PATH : ${
        lib.makeBinPath [
          mpv
          ffmpeg-full
          mkvtoolnix-cli
          xdg-utils
        ]
      }
  '';

  meta = {
    description = "Terminal media client for Plex and Jellyfin";
    homepage = "https://github.com/SuperCoolPencil/cue";
    license = lib.licenses.mit;
    mainProgram = "cue";
    platforms = lib.platforms.linux;
  };
}
