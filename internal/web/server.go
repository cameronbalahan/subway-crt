package web

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"

	"subway-crt/internal/transit"
)

type Server struct {
	port         string
	mu           sync.RWMutex
	currentFrame string
	snapshot     transit.TransitSnapshot
	listeners    map[chan string]bool
	themeChan    chan string
}

func NewServer(port string) *Server {
	if port == "" {
		port = "3000"
	}
	return &Server{
		port:      port,
		listeners: make(map[chan string]bool),
		themeChan: make(chan string, 16),
	}
}

// ThemeChan returns a channel that emits theme switches from the web UI
func (s *Server) ThemeChan() <-chan string {
	return s.themeChan
}

// UpdateFrame updates the latest rendered Lip Gloss ANSI frame and broadcasts to web clients
func (s *Server) UpdateFrame(frame string, snap transit.TransitSnapshot) {
	s.mu.Lock()
	s.currentFrame = frame
	s.snapshot = snap
	s.mu.Unlock()

	s.mu.RLock()
	defer s.mu.RUnlock()
	for ch := range s.listeners {
		select {
		case ch <- frame:
		default:
		}
	}
}

// Start begins the HTTP listener on the configured port
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// 1. Root CRT Web Interface
	mux.HandleFunc("/", s.handleIndex)

	// 2. SSE Stream for Live Charm TUI ANSI Frames
	mux.HandleFunc("/stream", s.handleStream)

	// 3. Current Frame Endpoint
	mux.HandleFunc("/api/frame", s.handleFrame)

	// 4. Transit Snapshot JSON API
	mux.HandleFunc("/api/transit/snapshot", s.handleSnapshot)

	// 5. Theme Switch Endpoint
	mux.HandleFunc("/api/theme", s.handleSetTheme)

	// 6. Pure Go Zip Codebase Downloader
	mux.HandleFunc("/api/download/codebase.zip", s.handleDownloadZip)

	server := &http.Server{
		Addr:    ":" + s.port,
		Handler: mux,
	}

	return server.ListenAndServe()
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	goVer := strings.TrimPrefix(runtime.Version(), "go")
	html := strings.ReplaceAll(crtHTMLTemplate, "{{GO_VERSION}}", "Go "+goVer)
	w.Write([]byte(html))
}

func (s *Server) handleSetTheme(w http.ResponseWriter, r *http.Request) {
	theme := r.URL.Query().Get("set")
	if theme != "" {
		select {
		case s.themeChan <- theme:
		default:
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ch := make(chan string, 4)
	s.mu.Lock()
	s.listeners[ch] = true
	initialFrame := s.currentFrame
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.listeners, ch)
		s.mu.Unlock()
		close(ch)
	}()

	// Send initial frame immediately
	if initialFrame != "" {
		fmt.Fprintf(w, "data: %s\n\n", escapeForSSE(initialFrame))
		flusher.Flush()
	}

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case frame := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", escapeForSSE(frame))
			flusher.Flush()
		}
	}
}

func (s *Server) handleFrame(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	frame := s.currentFrame
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(frame))
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	snap := s.snapshot
	s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snap)
}

