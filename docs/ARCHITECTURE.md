# DJ MorseCode Architecture

> "Every UI element must be backed by real data."

DJ MorseCode is a terminal-first music companion built around continuous discovery.

The application translates listening intent and explicit media selections into playback, delegates media transport to mpv, enriches playback with metadata and lyrics when available, records runtime evidence, derives meaningful listening history, and presents the resulting session through a Bubble Tea terminal interface.

The architecture intentionally separates source identity, programming decisions, playback, observations, enrichment, history, domain data, and presentation so each subsystem has a clear responsibility.

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
          +-------------------------------+-------------------------------+
          |                 |             |              |                |
          v                 v             v              v                v
     +---------+       +---------+   +---------+   +-------------+   +---------+
     |  radio  |       | source  |   | player  |   | observation |   |   ui    |
     |program- |       | identity|   |playback |   |   evidence  |   | render  |
     |  ming   |       |catalogs |   +----+----+   +------+------+   +---------+
     +----+----+       +----+----+        |               |
          |                 |             v               v
          |                 |        +---------+      +---------+
          |                 |        |player/  |      | history |
          |                 |        |  mpv    |      |projection|
          |                 |        +----+----+      +---------+
          |                 |             |
          +-----------------+-------------+
                            |
                            v
                           mpv
                            |
                            v
                        Playback
                            |
          +-----------------+------------------+
          |                                    |
          v                                    v
      metadata                            enrichment/context
                                               |
                              +----------------+----------------+
                              |                |                |
                              v                v                v
                         musicbrainz        lastfm           lrclib
                                                               |
                                                               v
                                                            library
                                                               |
                                                               v
                                                             lyrics
                                                               |
                                                               v
                                                             music
```

`cmd/dj` coordinates these systems.

It does not own the underlying rules for station matching, source identity, observation deduplication, history derivation, lyric parsing, enrichment matching, or mpv IPC.

---

# Core Architectural Distinctions

Several concepts that may look similar in the UI intentionally remain separate.

## Source Identity Is Not Playback Transport

A media source describes what an item is and where it came from.

mpv describes how that item is played.

```text
Source Media
    |
    | identity / provenance
    v
source.MediaItem
    |
    v
Player.Load()
    |
    v
mpv
```

A YouTube video, radio station, discovered Last.fm track, or future Jellyfin item can carry source identity without requiring a provider-specific playback engine inside the application model.

Network transport is not itself a media type.

In particular:

```text
network playback != radio
```

Finite network media and live radio have different lifecycle semantics even when mpv transports both.

## Navigation History Is Not Listening History

`radio.History` is browser-style station navigation state.

It supports:

- current station awareness
- backward navigation
- forward navigation
- recent-station avoidance
- station-selection policy

It contains a cursor and may truncate forward navigation when a new station is selected after navigating backward.

Listening history is different.

`history` derives chronological listening events from runtime observations.

```text
radio.History
    |
    +--> previous / next station navigation

observation
    |
    +--> history.Derive()
             |
             v
       chronological listening history
```

Navigating backward or forward does not redefine what was actually listened to.

## Observations Are Evidence, Not Presentation

The observation layer records runtime facts.

Examples include:

- station tune requested
- station tune confirmed
- station tune failed
- normalized radio playback metadata
- source-media playback

Observation History exposes this evidence for diagnostics.

Listening History projects only the evidence that currently represents a meaningful listening event.

For v1:

```text
StationTuneConfirmed
        |
        +--> listening-history station entry

MediaObservation
        |
        +--> listening-history media entry

PlaybackObservation
        |
        +--> retained as evidence
             not yet projected into listening history
```

Track-level history from radio playback is intentionally deferred until meaningful-listen semantics are defined.

---

# Primary Runtime Flows

## Radio Programming

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
Pending Tune
      |
      v
     mpv
```

A user may establish intent through:

- a vibe
- a mood family
- a genre

Direct station selection bypasses programming intent and gives control back to the user.

An active programmed intent survives station changes.

This allows manual next behavior and automatic station rotation to continue programming toward the same listening goal.

## Station Playback Reconciliation

A request to load a station is not treated as proof that playback succeeded.

```text
Player.Load()
      |
      v
StationTuneRequested
      |
      v
Pending Station Identity
      |
      +------ mpv becomes active ------+
      |                                |
      |                                v
      |                        Reconcile Station
      |                                |
      |                 +--------------+--------------+
      |                 |                             |
      |                 v                             v
      |          requested URL                 resolved/redirected
      |             matches                         URL differs
      |                 |                             |
      |                 +-------------+---------------+
      |                               |
      |                               v
      |                     StationTuneConfirmed
      |                               |
      |                    +----------+----------+
      |                    |                     |
      |                    v                     v
      |             radio.History          observation
      |
      +------ grace period expires ------> StationTuneFailed
                                               |
                                               v
                                      Exclude For Session
                                               |
                                               v
                                          Choose Again
```

