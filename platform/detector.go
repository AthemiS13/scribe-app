package platform

import (
	"runtime"
)

// Platform represents the different operating systems
type Platform int

const (
	Linux Platform = iota
	MacOS
	Windows
	Unknown
)

// String returns the string representation of the platform
func (p Platform) String() string {
	switch p {
	case Linux:
		return "linux"
	case MacOS:
		return "darwin"
	case Windows:
		return "windows"
	default:
		return "unknown"
	}
}

// DetectPlatform returns the current operating system platform
func DetectPlatform() Platform {
	switch runtime.GOOS {
	case "linux":
		return Linux
	case "darwin":
		return MacOS
	case "windows":
		return Windows
	default:
		return Unknown
	}
}

// IsSupported returns true if the platform is supported
func IsSupported(p Platform) bool {
	return p == Linux || p == MacOS || p == Windows
}
