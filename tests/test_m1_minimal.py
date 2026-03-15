import json
import os
import subprocess
import tempfile
import unittest


ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))


def run_cmd(args):
    result = subprocess.run(args, cwd=ROOT, capture_output=True, text=True, check=True)
    return result.stdout.strip()


class TestM1Minimal(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        run_cmd(
            [
                "python",
                "tools/compiler/compiler.py",
                "--root",
                "blueprint",
                "--db",
                "blueprint/index/blueprint.db",
                "--report-dir",
                "blueprint/index",
            ]
        )

    def test_compiler_reports_no_errors(self):
        with open(
            os.path.join(ROOT, "blueprint", "index", "validation-report.json"),
            "r",
            encoding="utf-8",
        ) as handle:
            report = json.load(handle)
        self.assertEqual(report.get("errors"), [])

        with open(
            os.path.join(ROOT, "blueprint", "index", "missing-reference-report.json"),
            "r",
            encoding="utf-8",
        ) as handle:
            report = json.load(handle)
        self.assertEqual(report.get("missing"), [])

    def test_route_target_pack_has_highest_priority(self):
        output = run_cmd(
            [
                "python",
                "tools/query-mcp/query_mcp.py",
                "route_query",
                "--db",
                "blueprint/index/blueprint.db",
                "--level",
                "L1",
                "--task",
                "unrelated task",
                "--target-pack",
                "security-permissions",
            ]
        )
        data = json.loads(output)
        self.assertEqual(data["candidates"][0]["id"], "L1.security.permissions")
        self.assertIn("target_pack match", data["candidates"][0]["reason"])

    def test_route_respects_target_domain(self):
        output = run_cmd(
            [
                "python",
                "tools/query-mcp/query_mcp.py",
                "route_query",
                "--db",
                "blueprint/index/blueprint.db",
                "--level",
                "L1",
                "--task",
                "permissions review for browser extension",
                "--target-domain",
                "wxt",
            ]
        )
        data = json.loads(output)
        self.assertTrue(data["candidates"])
        for candidate in data["candidates"]:
            self.assertEqual(candidate["id"].split(".")[1], "wxt")

    def test_route_target_pack_overrides_target_domain_conflict(self):
        output = run_cmd(
            [
                "python",
                "tools/query-mcp/query_mcp.py",
                "route_query",
                "--db",
                "blueprint/index/blueprint.db",
                "--level",
                "L1",
                "--task",
                "manifest permissions",
                "--target-pack",
                "security-permissions",
                "--target-domain",
                "wxt",
            ]
        )
        data = json.loads(output)
        self.assertEqual(data["candidates"][0]["id"], "L1.security.permissions")

    def test_route_antitrigger_exclusion(self):
        output = run_cmd(
            [
                "python",
                "tools/query-mcp/query_mcp.py",
                "route_query",
                "--db",
                "blueprint/index/blueprint.db",
                "--level",
                "L1",
                "--task",
                "tauri permission setup",
                "--target-domain",
                "wxt",
            ]
        )
        data = json.loads(output)
        self.assertEqual(data["candidates"], [])

    def test_frontmatter_summary_with_colon(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            l0_dir = os.path.join(temp_dir, "L0", "demo")
            os.makedirs(l0_dir, exist_ok=True)
            file_path = os.path.join(l0_dir, "overview.md")
            with open(file_path, "w", encoding="utf-8") as handle:
                handle.write(
                    """---
id: L0.demo
level: L0
domain: demo
subdomain: null
capability: null
title: Demo
summary: This summary has colon: valid content
aliases: []
triggers:
  - demo
anti_triggers: []
required_with: []
may_include: []
children: []
entry_conditions: []
stop_conditions: []
---

demo body
"""
                )

            db_path = os.path.join(temp_dir, "index", "blueprint.db")
            report_dir = os.path.join(temp_dir, "index")
            run_cmd(
                [
                    "python",
                    "tools/compiler/compiler.py",
                    "--root",
                    temp_dir,
                    "--db",
                    db_path,
                    "--report-dir",
                    report_dir,
                ]
            )

            with open(
                os.path.join(report_dir, "validation-report.json"),
                "r",
                encoding="utf-8",
            ) as handle:
                report = json.load(handle)
            self.assertEqual(report.get("errors"), [])

    def test_frontmatter_missing_required_keys_detected(self):
        with tempfile.TemporaryDirectory() as temp_dir:
            l0_dir = os.path.join(temp_dir, "L0", "demo")
            os.makedirs(l0_dir, exist_ok=True)
            file_path = os.path.join(l0_dir, "overview.md")
            with open(file_path, "w", encoding="utf-8") as handle:
                handle.write(
                    """---
id: L0.demo
level: L0
domain: demo
summary: missing required keys
---

demo body
"""
                )

            db_path = os.path.join(temp_dir, "index", "blueprint.db")
            report_dir = os.path.join(temp_dir, "index")
            run_cmd(
                [
                    "python",
                    "tools/compiler/compiler.py",
                    "--root",
                    temp_dir,
                    "--db",
                    db_path,
                    "--report-dir",
                    report_dir,
                ]
            )

            with open(
                os.path.join(report_dir, "validation-report.json"),
                "r",
                encoding="utf-8",
            ) as handle:
                report = json.load(handle)
            self.assertTrue(report.get("errors"))

    def test_activation_completed_path(self):
        output = run_cmd(
            [
                "python",
                "tools/activation-entry/activation_entry.py",
                "--db",
                "blueprint/index/blueprint.db",
                "--request",
                "fixtures/activation-request.sample.json",
            ]
        )
        data = json.loads(output)
        self.assertEqual(data["status"], "completed")

    def test_activation_handoff_path(self):
        request = {
            "request_id": "req-handoff",
            "task": "handoff to next pack after manifest review",
            "target_domain": "wxt",
            "bounded_context": {
                "selected_files": ["wxt.config.ts"],
                "config_fragments": ["manifest.permissions"],
            },
        }
        with tempfile.NamedTemporaryFile(
            "w", suffix=".json", delete=False, encoding="utf-8"
        ) as temp:
            json.dump(request, temp)
            temp_path = temp.name
        try:
            output = run_cmd(
                [
                    "python",
                    "tools/activation-entry/activation_entry.py",
                    "--db",
                    "blueprint/index/blueprint.db",
                    "--request",
                    temp_path,
                ]
            )
            data = json.loads(output)
            self.assertEqual(data["status"], "handoff")
        finally:
            os.remove(temp_path)

    def test_activation_partial_when_domain_known_but_no_candidate(self):
        request = {
            "request_id": "req-partial",
            "task": "totally unrelated text",
            "target_domain": "wxt",
            "bounded_context": {
                "selected_files": [],
                "config_fragments": [],
            },
        }
        with tempfile.NamedTemporaryFile(
            "w", suffix=".json", delete=False, encoding="utf-8"
        ) as temp:
            json.dump(request, temp)
            temp_path = temp.name
        try:
            output = run_cmd(
                [
                    "python",
                    "tools/activation-entry/activation_entry.py",
                    "--db",
                    "blueprint/index/blueprint.db",
                    "--request",
                    temp_path,
                ]
            )
            data = json.loads(output)
            self.assertEqual(data["status"], "partial")
        finally:
            os.remove(temp_path)

    def test_activation_failed_when_request_invalid(self):
        request = {
            "request_id": "req-failed",
            "task": "manifest",
        }
        with tempfile.NamedTemporaryFile(
            "w", suffix=".json", delete=False, encoding="utf-8"
        ) as temp:
            json.dump(request, temp)
            temp_path = temp.name
        try:
            output = run_cmd(
                [
                    "python",
                    "tools/activation-entry/activation_entry.py",
                    "--db",
                    "blueprint/index/blueprint.db",
                    "--request",
                    temp_path,
                ]
            )
            data = json.loads(output)
            self.assertEqual(data["status"], "failed")
        finally:
            os.remove(temp_path)


if __name__ == "__main__":
    unittest.main()
