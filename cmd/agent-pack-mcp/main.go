package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"the-agent-packs/internal/activation"
	"the-agent-packs/internal/compiler"
	"the-agent-packs/internal/query"
)

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(true)
	return enc.Encode(v)
}

func parseCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func cmdCompile(args []string) error {
	fs := flag.NewFlagSet("compile", flag.ContinueOnError)
	root := fs.String("root", "blueprint", "Blueprint root directory")
	db := fs.String("db", "blueprint/index/blueprint.db", "SQLite output path")
	reportDir := fs.String("report-dir", "blueprint/index", "Report output directory")
	if err := fs.Parse(args); err != nil {
		return err
	}
	errs, err := compiler.Compile(*root, *db, *reportDir)
	if err != nil {
		return err
	}
	if len(errs) > 0 {
		fmt.Println("compile completed with validation errors")
		return nil
	}
	fmt.Println("compile completed")
	return nil
}

func cmdRouteQuery(args []string) error {
	fs := flag.NewFlagSet("route_query", flag.ContinueOnError)
	db := fs.String("db", "blueprint/index/blueprint.db", "SQLite path")
	level := fs.String("level", "", "Blueprint level")
	task := fs.String("task", "", "Task text")
	targetPack := fs.String("target-pack", "", "Target pack")
	targetDomain := fs.String("target-domain", "", "Target domain")
	selectedFiles := fs.String("selected-files", "", "Comma-separated selected files")
	configFragments := fs.String("config-fragments", "", "Comma-separated config fragments")
	contextHints := fs.String("context-hints", "", "Comma-separated context hints")
	maxResults := fs.Int("max-results", 3, "Max candidates")
	if err := fs.Parse(args); err != nil {
		return err
	}
	conn, err := query.OpenDB(*db)
	if err != nil {
		return err
	}
	defer conn.Close()

	var targetPackPtr *string
	if strings.TrimSpace(*targetPack) != "" {
		targetPackPtr = targetPack
	}
	var targetDomainPtr *string
	if strings.TrimSpace(*targetDomain) != "" {
		targetDomainPtr = targetDomain
	}

	result, err := query.RouteQuery(
		conn,
		*level,
		*task,
		targetPackPtr,
		targetDomainPtr,
		parseCSV(*selectedFiles),
		parseCSV(*configFragments),
		parseCSV(*contextHints),
		*maxResults,
	)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func cmdReadNode(args []string) error {
	fs := flag.NewFlagSet("read_node", flag.ContinueOnError)
	db := fs.String("db", "blueprint/index/blueprint.db", "SQLite path")
	nodeID := fs.String("node-id", "", "Node ID")
	section := fs.String("section", "summary", "Section type")
	if err := fs.Parse(args); err != nil {
		return err
	}
	conn, err := query.OpenDB(*db)
	if err != nil {
		return err
	}
	defer conn.Close()

	result, err := query.ReadNode(conn, *nodeID, *section)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func cmdBuildBundle(args []string) error {
	fs := flag.NewFlagSet("build_context_bundle", flag.ContinueOnError)
	db := fs.String("db", "blueprint/index/blueprint.db", "SQLite path")
	nodeID := fs.String("node-id", "", "Main node ID")
	includeRequired := fs.Bool("include-required", false, "Include required nodes")
	includeMayInclude := fs.Bool("include-may-include", false, "Include may_include nodes")
	includeChildren := fs.Bool("include-children", false, "Include child nodes")
	if err := fs.Parse(args); err != nil {
		return err
	}
	conn, err := query.OpenDB(*db)
	if err != nil {
		return err
	}
	defer conn.Close()

	result, err := query.BuildContextBundle(conn, *nodeID, *includeRequired, *includeMayInclude, *includeChildren)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func cmdExpandNode(args []string) error {
	fs := flag.NewFlagSet("expand_node", flag.ContinueOnError)
	db := fs.String("db", "blueprint/index/blueprint.db", "SQLite path")
	nodeID := fs.String("node-id", "", "Node ID")
	edgeType := fs.String("edge-type", "child", "Edge type")
	if err := fs.Parse(args); err != nil {
		return err
	}
	conn, err := query.OpenDB(*db)
	if err != nil {
		return err
	}
	defer conn.Close()

	result, err := query.ExpandNode(conn, *nodeID, *edgeType)
	if err != nil {
		return err
	}
	return printJSON(result)
}

func cmdActivate(args []string) error {
	fs := flag.NewFlagSet("activate", flag.ContinueOnError)
	db := fs.String("db", "blueprint/index/blueprint.db", "SQLite path")
	requestPath := fs.String("request", "", "Activation request JSON path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*requestPath) == "" {
		return errors.New("--request is required")
	}
	conn, err := query.OpenDB(*db)
	if err != nil {
		return err
	}
	defer conn.Close()

	result, err := activation.Execute(conn, *requestPath)
	if err != nil {
		return err
	}
	return printJSON(result)
}

type mcpRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type mcpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type mcpResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id,omitempty"`
	Result  any       `json:"result,omitempty"`
	Error   *mcpError `json:"error,omitempty"`
}

func mcpWrite(resp mcpResponse) {
	b, _ := json.Marshal(resp)
	fmt.Println(string(b))
}

func mcpResult(id any, result any) {
	mcpWrite(mcpResponse{JSONRPC: "2.0", ID: id, Result: result})
}

func mcpErr(id any, code int, msg string) {
	mcpWrite(mcpResponse{JSONRPC: "2.0", ID: id, Error: &mcpError{Code: code, Message: msg}})
}

func parseParamString(raw json.RawMessage, key string) string {
	obj := map[string]any{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ""
	}
	v, ok := obj[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func parseParamBool(raw json.RawMessage, key string, def bool) bool {
	obj := map[string]any{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return def
	}
	v, ok := obj[key]
	if !ok || v == nil {
		return def
	}
	b, ok := v.(bool)
	if !ok {
		return def
	}
	return b
}

func parseParamInt(raw json.RawMessage, key string, def int) int {
	obj := map[string]any{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return def
	}
	v, ok := obj[key]
	if !ok || v == nil {
		return def
	}
	f, ok := v.(float64)
	if !ok {
		return def
	}
	return int(f)
}

func parseParamStringList(raw json.RawMessage, key string) []string {
	obj := map[string]any{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return []string{}
	}
	v, ok := obj[key]
	if !ok || v == nil {
		return []string{}
	}
	arr, ok := v.([]any)
	if !ok {
		return []string{}
	}
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out
}

func cmdMCP(args []string) error {
	fs := flag.NewFlagSet("mcp", flag.ContinueOnError)
	dbPath := fs.String("db", "blueprint/index/blueprint.db", "SQLite path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	conn, err := query.OpenDB(*dbPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	decoder := json.NewDecoder(os.Stdin)
	for {
		var req mcpRequest
		if err := decoder.Decode(&req); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			mcpErr(nil, -32700, "parse error")
			continue
		}

		switch req.Method {
		case "initialize":
			mcpResult(req.ID, map[string]any{
				"serverInfo": map[string]any{"name": "agent-pack-mcp", "version": "0.1.0"},
				"capabilities": map[string]any{
					"tools":     map[string]any{},
					"resources": map[string]any{},
					"prompts":   map[string]any{},
				},
			})
		case "tools/list":
			mcpResult(req.ID, map[string]any{"tools": []map[string]any{
				{"name": "route_query", "description": "Route task to blueprint node"},
				{"name": "read_node", "description": "Read node summary/body"},
				{"name": "build_context_bundle", "description": "Build minimal context bundle"},
				{"name": "expand_node", "description": "Expand node edges"},
				{"name": "rebuild_index", "description": "Rebuild blueprint sqlite index"},
			}})
		case "resources/list":
			mcpResult(req.ID, map[string]any{"resources": []map[string]any{
				{"uri": "blueprint://node/{id}", "name": "blueprint node"},
				{"uri": "blueprint://children/{id}", "name": "blueprint children"},
				{"uri": "blueprint://required/{id}", "name": "blueprint required"},
				{"uri": "blueprint://bundle/{bundle_id}", "name": "blueprint bundle"},
			}})
		case "prompts/list":
			mcpResult(req.ID, map[string]any{"prompts": []map[string]any{
				{"name": "route-task"},
				{"name": "expand-subdomain"},
				{"name": "debug-validator-failure"},
			}})
		case "tools/call":
			name := parseParamString(req.Params, "name")
			obj := map[string]any{}
			_ = json.Unmarshal(req.Params, &obj)
			argsRaw, _ := json.Marshal(obj["arguments"])

			switch name {
			case "route_query":
				level := parseParamString(argsRaw, "level")
				task := parseParamString(argsRaw, "task")
				targetPack := parseParamString(argsRaw, "target_pack")
				targetDomain := parseParamString(argsRaw, "target_domain")
				selectedFiles := parseParamStringList(argsRaw, "selected_files")
				configFragments := parseParamStringList(argsRaw, "config_fragments")
				contextHints := parseParamStringList(argsRaw, "context_hints")
				maxResults := parseParamInt(argsRaw, "max_results", 3)
				var tp *string
				if targetPack != "" {
					tp = &targetPack
				}
				var td *string
				if targetDomain != "" {
					td = &targetDomain
				}
				result, err := query.RouteQuery(conn, level, task, tp, td, selectedFiles, configFragments, contextHints, maxResults)
				if err != nil {
					mcpErr(req.ID, -32000, err.Error())
					continue
				}
				mcpResult(req.ID, map[string]any{"content": []map[string]any{{"type": "text", "text": toJSONString(result)}}})
			case "read_node":
				nodeID := parseParamString(argsRaw, "node_id")
				section := parseParamString(argsRaw, "section")
				if section == "" {
					section = "summary"
				}
				result, err := query.ReadNode(conn, nodeID, section)
				if err != nil {
					mcpErr(req.ID, -32000, err.Error())
					continue
				}
				mcpResult(req.ID, map[string]any{"content": []map[string]any{{"type": "text", "text": toJSONString(result)}}})
			case "build_context_bundle":
				nodeID := parseParamString(argsRaw, "main_node")
				includeRequired := parseParamBool(argsRaw, "include_required", true)
				includeMayInclude := parseParamBool(argsRaw, "include_may_include", false)
				includeChildren := parseParamBool(argsRaw, "include_children", false)
				result, err := query.BuildContextBundle(conn, nodeID, includeRequired, includeMayInclude, includeChildren)
				if err != nil {
					mcpErr(req.ID, -32000, err.Error())
					continue
				}
				mcpResult(req.ID, map[string]any{"content": []map[string]any{{"type": "text", "text": toJSONString(result)}}})
			case "expand_node":
				nodeID := parseParamString(argsRaw, "node_id")
				edgeType := parseParamString(argsRaw, "edge_type")
				if edgeType == "" {
					edgeType = "child"
				}
				result, err := query.ExpandNode(conn, nodeID, edgeType)
				if err != nil {
					mcpErr(req.ID, -32000, err.Error())
					continue
				}
				mcpResult(req.ID, map[string]any{"content": []map[string]any{{"type": "text", "text": toJSONString(result)}}})
			case "rebuild_index":
				root := parseParamString(argsRaw, "root")
				reportDir := parseParamString(argsRaw, "report_dir")
				if root == "" {
					root = "blueprint"
				}
				if reportDir == "" {
					reportDir = "blueprint/index"
				}
				errs, err := compiler.Compile(root, *dbPath, reportDir)
				if err != nil {
					mcpErr(req.ID, -32000, err.Error())
					continue
				}
				mcpResult(req.ID, map[string]any{"content": []map[string]any{{"type": "text", "text": toJSONString(map[string]any{"errors": errs})}}})
			default:
				mcpErr(req.ID, -32601, "unknown tool")
			}
		default:
			mcpErr(req.ID, -32601, "method not found")
		}
	}
}

func toJSONString(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func usage() {
	fmt.Println("agent-pack-mcp commands:")
	fmt.Println("  compile")
	fmt.Println("  route_query")
	fmt.Println("  read_node")
	fmt.Println("  build_context_bundle")
	fmt.Println("  expand_node")
	fmt.Println("  activate")
	fmt.Println("  mcp")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	command := os.Args[1]
	args := os.Args[2:]

	var err error
	switch command {
	case "compile":
		err = cmdCompile(args)
	case "route_query":
		err = cmdRouteQuery(args)
	case "read_node":
		err = cmdReadNode(args)
	case "build_context_bundle":
		err = cmdBuildBundle(args)
	case "expand_node":
		err = cmdExpandNode(args)
	case "activate":
		err = cmdActivate(args)
	case "mcp":
		err = cmdMCP(args)
	default:
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
