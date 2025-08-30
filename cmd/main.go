package cmd

import (
	"fmt"
	"os"
)

const version = "0.0.1"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "mount":
		handleMount()
	case "status":
		handleStatus()
	case "sync":
		handleSync()
	case "pull":
		handlePull()
	case "--version", "-v":
		fmt.Println("StorachaFS version:", version)
	default:
		fmt.Println("Unknown command:", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println("StorachaFS CLI")
	fmt.Println("Usage: storachafs <command>")
	fmt.Println("Commands:")
	fmt.Println("  mount   Mount a Storacha space at a local directory")
	fmt.Println("  status  Show local vs remote changes")
	fmt.Println("  sync    Sync local files to Storacha")
	fmt.Println("  pull    Fetch new or updated files from Storacha")
	fmt.Println("Options:")
	fmt.Println("  -v, --version   Show version")
}

func handleMount() {
	fmt.Println("[stub] Mounting Storacha space... (to be implemented)")
}

func handleStatus() {
	fmt.Println("[stub] Checking sync status... (to be implemented)")
}

func handleSync() {
	fmt.Println("[stub] Syncing local files to Storacha... (to be implemented)")
}

func handlePull() {
	fmt.Println("[stub] Pulling new/updated files from Storacha... (to be implemented)")
}
