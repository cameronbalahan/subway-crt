# subway-crt: NYC Subway Real-Time Tracker & System Map

An autonomous, retro-futuristic CRT departure board and animated transit schematic for the Bedford Ave (L) station, designed as a 1080p full-screen terminal kiosk with real-time MTA GTFS-RT ingestion, high-fidelity offline simulation, and a complete exportable Go/Charm TUI codebase for direct-to-console Raspberry Pi 3B deployment.

## User Review & Critical Decisions

> [!IMPORTANT]
> The following architectural decisions were confirmed during initial requirements clarification:

- **Deployment Mode**: Interactive Web CRT Kiosk simulating the 1080p direct-to-console display, paired with an integrated exporter providing the complete Go 1.27 + Bubble Tea / Lip Gloss / systemd source tree for physical Raspberry Pi 3B hardware.
- **Data Ingestion Strategy**: Dual-mode engine. Actively ingests live MTA Open Data GTFS-RT feeds (`nyct%2Fgtfs-l`, subway alerts, and intersecting trunk lines) via server-side proxy; automatically falls back to a deterministic, schedule-accurate kinematic simulation whenever MTA feeds timeout or network connectivity degrades.
- **Station Scope**: Dedicated Bedford Ave (L) commuter kiosk. The upper departure board is locked to Northbound (L10N: Manhattan-bound to 8 Av) and Southbound (L10S: Brooklyn-bound to Canarsie-Rockaway Pkwy), faithfully matching the physical hardware specification.

---

## 1. Overview & Core Concept

- **What It Does**: Transports transit commuters and retro tech enthusiasts to an authentic 1980s phosphor-CRT transit hub monitor. The top third provides live countdowns and train locations for Bedford Ave (L), accompanied by an alert banner for weekend closures or track incidents. The lower two-thirds renders an animated vector schematic of Manhattan and North Brooklyn transit trunks, showing real-time train beacons gliding along track paths.
- **Target Audience / Persona**: NYC commuters living off Bedford Ave, transit cartography enthusiasts, retro-computing makers running dedicated Pi displays, and subway riders needing glanceable headways.
- **Key Value**: Delivers zero-clutter, instant-read transit intelligence with phosphor aesthetics, zero heap bloat, and turn-key portability to physical Raspberry Pi hardware.

---

## 2. User Experience & Visual Design

### Key User Flows
1. **Always-On Kiosk View**: Default full-screen view partitioned into the Bedford Ave Departure Board (top 35%) and the Manhattan/Brooklyn Schematic Transit Map (bottom 65%). Live countdown timers tick down continuously; real-time train glyphs glide along track vectors.
2. **Alert Trigger & Incident Display**: When planned weekend closures or service disruptions occur on the L train (e.g., track work between 8 Av and Lorimer), an Outrun Magenta / Red high-priority warning banner flashes across the header with specific impact details.
3. **Display Controls & CRT Tuner**: Minimalist keyboard/touch drawer allows toggling scanlines, phosphor bloom intensity, green vs. amber phosphor modes, audio arrival chimes, and manual switching between Live MTA feeds and Kinematic Simulation.
4. **Pi 3B Hardware Exporter**: Commuters wishing to flash their physical Raspberry Pi 3B can open the Hardware Export drawer to inspect, copy, or download the full Go codebase (`main.go`, `mta/`, `transit/`, `tui/`, `systemd` service, `Makefile`, and ARM64 cross-compile script).

### Visual Identity & Theme
- **Aesthetic Direction**: High-contrast Outrun / Monochrome Phosphor CRT terminal. Evokes an authentic 1980s transit command center with phosphor persistence, subtle scanlines, and mechanical clarity.
- **Color Palette & Mood**:
  - Deep Void Canvas: `#0B0B10`
  - CRT Phosphor Green: `#00FF66` (primary nominal state & beacons)
  - CRT Warm Amber: `#FFB000` (secondary trunk lines & advisory alerts)
  - Outrun Magenta / Red Alert: `#FF007F` (service disruptions & warnings)
  - Cyber Cyan: `#00E5FF` (schematic rail vectors & East River boundary)
  - L Line Slate Silver: `#D4D4D8` (dedicated L track and station pips)
- **Typography & Hierarchy**:
  - Primary Monospace Display: Tabular figures with clean retro terminal proportions (`font-mono tabular-nums`).
  - Board Headings: Prominent station callouts (`BEDFORD AV (L)`, `MANHATTAN - 8 AV`, `CANARSIE - ROCKAWAY PKWY`).
  - Arrival Numbers: Massive, glanceable countdown numerals (`2 MIN`, `6 MIN`, `11 MIN`) paired with exact station locations (`At 1 Av`, `Entering Lorimer St`).
- **Interactive Feedback & Motion**:
  - Smooth 60fps interpolation for train beacons traveling along track curves and across the East River tube.
  - Subdued CRT scanline sweep and subtle phosphor persistence glow (customizable / toggleable for low-spec devices).
  - Flashing Outrun magenta beacon pulses on delayed or disrupted track segments.

---

## 3. Key Product Decisions & Trade-Offs

