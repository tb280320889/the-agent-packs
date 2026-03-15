import argparse
import json
import os
import sqlite3
import hashlib
from datetime import datetime, timezone


REQUIRED_KEYS = [
    "id",
    "level",
    "domain",
    "subdomain",
    "capability",
    "title",
    "summary",
    "aliases",
    "triggers",
    "anti_triggers",
    "required_with",
    "may_include",
    "children",
    "entry_conditions",
    "stop_conditions",
]


def parse_frontmatter(text):
    lines = text.splitlines()
    if not lines or lines[0].strip() != "---":
        raise ValueError("frontmatter missing or not starting with ---")
    end_idx = None
    for idx in range(1, len(lines)):
        if lines[idx].strip() == "---":
            end_idx = idx
            break
    if end_idx is None:
        raise ValueError("frontmatter not closed with ---")

    fm_lines = lines[1:end_idx]
    body = "\n".join(lines[end_idx + 1 :]).strip()
    data = {}
    current_key = None
    for raw in fm_lines:
        line = raw.rstrip()
        if not line:
            continue
        if line.lstrip().startswith("-") and current_key:
            item = line.lstrip()[1:].strip()
            data[current_key].append(item)
            continue
        if ":" in line:
            key, value = line.split(":", 1)
            key = key.strip()
            value = value.strip()
            current_key = key
            if value == "":
                data[key] = []
            elif value == "[]":
                data[key] = []
            elif value == "null":
                data[key] = None
            else:
                data[key] = value
            continue
        raise ValueError(f"invalid frontmatter line: {line}")

    return data, body


def derive_id(level, domain, subdomain, stem):
    if level == "L0":
        return f"{level}.{domain}"
    if level == "L1":
        return f"{level}.{domain}.{stem}"
    return f"{level}.{domain}.{subdomain}.{stem}"


def compute_parent_id(level, domain, subdomain):
    if level == "L1":
        return f"L0.{domain}"
    if level == "L2":
        return f"L1.{domain}.{subdomain}"
    return None


def checksum_text(text):
    return hashlib.sha256(text.encode("utf-8")).hexdigest()


def ensure_schema(conn):
    cursor = conn.cursor()
    cursor.execute(
        """
        CREATE TABLE IF NOT EXISTS nodes (
            id TEXT PRIMARY KEY,
            level TEXT,
            domain TEXT,
            subdomain TEXT,
            capability TEXT,
            title TEXT,
            summary TEXT,
            path TEXT,
            parent_id TEXT,
            body_md TEXT,
            entry_conditions_json TEXT,
            stop_conditions_json TEXT,
            checksum TEXT,
            updated_at TEXT
        )
        """
    )
    cursor.execute(
        """
        CREATE TABLE IF NOT EXISTS node_meta (
            node_id TEXT PRIMARY KEY,
            aliases TEXT,
            triggers TEXT,
            anti_triggers TEXT,
            tags TEXT
        )
        """
    )
    cursor.execute(
        """
        CREATE TABLE IF NOT EXISTS edges (
            source_id TEXT,
            target_id TEXT,
            edge_type TEXT
        )
        """
    )
    conn.commit()


def load_blueprint_files(root_dir):
    md_files = []
    for dirpath, _, filenames in os.walk(root_dir):
        for name in filenames:
            if name.endswith(".md"):
                md_files.append(os.path.join(dirpath, name))
    return sorted(md_files)


