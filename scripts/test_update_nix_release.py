import base64
from pathlib import Path
import runpy
import unittest


ROOT = Path(__file__).resolve().parents[1]


class UpdateNixReleaseTests(unittest.TestCase):
    def test_updates_version_and_both_archive_hashes(self):
        script = ROOT / "scripts/update-nix-release.py"
        self.assertTrue(script.is_file(), "Release updater is missing")
        update = runpy.run_path(str(script))["update_package"]
        assets = [
            {
                "name": f"cue_9.8.7_Linux_{arch}.tar.gz",
                "digest": "sha256:" + digest,
            }
            for arch, digest in [("x86_64", "01" * 32), ("arm64", "02" * 32)]
        ]
        release = {"tag_name": "v9.8.7", "draft": False, "prerelease": False, "assets": assets}
        updated = update((ROOT / "package.nix").read_text(), release)
        self.assertIn('version = "9.8.7";', updated)
        for asset in assets:
            expected = base64.b64encode(bytes.fromhex(asset["digest"][7:])).decode()
            self.assertIn(f'hash = "sha256-{expected}";', updated)

    def test_rejects_draft_and_prerelease_metadata(self):
        update = runpy.run_path(str(ROOT / "scripts/update-nix-release.py"))["update_package"]
        for flag in ["draft", "prerelease"]:
            with self.subTest(flag=flag):
                release = {
                    "tag_name": "v9.8.7", "draft": False, "prerelease": False,
                    "assets": [
                        {"name": f"cue_9.8.7_Linux_{arch}.tar.gz", "digest": "sha256:" + "01" * 32}
                        for arch in ["x86_64", "arm64"]
                    ],
                }
                release[flag] = True
                with self.assertRaises(ValueError):
                    update((ROOT / "package.nix").read_text(), release)

    def test_rejects_incomplete_or_invalid_updates(self):
        import copy

        update = runpy.run_path(str(ROOT / "scripts/update-nix-release.py"))["update_package"]
        original = (ROOT / "package.nix").read_text()
        base = {
            "tag_name": "v9.8.7", "draft": False, "prerelease": False,
            "assets": [
                {"name": f"cue_9.8.7_Linux_{arch}.tar.gz", "digest": "sha256:" + "01" * 32}
                for arch in ["x86_64", "arm64"]
            ],
        }
        for case in ["tag", "missing_asset", "duplicate_asset", "digest", "layout"]:
            with self.subTest(case=case):
                release = copy.deepcopy(base)
                source = original
                if case == "tag":
                    release["tag_name"] = "invalid-tag"
                elif case == "missing_asset":
                    release["assets"].pop()
                elif case == "duplicate_asset":
                    release["assets"].append(release["assets"][0])
                elif case == "digest":
                    release["assets"][0]["digest"] = "sha256:abcd"
                else:
                    source = source.replace("  version =", "  renamed_version =")
                with self.assertRaises(ValueError):
                    update(source, release)


if __name__ == "__main__":
    unittest.main()
