# DJ MorseCode Architecture

> "Every UI element must be backed by real data."

DJ MorseCode is a terminal-first music companion built around continuous discovery.

The application translates listening intent into radio programming, delegates media playback to mpv, enriches playback with metadata and synchronized lyrics when available, and presents the resulting session through a Bubble Tea terminal interface.

The architecture intentionally separates programming decisions, playback, enrichment, domain data, and presentation so each subsystem has a clear responsibility.

This document describes architectural intent rather than every implementation detail.

When implementation and this document disagree, implementation should be questioned before the document is changed.

---

# Architecture

```text
                              +----------------------+
                              |        cmd/dj        |
                              |  application model   |
                              +----------+-----------+
                                         |
                +------------------------+------------------------+
                |                        |                        |
                v                        v                        v
        +---------------+        +---------------+        +---------------+
        |     radio     |        |    player     |        |      ui       |
        | programming   |        |   playback    |        | presentation  |
        +-------+-------+        +-------+-------+        +---------------+
                |                        |
                |                        v
                |                +---------------+
                |                |  player/mpv   |
                |                | IPC + process |
                |                +-------+-------+
                |                        |
                |                        v
                |                      mpv
                |
                +-----------------------------+
                                              |
                                              v
                                     +----------------+
                                     |   metadata     |
                                     +----------------+

              +-------------+     +-------------+     +-------------+
              |   library   | --> |   lrclib    | --> |   lyrics    |
              | lyric cache |     | remote data |     | LRC parsing |
              +-------------+     +-------------+     +------+------+
                                                               |
                                                               v
                                                        +-------------+
                                                        |    music    |
                                                        | song / cues |
                                                        +-------------+
```

`cmd/dj` coordinates these systems.

It does not own the underlying rules for station matching, lyric parsing, or mpv IPC.

---

# Primary Runtime Flows

## Programming

DJ MorseCode programs radio from user intent.

```text
User Selection
      |
      v
    Intent
      |
      v
   Criteria
      |
      v
    Match
      |
      v
Eligible Stations
      |
      v
History / Failure Filtering
      |
      v
Candidate Selection
      |
      v
   Station
      |
      v
Player.Load()
      |
      v
     mpv
      |
      v
    Audio
```

A user may establish intent through:

- a vibe
- a mood family
- a genre

Direct station selection bypasses programming intent and gives control back to the user.

An active intent survives station changes.

This allows manual `next` behavior and automatic station rotation to continue programming toward the same listening goal.

---

## Playback Reconciliation

A request to load a station is not treated as proof that playback succeeded.

```text
Player.Load()
      |
      v
Pending Tune
      |
      +------ mpv becomes active ------> Confirm Playback
      |                                      |
      |                                      v
      |                                Record History
      |
      +------ grace period expires ----> Failed Attempt
                                             |
                                             v
                                      Exclude Station
                                             |
                                             v
                                       Choose Again
```

mpv remains authoritative for actual playback state.

DJ MorseCode records successful station history only after playback is observed as active.

Failed stations are excluded from the current programmed session so recovery cannot immediately select the same failed endpoint again.

---

## Lyrics and Enrichment

Lyrics are enrichment, not a playback dependency.

```text
Playback Metadata
      |
      v
Track Resolution
      |
      v
Local Lyric Cache
      |
      +------ hit ------> Parsed Timeline
      |
      +------ miss -----> LRCLIB
                              |
                              v
                        Parsed Timeline
                              |
                              v
                         Cache Result
                              |
                              v
                              UI
```

Playback continues when lyrics are unavailable.

Remote lyric retrieval is asynchronous so network work does not block the terminal application.

Responses are associated with track identity so stale lyric results cannot replace the current track's timeline.

---

# Package Responsibilities

## music

The core playback-domain model.

Contains concepts such as:

- Song
- Cue
- CueType

`music` represents parsed musical information without knowing how that information was retrieved or displayed.

The package knows nothing about:

- mpv
- radio selection
- Bubble Tea
- terminal layout
- network retrieval

---

## radio

Owns radio programming and station-domain behavior.

Responsibilities include:

- curated station catalog
- station lookup
- semantic station metadata
- genres
- tags
- moods
- contexts
- energy
- dynamic catalog facets
- vibe presets
- mood families
- listening criteria
- station matching
- session intent
- tune history
- recent-station avoidance
- explicit station exclusion
- candidate selection
- controlled variety
- selection relaxation

The package answers questions such as:

> Which stations satisfy this listening intent?

and:

> Given these eligible stations and this listening history, which station should be selected next?

The package does not load media.

Selection returns a station; the application decides when to ask the player to load it.

---

## player

Provides the playback-facing abstraction used by the application.

Responsibilities include:

- load media
- expose current media metadata
- expose playback position and duration
- identify network playback
- expose playback idle state
- determine the active timeline cue
- format playback information

The application should not need to understand mpv IPC commands or property names.

---

## player/mpv

Owns the mpv integration.

Responsibilities include:

- mpv IPC communication
- property retrieval
- media loading
- managed mpv process startup
- IPC readiness
- managed process shutdown
- IPC socket cleanup

DJ MorseCode launches mpv idle.

```text
DJ MorseCode starts
      |
      v
Start mpv
      |
      v
Wait for IPC readiness
      |
      v
Connect Player
      |
      v
Bubble Tea starts
```

On normal application shutdown, DJ MorseCode terminates the mpv process it owns and removes the IPC socket.

mpv owns media transport and audio playback.

DJ MorseCode owns programming behavior.

---

## metadata

Resolves playback metadata into track identity.

Network streams frequently expose changing ICY/media-title data, while local media may provide structured artist and title tags.

Metadata precedence is intentionally different for local and network playback so static stream metadata cannot overwrite changing radio-track metadata.

---

## lyrics

Parses LRC documents into domain objects.

Responsibilities include:

- parse metadata
- parse timestamps
- parse synchronized lyric cues
- construct lyric timelines

The parser does not perform network retrieval.

It parses documents.

---

## lrclib

Owns LRCLIB integration.

Responsibilities include:

- search remote lyric records
- identify an appropriate result
- translate remote data into the music domain

LRCLIB is optional enrichment.

Failure to retrieve lyrics must never stop playback.

---

## library

Owns persistent local lyric cache behavior.

Remote lyric results may be stored locally and reused on later playback.

The cache is application data, not repository-owned demo media.

The application should not require bundled songs or lyrics to start.

---

## ui

Presentation only.

Responsibilities include:

- idle state
- current playback
- station information
- active session intent
- synchronized lyric timeline
- station picker
- vibe picker
- mood picker
- genre picker
- station history
- controls and status

The UI never parses LRC syntax.

The UI never decides which station matches an intent.

The UI renders application state.

---

## cmd/dj

The application orchestration layer.

Responsibilities include:

- Bubble Tea model and update loop
- keyboard interaction
- picker state
- active session coordination
- pending tune state
- station failure recovery
- automatic rotation timing
- mpv lifecycle coordination
- connecting domain decisions to playback
- connecting playback state to presentation

`cmd/dj` may coordinate packages, but domain rules should remain in the packages that own them.

For example:

```text
cmd/dj decides WHEN to choose another station.

radio decides WHICH stations are valid choices.

player decides HOW to request playback.

mpv decides HOW media is transported and rendered as audio.
```

---

# State Ownership

Clear state ownership prevents UI observations from becoming accidental domain rules.

## mpv owns

- actual media playback
- playback position
- media duration
- current media path
- playback idle state
- stream metadata

## radio owns

- station definitions
- station semantics
- matching rules
- selection rules
- intent representation
- tune-history behavior

## cmd/dj owns

- current application session
- active intent
- pending tune attempt
- failed stations for the active session
- automatic rotation timing
- picker visibility and selection
- orchestration between radio and player

## ui owns

- presentation
- layout
- terminal styles

