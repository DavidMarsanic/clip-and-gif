# Clip & GIF

Download a full video, an exact clip, or a high-quality GIF from any
yt-dlp-supported URL — paste a link, drag a timeline, save. Opens as its
own window.

This app doesn't reimplement anything: it imports
[video-clipper](https://github.com/DavidMarsanic/video-clipper)'s engine
(yt-dlp/ffmpeg — fetch, clip, extract audio) and
[gif-maker](https://github.com/DavidMarsanic/gif-maker)'s engine
(ffmpeg | gifski) as real Go module dependencies and composes them: fetch
the requested range with video-clipper's engine, then hand the resulting
local file to gif-maker's engine for GIF conversion. If you only need
video/audio, use video-clipper directly; if you only need to GIF a video
file you already have, use gif-maker directly. This app is for when you
want a URL to GIF in one step.

## Requirements

Three external tools, all expected on `PATH`, none bundled:

- [`yt-dlp`](https://github.com/yt-dlp/yt-dlp#installation) — extraction
- [`ffmpeg`](https://ffmpeg.org/download.html) — trimming, remuxing, encoding, GIF decoding
- [`gifski`](https://github.com/ImageOptim/gifski#download-and-install) — GIF encoding

On macOS: `brew install yt-dlp ffmpeg gifski`.

**A Chromium-based browser already installed**: Google Chrome, Chromium,
Brave, Microsoft Edge, or Arc — renders the app's own UI window.

If anything is missing, the app still opens — it'll tell you what's
missing the moment you try to use it.

## Use

Paste a URL, wait for the preview to load, then either download the whole
thing or drag the timeline to pick a range. Pick a format — video,
audio-only, or GIF — and for GIF, a width (the single biggest lever for
file size). Save.

## License

MIT — see [LICENSE](LICENSE).
