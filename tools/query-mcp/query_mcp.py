import argparse
import json
import os
import sqlite3


PACK_NODE_MAP = {
    "wxt-manifest": "L1.wxt.manifest",
    "security-permissions": "L1.security.permissions",
    "release-store-review": "L1.release.store-review",
}


def open_db(db_path):
    if not os.path.exists(db_path):
        raise FileNotFoundError(f"db not found: {db_path}")
    return sqlite3.connect(db_path)


def fetch_nodes(conn, level=None):
    cursor = conn.cursor()
    if level:
        cursor.execute(
            "SELECT id, level, domain, subdomain, title, summary FROM nodes WHERE level = ?",
            (level,),
        )
    else:
        cursor.execute("SELECT id, level, domain, subdomain, title, summary FROM nodes")
    return cursor.fetchall()


def fetch_meta(conn, node_id):
    cursor = conn.cursor()
    cursor.execute(
        "SELECT aliases, triggers, anti_triggers FROM node_meta WHERE node_id = ?",
        (node_id,),
    )
    row = cursor.fetchone()
    if not row:
        return {"aliases": [], "triggers": [], "anti_triggers": []}
    return {
        "aliases": json.loads(row[0] or "[]"),
        "triggers": json.loads(row[1] or "[]"),
        "anti_triggers": json.loads(row[2] or "[]"),
    }


def fetch_edges(conn, source_id, edge_type):
    cursor = conn.cursor()
    cursor.execute(
        "SELECT target_id FROM edges WHERE source_id = ? AND edge_type = ?",
        (source_id, edge_type),
    )
    return [row[0] for row in cursor.fetchall()]


def tokenize(text):
    return text.lower()


def list_token_score(task_text, values, weight, reason_prefix):
    score = 0.0
    reason = []
    for value in values or []:
        token = value.lower()
        if token and token in task_text:
            score += weight
            reason.append(f"{reason_prefix}:{value}")
    return score, reason


def evidence_score(task_text, selected_files, config_fragments, context_hints):
    score = 0.0
    reason = []
    for file_path in selected_files or []:
        file_token = os.path.basename(file_path).lower()
        if file_token and file_token in task_text:
            score += 0.4
            reason.append(f"selected_file:{file_token}")
    for fragment in config_fragments or []:
        token = str(fragment).lower()
        if token and token in task_text:
            score += 0.3
            reason.append(f"config_fragment:{fragment}")
    for hint in context_hints or []:
        token = str(hint).lower()
        if token and token in task_text:
            score += 0.2
            reason.append(f"context_hint:{hint}")
    return score, reason


def score_candidate(
    task,
    candidate,
    meta,
    target_domain=None,
    selected_files=None,
    config_fragments=None,
    context_hints=None,
):
    score = 0.0
    reason = []
    task_text = tokenize(task)
    if target_domain and candidate[2] == target_domain:
        score += 3.0
        reason.append("target_domain match")
    trigger_score, trigger_reason = list_token_score(
        task_text, meta["triggers"], 1.0, "trigger"
    )
    alias_score, alias_reason = list_token_score(
        task_text, meta["aliases"], 0.5, "alias"
    )
    evidence, evidence_reason = evidence_score(
        task_text, selected_files, config_fragments, context_hints
    )
    score += trigger_score + alias_score + evidence
    reason.extend(trigger_reason)
    reason.extend(alias_reason)
    reason.extend(evidence_reason)
    return score, reason


def route_query(
    conn,
    level,
    task,
    target_pack=None,
    target_domain=None,
    selected_files=None,
    config_fragments=None,
    context_hints=None,
    max_results=3,
):
    if target_pack:
        mapped_node = PACK_NODE_MAP.get(target_pack)
        if mapped_node:
            direct = read_node(conn, mapped_node, section="summary")
            if direct and direct.get("level") == level:
                return {
                    "candidates": [
                        {
                            "id": direct["id"],
                            "title": direct["title"],
                            "summary": direct["summary"],
                            "score": 99.0,
                            "reason": ["target_pack match"],
                        }
                    ],
                    "must_include": fetch_edges(conn, direct["id"], "required_with"),
                }

    candidates = []
    for row in fetch_nodes(conn, level=level):
        if target_domain and row[2] != target_domain:
            continue
        meta = fetch_meta(conn, row[0])
        if any(trigger.lower() in tokenize(task) for trigger in meta["anti_triggers"]):
            continue
        score, reason = score_candidate(
            task,
            row,
            meta,
            target_domain=target_domain,
            selected_files=selected_files,
            config_fragments=config_fragments,
            context_hints=context_hints,
        )
        if score > 0:
            candidates.append(
                {
                    "id": row[0],
                    "title": row[4],
                    "summary": row[5],
                    "score": score,
                    "reason": reason,
                }
            )

    candidates.sort(key=lambda item: (-item["score"], item["id"]))
    top = candidates[: max_results or 3]
    must_include = []
    if top:
        must_include = fetch_edges(conn, top[0]["id"], "required_with")
    return {
        "candidates": top,
        "must_include": must_include,
    }


