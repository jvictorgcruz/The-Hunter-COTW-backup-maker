package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"backup-maker/internal/onedrive"
)

const graphBaseURL = "https://graph.microsoft.com/v1.0"

type OnedriveProvider struct{}

type graphItem struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	CreatedDateTime time.Time `json:"createdDateTime"`
}

type graphListResponse struct {
	Value []graphItem `json:"value"`
}

func (lp *OnedriveProvider) ID() string {
	return "onedrive"
}

func (lp *OnedriveProvider) IsConfigured() bool {
	return onedrive.IsAuthenticated()
}

func (lp *OnedriveProvider) Send(localZipPath string, filename string) error {
	ctx := context.Background()
	client, err := onedrive.GetClient(ctx)
	if err != nil {
		return err
	}

	file, err := os.Open(localZipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	uploadURL := fmt.Sprintf("%s/me/drive/special/approot:/%s:/content", graphBaseURL, filename)
	req, err := http.NewRequestWithContext(ctx, "PUT", uploadURL, file)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/zip")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("onedrive upload failed with status: %s", resp.Status)
	}

	return nil
}

func (lp *OnedriveProvider) Cleanup(limit int) error {
	ctx := context.Background()
	client, err := onedrive.GetClient(ctx)
	if err != nil {
		return err
	}

	listURL := graphBaseURL + "/me/drive/special/approot/children"
	req, err := http.NewRequestWithContext(ctx, "GET", listURL, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list onedrive files: %s", resp.Status)
	}

	var listResp graphListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return err
	}

	var backupItems []graphItem
	for _, item := range listResp.Value {
		if strings.HasPrefix(item.Name, backupPrefix) && strings.HasSuffix(item.Name, backupSuffix) {
			backupItems = append(backupItems, item)
		}
	}

	if len(backupItems) <= limit {
		return nil
	}

	sort.Slice(backupItems, func(i, j int) bool {
		return backupItems[i].Name < backupItems[j].Name
	})

	itemsToDelete := backupItems[:len(backupItems)-limit]
	for _, item := range itemsToDelete {
		deleteURL := fmt.Sprintf("%s/me/drive/items/%s", graphBaseURL, item.ID)
		delReq, err := http.NewRequestWithContext(ctx, "DELETE", deleteURL, nil)
		if err != nil {
			return err
		}

		delResp, err := client.Do(delReq)
		if err != nil {
			return err
		}
		delResp.Body.Close()

		if delResp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("failed to delete onedrive file %s: %s", item.Name, delResp.Status)
		}
	}

	return nil
}

func init() {
	RegisterProvider(&OnedriveProvider{})
}
