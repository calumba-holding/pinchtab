# Snapshot

Get an accessibility snapshot of the current page, including element refs that can be reused by action commands.

Iframe content is detected automatically during snapshot capture. Same-origin iframe descendants are included beneath the iframe owner element, and their refs can be reused directly with action commands. Cross-origin iframes currently remain as owner nodes only.

Selector scoping is explicit. `selector=...` only searches the current frame scope, which defaults to `main`. To scope selector-based snapshots into an iframe, set the frame first with [`/frame`](./frame.md) or `pinchtab frame`.

```bash
curl "http://localhost:9867/snapshot?filter=interactive"
# CLI Alternative (defaults to compact text output)
pinchtab snap -i
# Output
[e5] link "More information..."

# Use --full or --compact=false for JSON
pinchtab snap --full
```

## CLI Flags

| Flag | Description |
|------|-------------|
| `-i`, `--interactive` | Filter to interactive elements + headings (default: true) |
| `-c`, `--compact` | Compact text output (default: true) |
| `-d`, `--diff` | Show diff from previous snapshot |
| `--full` | Full JSON output (shorthand for `--interactive=false --compact=false`) |
| `--text` | Text output format |
| `-s`, `--selector` | CSS selector to scope snapshot |
| `--max-tokens` | Maximum token budget |
| `--depth` | Tree depth limit |
| `--tab` | Target specific tab |

## Examples

```bash
pinchtab snap                           # Interactive compact (default)
pinchtab snap -i -c                     # Same as above
pinchtab snap --full                    # Full JSON with all nodes
pinchtab snap -d                        # Show changes since last snapshot
pinchtab snap --selector "#main"        # Scope to element
pinchtab snap --max-tokens 2000         # Limit output size
```

## API Parameters

| Parameter | Description |
|-----------|-------------|
| `filter` | `interactive` for interactive + headings |
| `format` | `compact`, `text`, `yaml`, or default JSON |
| `diff` | `true` for diff mode |
| `selector` | CSS selector to scope |
| `maxTokens` | Token budget limit |
| `depth` | Tree depth limit |

## Related Pages

- [Click](./click.md)
- [Frame](./frame.md)
- [Tabs](./tabs.md)