mpv remains authoritative for whether playback is active.

The pending station identity remains authoritative for a station DJ MorseCode explicitly requested when mpv resolves or redirects the stream to another URL.

Exact stream-URL equality is therefore not required to confirm a requested station.

Failed stations are excluded from the current programmed session so recovery cannot immediately select the same failed endpoint again.

Manual station selections are not silently redirected to unrelated stations after failure.

## Source-Media Playback

Finite source media follows a different lifecycle from radio.

```text
Catalog Selection
      |
      v
source.MediaItem
      |
      v
Player.Load()
      |
      v
     mpv
      |
      v
Playback Identity
      |
      +--> track initialization
      |
      +--> MediaObservation
      |
      +--> enrichment
      |
      +--> lyrics
```

A successful `Player.Load()` request does not itself create a media observation.

Source-media playback is observed once mpv reports active playback.

Repeated application ticks are deduplicated by the observation recorder.

A genuine revisit remains a new event:

```text
Media A
   |
   v
Media B
   |
   v
Media A
```

and:

```text
Media A
   |
   v
Radio
   |
   v
Media A
```

both preserve the second visit to Media A.

## Listening History

Listening history is derived rather than independently mutated.

```text
Station Observations -----+
                          |
                          +--> history.Derive() --> []history.Entry --> UI
                          |
Media Observations -------+
```

The current projection includes:

- confirmed station tunes
- valid source-media playback observations

It excludes:

- tune requests
- tune failures
- raw radio playback observations

Entries are sorted chronologically and revisits are preserved.

This keeps history deterministic and prevents navigation behavior from rewriting listening chronology.

## Lyrics and Enrichment

Lyrics and metadata enrichment are optional playback enhancements.

```text
Playback Identity
      |
      +-------------------+
      |                   |
      v                   v
  Enrichment            Lyrics
      |                   |
      v                   v
 MusicBrainz         Local Cache
 Last.fm                  |
                          +------ hit ------> Parsed Lyrics
                          |
                          +------ miss -----> LRCLIB
                                                |
                                                v
                                           Best Match
                                                |
                                                v
                                           Cache Result
                                                |
                                                v
                                          Parsed Lyrics
```

Playback continues when enrichment or lyrics are unavailable.

Remote work is asynchronous so network activity does not block the terminal application.

Responses are associated with track identity so stale results cannot replace the current track's state.

Finite network media may initially expose identity before useful duration becomes available. Duration-sensitive lyric matching is deferred until playback duration is known.

---

# Package Responsibilities

## source

Owns provider-neutral media identity, provenance, and source capabilities.

Core concepts include:

- `Kind`
- `Source`
- `ItemRef`
- `MediaKind`
- `MediaItem`
- catalog providers
- discoverers
- typed discovery requests
- discovery relevance

The package answers:

> What kind of media is this?

> Where did it come from?

> What stable provider identity does it carry?

Source identity does not determine playback transport.

## m3u

Owns M3U playlist ingestion.

Responsibilities include:

- basic M3U parsing
- extended M3U parsing
- `EXTINF` metadata handling
- playlist entry identity
- station proposal creation
- source provenance
- exposing station media through source capabilities

M3U parsing does not directly mutate the built-in radio catalog.

Imported entries become structured proposals or media items that other layers may review or consume.

## youtube

Owns YouTube media identity and the current curated YouTube catalog.

Responsibilities include:

- supported YouTube URL parsing
- video identity
- playlist identity
- canonical media references
- curated media labels
- curated artist/title identity

YouTube does not own playback.

The selected URL is handed to the player, and mpv/yt-dlp handle transport.

Current YouTube playback is audio-only.

## lastfm

Owns Last.fm integration.

Responsibilities include:

- track information
- track tags
- artist tags
- contextual tag evidence
- similar-track discovery
- discovery relevance
- provider media identity

Last.fm currently acts as an enrichment/context and discovery provider, not as the playback engine.

## musicbrainz

Owns MusicBrainz integration.

Responsibilities include:

- recording search
- recording lookup
- canonical recording candidates
- identifiers such as MusicBrainz recording IDs and ISRCs
- release evidence
- enrichment candidates

MusicBrainz results are matched against observed playback identity before being accepted.

## enrichment

Coordinates canonical track enrichment.

Responsibilities include:

- provider enrichment requests
- accepted/ambiguous/no-match status
- enrichment caching
- associating enrichment with playback identity

Enrichment must not block playback.

## context

Coordinates track-context providers.

Context is separate from canonical identity.

A track may have a stable canonical identity while multiple providers contribute contextual information such as tags or release evidence.

