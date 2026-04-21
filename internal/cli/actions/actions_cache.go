package actions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pinchtab/pinchtab/internal/cli/apiclient"
	"github.com/spf13/cobra"
)

// CacheClear clears the browser's HTTP disk cache.
func CacheClear(client *http.Client, base, token string, cmd *cobra.Command) {
	result := apiclient.DoPostQuiet(client, base, token, "/cache/clear", nil)
	if result == nil {
		fmt.Fprintln(os.Stderr, "ERROR: cache: clear failed")
		os.Exit(2)
	}

	jsonOut, _ := cmd.Flags().GetBool("json")
	if jsonOut {
		out, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(out))
	} else {
		fmt.Println("OK")
	}
}

// CacheStatus checks if the browser cache can be cleared.
func CacheStatus(client *http.Client, base, token string, cmd *cobra.Command) {
	result := apiclient.DoGetRaw(client, base, token, "/cache/status", nil)
	if result == nil {
		fmt.Fprintln(os.Stderr, "ERROR: cache: status check failed")
		os.Exit(2)
	}

	var buf map[string]any
	if err := json.Unmarshal(result, &buf); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cache: %v\n", err)
		os.Exit(2)
	}

	jsonOut, _ := cmd.Flags().GetBool("json")
	if jsonOut {
		out, _ := json.MarshalIndent(buf, "", "  ")
		fmt.Println(string(out))
	} else {
		if canClear, ok := buf["canClear"].(bool); ok && canClear {
			fmt.Println("can-clear")
		} else {
			fmt.Println("cache-empty")
		}
	}
}
