![GitHub license](https://img.shields.io/github/license/neptunehub/AudioMuse-AI-NV-plugin.svg)
![Latest Tag](https://img.shields.io/github/v/tag/neptunehub/AudioMuse-AI-NV-plugin?label=latest-tag)
![Media Server Support: Navidrome 0.61.0](https://img.shields.io/badge/Media%20Server-Navidrome%200.61.0-blue?style=flat-square&logo=server&logoColor=white)

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

**Front-End tested with the plugin are:**
- Navidrome integrated web frontend
- [Substreamer](https://github.com/ghenry22/substreamer) - iOS/Android opensource mobile frontend, more information here: https://github.com/ghenry22/substreamer/issues/58
- [Tempus](https://github.com/eddyizm/tempus) - Android opensource mobile frontend, more information here: https://github.com/eddyizm/tempus/issues/410
- [Symfonium](https://symfonium.app/) - Androind closed source mobile frontend. You need to enable in the configuration `Use similar tracks for Radio Mix`. Additional functionality that use AudioMuse-AI API are actually implemented only using Jellyfin. More information here: https://support.symfonium.app/t/implement-navidrome-audiomuse-ai-plugin-to-symfonium/12238/12
- [Feishin](https://github.com/jeffvli/feishin/issues/1675) - Web opensource frontend, more information here: https://github.com/jeffvli/feishin/issues/1675

Other frontnend not in this list could also work by using those API even.

**NEWS**
> * From AudioMuse-AI-NV-plugin v8 the Sonic Similarity API extension are supported on top of the previous one. The API are [documented here](https://opensubsonic.netlify.app/docs/extensions/sonicsimilarity/) and supported by [this Navidrome PR](https://github.com/navidrome/navidrome/pull/5419)
> * InstantMix support in Navidrome start from v0.60.0: https://github.com/navidrome/navidrome/releases/tag/v0.60.0


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
