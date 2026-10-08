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

DJ MorseCode is a terminal-based listening companion that turns intent into continuous music while keeping discovery at the center of the experience.

Internet radio remains the heart of the DJ, but playback is no longer architecturally tied to radio. DJ MorseCode can also play curated source media through the same playback, enrichment, lyric, observation, and history pipelines.

Instead of building playlists, choose how you want to listen:

- **Vibes** describe the session: Focus, Discovery, Energy, or Wind Down.
- **Moods** describe how the music should feel.
- **Genres** describe the musical neighborhood.
- **Stations** let you take direct control when you already know what you want.
- **Source media** provides explicitly curated on-demand listening when appropriate.

DJ MorseCode selects appropriate internet radio stations, remembers where it has been, introduces variety, rotates stations during programmed sessions, and recovers from failed streams without abandoning the listening intent.

mpv handles audio playback while DJ MorseCode handles programming, source identity, enrichment, observations, and listening history.

## Quick Start

### Requirements

DJ MorseCode currently requires:

- Go 1.26.5
- mpv
- yt-dlp for YouTube playback
- A terminal with Unicode and color support
- Internet access for streams, source media, enrichment, and remote lyric retrieval

The currently tested mpv version is 0.41.0.

On macOS with Homebrew:

```bash
brew install mpv yt-dlp
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

DJ MorseCode uses contextual controls. Available commands depend on the current control mode.

Core interactions include:

- choosing a station
- choosing a vibe
- choosing a mood
- choosing a genre
- choosing curated source media
- moving backward and forward through station navigation history
- viewing chronological listening history
- viewing raw observation history
- pausing and resuming playback
- muting and unmuting
- changing volume
- showing or hiding lyrics
- signing off

Pickers support arrow keys and `j` / `k` navigation, `Enter` to select, and `Esc` to return.

The terminal footer displays the controls available in the current context.

## Listening Modes

### Vibes

Vibes describe what you're doing rather than prescribing a genre.

Current presets include:

- Focus
- Discovery
- Energy
- Wind Down

A vibe becomes the active session intent. DJ MorseCode continues programming within that intent until you select something else or manually take control.

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

### Source Media

DJ MorseCode can play finite network media without treating it as live radio.

The current implementation includes a manually curated YouTube catalog. YouTube playback is audio-only and participates in the same track enrichment, lyric lookup, observation, and listening-history pipelines used elsewhere in the application.

The source architecture distinguishes media identity and provenance from playback transport, providing a foundation for additional catalogs and providers without coupling the terminal UI to a specific service.

## Autonomous Sessions

Vibe, mood, and genre selections create an active session intent.

During an active session, DJ MorseCode:

1. Finds stations matching the requested criteria.
2. Avoids the currently playing and recently tuned stations when possible.
3. Introduces controlled variety among eligible stations.
4. Remembers successful station transitions.
5. Rotates to another appropriate station after the configured dwell interval.
6. Preserves the original listening intent across rotations.

Manual next behavior uses the same selection machinery as automatic rotation.

### Stream Recovery

Internet radio is messy.

Streams disappear, stall, redirect, buffer, and occasionally just decide today is not their day.

DJ MorseCode treats a requested tune as pending until playback is observed. If the stream remains idle beyond the grace period, that station is excluded for the current programmed session and DJ MorseCode selects another station satisfying the same intent.

Resolved or redirected stream URLs do not erase the identity of a station DJ MorseCode explicitly requested.

Failed tune attempts are not added to listening history.

## Listening History

DJ MorseCode maintains chronological listening history for the current session.

Listening history is derived from successful playback evidence rather than navigation state:

- confirmed radio station tunes become station entries
- explicitly played source media become media entries
- revisits are preserved chronologically
- failed tune attempts are excluded

Radio station back/forward navigation maintains its own browser-style history and cursor. Navigating through previous stations does not redefine listening history.

Raw station, playback, and source-media observations remain available separately for diagnostics.

Track-level history from songs encountered within live radio streams is intentionally deferred until meaningful-listen semantics are defined.

## Lyrics

DJ MorseCode supports plain and synchronized lyrics when lyrics are available.

The current lyric path includes:

- locally cached lyrics
- LRCLIB lookup
- persistent LRCLIB results in the DJ MorseCode cache
- synchronized timeline cues when available

Lyrics are enrichment, not a playback requirement.

Missing lyrics never prevent the music from continuing.

## Metadata and Enrichment

Playback identity can be enriched independently of the playback source.

Current enrichment and context sources include:

- MusicBrainz recording identity and release metadata
- Last.fm contextual tags
- LRCLIB lyrics

Provider results are translated into application-domain models before they reach the UI.

Enrichment is asynchronous and must never block playback.

## Media Sources

DJ MorseCode separates media identity from playback transport.

Current source capabilities include:

- built-in radio stations
- M3U playlist parsing and station proposals
- Last.fm similar-track discovery
- curated YouTube video playback

A source may provide catalog entries, discovery results, provenance, or playable media identity without becoming the playback engine itself.

mpv remains the playback engine.

## Station Intelligence

Stations carry semantic metadata used by the selection engine:

- Genre
- Tags
- Moods
- Contexts
- Energy

Matching determines which stations are appropriate for an intent.

Selection then considers navigation history, recent stations, and failed stations before choosing among eligible candidates.

This keeps the decision process deterministic where rules matter and varied where multiple answers are equally valid.

## Architecture

At a high level:

```text
                    User Intent / Selection
                             |
                 +-----------+-----------+
                 |                       |
                 v                       v
          Radio Programming        Source Catalogs
                 |                       |
                 v                       v
          Station Selection         Media Identity
                 |                       |
                 +-----------+-----------+
                             |
                             v
                           Player
                             |
                             v
                            mpv
                             |
                 +-----------+-----------+
                 |                       |
                 v                       v
              Playback              Observations
                 |                       |
                 v              +--------+--------+
             Enrichment          |                 |
                                 v                 v
                           Diagnostics       Listening History
```

DJ MorseCode owns intent, source identity, programming decisions, enrichment, observations, and listening history.

mpv owns media transport and audio playback.

The source layer describes what media is and where it came from without prescribing how it is played.

The observation layer records runtime evidence.

The history layer derives meaningful listening chronology from that evidence.

`radio.History` is a separate browser-style navigation model for station back/forward controls.

Bubble Tea owns the terminal interaction model.

See [ARCHITECTURE.md](/docs/ARCHITECTURE.md) for the detailed architecture.

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

Useful inspection commands also live under `cmd/` for focused MusicBrainz, Last.fm, and LRCLIB diagnostics.

## Project Status

DJ MorseCode currently provides a complete terminal-first listening loop centered on discovery.

The implemented foundation includes:

- managed mpv lifecycle
- audio-only playback
- curated internet radio
- direct station tuning
- vibe-based programming
- mood-based programming
- genre-based programming
- persistent session intent
- history-aware station selection
- browser-style station back/forward navigation
- controlled selection variety
- manual and automatic station rotation
- redirected-stream reconciliation
- failed-stream recovery
- provider-neutral media identity and provenance
- M3U playlist ingestion and station proposals
- Last.fm media discovery
- curated YouTube audio playback
- MusicBrainz and Last.fm enrichment
- runtime station, playback, and media observations
- chronological session listening history
- synchronized lyric support
- persistent lyric caching
- tmux-friendly terminal UI

Future directions live in [ROADMAP.md](ROADMAP.md).

## License

DJ MorseCode is available under the [MIT License](LICENSE).