# DJ MorseCode Architecture

> "Every UI element must be backed by real data."

DJ MorseCode is a terminal-first music companion built around
a synchronized musical timeline.

The application intentionally separates domain logic,
playback, parsing, and presentation so each subsystem has a
single responsibility.

This document describes architectural intent rather than
implementation details.

When implementation and this document disagree,
implementation should be questioned before the document is
changed.

---

# Architecture

                    +----------------------+
                    |      main.go         |
                    +----------+-----------+
                               |
      +------------------------+------------------------+
      |                        |                        |
      ▼                        ▼                        ▼
+--------------+       +---------------+       +---------------+
|    player    |       |    lyrics     |       |    radio      |
+--------------+       +---------------+       +---------------+
       |                       |                       |
       +-----------+-----------+-----------+-----------+
                   |                       |
                   ▼                       ▼
             +-------------------------------+
             |            music              |
             +-------------------------------+
                           |
                           ▼
                    +---------------+
                    |      ui       |
                    +---------------+

---

# Package Responsibilities

## music

The domain model.

Defines the language of DJ MorseCode.

Contains:

- Song
- Cue
- CueType

The music package knows nothing about:

- playback
- parsing
- rendering
- Bubble Tea

---

## lyrics

Converts LRC documents into domain objects.

Responsibilities:

- Parse metadata
- Parse timeline
- Parse timestamps
- Parse duration

The parser never constructs applications.

It only parses documents.

---

## player

Owns musical time.

Responsibilities:

- Measure playback time
- Determine the active cue
- Format playback durations
- Calculate playback progress

The player knows nothing about:

- Bubble Tea
- rendering
- terminal layout

---

## radio

Owns station metadata.

Responsibilities:

- Station catalog
- Station lookup

Future responsibilities may include:

- station groups
- genres
- metadata providers

---

## ui

Presentation only.

Responsibilities:

- Render panels
- Render timeline
- Render transport
- Render status

The UI never parses data.

The UI never calculates playback.

The UI renders whatever it is given.

---

# Guiding Principles

## Tell the Truth

Every UI element must be backed by real data.

No simulated transports.

No simulated playback.

No fake metadata.

---

## Parse Once

Documents are parsed into domain objects.

The UI never interprets LRC syntax.

---

## Timeline First

Songs are represented as musical timelines.

Lyrics are one kind of cue.

Instrumental breaks are another.

Future cue types may include:

- trivia
- station IDs
- DJ commentary

---

## Bubble Tea is the Engine

Bubble Tea owns terminal interaction.

DJ MorseCode owns behavior.

The application should be able to change UI frameworks
without changing its domain model.

---

## Iterate Before Abstracting

Abstractions are introduced only after at least two
concrete implementations exist.

Examples:

- playback adapters
- metadata providers
- layout engines

---

# Roadmap

Current

✓ Timeline Engine

✓ Playback Clock

✓ Transport

✓ Progress

Future

□ MPV Adapter

□ LRCLIB

□ Live Radio

□ Responsive Panels

□ Trivia Timeline

□ AI DJ Commentary