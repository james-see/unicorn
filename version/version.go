package version

// Version is set at build time via -ldflags (e.g. -X github.com/jamesacampbell/unicorn/version.Version=v3.31.1).
// Default matches latest CHANGELOG release when building locally.
var Version = "3.33.0"

// ReleaseURL is the base URL for release info and downloads.
const ReleaseURL = "https://github.com/james-see/unicorn/releases"

// String returns the version string, ensuring a leading 'v' for display.
func String() string {
	if Version == "" {
		return "dev"
	}
	if len(Version) > 0 && Version[0] != 'v' {
		return "v" + Version
	}
	return Version
}

// ReleaseInfoURL returns the URL to the release page for the current version (latest for tagged builds).
func ReleaseInfoURL() string {
	return ReleaseURL + "/latest"
}