## observation

Owns runtime evidence.

Observation types include:

- station observations
- radio playback observations
- source-media observations

The recorder owns transition deduplication.

The service coordinates recorder and sink behavior.

The current memory sink stores session evidence and exposes defensive copies for inspection and projection.

Observation semantics should describe what happened, not how a UI intends to display it.

## history

Owns meaningful listening chronology.

`history.Derive()` projects observation evidence into chronological entries.

Current entry kinds include:

- station
- media

The package:

- includes confirmed station tunes
- includes valid source-media observations
- preserves revisits
- sorts entries chronologically
- ignores failed/requested station events
- intentionally excludes radio `PlaybackObservation` from v1 history

The history package does not own station navigation.

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
- browser-style station navigation history
- recent-station avoidance
- explicit station exclusion
- candidate selection
- controlled variety
- selection relaxation
- station proposals

The package answers questions such as:

> Which stations satisfy this listening intent?

and:

> Given these eligible stations and this navigation history, which station should be selected next?

The package does not load media.

Selection returns a station; the application decides when to ask the player to load it.

## player

Provides the playback-facing abstraction used by the application.

Responsibilities include:

- load media
- expose current media metadata
- expose playback position and duration
- expose current media path
- expose playback idle state
- distinguish playback state used by application reconciliation
- determine the active timeline cue
- format playback information

The application should not need to understand mpv IPC commands or property names.

## player/mpv

Owns the mpv integration.

Responsibilities include:

- mpv IPC communication
- property retrieval
- media loading
- managed mpv process startup
- audio-only playback policy
- IPC readiness
- managed process shutdown
- IPC socket cleanup

DJ MorseCode launches mpv idle and globally disables video for the current product experience.

```text
DJ MorseCode starts
      |
      v
Start mpv --no-video
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

DJ MorseCode owns programming behavior and source semantics.

## metadata

Resolves playback metadata into track identity.

Radio streams frequently expose changing ICY/media-title data, while finite media may provide structured tags or curated source identity.

Metadata precedence depends on playback semantics.

Static stream metadata must not overwrite changing radio-track metadata.

Curated source identity may provide canonical artist/title information for finite media while mpv remains authoritative for playback state and duration.

## music

The core parsed-music model.

Contains concepts such as:

- Song
- Cue
- CueType

`music` represents parsed musical information without knowing how that information was retrieved or displayed.

The package knows nothing about:

- mpv IPC
- radio selection
- source providers
- Bubble Tea
- terminal layout
- network retrieval

## lyrics

Parses lyric documents into domain objects.

Responsibilities include:

- parse metadata
- parse timestamps
- parse synchronized lyric cues
- construct lyric timelines

The parser does not perform network retrieval.

It parses documents.

## lrclib

Owns LRCLIB integration.

Responsibilities include:

- search remote lyric records
- identify an appropriate duration-aware result
- translate remote data into the music domain

LRCLIB is optional enrichment.

Failure to retrieve lyrics must never stop playback.

## library

Owns persistent local lyric cache behavior.

Remote lyric results may be stored locally and reused on later playback.

The cache is application data, not repository-owned demo media.

The application should not require bundled songs or lyrics to start.

## ui

Presentation only.

Responsibilities include:

- idle state
- current playback
- station information
- active session intent
- lyric state and synchronized timeline
- station picker
- source-media picker
- vibe picker
- mood picker
- genre picker
- chronological listening history
- raw observation history
- controls and status

The UI never parses LRC syntax.

The UI never decides which station matches an intent.

The UI never decides which observations constitute listening history.

The UI renders application state and domain projections.

## controls

Owns semantic control resolution.

Keyboard input is translated into application actions based on the active control mode.

This keeps literal key bindings from becoming application-domain behavior.

## cmd/dj

The application orchestration layer.

Responsibilities include:

- Bubble Tea model and update loop
- keyboard interaction
- picker state
- active session coordination
- pending tune state
- station identity reconciliation
- station failure recovery
- automatic rotation timing
- source-media selection state
- mpv lifecycle coordination
- connecting domain decisions to playback
- connecting playback state to observations
- coordinating enrichment and lyrics
- connecting history projections to presentation

`cmd/dj` may coordinate packages, but domain rules should remain in the packages that own them.

For example:

```text
cmd/dj decides WHEN to choose another station.

radio decides WHICH stations are valid choices.

source describes WHAT media is and WHERE it came from.

player decides HOW to request playback.

mpv decides HOW media is transported and rendered as audio.

observation decides WHETHER a runtime fact is a new observation.

history decides WHICH observations represent listening chronology.

ui decides HOW application state is presented.
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
- raw stream/media metadata

## source owns

- media kinds
- source kinds
- source provenance
- provider item references
- provider-neutral media identity
- source capability contracts