def read_node(conn, node_id, section="summary"):
    cursor = conn.cursor()
    cursor.execute(
        "SELECT id, title, summary, body_md, level FROM nodes WHERE id = ?", (node_id,)
    )
    row = cursor.fetchone()
    if not row:
        return None
    if section == "summary":
        return {"id": row[0], "title": row[1], "summary": row[2], "level": row[4]}
    return {"id": row[0], "title": row[1], "body": row[3], "level": row[4]}


def build_context_bundle(
    conn,
    main_node,
    include_required=True,
    include_may_include=False,
    include_children=False,
):
    bundle = {
        "main": None,
        "required": [],
        "execution_children": [],
        "deferred": [],
        "recommended_validators": [],
        "recommended_artifacts": [],
    }

    main = read_node(conn, main_node, section="summary")
    if not main:
        return bundle
    bundle["main"] = main

    def append_node(target_id, collection):
        node = read_node(conn, target_id, section="summary")
        if not node:
            return
        if node.get("level") == "L3":
            bundle["deferred"].append(node)
            return
        collection.append(node)

    if include_required:
        for target_id in fetch_edges(conn, main_node, "required_with"):
            append_node(target_id, bundle["required"])
    if include_children:
        for target_id in fetch_edges(conn, main_node, "child"):
            append_node(target_id, bundle["execution_children"])
    if include_may_include:
        for target_id in fetch_edges(conn, main_node, "may_include"):
            append_node(target_id, bundle["execution_children"])

    return bundle


def expand_node(conn, node_id, edge_type="child"):
    return fetch_edges(conn, node_id, edge_type)


def main():
    parser = argparse.ArgumentParser(description="Blueprint Query MCP (M1 minimal)")
    parser.add_argument(
        "command",
        choices=["route_query", "read_node", "build_context_bundle", "expand_node"],
    )
    parser.add_argument("--db", default="blueprint/index/blueprint.db")
    parser.add_argument("--level", default=None)
    parser.add_argument("--task", default="")
    parser.add_argument("--target-pack", default=None)
    parser.add_argument("--target-domain", default=None)
    parser.add_argument("--selected-files", default="")
    parser.add_argument("--config-fragments", default="")
    parser.add_argument("--context-hints", default="")
    parser.add_argument("--max-results", type=int, default=3)
    parser.add_argument("--node-id", default=None)
    parser.add_argument("--edge-type", default="child")
    parser.add_argument("--include-required", action="store_true")
    parser.add_argument("--include-may-include", action="store_true")
    parser.add_argument("--include-children", action="store_true")
    args = parser.parse_args()

    selected_files = [item for item in args.selected_files.split(",") if item]
    config_fragments = [item for item in args.config_fragments.split(",") if item]
    context_hints = [item for item in args.context_hints.split(",") if item]

    conn = open_db(args.db)
    if args.command == "route_query":
        result = route_query(
            conn,
            level=args.level,
            task=args.task,
            target_pack=args.target_pack,
            target_domain=args.target_domain,
            selected_files=selected_files,
            config_fragments=config_fragments,
            context_hints=context_hints,
            max_results=args.max_results,
        )
    elif args.command == "read_node":
        result = read_node(conn, args.node_id, section="summary")
    elif args.command == "build_context_bundle":
        result = build_context_bundle(
            conn,
            args.node_id,
            include_required=args.include_required,
            include_may_include=args.include_may_include,
            include_children=args.include_children,
        )
    else:
        result = expand_node(conn, args.node_id, edge_type=args.edge_type)

    conn.close()
    print(json.dumps(result, indent=2, ensure_ascii=True))


if __name__ == "__main__":
    main()
