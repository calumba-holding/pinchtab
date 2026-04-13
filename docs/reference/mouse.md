# Mouse

Low-level pointer controls for cases where DOM-native click or hover behavior is not enough.

```bash
curl -X POST http://localhost:9867/action \
  -H "Content-Type: application/json" \
  -d '{"kind":"mousemove","ref":"e5"}'

curl -X POST http://localhost:9867/action \
  -H "Content-Type: application/json" \
  -d '{"kind":"mousedown","ref":"e5","button":"left"}'

curl -X POST http://localhost:9867/action \
  -H "Content-Type: application/json" \
  -d '{"kind":"mouseup","ref":"e5","button":"left"}'

curl -X POST http://localhost:9867/action \
  -H "Content-Type: application/json" \
  -d '{"kind":"mousewheel","x":400,"y":320,"wheelDeltaY":240}'

# CLI Alternatives
pinchtab mouse move e5
pinchtab mouse down e5 --button left
pinchtab mouse up e5 --button left
pinchtab mouse wheel --x 400 --y 320 --wheel-delta-y 240
```

Notes:

- mouse actions accept the same targeting modes as other action commands: `ref`, `selector`, `nodeId`, or `x`/`y`
- `mousedown` and `mouseup` accept `button=left|right|middle`
- `mousewheel` accepts `wheelDeltaX` and `wheelDeltaY`
- use these for drag handles, canvas controls, precise hover choreography, or sites that require exact pointer sequencing

## Related Pages

- [Click](./click.md)
- [Hover](./hover.md)
- [Scroll](./scroll.md)
