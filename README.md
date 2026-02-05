# AudioMuse-AI Navidrome Plugin

<p align="center">
  <img src="https://github.com/NeptuneHub/audiomuse-ai-plugin/blob/master/audiomuseai.png?raw=true" alt="AudioMuse-AI Logo" width="480">
</p>


**AudioMuse-AI-NV-Plugin** the a Navidrome plugin that integrates core AudioMuse-AI features into the Navidrome frontend.

Actually this is the list of integrated functionality:
- Instant Mix - Song similarity
- Radio - Artist Similarity
- Artist Info - It return similar artist

For Mobile app that want to map this functionality they need to implement the `getSimilarSongs2` / `getSimilarSongs` and `getArtistInfo` API.


> **IMPORTANT** InstantMix support in Navidrome start from v0.60.0: https://github.com/navidrome/navidrome/releases/tag/v0.60.0

**The full list or AudioMuse-AI related repository are:** 
  > * [AudioMuse-AI](https://github.com/NeptuneHub/AudioMuse-AI): the core application, it run Flask and Worker containers to actually run all the feature;
  > * [AudioMuse-AI Helm Chart](https://github.com/NeptuneHub/AudioMuse-AI-helm): helm chart for easy installation on Kubernetes;
  > * [AudioMuse-AI Plugin for Jellyfin](https://github.com/NeptuneHub/audiomuse-ai-plugin): Jellyfin Plugin;
  > * [AudioMuse-AI Plugin for Navidrome](https://github.com/NeptuneHub/AudioMuse-AI-NV-plugin): Navidrome Plugin;
  > * [AudioMuse-AI MusicServer](https://github.com/NeptuneHub/AudioMuse-AI-MusicServer): Open Subosnic like Music Sever with integrated sonic functionality.

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
      - ND_AGENTS=audiomuseai,lastfm,spotify,deezer
      - ND_DEVARTISTINFOTIMETOLIVE=1s
    volumes:
      - ./data:/data
      - /path/to/music:/music:ro
```

- Then you need to put `audiomuseai.ndp` in Navidrome data plugins folder (default: `/data/plugins`).
- Restart Navidrome, go to UI -> Plugins, enable **AudioMuse-AI**, set **AudioMuse-AI API URL** and other configuration parameter.

> Note: the audiomuseai.npd can be found attached to the release: https://github.com/NeptuneHub/AudioMuse-AI-NV-plugin/releases

See the official [Navidrome Documentation](https://www.navidrome.org/docs/usage/features/plugins/#managing-plugins-in-the-web-ui) for more information on how the plugin works.

## HOW-TO build

- Requirements (Ubuntu / macOS): Go, TinyGo.
- Build:

```bash
make build    # -> audiomuseai.wasm
make package  # -> audiomuseai.ndp
```

Full stop.
