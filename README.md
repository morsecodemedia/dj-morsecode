# DJ MorseCode

**A terminal DJ for continuous musical discovery.**

There was something magical about never knowing what came next.

Modern streaming services optimize for control.

DJ MorseCode optimizes for discovery.

It should feel like sitting beside a great late-night radio DJ that somehow always understands your mood.

## Philosophy

The goal is not to play exactly the songs you ask for.

The goal is to keep you in the groove.

You shouldn't think about playlists.

You shouldn't think about albums.

You shouldn't even have to think about songs.

You think about a feeling.

DJ MorseCode handles the rest.

### Continuous Discovery

Every song should feel like it belongs.

Every song should also feel like a pleasant surprise.

### Infinite Radio

Music keeps moving without playlist management, queue management, or decision fatigue.

You get what you get.

If you're not feeling it, move on.

If it hits, enjoy the ride.

### Ambient Companion

DJ MorseCode should never demand attention.

It quietly enhances your workspace.

It should feel like another pane in tmux rather than another application.

### Joy

Every feature should make someone smile.

---

## What DJ MorseCode Does

DJ MorseCode is a terminal-based radio companion that turns listening intent into continuous music.

Instead of building playlists, choose how you want to listen:

- **Vibes** describe the session: Focus, Discovery, Energy, or Wind Down.
- **Moods** describe how the music should feel.
- **Genres** describe the musical neighborhood.
- **Stations** let you take direct control when you already know what you want.

DJ MorseCode selects an appropriate internet radio station, remembers where it has been, introduces variety, rotates stations during programmed sessions, and recovers from failed streams without abandoning the listening intent.

MPV handles audio playback while DJ MorseCode handles the programming.

## Quick Start

### Requirements

DJ MorseCode currently requires:

- Go 1.26.5
- mpv
- A terminal with Unicode and color support
- Internet access for radio streams and remote lyric retrieval

The currently tested mpv version is 0.41.0.

On macOS with Homebrew:

```bash
brew install mpv
```

### Run

Clone the repository and start DJ MorseCode:

```bash
git clone git@github.com:morsecodemedia/dj-morsecode.git
cd dj-morsecode
go run ./cmd/dj
```

That's it.

DJ MorseCode starts and manages its own idle mpv process. You do not need to launch mpv separately.

When DJ MorseCode exits normally, it also shuts down the mpv process and removes its IPC socket.

## Controls

| Key | Action |
| --- | --- |
| `s` | Choose a station |
| `v` | Choose a vibe |
| `m` | Choose a mood |
| `g` | Choose a genre |
| `n` | Choose the next station within the active session |
| `h` | View station history |
| `q` | Sign off |
| `Ctrl-C` | Sign off |
| `↑` / `↓` | Navigate a picker |
| `j` / `k` | Navigate a picker |
| `Enter` | Select |
| `Esc` | Close a picker |

`n` only has meaning during an active vibe, mood, or genre session. Choosing a station manually with `s` ends the active programmed session and returns control to you.

## Listening Modes

### Vibes

Vibes describe what you're doing rather than prescribing a genre.

Current presets include:

- **Focus**
- **Discovery**
- **Energy**
- **Wind Down**

A vibe becomes the active session intent. DJ MorseCode continues programming within that intent until you select something else or manually tune a station.

### Moods

Moods describe how the session should feel.

DJ MorseCode groups its richer station metadata into human-friendly mood families such as:

- Calm
- Dreamy
- Energetic
- Uplifting
- Adventurous
- Edgy
- Nostalgic
- Playful
- Thoughtful

### Genres

Genres are derived dynamically from the curated station catalog.

Adding a station with a new genre automatically makes that genre available to the genre picker.

### Stations

The station picker provides direct access to the curated radio catalog.

Manual station selection clears any active vibe, mood, or genre intent. From that point DJ MorseCode stays on the station you chose until you make another selection.

## Autonomous Sessions

Vibe, mood, and genre selections create an active session intent.

During an active session, DJ MorseCode:

1. Finds stations matching the requested criteria.
2. Avoids the currently playing and recently tuned stations when possible.
3. Introduces controlled variety among eligible stations.
4. Remembers successful station transitions.
5. Rotates to another appropriate station after the configured dwell interval.
6. Preserves the original listening intent across rotations.

Press `n` to trigger the same rotation behavior manually.

### Stream Recovery

Internet radio is messy.

Streams disappear, stall, redirect, buffer, and occasionally just decide today is not their day.

DJ MorseCode gives an intent-driven station time to establish playback. If the stream remains idle beyond the grace period, that station is excluded for the current session and DJ MorseCode selects another station satisfying the same intent.

Failed tune attempts are not added to listening history.

## Lyrics

DJ MorseCode supports synchronized lyric timelines when lyrics are available.

The current lyric path includes:

- locally cached lyrics
- LRCLIB lookup
- persistent LRCLIB results in the DJ MorseCode cache

Lyrics are enrichment, not a playback requirement. Missing lyrics never prevent the music from continuing.

## Station Intelligence

Stations carry semantic metadata used by the selection engine:

- Genre
- Tags
- Moods
- Contexts
- Energy

Matching determines which stations are appropriate for an intent.

Selection then considers session history and failed stations before choosing among eligible candidates.

This keeps the decision process deterministic where rules matter and varied where multiple answers are equally valid.

## Architecture

At a high level:

```text
                     ┌──────────────┐
                     │ User Intent  │
                     │ vibe / mood  │
                     │ genre        │
                     └──────┬───────┘
                            │
                            ▼
                       Criteria
                            │
                            ▼
                     Station Match
                            │
                            ▼
                History / Failure Filters
                            │
                            ▼
                    Candidate Selection
                            │
                            ▼
                         Station
                            │
                            ▼
                      Player.Load()
                            │
                            ▼
                           mpv
                            │
                            ▼
                         Audio
```

DJ MorseCode owns session intent and programming decisions.

mpv owns media transport and audio playback.

Bubble Tea owns the terminal interaction model.

## Development

Format:

```bash
gofmt -w ./cmd ./internal
```

Test:

```bash
go test ./...
```

Vet:

```bash
go vet ./...
```

Release sanity check:

```bash
go test ./...
go vet ./...
git diff --check
```

Run:

```bash
go run ./cmd/dj
```

## Project Status

DJ MorseCode v1 focuses on the smallest complete listening experience:

> Start the application, describe what you want to hear, and let the DJ handle the rest.

The v1 feature set includes:

- managed mpv lifecycle
- curated internet radio
- direct station tuning
- vibe-based programming
- mood-based programming
- genre-based programming
- persistent session intent
- history-aware station selection
- controlled selection variety
- manual and automatic station rotation
- failed-stream recovery
- synchronized lyric support
- persistent lyric caching
- tmux-friendly terminal UI

Future directions live in [ROADMAP.md](ROADMAP.md).

## License

DJ MorseCode is available under the [MIT License](LICENSE).