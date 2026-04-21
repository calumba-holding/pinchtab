# Navigate

Open a new tab and navigate it to a URL, or reuse a tab when a tab ID is provided.

```bash
curl -X POST http://localhost:9867/navigate \
  -H "Content-Type: application/json" \
  -d '{"url":"https://pinchtab.com"}'
# CLI Alternative
pinchtab nav https://pinchtab.com
# Response (default is tab ID; use --json for full JSON)
8f9c7d4e1234567890abcdef12345678
```

## CLI Flags

| Flag | Description |
|------|-------------|
| `--tab` | Reuse existing tab |
| `--new-tab` | Force new tab |
| `--block-images` | Block image loading |
| `--block-ads` | Block ads |
| `--snap` | Output snapshot after navigation |
| `--snap-diff` | Output snapshot diff after navigation |
| `--print-tab-id` | Print only tab ID (auto when piped) |
| `--json` | Full JSON response |

## Examples

```bash
pinchtab nav https://example.com              # Navigate, print tab ID
pinchtab nav https://example.com --snap       # Navigate and snapshot
TAB=$(pinchtab nav https://example.com)       # Capture tab ID for reuse
pinchtab nav https://other.com --tab "$TAB"   # Reuse tab
pinchtab nav https://example.com --block-images  # Skip images
```

## API Body Fields

| Field | Description |
|-------|-------------|
| `url` | Target URL (required) |
| `tabId` | Reuse existing tab |
| `newTab` | Force new tab |
| `blockImages` | Block image loading |
| `blockAds` | Block ads |
| `timeout` | Navigation timeout |
| `waitFor` | Wait condition |
| `waitSelector` | Wait for selector |

## Related Pages

- [Snapshot](./snapshot.md)
- [Tabs](./tabs.md)