func (s *Server) handleDownloadZip(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	filesToZip := []string{
		"cmd/subway-crt/main.go",
		"internal/config/config.go",
		"internal/mta/client.go",
		"internal/transit/state.go",
		"internal/transit/state_test.go",
		"internal/tui/styles.go",
		"internal/tui/board.go",
		"internal/tui/board_test.go",
		"internal/tui/map.go",
		"internal/tui/tui.go",
		"internal/web/server.go",
		"scripts/subway-crt.service",
		"scripts/build-arm64.sh",
		"Makefile",
		"go.mod",
		"go.sum",
	}

	for _, file := range filesToZip {
		if content, err := os.ReadFile(file); err == nil {
			f, err := zw.Create(file)
			if err == nil {
				f.Write(content)
			}
		}
	}

	// Add README.md
	if f, err := zw.Create("README.md"); err == nil {
		readme := "# subway-crt: Pure Go NYC Subway Departure Board & Schematic Map\n\n" +
			"An autonomous, retro-futuristic CRT departure board and live animated transit schematic for the Bedford Ave (L) station. Built exclusively in Pure Go with Charm Bracelet libraries (Bubble Tea, Lip Gloss) and native GTFS-Realtime Protobuf feeds.\n\n" +
			"## Hardware Target\n" +
			"Raspberry Pi (3B / 4 / 5 / Zero 2W), connected to CRT / HDMI display or /dev/tty1 direct console.\n\n" +
			"## Fixing \"Exec format error\" on Raspberry Pi\n" +
			"If running the binary gives **\"Exec format error\"**, your Raspberry Pi is running a **32-bit OS** (e.g. standard Raspberry Pi OS 32-bit / `armv7l`).\n\n" +
			"Check your Pi's kernel architecture by running:\n" +
			"```bash\n" +
			"uname -m\n" +
			"```\n" +
			"- If output is **`aarch64`**: Run `./subway-crt-arm64` (64-bit)\n" +
			"- If output is **`armv7l`** or **`armv6l`**: Run `./subway-crt-armv7` (32-bit)\n\n" +
			"## Cross-Compiling for Raspberry Pi\n" +
			"```bash\n" +
			"# Build 64-bit ARM (Raspberry Pi OS 64-bit / aarch64)\n" +
			"make build-arm64\n\n" +
			"# Build 32-bit ARM (Raspberry Pi OS 32-bit / armv7l)\n" +
			"make build-armv7\n\n" +
			"# Or build both automatically:\n" +
			"./scripts/build-arm64.sh\n" +
			"```\n\n" +
			"## Running on the Pi\n" +
			"```bash\n" +
			"# On 64-bit OS:\n" +
			"./subway-crt-arm64\n\n" +
			"# On 32-bit OS:\n" +
			"./subway-crt-armv7\n" +
			"```\n"
		f.Write([]byte(readme))
	}

	zw.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\"subway-crt-pure-go.zip\"")
	io.Copy(w, buf)
}

func escapeForSSE(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}

const crtHTMLTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>subway-crt: NYC Subway Real-Time Tracker &amp; System Map (Pure Go / Charm TUI)</title>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Share+Tech+Mono&family=JetBrains+Mono:wght@400;700;800&family=VT323&display=swap" rel="stylesheet">
  <style>
    :root {
      --crt-pink: #FF007F;
      --crt-cyan: #00E5FF;
      --crt-yellow: #FFE600;
      --crt-violet: #9D00FF;
      --crt-green: #00FF66;
      --crt-amber: #FFB000;
      --crt-bg: #07070F;
      --theme-color: #00E5FF;
      --glow-color: rgba(0, 229, 255, 0.45);
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background: #040408;
      color: var(--theme-color);
      font-family: 'Share Tech Mono', 'JetBrains Mono', monospace;
      overflow: hidden;
      height: 100vh;
      display: flex;
      flex-direction: column;
      user-select: none;
    }

    /* Outrun Synthwave Bezel Frame */
    .crt-bezel {
      flex: 1;
      position: relative;
      margin: 8px 10px 4px 10px;
      border: 5px solid #18122B;
      border-radius: 14px;
      background: var(--crt-bg);
      box-shadow: 
        inset 0 0 100px rgba(0, 0, 0, 0.95),
        0 0 0 1px #2E1A47,
        0 0 35px rgba(157, 0, 255, 0.25),
        0 0 15px rgba(255, 0, 127, 0.2);
      overflow: hidden;
      display: flex;
      flex-direction: column;
    }

    /* Outrun Scanlines */
    .crt-scanlines {
      position: absolute;
      top: 0; left: 0; right: 0; bottom: 0;
      background: linear-gradient(rgba(18, 16, 26, 0) 50%, rgba(0, 0, 0, 0.42) 50%),
                  linear-gradient(90deg, rgba(255, 0, 127, 0.03), rgba(0, 229, 255, 0.02), rgba(157, 0, 255, 0.03));
      background-size: 100% 3px, 6px 100%;
      pointer-events: none;
      z-index: 40;
      opacity: 0.9;
    }

    /* Phosphor screen curvature vignette */
    .crt-vignette {
      position: absolute;
      top: 0; left: 0; right: 0; bottom: 0;
      background: radial-gradient(circle at center, transparent 60%, rgba(4, 4, 10, 0.8) 100%);
      pointer-events: none;
      z-index: 41;
    }

    /* Glass glare sweep */
    .crt-glare {
      position: absolute;
      top: 0; left: 0; right: 0; height: 35%;
      background: linear-gradient(180deg, rgba(255, 255, 255, 0.04) 0%, rgba(255, 255, 255, 0) 100%);
      pointer-events: none;
      z-index: 42;
    }

    /* Terminal Output Container */
    #terminal-display {
      flex: 1;
      padding: 14px 20px;
      font-size: 13.5px;
      line-height: 1.28;
      white-space: pre;
      overflow: auto;
      text-shadow: 0 0 5px var(--glow-color), 0 0 12px var(--glow-color);
      font-weight: 600;
      z-index: 10;
    }

    /* Synthwave Control Bar */
    .crt-footer {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 7px 16px;
      background: #080712;
      border-top: 1px solid rgba(157, 0, 255, 0.3);
      font-size: 11.5px;
      z-index: 50;
    }

    .btn {
      background: #120F24;
      color: #9D9DAF;
      border: 1px solid #2B2547;
      padding: 5px 12px;
      border-radius: 4px;
      font-family: inherit;
      font-size: 11px;
      font-weight: bold;
      cursor: pointer;
      transition: all 0.15s ease;
      letter-spacing: 0.5px;
    }
    .btn:hover {
      background: #1E193B;
      color: #FFF;
      border-color: var(--crt-cyan);
      box-shadow: 0 0 10px rgba(0, 229, 255, 0.35);
    }
    .btn.active {
      background: var(--crt-pink);
      color: #FFF;
      border-color: var(--crt-pink);
      box-shadow: 0 0 12px rgba(255, 0, 127, 0.6);
    }
    .btn.active-cyan {
      background: var(--crt-cyan);
      color: #000;
      border-color: var(--crt-cyan);
      box-shadow: 0 0 12px rgba(0, 229, 255, 0.6);
    }
    .btn.active-green {
      background: var(--crt-green);
      color: #000;
      border-color: var(--crt-green);
      box-shadow: 0 0 12px rgba(0, 255, 102, 0.6);
    }
    .btn.active-amber {
      background: var(--crt-amber);
      color: #000;
      border-color: var(--crt-amber);
      box-shadow: 0 0 12px rgba(255, 176, 0, 0.6);
    }

    .tag-purego {
      background: #00ADD8;
      color: #FFF;
      padding: 2px 7px;
      border-radius: 3px;
      font-weight: 800;
      font-size: 10px;
      letter-spacing: 0.5px;
      display: inline-flex;
      align-items: center;
      gap: 4px;
    }

    .live-dot {
      width: 7px;
      height: 7px;
      border-radius: 50%;
      background: #00FF66;
      box-shadow: 0 0 8px #00FF66;
      animation: pulse 1.8s infinite;
    }
    @keyframes pulse {
      0%, 100% { opacity: 1; transform: scale(1); }
      50% { opacity: 0.4; transform: scale(0.85); }
    }
  </style>
