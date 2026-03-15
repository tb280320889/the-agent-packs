import argparse
import importlib.util
import json
import os


def load_query_mcp():
    module_path = os.path.join(
        os.path.dirname(__file__), "..", "query-mcp", "query_mcp.py"
    )
    spec = importlib.util.spec_from_file_location("query_mcp", module_path)
    if spec is None or spec.loader is None:
        raise RuntimeError("failed to load query_mcp module")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def load_request(path):
    with open(path, "r", encoding="utf-8") as handle:
        return json.load(handle)


def activation_result(request, route, bundle, status):
    return {
        "request_id": request.get("request_id"),
        "status": status,
        "main_pack": route.get("main_pack"),
        "artifacts": [],
        "validation_results": [],
        "handoff": None,
        "summary": route.get("route_reason", ""),
    }


def merge_hints(primary, secondary):
    merged = []
    for source in (primary or [], secondary or []):
        for item in source:
            if item not in merged:
                merged.append(item)
    return merged


def main():
    parser = argparse.ArgumentParser(description="Activation Entry (M1 minimal)")
    parser.add_argument("--db", default="blueprint/index/blueprint.db")
    parser.add_argument("--request", required=True, help="Activation request JSON path")
    args = parser.parse_args()

    query_mcp = load_query_mcp()
    request = load_request(args.request)
    if (
        not request.get("request_id")
        or not request.get("task")
        or "bounded_context" not in request
    ):
        print(
            json.dumps(
                {"status": "failed", "reason": "invalid activation request"},
                indent=2,
                ensure_ascii=True,
            )
        )
        return

    task = request.get("task", "")
    target_pack = request.get("target_pack")
    target_domain = request.get("target_domain")
    bounded = request.get("bounded_context") or {}
    selected_files = merge_hints(
        request.get("selected_files"), bounded.get("selected_files")
    )
    config_fragments = merge_hints(
        request.get("config_fragments"), bounded.get("config_fragments")
    )
    context_hints = merge_hints(request.get("context_hints"), bounded.get("host_hints"))
    conn = query_mcp.open_db(args.db)
    route = query_mcp.route_query(
        conn,
        level="L1",
        task=task,
        target_pack=target_pack,
        target_domain=target_domain,
        selected_files=selected_files,
        config_fragments=config_fragments,
        context_hints=context_hints,
        max_results=1,
    )
    if not route["candidates"]:
        conn.close()
        if target_domain:
            result = activation_result(
                request,
                {
                    "main_pack": None,
                    "route_reason": "insufficient evidence for L1, fallback to L0 recommended",
                },
                {},
                "partial",
            )
        else:
            result = activation_result(
                request,
                {"main_pack": None, "route_reason": "no route candidate"},
                {},
                "failed",
            )
        print(json.dumps(result, indent=2, ensure_ascii=True))
        return

    main_node = route["candidates"][0]["id"]
    bundle = query_mcp.build_context_bundle(conn, main_node, include_required=True)
    conn.close()

    if not bounded.get("selected_files") and not bounded.get("config_fragments"):
        result = activation_result(
            request,
            {
                "main_pack": "wxt-manifest",
                "route_reason": "bounded context missing required evidence",
            },
            bundle,
            "partial",
        )
    elif "handoff" in task.lower():
        result = activation_result(
            request,
            {
                "main_pack": "wxt-manifest",
                "route_reason": "handoff requested by task boundary",
            },
            bundle,
            "handoff",
        )
    else:
        result = activation_result(
            request,
            {
                "main_pack": "wxt-manifest",
                "route_reason": "route to L1.wxt.manifest with required cross-cutting lines",
            },
            bundle,
            "completed",
        )

    print(json.dumps(result, indent=2, ensure_ascii=True))


if __name__ == "__main__":
    main()