## radio owns

- station definitions
- station semantics
- matching rules
- selection rules
- intent representation
- browser-style station navigation history

## observation owns

- runtime observation types
- transition deduplication
- session evidence storage contracts

## history owns

- listening-history entry semantics
- chronological projection from observations

## enrichment/context owns

- canonical enrichment coordination
- contextual provider coordination

## cmd/dj owns

- current application session
- active intent
- pending station tune
- failed stations for the active session
- automatic rotation timing
- picker visibility and selection
- currently selected source media
- orchestration between domains and player

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

Selecting explicit source media also leaves radio programming intent.

This makes direct selection an explicit manual takeover.

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

## Navigation History

`radio.History` records station navigation state.

It supports:

- current station awareness
- backward navigation
- forward navigation
- recent-station avoidance
- station-selection policy

Consecutive additions of the same current station do not create duplicate navigation entries.

Navigating backward and forward moves a cursor rather than creating new tune events.

This structure is intentionally not the user-facing listening history.

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

# Listening History

Listening history represents successful listening events rather than navigation state.

Current inputs are:

```text
StationTuneConfirmed
MediaObservation
```

The projection is chronological and preserves revisits.

For example:

```text
23:27  Beastie Boys - Intergalactic
23:27  A Tribe Called Quest - Can I Kick It?
23:28  Hard Rock Radio FM
23:29  Beastie Boys - Intergalactic
```

The second Beastie Boys entry is intentional.

The listener left that media, listened to something else, and returned.

Raw radio track observations are not currently projected into listening history.

Before adding them, DJ MorseCode should define meaningful-listen semantics such as whether a track must play for a minimum duration or proportion before it counts as heard.

---

# Automatic Programming

During an active intent, stations rotate after a configured dwell interval.

Automatic rotation uses the same station-selection machinery as manual next behavior.

```text
manual next ----------------+
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

DJ MorseCode treats a requested station load as pending until playback is confirmed.

If mpv remains idle beyond the tune grace period:

1. the pending tune becomes a failed attempt
2. the station is excluded for the active session
3. DJ MorseCode selects another station satisfying the same intent
4. the replacement becomes the new pending tune

If every eligible station has failed, selection stops rather than retrying indefinitely or leaving the requested musical intent.

A resolved or redirected stream URL may differ from the catalog URL.

When DJ MorseCode initiated the tune and mpv becomes active, the pending station identity can confirm that playback without requiring exact URL equality.

Manual station selections are not silently redirected to unrelated stations after failure.

---

# Guiding Principles

## Tell the Truth

Every UI element must be backed by real data.

No simulated transports.

No simulated playback.

No fake metadata.

Silence is represented as an intentional `STANDING BY` state rather than fake initial media.

## Playback Is Authoritative

Requesting playback is not the same as successfully playing media.

Application state should reconcile against actual player state before claiming that a station or source media played successfully.

Provider identity may describe what was requested.

mpv remains authoritative for whether playback actually happened.

## Identity Is Not Transport

Source identity and playback transport are separate concerns.

A provider may identify or discover media without becoming the playback engine.

Network transport must not be treated as synonymous with live radio.

## Evidence Before Intelligence

Observations record facts before higher-level features interpret them.

Listening history, future station health, session intelligence, and personalization should derive from evidence rather than creating competing state models.

## Intent Over Tracks

DJ MorseCode is not primarily a playlist manager.

The central user concept remains listening intent:

> What should this session feel like?

Explicit source media may provide useful manual control without replacing discovery as the product's center of gravity.

## Parse Once

Documents are parsed into domain objects.

The UI never interprets LRC or playlist syntax.

## Enrichment Must Not Block Playback

Lyrics, metadata enrichment, contextual providers, and future DJ Notes are optional.

Music should continue when enrichment is unavailable.

## Bubble Tea Is the Interaction Engine

Bubble Tea owns terminal interaction.

DJ MorseCode owns behavior.

The domain model should not depend on Bubble Tea.

## mpv Is the Playback Engine

DJ MorseCode should not reimplement a media player.

mpv owns playback.

DJ MorseCode owns programming, source semantics, session behavior, observations, and the listening experience.

## Iterate Before Abstracting

Abstractions are introduced after concrete behavior demonstrates the need for them.

Current abstractions such as source capabilities and observation/history separation were introduced after concrete integrations demonstrated the need.

Potential future abstractions include:

- additional playback adapters
- additional metadata providers
- additional lyric providers
- durable observation storage
- user-managed catalogs

Do not generalize a subsystem merely because it might someday have another implementation.

## Ambient Companion

DJ MorseCode should never demand attention.

The application should remain useful as another pane in tmux rather than becoming another application that requires constant management.

## Joy

Every feature should make someone smile.