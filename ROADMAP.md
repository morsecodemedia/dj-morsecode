# DJ MorseCode Roadmap

DJ MorseCode optimizes for discovery.

The roadmap is organized around product directions rather than promised release numbers. Features move into a release when they make the listening experience better without making the DJ demand more attention.

## Current Foundation

DJ MorseCode currently provides:

- live internet-radio playback through mpv
- audio-only finite network-media playback
- direct station tuning
- mood, genre, and vibe-based station discovery
- browser-style station back/forward navigation
- chronological session listening history
- provider-neutral source identity and provenance
- M3U playlist parsing and station proposals
- curated YouTube media playback
- Last.fm similar-track discovery
- playback controls for pause, mute, and volume
- tiered contextual keyboard controls
- normalized provider-specific stream metadata
- redirected-stream station reconciliation
- MusicBrainz identity and release enrichment
- Last.fm contextual tags
- station, playback, and source-media observations
- LRCLIB lyric discovery with local caching
- plain and synchronized lyric content

## Listening History

The current session history records confirmed station tunes and explicitly played source media while preserving revisits chronologically.

Listening history is intentionally separate from browser-style station navigation and from raw observation history.

Future work includes:

- persistence across application sessions
- track-level history derived from radio playback observations
- meaningful-listen thresholds before a track counts as heard
- filtering and grouping
- source and provenance display
- session summaries
- listening intelligence

Track-level history should not be introduced until DJ MorseCode has an explicit definition of what constitutes a meaningful listen.

## Session Intelligence

### Keep Cookin'

Continue the current vibe.

### Surprise Me

Take an unexpected turn without breaking the mood.

### Go Deeper

Move toward less obvious selections, deep cuts, and musical neighborhoods outside the usual path.

### Take Me Sideways

Keep similar energy while moving into a different genre.

### Cool It Down

Reduce the session's intensity.

### Go Harder

Increase the session's intensity.

### Encore

Stay in the current musical neighborhood without simply replaying what just happened.

## Personalization

Potential directions include:

- favorite stations
- disliked stations
- persistent listening history across sessions
- track-level listening history with meaningful-listen semantics
- preference learning
- session memory
- time-of-day behavior
- frequently used vibes and moods

The DJ should learn without turning listening into configuration work.

## Station Health

Persist station tune failures across sessions.

A failed tune attempt should be recorded as an attempt-level event, not counted repeatedly from polling ticks.

Potential lifecycle:

```text
Active
  |
  v
failed tune attempts
  |
  v
warning threshold
  |
  v
flagged for decommission
  |
  v
decommissioned
```

A station that reaches a threshold such as five independent failed tune attempts may be removed from the effective active catalog and retained in a decommissioned catalog for inspection or rehabilitation.

A temporary outage should not permanently condemn a station.

Station health should build on durable observations rather than inventing a second failure-tracking model.

## DJ Notes

Inspired by VH1 Pop-Up Video.

Occasionally surface small pieces of context without interrupting playback.

### Examples:

> This guitar solo was recorded in one take.

> This band originally opened for Soundgarden.

> You haven't heard this station in eight months.

Notes should disappear automatically and never become walls of text.

## Playback Sources

Playback should not be permanently tied to one source.

The source layer now provides provider-neutral media identity, provenance, catalog capabilities, and discovery capabilities.

Current integrations include:

- built-in radio stations
- M3U station catalogs and proposals
- curated YouTube media
- Last.fm track discovery

Possible future playback or catalog adapters include:

- local libraries
- Jellyfin
- Navidrome
- Spotify
- Apple Music
- YouTube Music
- additional radio catalogs

The terminal renderer should care about what is playing, not where it came from.

Source providers should describe media and discovery.

Playback adapters should describe how media is transported.

Those responsibilities should remain independent.

## Lyrics and Enrichment

Lyrics should remain provider-agnostic.

Reliable lyric acquisition and readable lyric content remain the priority.

Current enrichment includes:

- MusicBrainz recording identity
- MusicBrainz release information
- Last.fm contextual tags
- LRCLIB plain and synchronized lyrics

Potential future enhancements include:

- additional lyric providers
- user-managed local LRC discovery
- richer lyric provenance
- transcripts for non-music media
- richer track and artist context
- enrichment confidence and provenance in the UI

Lyrics availability and lyrics visibility should remain separate concepts.

The lyrics panel should default to hidden while still exposing whether lyrics are available.

## Catalogs

Potential station-catalog improvements include:

- user-managed station configuration
- user-managed M3U catalog configuration
- multiple endpoints per station
- preferred codec/quality selection
- endpoint failover
- catalog health tooling
- persistent imported catalogs

## Catalog Intelligence

DJ MorseCode should make adding and maintaining stations substantially easier than manually researching and classifying them.

### Assisted Station Import

Given a stream URL and optional homepage, inspect available source evidence and produce a structured station proposal.

Potential evidence includes:

- stream metadata and ICY headers
- station name and description
- declared genres
- homepage content
- current or sample track metadata
- existing DJ MorseCode catalog taxonomy

An LLM-backed classifier may use this evidence to propose:

- station name
- genre
- tags
- moods
- contexts
- energy
- other catalog metadata

Generated classifications must conform to DJ MorseCode's existing taxonomy rather than inventing arbitrary categories.

Imports should be reviewable before becoming part of the active catalog.

Generated values should retain enough provenance to understand why they were selected.

Future ingestion sources may include:

- individual stream URLs
- station homepages
- M3U playlists
- external radio catalogs
- user-maintained catalog files

## Observability and Diagnostics

Observation History provides raw runtime evidence, but debugging should not require exiting DJ MorseCode or adding temporary logging.

Add a live, read-only playback and metadata inspector under the Info controls.

Useful diagnostic state includes:

- current mpv path
- idle state
- playback position and duration
- raw mpv metadata
- normalized playback metadata
- source media identity and provenance
- current and pending station identity
- resolved station identity
- track identity
- enrichment status and provider
- lyric state and lookup duration

The inspector should update through the normal application tick and must not become another state machine.

## Interface

The idle screen may eventually evolve into a lightweight home view containing useful context such as:

- quick-start intents
- recent listening history
- favorite stations
- current station health

Additional interface work includes:

- lyrics hidden by default
- live playback and metadata inspector
- DJ Notes
- optional Follow Mode
- broader UI cleanup
- history filtering and grouping

It should still feel like another pane in tmux rather than another application.

## Other Media

The architecture may eventually support listening contexts beyond live music radio:

- podcasts
- audiobooks
- spoken-word streams
- local media

Those modes should earn their way into the product rather than complicating the radio experience prematurely.

## Joy

Every feature should still make someone smile.

## MISC TODO

### TODO: Persistence/session foundation
  → durable observations
  → persist listening history across sessions
  → track-level history / meaningful-listen semantics
  → station health

### TODO: Catalog evolution
  → stations become data rather than compiled Go
  → user-managed M3U catalogs
  → assisted URL/homepage importer
  → taxonomy-constrained LLM classification
  → review/provenance workflow

### TODO: Session Intelligence
  → Keep Cookin'
  → Surprise Me
  → Go Deeper
  → Take Me Sideways
  → Go Harder
  → Cool It Down
  → Encore

### TODO: Experience/UI
  → DJ Notes
  → lyrics panel hidden by default
  → live playback/metadata debug inspector
  → optional Follow Mode
  → broader UI cleanup