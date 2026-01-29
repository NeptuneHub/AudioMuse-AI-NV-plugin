# AudioMuse-AI Navidrome Plugin

A Navidrome plugin that reimplements the "Instant Mix" button to use AudioMuse-AI for similar track recommendations.

**IMPORTANT** InstantMix support in Navidrome is still not released in the stable image, you can find only in the develop image

## HOW-TO Install

- The ENV var ND_PLUGINS_ENABLED, ND_PLUGINS_AUTORELOAD and ND_AGENTS are important, assuming that you deploy with docker compose you should use something like this:

```yaml
version: '3'
services:
  navidrome:
    image: deluan/navidrome:latest
    ports:
      - '4533:4533'
    environment:
      - ND_PLUGINS_ENABLED=true
      - ND_PLUGINS_AUTORELOAD=true
      - ND_AGENTS=audiomuseai
    volumes:
      - ./data:/data
      - /path/to/music:/music:ro
```

- Then you need to put `audiomuseai.ndp` in Navidrome data plugins folder (default: `/data/plugins`).
- Restart Navidrome, go to UI -> Plugins, enable **AudioMuse-AI**, set **AudioMuse-AI API URL** and other configuration parameter.

## HOW-TO build

- Requirements (Ubuntu / macOS): Go, TinyGo.
- Build:

```bash
make build    # -> audiomuseai.wasm
make package  # -> audiomuseai.ndp
```

Full stop.