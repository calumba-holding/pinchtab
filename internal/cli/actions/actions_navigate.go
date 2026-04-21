package actions

import (
	"fmt"
	"net/http"

	"github.com/pinchtab/pinchtab/internal/cli/apiclient"
	"github.com/pinchtab/pinchtab/internal/cli/output"
	"github.com/spf13/cobra"
)

// Back navigates the current (or specified) tab back in history.
func Back(client *http.Client, base, token string, cmd *cobra.Command) {
	tabID, _ := cmd.Flags().GetString("tab")
	path := "/back"
	if tabID != "" {
		path = "/tabs/" + tabID + "/back"
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		apiclient.DoPost(client, base, token, path, nil)
		return
	}

	result := apiclient.DoPostQuiet(client, base, token, path, nil)
	if url, ok := result["url"].(string); ok {
		output.Value(url)
	} else {
		output.Success()
	}

	snap, _ := cmd.Flags().GetBool("snap")
	snapDiff, _ := cmd.Flags().GetBool("snap-diff")
	if snap || snapDiff {
		fetchAndPrintSnapshot(client, base, token, tabID, snapDiff)
	}
	text, _ := cmd.Flags().GetBool("text")
	if text {
		fetchAndPrintText(client, base, token, tabID)
	}
}

// Forward navigates the current (or specified) tab forward in history.
func Forward(client *http.Client, base, token string, cmd *cobra.Command) {
	tabID, _ := cmd.Flags().GetString("tab")
	path := "/forward"
	if tabID != "" {
		path = "/tabs/" + tabID + "/forward"
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		apiclient.DoPost(client, base, token, path, nil)
		return
	}

	result := apiclient.DoPostQuiet(client, base, token, path, nil)
	if url, ok := result["url"].(string); ok {
		output.Value(url)
	} else {
		output.Success()
	}

	snap, _ := cmd.Flags().GetBool("snap")
	snapDiff, _ := cmd.Flags().GetBool("snap-diff")
	if snap || snapDiff {
		fetchAndPrintSnapshot(client, base, token, tabID, snapDiff)
	}
	text, _ := cmd.Flags().GetBool("text")
	if text {
		fetchAndPrintText(client, base, token, tabID)
	}
}

// Reload reloads the current (or specified) tab.
func Reload(client *http.Client, base, token string, cmd *cobra.Command) {
	tabID, _ := cmd.Flags().GetString("tab")
	path := "/reload"
	if tabID != "" {
		path = "/tabs/" + tabID + "/reload"
	}

	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		apiclient.DoPost(client, base, token, path, nil)
		return
	}

	apiclient.DoPostQuiet(client, base, token, path, nil)
	output.Success()

	snap, _ := cmd.Flags().GetBool("snap")
	snapDiff, _ := cmd.Flags().GetBool("snap-diff")
	if snap || snapDiff {
		fetchAndPrintSnapshot(client, base, token, tabID, snapDiff)
	}
	text, _ := cmd.Flags().GetBool("text")
	if text {
		fetchAndPrintText(client, base, token, tabID)
	}
}

func Navigate(client *http.Client, base, token string, url string, cmd *cobra.Command) string {
	body := map[string]any{"url": url}
	newTab, _ := cmd.Flags().GetBool("new-tab")
	if newTab {
		body["newTab"] = true
	}
	if v, _ := cmd.Flags().GetBool("block-images"); v {
		body["blockImages"] = true
	}
	if v, _ := cmd.Flags().GetBool("block-ads"); v {
		body["blockAds"] = true
	}
	tabID, _ := cmd.Flags().GetString("tab")
	path := "/navigate"
	// Don't use tab-specific path when creating a new tab
	if tabID != "" && !newTab {
		path = "/tabs/" + tabID + "/navigate"
	}

	// JSON output mode
	jsonOutput, _ := cmd.Flags().GetBool("json")
	if jsonOutput {
		result := apiclient.DoPost(client, base, token, path, body)
		apiclient.SuggestNextAction("navigate", result)
		return ""
	}

	// Terse mode: emit only the tabId (default, or when --print-tab-id / pipe)
	result := apiclient.DoPostQuiet(client, base, token, path, body)
	resultTabID := ""
	if tid, ok := result["tabId"].(string); ok && tid != "" {
		resultTabID = tid
		fmt.Println(tid)
	}

	// If --snap or --snap-diff flag is set, fetch and output snapshot
	snap, _ := cmd.Flags().GetBool("snap")
	snapDiff, _ := cmd.Flags().GetBool("snap-diff")
	if snap || snapDiff {
		fetchAndPrintSnapshot(client, base, token, resultTabID, snapDiff)
	}

	return resultTabID
}