- **Decision 1: Dual-Mode Transit Ingestion (Live MTA + Kinematic Simulator)**
  - *Chosen Approach*: Server proxy ingests official MTA GTFS-RT protobuf feeds (`nyct/gtfs-l` and `camsys/subway-alerts`). If the MTA feed is unreachable, returns 429/timeout, or if offline, the system seamlessly transitions to a deterministic schedule simulator calculating train movements based on actual L train headways.
  - *Why*: Guarantees the kiosk is immediately alive, functional, and visually dynamic under all network conditions while providing real-time data when live.
  - *Alternatives Considered*: Live-only (risks blank screens or static error states when MTA endpoints throttle requests) vs. static simulation-only (lacks real-world utility).

- **Decision 2: Integrated SVG/Vector Schematic Map vs Static ASCII Canvas**
  - *Chosen Approach*: Scalable SVG vector schematic mapped to terminal coordinate grid nodes with CRT phosphor styling. Trains are represented as glowing vector beacons that interpolate along track splines.
  - *Why*: Looks razor-sharp across all resolutions (from 1080p television monitors down to laptop screens) while faithfully rendering the geometry of 14th Street, the East River crossing, Bedford Ave, and intersecting trunks (A/C/E, 1/2/3, B/D/F/M, N/Q/R/W, 4/5/6, G, J/M/Z).
  - *Alternatives Considered*: Pure text character grid (limits beacon animation smoothness) vs raster canvas (poor scaling and blurry fonts).

- **Decision 3: Complete Go 1.27 Raspberry Pi Source Code Export**
  - *Chosen Approach*: Bundle the complete Go source tree adhering to the Charm Bubbletea/Lipgloss architecture, protobuf schemas, systemd service, and ARM64 Makefile in a dedicated, interactive export panel with direct code inspection and one-click zip download.
  - *Why*: Delivers 100% fidelity to the user's technical specification for running on physical Pi 3B hardware while providing the instant web experience.

---

## 4. Technical Architecture & Data Strategy

```
┌────────────────────────────────────────────────────────────────────────┐
│                        subway-crt Kiosk System                         │
└────────────────────────────────────────────────────────────────────────┘
                                    │
          ┌─────────────────────────┴─────────────────────────┐
          ▼                                                   ▼
┌───────────────────────────────┐           ┌────────────────────────────────┐
│   Server Proxy & Ingestion    │           │    Client Kiosk Application    │
│  - /api/transit/live          │           │   (React 19 + Tailwind v4)     │
│    * MTA GTFS-RT Protobuf     │           └────────────────────────────────┘
│    * MTA Subway Alerts        │                           │
│  - /api/transit/status        │             ┌─────────────┴─────────────┐
│  - Realistic Fallback Model   │             ▼                           ▼
└───────────────────────────────┘   ┌────────────────────┐      ┌────────────────────┐
                                    │  Bedford Board     │      │  Schematic Map     │
                                    │  - L10N Countdowns │      │  - 14th St Trunk   │
                                    │  - L10S Countdowns │      │  - East River Tube │
                                    │  - Train Positions │      │  - Brooklyn Loops  │
                                    │  - Flashing Alerts │      │  - Intersecting N  │
                                    └────────────────────┘      └────────────────────┘
                                              │                           │
                                              └─────────────┬─────────────┘
                                                            ▼
                                                ┌───────────────────────┐
                                                │   CRT Effect Shader   │
                                                │  - Phosphor Bloom     │
                                                │  - Scanlines & Glow   │
                                                │  - Amber / Green Mode │
                                                │  - Pi Source Exporter │
                                                └───────────────────────┘
```

### Data Model & State
- **TransitSnapshot**:
  - `timestamp`: Milliseconds timestamp of the current snapshot.
  - `feedMode`: `'live' | 'simulated' | 'degraded'`.
  - `manhattanTrains`: List of upcoming L10N departures with `routeId`, `destination` ("8 Av"), `minutesAway`, and `currentStatus` (e.g., "At 1 Av").
  - `canarsieTrains`: List of upcoming L10S departures with `routeId`, `destination` ("Rockaway Pkwy"), `minutesAway`, and `currentStatus` (e.g., "Approaching Bedford Av").
  - `alerts`: List of active alerts with `severity`, `title`, `description`, and `affectsL`.
  - `trainPositions`: Real-time vehicle positions with `trainId`, `direction`, `progress` (0.0 to 1.0 along line segment), and `status`.

### Interactive Component & State Mapping
- **Board Controller**: Updates countdowns every second, manages alert banners with priority sorting, and pulses warning colors if delays occur.
- **Map Vector Renderer**: Smoothly steps vehicle coordinates along track curves using requestAnimationFrame, highlighting station stops and transfer hubs.
- **CRT Controls Overlay**: Allows switching color profile (Phosphor Green, Amber Monochrome, Cyber Synthwave), adjusting CRT curve/scanline density, and triggering test alerts (e.g. simulated weekend track work).
- **Codebase Exporter**: Provides multi-tab viewing and download of the Go project files (`main.go`, `board.go`, `map.go`, `styles.go`, `subway-crt.service`, `Makefile`).
