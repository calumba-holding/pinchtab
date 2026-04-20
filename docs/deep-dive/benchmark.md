# Benchmark: Deep Dive

This document compares two agent-driven browser surfaces — **PinchTab** and
**agent-browser** — on the same 10-step benchmark task set, driven by
`claude-haiku-4-5-20251001` through our API runner. The goal is to measure the
real, end-to-end token cost of using each tool with an LLM agent, not an
abstract "tool cost".

The numbers below come from ten fresh runs recorded on 2026-04-20:
`lp1`–`lp5` (PinchTab lane) and `la1`–`la5` (agent-browser lane). Raw logs
are attached at the bottom of this document.

## How the Test is Conducted

Both lanes complete the **same 10 steps** — 5 steps in Group 0 (navigate,
snapshot, read text, click a link, extract specific data) plus 5 steps in
Group 1 (fill a text input, click a button, use search functionality, handle
dynamic content, navigate back). Both lanes read the **shared task file** no
lane-specific task overrides.

### Environment

Every run is executed inside Docker Compose (`tests/tools/docker/compose.yml`)
with three services:

- `fixtures` — the benchmark web server hosting the test pages
  (`/`, `/wiki.html`, `/articles.html`, `/search.html`, `/form.html`,
  `/dashboard.html`, `/ecommerce.html`, `/spa.html`, `/login.html`, etc.).
