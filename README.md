# DJ MorseCode

> "The DJ that quietly codes with you."

---

## Overview

DJ MorseCode is not a music player.

It is not another streaming service.

It is not another playlist manager.

DJ MorseCode is a coding companion.

It recreates the feeling of late-night terrestrial radio—where you never quite knew what song was coming next—but removes everything that got in the way.

No twenty-minute commercial breaks.

No DJs talking over the intro.

No endlessly replaying the same twenty songs.

Just music.

Just discovery.

Just flow.

---

## The Feeling

I grew up waiting beside the radio with a cassette recorder.

When the DJ finally played *that* song, you hit RECORD and hoped they wouldn't talk over the intro.

There was something magical about never knowing what came next.

Modern streaming services optimize for control.

DJ MorseCode optimizes for discovery.

It should feel like sitting beside a great late-night radio DJ that somehow always understands your mood.

---

## Philosophy

The goal is not to play exactly the songs you ask for.

The goal is to keep you in the groove.

You shouldn't think about playlists.

You shouldn't think about albums.

You shouldn't even think about songs.

You think about a feeling.

DJ MorseCode handles the rest.

---

## Core Principles

### Continuous Discovery

Every song should feel like it belongs.

Every song should also feel like a pleasant surprise.

---

### Infinite Radio

Music never stops.

No playlist management.

No queue management.

No decision fatigue.

You get what you get.

If you're not feeling it...

Skip.

If it hits...

Enjoy the ride.

---

### Ambient Companion

DJ MorseCode should never demand attention.

It quietly enhances your workspace.

It should feel like another pane in tmux rather than another application.

---

### Joy

Every feature should make someone smile.

---

# Version 0.1

The smallest useful version.

Features

- Display currently playing song
- Display synchronized lyrics
- Karaoke-style highlighting
- Optimized for tmux
- Minimal distractions

Nothing else.

Ship it.

---

# Version 0.2

Music Context

Display

- Artist
- Album
- Year
- Genre
- Elapsed Time

---

# Version 0.3

DJ Notes

Inspired by VH1 Pop-Up Video.

Occasionally display tiny contextual facts.

Examples

> This guitar solo was recorded in one take.

> This band originally opened for Soundgarden.

> You've listened to this song 37 times.

> Last played eight months ago.

These disappear automatically after a few seconds.

No interruptions.

No walls of text.

---

# Version 0.4

Station Controls

Instead of technical controls...

Think like a DJ.

Examples

Keep Cookin'

Continue the current vibe.

---

Surprise Me

Take an unexpected turn without breaking the mood.

---

Go Deeper

Less popular tracks.

B-sides.

Deep cuts.

---

Take Me Sideways

Same energy.

Different genre.

---

Cool It Down

Reduce intensity.

---

Go Harder

Increase intensity.

---

Encore

Stay in this musical neighborhood.

Not replay.

Continue the feeling.

---

# Version 0.5

Playback Adapters

Playback should never be tied to one platform.

Possible adapters include

- Pandora
- Spotify
- Apple Music
- YouTube Music
- Jellyfin
- Navidrome
- Local Library

The renderer should not know where the music came from.

It only knows what is currently playing.

---

# Lyrics Providers

Lyrics should be provider-agnostic.

Preferred format:

Timed lyrics (LRC or equivalent)

Renderer responsibilities

- Previous line
- Current highlighted line
- Next line

Always centered.

Never scrolling walls of text.

---

# Renderer

Designed specifically for developers.

tmux first.

Terminal first.

Keyboard first.

Everything else is secondary.

---

Example

────────────────────────────────────────────

♫ Interstate Love Song

Stone Temple Pilots

Alternative • 1994

────────────────────────────────────────────

      Leaving on a southern train...

████████████ Only yesterday you lied...

      Promises of what I seemed...

────────────────────────────────────────────

DJ NOTE

This song was recorded in one take.

────────────────────────────────────────────

---

# Future

Mood Engine

Instead of choosing playlists...

Choose feelings.

Examples

Late-night coding

Progressive rabbit hole

Forgotten 90s alternative

Dark synthwave

Jazz while debugging

Sunday morning coffee

Coding through a thunderstorm

The DJ builds the station.

---

# Long-Term Vision

DJ MorseCode is not trying to replace Spotify.

It isn't trying to replace Pandora.

It isn't trying to replace radio.

It's trying to recreate something we've quietly lost.

The joy of discovery.

The comfort of an excellent DJ.

The feeling that someone else is curating the soundtrack while you disappear into your work.

---

# One Design Rule

Every feature must answer one question.

> Does this make coding more enjoyable?

If the answer is no...

It doesn't belong.

---

# Motto

Don't build a music player.

Build the DJ that quietly codes with you.
