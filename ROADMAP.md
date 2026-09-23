# DJ MorseCode Roadmap

DJ MorseCode optimizes for discovery.

The roadmap is organized around product directions rather than promised release numbers. Features move into a release when they make the listening experience better without making the DJ demand more attention.

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
- longer-term listening history
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
  ↓
failed tune attempts
  ↓
warning threshold
  ↓
flagged for decommission
  ↓
decommissioned
```

A station that reaches a threshold such as five independent failed tune attempts may be removed from the effective active catalog and retained in a decommissioned catalog for inspection or rehabilitation.

A temporary outage should not permanently condemn a station.

## DJ Notes

Inspired by VH1 Pop-Up Video.

Occasionally surface small pieces of context without interrupting playback.

Examples:

> This guitar solo was recorded in one take.

> This band originally opened for Soundgarden.

> You haven't heard this station in eight months.

Notes should disappear automatically and never become walls of text.

## Playback Sources

Playback should not be permanently tied to one source.

Possible future adapters include:

* local libraries
* Jellyfin
* Navidrome
* Spotify
* Apple Music
* YouTube Music
* Last.FM
* other radio catalogs

The terminal renderer should care about what is playing, not where it came from.

## Lyrics and Enrichment

Lyrics should remain provider-agnostic.

Potential directions include:

* additional lyric providers
* user-managed local LRC discovery
* richer lyric provenance
* transcripts for non-music media
* track and artist context

## Catalogs

Potential station-catalog improvements include:

* user-managed station configuration
* M3U catalog imports
* multiple endpoints per station
* preferred codec/quality selection
* endpoint failover
* catalog health tooling

## Interface

The idle screen may eventually evolve into a lightweight home view containing useful context such as:

* quick-start intents
* recently tuned stations
* favorite stations
* current station health

It should still feel like another pane in tmux rather than another application.

## Other Media

The architecture may eventually support listening contexts beyond live music radio:

* podcasts
* audiobooks
* spoken-word streams
* local media

Those modes should earn their way into the product rather than complicating the radio experience prematurely.

## Joy

Every feature should still make someone smile.