def validate_and_collect(root_dir):
    nodes = []
    node_meta = []
    edges = []
    errors = []

    for path in load_blueprint_files(root_dir):
        if os.path.basename(path) in (
            "README.md",
            "frontmatter-examples.md",
            "schema.md",
        ):
            continue
        rel_path = os.path.relpath(path, root_dir)
        parts = rel_path.split(os.sep)
        if len(parts) < 3:
            errors.append({"path": path, "error": "invalid path structure"})
            continue
        level_dir = parts[0]
        domain_dir = parts[1]
        stem = os.path.splitext(parts[-1])[0]
        subdomain_dir = None
        if level_dir in ("L2", "L3"):
            subdomain_dir = parts[2]
        if level_dir == "L1":
            subdomain_dir = stem
        if level_dir == "L0":
            subdomain_dir = None

        with open(path, "r", encoding="utf-8") as handle:
            content = handle.read()
        try:
            fm, body = parse_frontmatter(content)
        except ValueError as exc:
            errors.append({"path": path, "error": str(exc)})
            continue

        missing = [key for key in REQUIRED_KEYS if key not in fm]
        if missing:
            errors.append({"path": path, "error": f"missing keys: {missing}"})
            continue

        derived_id = derive_id(level_dir, domain_dir, subdomain_dir, stem)
        if fm["id"] != derived_id:
            errors.append(
                {"path": path, "error": f"id mismatch: {fm['id']} != {derived_id}"}
            )
            continue
        if fm["level"] != level_dir:
            errors.append({"path": path, "error": "level mismatch"})
            continue
        if fm["domain"] != domain_dir:
            errors.append({"path": path, "error": "domain mismatch"})
            continue
        if level_dir != "L0" and fm["subdomain"] != subdomain_dir:
            errors.append({"path": path, "error": "subdomain mismatch"})
            continue

        parent_id = compute_parent_id(level_dir, domain_dir, subdomain_dir)
        nodes.append(
            {
                "id": fm["id"],
                "level": fm["level"],
                "domain": fm["domain"],
                "subdomain": fm["subdomain"],
                "capability": fm["capability"],
                "title": fm["title"],
                "summary": fm["summary"],
                "path": os.path.join("blueprint", rel_path).replace("\\", "/"),
                "parent_id": parent_id,
                "body_md": body,
                "entry_conditions_json": json.dumps(fm["entry_conditions"]),
                "stop_conditions_json": json.dumps(fm["stop_conditions"]),
                "checksum": checksum_text(content),
                "updated_at": datetime.now(timezone.utc).isoformat(),
            }
        )
        node_meta.append(
            {
                "node_id": fm["id"],
                "aliases": json.dumps(fm["aliases"]),
                "triggers": json.dumps(fm["triggers"]),
                "anti_triggers": json.dumps(fm["anti_triggers"]),
                "tags": json.dumps([]),
            }
        )
        for target in fm["children"]:
            edges.append(
                {"source_id": fm["id"], "target_id": target, "edge_type": "child"}
            )
        for target in fm["required_with"]:
            edges.append(
                {
                    "source_id": fm["id"],
                    "target_id": target,
                    "edge_type": "required_with",
                }
            )
        for target in fm["may_include"]:
            edges.append(
                {"source_id": fm["id"], "target_id": target, "edge_type": "may_include"}
            )

    return nodes, node_meta, edges, errors


def write_index(db_path, nodes, node_meta, edges):
    os.makedirs(os.path.dirname(db_path), exist_ok=True)
    conn = sqlite3.connect(db_path)
    ensure_schema(conn)
    cursor = conn.cursor()
    cursor.execute("DELETE FROM nodes")
    cursor.execute("DELETE FROM node_meta")
    cursor.execute("DELETE FROM edges")
    for node in nodes:
        cursor.execute(
            """
            INSERT INTO nodes (id, level, domain, subdomain, capability, title, summary, path, parent_id, body_md,
                               entry_conditions_json, stop_conditions_json, checksum, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
            (
                node["id"],
                node["level"],
                node["domain"],
                node["subdomain"],
                node["capability"],
                node["title"],
                node["summary"],
                node["path"],
                node["parent_id"],
                node["body_md"],
                node["entry_conditions_json"],
                node["stop_conditions_json"],
                node["checksum"],
                node["updated_at"],
            ),
        )
    for meta in node_meta:
        cursor.execute(
            """
            INSERT INTO node_meta (node_id, aliases, triggers, anti_triggers, tags)
            VALUES (?, ?, ?, ?, ?)
            """,
            (
                meta["node_id"],
                meta["aliases"],
                meta["triggers"],
                meta["anti_triggers"],
                meta["tags"],
            ),
        )
    for edge in edges:
        cursor.execute(
            "INSERT INTO edges (source_id, target_id, edge_type) VALUES (?, ?, ?)",
            (edge["source_id"], edge["target_id"], edge["edge_type"]),
        )
    conn.commit()
    conn.close()


def write_reports(report_dir, errors, edges, nodes):
    os.makedirs(report_dir, exist_ok=True)
    missing = []
    node_ids = {node["id"] for node in nodes}
    for edge in edges:
        if edge["target_id"] not in node_ids:
            missing.append(edge)
    with open(
        os.path.join(report_dir, "validation-report.json"), "w", encoding="utf-8"
    ) as handle:
        json.dump({"errors": errors}, handle, indent=2, ensure_ascii=True)
    with open(
        os.path.join(report_dir, "missing-reference-report.json"), "w", encoding="utf-8"
    ) as handle:
        json.dump({"missing": missing}, handle, indent=2, ensure_ascii=True)


def main():
    parser = argparse.ArgumentParser(description="Blueprint compiler (M1 minimal)")
    parser.add_argument("--root", default="blueprint", help="Blueprint root directory")
    parser.add_argument(
        "--db", default="blueprint/index/blueprint.db", help="SQLite output path"
    )
    parser.add_argument(
        "--report-dir", default="blueprint/index", help="Report output directory"
    )
    args = parser.parse_args()

    nodes, node_meta, edges, errors = validate_and_collect(args.root)
    write_index(args.db, nodes, node_meta, edges)
    write_reports(args.report_dir, errors, edges, nodes)

    if errors:
        print("compile completed with validation errors")
    else:
        print("compile completed")


if __name__ == "__main__":
    main()