</head>
<body>
  <div class="crt-bezel">
    <div class="crt-glare"></div>
    <div class="crt-scanlines" id="scanlines-layer"></div>
    <div class="crt-vignette"></div>
    <pre id="terminal-display">INITIALIZING SUBWAY-CRT (PURE GO + CHARM BRACELET ENGINE)...</pre>
  </div>

  <div class="crt-footer">
    <div style="display:flex; gap:6px; align-items:center;">
      <span style="color:#777; font-weight:bold; font-size:10px;">PALETTE:</span>
      <button class="btn active" id="btn-outrun" onclick="setPalette('outrun', '#00E5FF')">⚡ SYNTHWAVE</button>
      <button class="btn" id="btn-phosphor" onclick="setPalette('phosphor', '#00FF66')">P1 GREEN</button>
      <button class="btn" id="btn-amber" onclick="setPalette('amber', '#FFB000')">P3 AMBER</button>
    </div>

    <div style="display:flex; gap:8px; align-items:center;">
      <button class="btn" onclick="toggleScanlines()">SCANLINES: ON</button>
      <button class="btn" onclick="toggleFullscreen()">⛶ 1080p KIOSK</button>
      <a href="/api/download/codebase.zip" class="btn" style="text-decoration:none; display:inline-flex; align-items:center; gap:5px; color:#00E5FF; border-color:#00E5FF;">
        📥 DOWNLOAD PURE GO CODEBASE (.ZIP)
      </a>
    </div>

    <div style="display:flex; gap:10px; align-items:center;">
      <span class="tag-purego">{{GO_VERSION}}</span>
      <span style="color:#666; font-size:11px;">CHARM BUBBLE TEA</span>
      <span class="live-dot" title="Live MTA GTFS-RT Streaming"></span>
    </div>
  </div>

  <script>
    const term = document.getElementById('terminal-display');
    const scanlines = document.getElementById('scanlines-layer');

    function setPalette(name, glowHex) {
      document.documentElement.style.setProperty('--theme-color', glowHex);
      document.documentElement.style.setProperty('--glow-color', glowHex + '66');
      document.querySelectorAll('#btn-outrun, #btn-phosphor, #btn-amber').forEach(b => {
        b.className = 'btn';
      });
      const activeBtn = document.getElementById('btn-' + name);
      if (activeBtn) {
        if (name === 'outrun') activeBtn.className = 'btn active';
        else if (name === 'phosphor') activeBtn.className = 'btn active-green';
        else if (name === 'amber') activeBtn.className = 'btn active-amber';
      }
      fetch('/api/theme?set=' + name).catch(() => {});
    }

    let scanlinesEnabled = true;
    function toggleScanlines() {
      scanlinesEnabled = !scanlinesEnabled;
      scanlines.style.display = scanlinesEnabled ? 'block' : 'none';
      event.target.textContent = 'SCANLINES: ' + (scanlinesEnabled ? 'ON' : 'OFF');
    }

    function toggleFullscreen() {
      if (!document.fullscreenElement) {
        document.documentElement.requestFullscreen().catch(() => {});
      } else {
        document.exitFullscreen().catch(() => {});
      }
    }

    // Comprehensive ANSI TrueColor parser converting Charm Lip Gloss SGR escapes to styled HTML
    function ansiToHtml(text) {
      if (!text) return '';

      // Escape HTML special chars
      let s = text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');

      // Parse ANSI escape sequences
      // Matches \x1b[ ... m
      const ansiRegex = /\x1b\[([0-9;]*)m/g;
      let html = '';
      let lastIndex = 0;
      let activeSpans = 0;

      let match;
      while ((match = ansiRegex.exec(s)) !== null) {
        html += s.substring(lastIndex, match.index);
        lastIndex = match.index + match[0].length;

        const codes = match[1] ? match[1].split(';').map(Number) : [0];

        if (codes[0] === 0 || match[1] === '') {
          // Reset
          while (activeSpans > 0) {
            html += '</span>';
            activeSpans--;
          }
        } else {
          let style = '';
          let i = 0;
          while (i < codes.length) {
            const c = codes[i];
            if (c === 1) {
              style += 'font-weight:800;';
              i++;
            } else if (c === 38 && codes[i+1] === 2) {
              // 24-bit TrueColor Foreground: 38;2;R;G;B
              const r = codes[i+2], g = codes[i+3], b = codes[i+4];
              style += 'color:rgb(' + r + ',' + g + ',' + b + ');';
              i += 5;
            } else if (c === 48 && codes[i+1] === 2) {
              // 24-bit TrueColor Background: 48;2;R;G;B
              const r = codes[i+2], g = codes[i+3], b = codes[i+4];
              style += 'background-color:rgb(' + r + ',' + g + ',' + b + ');';
              i += 5;
            } else if (c >= 30 && c <= 37) {
              const stdColors = ['#000','#c00','#0c0','#c60','#00c','#c0c','#0cc','#ccc'];
              style += 'color:' + stdColors[c-30] + ';';
              i++;
            } else if (c >= 40 && c <= 47) {
              const stdBg = ['#000','#c00','#0c0','#c60','#00c','#c0c','#0cc','#ccc'];
              style += 'background-color:' + stdBg[c-40] + ';';
              i++;
            } else {
              i++;
            }
          }
          if (style) {
            html += '<span style="' + style + '">';
            activeSpans++;
          }
        }
      }

      html += s.substring(lastIndex);
      while (activeSpans > 0) {
        html += '</span>';
        activeSpans--;
      }
      return html;
    }

    // Connect to pure Go Server-Sent Events stream
    let eventSource;
    function connectStream() {
      eventSource = new EventSource('/stream');
      eventSource.onmessage = function(e) {
        try {
          const raw = JSON.parse(e.data);
          term.innerHTML = ansiToHtml(raw);
        } catch {
          term.innerHTML = ansiToHtml(e.data);
        }
      };
      eventSource.onerror = function() {
        eventSource.close();
        setTimeout(connectStream, 1500);
      };
    }

    connectStream();
  </script>
</body>
</html>
`
