// Command clip-and-gif downloads a video (or an exact range from one)
// from any yt-dlp-supported URL, as video, audio, or a high-quality GIF —
// combining video-clipper and gif-maker's engines, imported as real Go
// module dependencies, into one app. Bare invocation opens a local
// browser UI.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	gifengine "github.com/DavidMarsanic/gif-maker/engine"
	videoengine "github.com/DavidMarsanic/video-clipper/engine"

	"github.com/DavidMarsanic/clip-and-gif/internal/browser"
	"github.com/DavidMarsanic/clip-and-gif/internal/paths"
	"github.com/DavidMarsanic/clip-and-gif/internal/server"
)

const version = "0.1.0"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	fs := flag.NewFlagSet("clip-and-gif", flag.ContinueOnError)

	urlFlag := fs.String("url", "", "video URL to preload")
	output := fs.String("output", "", "output directory (default: your Downloads folder)")
	port := fs.Int("port", 0, "local UI server port (default: automatic)")
	showVersion := fs.Bool("version", false, "print the version and exit")
	fs.Usage = func() { printUsage(fs) }

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 1
	}

	if *showVersion {
		fmt.Println("clip-and-gif " + version)
		return 0
	}

	widenPATH()

	outputDir, err := paths.ResolveDownloadsDir(*output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}

	videoEng := videoengine.New(outputDir)
	gifEng := gifengine.New()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(ctx, videoEng, gifEng, outputDir)
	addr, err := srv.Start(*port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return 1
	}

	target := addr + "/"
	if *urlFlag != "" {
		target += "?url=" + url.QueryEscape(*urlFlag)
	}

	fmt.Fprintln(os.Stderr, "Clip & GIF running at", addr, "— press Ctrl+C to quit")

	// When a host process (securexe-launcher) is the one showing the UI —
	// in its own native window, so it can get a real Dock identity instead
	// of a spawned Chrome window — it sets this before starting us and
	// watches this same stderr line to discover the URL. Opening our own
	// Chrome window too would just leave a second, redundant one.
	if os.Getenv("SECUREXE_HOSTED") == "" {
		if err := browser.OpenAppWindow(target); err != nil {
			fmt.Fprintln(os.Stderr, "couldn't open a window automatically:", err)
			fmt.Fprintln(os.Stderr, "open this URL manually:", target)
		}
	}

	<-ctx.Done()
	return 0
}

// widenPATH adds common tool-install directories that a GUI-launched
// process often lacks — see the identical function in video-clipper for
// the full explanation. Duplicated rather than imported: it's a
// process-startup concern specific to this binary, not shared domain
// logic, so it doesn't belong in either imported engine package.
func widenPATH() {
	home, _ := os.UserHomeDir()
	candidates := []string{
		"/opt/homebrew/bin", "/opt/homebrew/sbin",
		"/usr/local/bin", "/usr/local/sbin",
		filepath.Join(home, ".local", "bin"),
	}

	current := os.Getenv("PATH")
	existing := map[string]bool{}
	for _, p := range filepath.SplitList(current) {
		existing[p] = true
	}

	var toAdd []string
	for _, dir := range candidates {
		if dir == "" || existing[dir] {
			continue
		}
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			toAdd = append(toAdd, dir)
		}
	}
	if len(toAdd) == 0 {
		return
	}
	toAdd = append(toAdd, current)
	os.Setenv("PATH", strings.Join(toAdd, string(os.PathListSeparator)))
}

func printUsage(fs *flag.FlagSet) {
	fmt.Fprint(os.Stderr, `clip-and-gif — download a video, or an exact range from one, as video,
audio, or a high-quality GIF, from any yt-dlp-supported URL.

Bare invocation opens a local browser UI: paste a URL, preview it, drag a
timeline to pick a range, and save as video, audio, or GIF.

Usage:
  clip-and-gif                open the browser UI
  clip-and-gif -url <url>      open the UI with that URL preloaded

Flags:
`)
	fs.PrintDefaults()
	fmt.Fprint(os.Stderr, `
Requires yt-dlp, ffmpeg, and gifski on PATH. None are bundled.
`)
}
