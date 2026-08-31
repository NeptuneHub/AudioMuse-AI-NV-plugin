![GitHub license](https://img.shields.io/github/license/neptunehub/AudioMuse-AI-NV-plugin.svg)
![Latest Tag](https://img.shields.io/github/v/tag/neptunehub/AudioMuse-AI-NV-plugin?label=latest-tag)
![Media Server Support: Navidrome 0.62.0](https://img.shields.io/badge/Media%20Server-Navidrome%200.62.0-blue?style=flat-square&logo=server&logoColor=white)
<a href="https://liberapay.com/NeptuneHub/donate"><img alt="Donate using Liberapay" src="https://liberapay.com/assets/widgets/donate.svg"></a>

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
- [Wavio](https://wavio-app.vercel.app) - Android opensource mobile frontend, more information here: https://github.com/Joel-Mercier/wavio

Other frontnend not in this list could also work by using those API even.

**NEWS**
> * From AudioMuse-AI-NV-plugin v10 different AudioMuse-AI functionality can be selected for the InstantMix, chose it in the plugin configuration
> * From AudioMuse-AI-NV-plugin v8 the Sonic Similarity API extension are supported on top of the previous one. The API are [documented here](https://opensubsonic.netlify.app/docs/extensions/sonicsimilarity/) and supported by [this Navidrome PR](https://github.com/navidrome/navidrome/pull/5419)
> * InstantMix support in Navidrome start from v0.60.0: https://github.com/navidrome/navidrome/releases/tag/v0.60.0


**The full list or AudioMuse-AI related repository are:** 
  > * [AudioMuse-AI](https://github.com/NeptuneHub/AudioMuse-AI): the core application, it run Flask and Worker containers to actually run all the feature;
  > * [AudioMuse-AI Helm Chart](https://github.com/NeptuneHub/AudioMuse-AI-helm): helm chart for easy installation on Kubernetes;
  > * [AudioMuse-AI Plugin for Jellyfin](https://github.com/NeptuneHub/audiomuse-ai-plugin): Jellyfin Plugin;
  > * [AudioMuse-AI Plugin for Navidrome](https://github.com/NeptuneHub/AudioMuse-AI-NV-plugin): Navidrome Plugin;
  > * [AudioMuse-AI MusicServer](https://github.com/NeptuneHub/AudioMuse-AI-MusicServer): Open Subsonic-like Music Server with integrated sonic functionality.

## HOW-TO Install

> IMPORTANT: Before start we suggest to have the last version of Navidrome, AudioMuse-AI-NV-plugin and AudioMuse-AI core containers. If the version are not aligned, some error could happen. Keep this in mind also for future update.

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
      - ND_AGENTS=audiomuseai,lastfm,deezer
      - ND_DEVARTISTINFOTIMETOLIVE=1s
    volumes:
      - ./data:/data
      - /path/to/music:/music:ro
```

- Then you need to put `audiomuseai.ndp` in Navidrome data plugins folder (default: `/data/plugins`).
- Restart Navidrome, go to UI -> Plugins, enable **AudioMuse-AI**, set **AudioMuse-AI API URL** and other configuration parameter.
- **The order of ND_AGENTS is important.** Navidrome will use the first listed agent supporting sonic similarity.

> [!IMPORTANT]
> - The audiomuseai.npd can be found attached to the release: https://github.com/NeptuneHub/AudioMuse-AI-NV-plugin/releases.
> - If you had configured authentication on AudioMuse-AI, you should also to create an 'API Token' on AudioMuse-AI core container and put in the plugin configuration.
> - If AudioMuse-Ai has multiple Music Server configured, you need to configure the 'Server Name' in the plugin to match the one in AudioMuse-AI
> - The **Instant Mix Functionality** select choose which AudioMuse-AI engine answer Instant Mix and Sonic Similarity:
>   - **Similar Song** (default) - `GET /api/similar_tracks`, the historical behaviour;
>   - **Lyrics by Song** - `POST /api/sem_grove/search`, it need the SemGrove index (lyrics + audio analysis) built in AudioMuse-AI;
>   - **Hyperbolic Similarity** - `POST /api/hyperbolic/similar` in its standard mode, it need the hyperbolic projection built in AudioMuse-AI.
>
>   'Eliminate Duplicates' and 'Use Radius Similarity' apply only to **Similar Song**, the other two engine use the AudioMuse-AI default.

See the official [Navidrome Documentation](https://www.navidrome.org/docs/usage/features/plugins/#managing-plugins-in-the-web-ui) for more information on how the plugin works.

## HOW-TO test

In the AudioMuse-AI integrated UI, use the **Similar Song** functionality on a track.

In the Navidrome Web UI, use **InstantMix** on the same track.

Both actions should generate similar requests in the AudioMuse-AI Flask logs. For the Navidrome request, you should also see `plugin=audiomuseai` in multiple line of the Navidrome logs with no related errors.

If both tests pass, supported third-party clients should work correctly. Any issue on third party client MUST be addressed on their side.

> [!NOTE]
> - Different frontends may return a different number of similar songs, depending on their configuration. This plugin does not enforce a default limit.
> - A 401 Unauthorized error in the Navidrome logs typically indicates that the apiToken is missing or incorrectly configured.
> - Navidrome must be able to reach the AudioMuse-AI core container. Ensure the plugin is configured with the correct host/IP address. If requests do not appear in the AudioMuse-AI Flask logs, the issue is most likely network connectivity or container routing.

## HOW-TO build

- Requirements (Ubuntu / macOS): Go, TinyGo.
- Build:

```bash
make build    # -> audiomuseai.wasm
make package  # -> audiomuseai.ndp
```

Full stop.

## **Code Mirror**

[AudioMuse-AI-NV-plugin](https://github.com/NeptuneHub/AudioMuse-AI-NV-plugin) repository code is mirrored here:
- https://codeberg.org/NeptuneHub/AudioMuse-AI-NV-plugin
  
DO **NOT** USE MIRROR TO RAISE ISSUE, PR OTHER ACTION DIFFERENT FROM GET THE CODE