- `pinchtab` or `agent-browser` — the browser surface being measured.
  PinchTab is built from `tests/tools/config/pinchtab-benchmark.json`
  (IDPI `wrapContent=false`, to match agent-browser's unwrapped output).
- `runner` — the Go program at `tests/tools/runner/` that drives the LLM
  agent loop.

### Agent loop

For each step:

1. The runner sends the task description plus the accumulated transcript to
   Anthropic.
2. The model emits a shell command (one of `./scripts/pt ...`, `./scripts/ab
   ...`, or `./scripts/runner step-end ...`).
3. The runner executes the command inside the lane's container and appends
   the result to the transcript.
4. Once the model answers, `./scripts/runner step-end` records the answer and
   verifies it against the oracle.
5. The loop continues until all steps are answered or `--max-turns` is hit.

Token usage is read directly off Anthropic's `usage` object per response and
summed across the run — no self-reporting by the model.

### What the skill does

Both lanes ship a "skill" — a markdown pack that teaches the model how to
drive the tool: command shapes, ref syntax, snapshot format, known gotchas,
recovery patterns. The skill sits in the **cached prefix** of every request
and is charged at `cache_read` rates (10% of uncached input on Anthropic),
so it costs a predictable, small amount per turn regardless of length.

In a short 10-step run the skill does three things that matter: kills
exploration turns (without the skill the agent has to poke at `--help`,
print `README.md`, try a wrong flag and recover), anchors ref syntax (what
`e5`, `@e5`, or a `[~]` marker means), and pre-empts known gotchas (the
navigation 409 on click, the `--snap-diff` format on PinchTab; `@ref` prefix
and session semantics on agent-browser).

The skill does not do the work; it lets the agent do the work in fewer
turns. Token savings come almost entirely from turn-count reduction.

### Why we ran with a *partial* skill

Both lanes run against a **trimmed subset** of the full shipped skill, not
the complete skill pack the product would load in daily use. This is
deliberate, to keep the comparison fair:

- **PinchTab full skill** = `skills/pinchtab/SKILL.md` (~14.5 KB) + six
  reference files under `skills/pinchtab/references/` (api.md, commands.md,
  env.md, mcp.md, profiles.md, agent-optimization.md) totalling ~44 KB. Full
  size ≈ **58.5 KB**.
- **PinchTab in the benchmark** = just `SKILL.md` (~14.5 KB). The reference
  subfolder is not inlined.
- **agent-browser full skill** comes from `agent-browwser skills get
  agent-browser --full`. In the benchmark the runner extracts only the
  header plus `references/commands.md` and `references/snapshot-refs.md`
  and drops the rest (see
  `tests/tools/runner/prompt.go:DownloadAgentBrowserSkill`).

In a 10-step test, the full reference bundles would dominate `cache_read`
tokens on the lane that happens to ship more reference content. Charging
both lanes for content the agent never consults in a 10-step run would
amplify whichever lane writes more documentation, not whichever tool is
more efficient at the actual task. Trimming both to "header + the one
reference file the agent actually reaches for" isolates the tool-surface
comparison.

If you want to measure *cost in daily production use* rather than *cost of
the tool surface itself*, re-run with `--full` skills on both lanes.

## Results: 2026-04-20 runs (n=5 each)

All ten runs scored 10/10 passes on the same 10-step set. Anthropic
`claude-haiku-4-5-20251001`, `--max-turns 120`.

### Raw per-run totals

| Run | Requests | Uncached input | Cache create | Cache read | Output | Total tokens |
|-----|---------:|---------------:|-------------:|-----------:|-------:|-------------:|
| lp1 | 33 | 68,028 | 6,232 | 199,424 | 4,008 | 277,692 |
| lp2 | 27 | 47,287 | 6,232 | 162,032 | 3,385 | 218,936 |
| lp3 | 31 | 55,872 | 6,232 | 186,960 | 3,991 | 253,055 |
| lp4 | 28 | 57,958 | 6,232 | 168,264 | 3,659 | 236,113 |
| lp5 | 32 | 56,603 | 6,232 | 193,192 | 4,217 | 260,244 |
| **lp avg** | **30.2** | **57,150** | **6,232** | **181,974** | **3,852** | **249,208** |
| la1 | 38 | 57,162 | 6,047 | 223,739 | 4,358 | 291,306 |
| la2 | 36 | 52,276 | 6,047 | 217,692 | 4,308 | 280,323 |
| la3 | 46 | 74,104 | 6,047 | 272,115 | 5,305 | 357,571 |
| la4 | 38 | 57,352 | 6,047 | 229,786 | 4,372 | 297,557 |
| la5 | 38 | 57,188 | 6,047 | 223,739 | 4,337 | 291,311 |
| **la avg** | **39.2** | **59,616** | **6,047** | **233,414** | **4,536** | **303,614** |

### Cost (Anthropic Haiku 4.5, 2026-04-20 pricing)

Pricing per 1M tokens: `$1.00` uncached input, `$1.25` cache-create (1.25×),
`$0.10` cache-read (0.1×), `$5.00` output (5×).

| Run | Cost |
|-----|-----:|
| lp1 | $0.1158 |
| lp2 | $0.0882 |
| lp3 | $0.1023 |
| lp4 | $0.1009 |
| lp5 | $0.1048 |
| **lp avg** | **$0.1024** |
| la1 | $0.1089 |
| la2 | $0.1031 |
| la3 | $0.1354 |
| la4 | $0.1097 |
| la5 | $0.1088 |
| **la avg** | **$0.1132** |

### Averages (the comparison)

| Lane | Avg requests | Avg uncached input | Avg cache create | Avg cache read | Avg output | Avg total tokens | Avg cost |
|------|-------------:|-------------------:|-----------------:|---------------:|-----------:|-----------------:|---------:|
| PinchTab (lp1–lp5) | **30.2** | 57,150 | 6,232 | 181,974 | 3,852 | **249,208** | **$0.1024** |
| agent-browser (la1–la5) | **39.2** | 59,616 | 6,047 | 233,414 | 4,536 | **303,614** | **$0.1132** |
| Δ (la − lp) | +9.0 | +2,466 | −185 | **+51,440** | +684 | +54,406 | +$0.0108 |
| Δ (% vs lp) | +30% | +4% | −3% | **+28%** | +18% | +22% | **+10.5%** |

**Per-run averages:**

- PinchTab is **~10% cheaper** per run on cost ($0.1024 vs $0.1132).
- PinchTab uses **~30% fewer API requests** (30.2 vs 39.2).
- PinchTab's total tokens are **~22% lower** on average, driven almost
  entirely by **cache-read** (+51k on la) — the click→snapshot ping-pong
  inflates re-reads of the cached prefix.

**Per-step** (10 steps per run): PinchTab $0.01024/step, agent-browser
$0.01132/step.

The token gap (~22%) is bigger than the dollar gap (~10%) because most of
agent-browser's extra tokens are `cache_read` at $0.10/1M — cheap per token
but many of them due to the click→snapshot pattern.

### Variance

- **PinchTab** spread: $0.0882 → $0.1158 (~$0.028, 31% of mean).
- **agent-browser** spread: $0.1031 → $0.1354 (~$0.032, 29% of mean),
  dominated by run 3 (46 requests, $0.1354). Without la3 the la mean drops
  to $0.1076 and the gap narrows to ~5%.

Variance is still high enough that single-run conclusions should be
treated with caution; n=5 gives a usable central tendency but the
confidence interval is wide.

## Extended Scope: 2026-04-20 runs (n=3 each, 24 steps)

To check whether the ~10% basic-scope gap holds up at longer workloads, we
re-ran both lanes against **6 groups (0, 1, 2, 3, 4, 5)** — 24 steps per
run — with `--max-turns 250`. Three runs per lane, same model
(`claude-haiku-4-5-20251001`). Logs prefixed `lae*` (agent-browser) and
`lpe*` (PinchTab).

### Raw per-run totals

| Run | Requests | Uncached input | Cache create | Cache read | Output | Total tokens | Answered | Passed | Pass rate |
|-----|---------:|---------------:|-------------:|-----------:|-------:|-------------:|---------:|-------:|----------:|
| lpe1 | 97 | 225,185 | 7,074 | 679,104 | 12,574 | 923,937 | 24 | 23 | 95.8% |
| lpe2 | 89 | 200,720 | 7,074 | 622,512 | 11,431 | 841,737 | 24 | 24 | 100% |
| lpe3 | 91 | 224,544 | 7,074 | 636,660 | 12,780 | 881,058 | 24 | 24 | 100% |
| **lpe avg** | **92.3** | **216,816** | **7,074** | **646,092** | **12,262** | **882,244** | **24.0** | **23.7** | **98.6%** |
| lae1 | 130 | 240,784 | 6,873 | 886,617 | 15,181 | 1,149,455 | 23 | 22 | 95.6% |
| lae2 | 112 | 231,525 | 6,873 | 762,903 | 13,605 | 1,014,906 | 24 | 24 | 100% |
| lae3 | 160 | 303,182 | 6,873 | 1,092,807 | 18,427 | 1,421,289 | 24 | 24 | 100% |
| **lae avg** | **134.0** | **258,497** | **6,873** | **914,109** | **15,738** | **1,195,217** | **23.7** | **23.3** | **98.5%** |

### Cost

| Run | Cost |
|-----|-----:|
| lpe1 | $0.3648 |
| lpe2 | $0.3290 |
| lpe3 | $0.3610 |
| **lpe avg** | **$0.3516** |
| lae1 | $0.4139 |
| lae2 | $0.3844 |
| lae3 | $0.5132 |
| **lae avg** | **$0.4372** |

### Averages (the comparison)

| Lane | Avg requests | Avg uncached input | Avg cache create | Avg cache read | Avg output | Avg total tokens | Avg cost |
|------|-------------:|-------------------:|-----------------:|---------------:|-----------:|-----------------:|---------:|
| PinchTab (lpe1–lpe3) | **92.3** | 216,816 | 7,074 | 646,092 | 12,262 | **882,244** | **$0.3516** |
| agent-browser (lae1–lae3) | **134.0** | 258,497 | 6,873 | 914,109 | 15,738 | **1,195,217** | **$0.4372** |
| Δ (lae − lpe) | +41.7 | +41,681 | −201 | **+268,017** | +3,476 | +312,973 | +$0.0856 |
| Δ (% vs lpe) | +45% | +19% | −3% | **+41%** | +28% | +35% | **+24.4%** |

**Per-run averages:**

- PinchTab is **~24% cheaper** per run on cost ($0.3516 vs $0.4372).
- PinchTab uses **~45% fewer API requests** (92.3 vs 134.0).
- PinchTab's total tokens are **~35% lower** on average, again dominated
  by **cache-read** (+268k on lae) as the click→snapshot pattern compounds
  with more steps.

**Per-step** (24 steps per run): PinchTab $0.01465/step, agent-browser
$0.01822/step. Δ = $0.00357/step (+24% on la).

### Reliability

Both lanes passed 72 total verifications except for one step each
(lpe1: 23/24, lae1: 22/23) — identical error rate of ~1 missed step per
24. Both remain in the mid-90s% range.

### How the gap scales with workload

| Scope | Cost/step (pt) | Cost/step (la) | Gap |
|-------|---------------:|---------------:|----:|
| Basic (10 steps, n=5) | $0.01024 | $0.01132 | +10.5% |
| Extended (24 steps, n=3) | $0.01465 | $0.01822 | +24.4% |

Two things move:

1. **The gap widens at longer scope.** On 10 steps agent-browser is ~10%
   more expensive per run; on 24 steps it's ~24% more expensive. The
   click→snapshot ping-pong pattern compounds — every additional step
   that involves a post-action snapshot costs agent-browser an extra
   round trip, while PinchTab packs it into one via `--snap-diff`.
2. **Per-step cost rises for both lanes** as groups 2–5 are added
   (pt $0.0102 → $0.0147/step, la $0.0113 → $0.0182/step). The later
   groups contain structurally harder steps — dashboards with dynamic
   state, e-commerce flows, SPA interactions — that take more turns on
   either surface.

**agent-browser variance is larger at extended scope.** lae3 (160
requests, $0.5132) is a clear outlier; without it la mean drops to
$0.3992 and the gap narrows to ~14%. PinchTab's lpe runs stay in a
tight $0.329–$0.365 band.

### Takeaway

The extended-scope runs reinforce the basic-scope story but at a larger
magnitude: at production-realistic workloads the PinchTab lane is
meaningfully cheaper (roughly one-quarter less) and uses substantially
fewer API round trips (roughly 45% fewer). The 10% number from the basic
suite understates the gap at scale.

For headline citations: use the basic number as the minimum-noise
apples-to-apples anchor (tight groups, high replicate count), and cite
the extended number when talking about cost at realistic workload size.

## Fairness Caveats

### 1. Task-suite bias

The 10-step task set was designed alongside PinchTab's development. Steps
that are awkward or multi-call in agent-browser — those needing an
explicit post-action snapshot after click / fill / submit / back /
dynamic-content interactions — make up most of Group 1 and all of
Groups 2–5.

A stronger future comparison would co-design the task set with both teams,
or run a much larger task set (the full 39-group benchmark lives in
`tests/benchmark/`) so idiosyncratic per-task biases average out.

### 2. Partial-skill configuration

As described above, both lanes run against a matched *subset* of their
full skills. This is fair for measuring the tool surface, but it's not
representative of production use. A different skill configuration may
have helped agent-browser.

## Methodology Notes

### How token usage is measured

The runner reads usage directly off Anthropic's `usage` object per
response and sums across the run. Fields captured per response:

- `input_tokens` — fresh uncached input
- `cache_creation_input_tokens` — prompt-cache write
- `cache_read_input_tokens` — prompt-cache hit
- `output_tokens` — output
- `request_count` — API requests for the run

### What the token total represents

This is the cost of the **entire agent loop** — system prompt, skill,
tool calls, tool outputs, reasoning, and retries — not a pure "browser
tool only" number. It does not include Docker CPU time or local shell
execution cost.

### Context compaction

To prevent the runner from resending the entire conversation on every turn
(the naive default), we:

1. Truncate tool output before feeding it back to the model (2 preview
   lines per call).
2. Compact old history into a short progress summary derived from the
   benchmark report.
3. Inline the lane setup + skill into the cached prefix so the agent
   doesn't spend uncached turns running `cat setup-*.md`.

### The step-end collapse

Each completed step used to cost two bookkeeping turns: one to record the
answer, one to verify it against the oracle. Those were collapsed into a
single `./scripts/runner step-end` invocation. The n=5 runs above confirm
10/10 adoption on both lanes; this alone saves ~10 turns per PinchTab run
and ~13 per agent-browser run vs earlier baselines.

## Reproducing This Benchmark

```bash
cd tests/benchmark

# 1. Start Docker environment
docker compose up -d --build

# 2. Initialise reports
./scripts/run-optimization.sh

# 3. Baseline lane (deterministic, ~30 seconds)
./scripts/baseline.sh

# 4. PinchTab lane (requires Anthropic key)
ANTHROPIC_API_KEY=... ./scripts/run-api-benchmark.ts --lane pinchtab --groups 0,1

# 5. agent-browser lane (requires Anthropic key)
ANTHROPIC_API_KEY=... ./scripts/run-api-benchmark.ts --lane agent-browser --groups 0,1

# 6. Inspect run-level usage
jq '.run_usage' results/pinchtab_benchmark_*.json
jq '.run_usage' results/agent_browser_benchmark_*.json
```

## Report Files

Results are written to `tests/benchmark/results/`:

| File Pattern | Contents |
|--------------|----------|
| `baseline_YYYYMMDD_HHMMSS.json` | Baseline run |
| `pinchtab_benchmark_YYYYMMDD_HHMMSS.json` | PinchTab results |
| `agent_browser_benchmark_YYYYMMDD_HHMMSS.json` | agent-browser results |
| `agent_browser_commands.ndjson` | Tool-call log for agent-browser |

## Attached Raw Logs

The transcripts behind the sixteen runs in this document.

### Basic scope (10 steps, groups 0+1, n=5 per lane)

- [lp1.txt](./logs/lp1.txt) — PinchTab run 1
- [lp2.txt](./logs/lp2.txt) — PinchTab run 2
- [lp3.txt](./logs/lp3.txt) — PinchTab run 3
- [lp4.txt](./logs/lp4.txt) — PinchTab run 4
- [lp5.txt](./logs/lp5.txt) — PinchTab run 5
- [la1.txt](./logs/la1.txt) — agent-browser run 1
- [la2.txt](./logs/la2.txt) — agent-browser run 2 *(cache-create patched)*
- [la3.txt](./logs/la3.txt) — agent-browser run 3
- [la4.txt](./logs/la4.txt) — agent-browser run 4 *(cache-create patched)*
- [la5.txt](./logs/la5.txt) — agent-browser run 5

### Extended scope (24 steps, groups 0–5, n=3 per lane)

- [lpe1.txt](./logs/lpe1.txt) — PinchTab extended run 1
- [lpe2.txt](./logs/lpe2.txt) — PinchTab extended run 2
- [lpe3.txt](./logs/lpe3.txt) — PinchTab extended run 3
- [lae1.txt](./logs/lae1.txt) — agent-browser extended run 1
- [lae2.txt](./logs/lae2.txt) — agent-browser extended run 2
- [lae3.txt](./logs/lae3.txt) — agent-browser extended run 3

Each log contains the full agent conversation, every tool call with
arguments, timing, and the `[run-usage]` line at the bottom.
Machine-specific paths have been replaced with `<repo>` for portability.

## Limitations

- Task-suite bias (see Fairness Caveats §1)
- Partial-skill configuration (Fairness Caveats §2)
- Basic scope: n=5 per lane; extended scope: n=3 per lane. Run-to-run
  variance remains ~25–30% of mean; agent-browser has one outlier in
  each scope
- Single model (Haiku 4.5); no Sonnet/Opus comparison
- Fixed Docker environment adds per-call overhead roughly equal across
  lanes, but absolute times are not production-representative
- Score is pass-count, not answer quality or time-to-complete

## Future Work

- Fix the two measurement bugs (pinchtab tool-calls, cache-create drop) in
  `tests/tools/runner/recordstep.go`
- 10+ runs per lane at extended scope to tighten the confidence interval
  and characterise the agent-browser outlier rate
- Model comparison (Haiku vs Sonnet vs Opus)
- Full-skill (not partial-skill) re-run for a production-realistic
  comparison
- Co-designed task set with both tool teams to reduce task-suite bias
- Per-step token tracking (currently run-level only)
- Retry rates and error-recovery patterns