---

# Session Intent

Vibes, moods, and genres are different user interfaces to the same programming abstraction.

```text
Vibe -----+
          |
Mood -----+----> Intent ----> Criteria
          |
Genre ----+
```

Only one programmed intent is active at a time.

Selecting a new vibe, mood, or genre replaces the previous intent.

Selecting a station manually clears the active intent.

This makes direct station selection an explicit manual takeover.

---

# Station Selection

Station selection is intentionally separated into stages.

## Match

`Match()` determines eligibility from criteria.

Across dimensions, criteria are combined as AND conditions.

Within a dimension, values are alternatives.

Conceptually:

```text
(context = coding OR focus)
AND
(mood = calm OR dreamy)
AND
energy <= 2
```

## History

Tune history records successful station transitions.

Consecutive observations of the same station do not create duplicate tune events.

History enables:

- current station awareness
- previous station awareness
- recent-station avoidance
- future listening analysis

## Select

Selection operates only on eligible candidates.

History filtering occurs before candidate choice.

A candidate chooser may introduce variety among equally valid surviving stations.

Randomness never determines whether a station satisfies the listening criteria.

## Choose

`Choose()` combines matching and selection policy.

It:

1. matches criteria
2. applies explicit exclusions
3. excludes the current station
4. avoids recent stations
5. progressively relaxes recent-history avoidance when necessary
6. chooses among the remaining valid candidates

Current-station exclusion does not relax.

---

# Automatic Programming

During an active intent, stations rotate after a configured dwell interval.

Automatic rotation uses the same operation as the manual `n` command.

```text
manual n -------------------+
                            |
dwell interval expires -----+----> choose next station
                                      |
                                      v
                                 Player.Load()
```

There is one station-selection path rather than separate manual and automatic algorithms.

---

# Stream Failure Recovery

Internet radio endpoints are unreliable by nature.

DJ MorseCode treats an intent-driven load as pending until playback is confirmed.

If mpv remains idle beyond the tune grace period:

1. the pending tune becomes a failed attempt
2. the station is excluded for the active session
3. DJ MorseCode selects another station satisfying the same intent
4. the replacement becomes the new pending tune

If every eligible station has failed, selection stops rather than retrying indefinitely or leaving the requested musical intent.

Manual station selections are not silently redirected to unrelated stations.

---

# Guiding Principles

## Tell the Truth

Every UI element must be backed by real data.

No simulated transports.

No simulated playback.

No fake metadata.

Silence is represented as an intentional `STANDING BY` state rather than fake initial media.

---

## Playback Is Authoritative

Requesting playback is not the same as successfully playing media.

Application state should reconcile against actual player state before claiming that a station played successfully.

---

## Intent Over Tracks

DJ MorseCode is not primarily a playlist manager.

The central user concept is listening intent:

> What should this session feel like?

Tracks and stations are implementation details of satisfying that intent.

---

## Parse Once

Documents are parsed into domain objects.

The UI never interprets LRC syntax.

---

## Enrichment Must Not Block Playback

Lyrics, metadata enrichment, and future contextual features are optional.

Music should continue when enrichment is unavailable.

---

## Bubble Tea Is the Interaction Engine

Bubble Tea owns terminal interaction.

DJ MorseCode owns behavior.

The domain model should not depend on Bubble Tea.

---

## mpv Is the Playback Engine

DJ MorseCode should not reimplement a media player.

mpv owns playback.

DJ MorseCode owns programming, session behavior, and the listening experience.

---

## Iterate Before Abstracting

Abstractions are introduced after concrete behavior demonstrates the need for them.

Potential future abstractions include:

- additional playback adapters
- additional metadata providers
- additional lyric providers
- user-managed station catalogs

Do not generalize a subsystem merely because it might someday have another implementation.

---

## Ambient Companion

DJ MorseCode should never demand attention.

The application should remain useful as another pane in tmux rather than becoming another application that requires constant management.

---

## Joy

Every feature should make someone smile.