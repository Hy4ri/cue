"""Refresh the Nix package from the latest stable GitHub release."""

import base64
import json
from pathlib import Path
import re
import subprocess


def update_package(source, release):
    if release["draft"] or release["prerelease"]:
        raise ValueError("Only published stable releases can update the Nix package")
    tag = re.fullmatch(r"v(\d+\.\d+\.\d+)", release["tag_name"])
    if tag is None:
        raise ValueError("Expected a release tag in the form vMAJOR.MINOR.PATCH")
    version = tag[1]
    updated, count = re.subn(
        r'^  version = "[^"]+";', f'  version = "{version}";', source, flags=re.MULTILINE
    )
    if count != 1:
        raise ValueError("Expected exactly one Nix version pin")
    for system, arch in [("x86_64-linux", "x86_64"), ("aarch64-linux", "arm64")]:
        name = f"cue_{version}_Linux_{arch}.tar.gz"
        assets = [asset for asset in release["assets"] if asset["name"] == name]
        if len(assets) != 1:
            raise ValueError(f"Expected exactly one release asset named {name}")
        digest = assets[0].get("digest") or ""
        if re.fullmatch(r"sha256:[0-9a-f]{64}", digest) is None:
            raise ValueError(f"Missing or invalid GitHub SHA-256 digest for {name}")
        sri = "sha256-" + base64.b64encode(bytes.fromhex(digest[7:])).decode()
        pattern = rf'({re.escape(system)} = \{{\s+arch = "{arch}";\s+hash = ")[^"]+(";)'
        updated, count = re.subn(pattern, lambda match: match[1] + sri + match[2], updated)
        if count != 1:
            raise ValueError(f"Expected exactly one archive hash for {system}")
    return updated


if __name__ == "__main__":
    release = json.loads(subprocess.check_output(
        ["gh", "api", "repos/SuperCoolPencil/cue/releases/latest"], text=True
    ))
    package = Path(__file__).resolve().parents[1] / "package.nix"
    source = package.read_text()
    updated = update_package(source, release)
    if updated != source:
        package.write_text(updated)
        print(f"Updated Nix package to {release['tag_name']}")
    else:
        print("Nix package is already up to date")
