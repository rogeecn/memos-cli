package cli

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/rogeecn/memos-cli/internal/config"
	"github.com/rogeecn/memos-cli/internal/input"
	"github.com/rogeecn/memos-cli/internal/memos"
	"github.com/rogeecn/memos-cli/internal/output"
	"github.com/spf13/cobra"
)

func newMemoCreateCommand() *cobra.Command {
	var visibility string
	var tags []string
	var images []string

	cmd := &cobra.Command{
		Use:   "create <content>",
		Short: "Create a memo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.LoadFromEnv()
			if err := cfg.Validate(); err != nil {
				return err
			}
			client := memos.NewClient(cfg.BaseURL, cfg.APIKey, cfg.AdminAPIKey)
			allTags := append([]string{}, tags...)
			if cfg.DefaultTag != "" {
				allTags = append(allTags, cfg.DefaultTag)
			}
			payload := memos.MemoPayload{
				Content:    input.MergeTags(args[0], allTags),
				Visibility: visibility,
			}
			memo, err := client.CreateMemo(payload)
			if err != nil {
				return err
			}
			if err := uploadMemoImages(client, memo.Name, images); err != nil {
				return err
			}
			if getCommandContext(cmd.Context()).jsonOutput {
				return output.WriteJSON(cmd.OutOrStdout(), memo)
			}
			return output.WriteMemoDetail(cmd.OutOrStdout(), memo)
		},
	}
	cmd.Flags().StringVar(&visibility, "visibility", "PRIVATE", "Memo visibility")
	cmd.Flags().StringSliceVar(&tags, "tag", nil, "Tags to append")
	cmd.Flags().StringSliceVar(&images, "image", nil, "Local image paths to upload")
	return cmd
}

func uploadMemoImages(client *memos.Client, memoName string, imagePaths []string) error {
	if len(imagePaths) == 0 {
		return nil
	}

	attachments := make([]memos.Attachment, 0, len(imagePaths))
	for _, imagePath := range imagePaths {
		content, err := os.ReadFile(imagePath)
		if err != nil {
			return fmt.Errorf("read image %q: %w", imagePath, err)
		}
		filename := filepath.Base(imagePath)
		mimeType := mime.TypeByExtension(filepath.Ext(filename))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
		attachment, err := client.CreateAttachment(memos.Attachment{
			Filename: filename,
			Content:  content,
			Type:     mimeType,
			Memo:     memoName,
		})
		if err != nil {
			return err
		}
		attachments = append(attachments, attachment)
	}

	return client.SetMemoAttachments(strings.TrimPrefix(memoName, "memos/"), memos.SetMemoAttachmentsPayload{
		Name:        memoName,
		Attachments: attachments,
	})
}

func newMemoUpdateCommand() *cobra.Command {
	var content string
	var visibility string
	cmd := &cobra.Command{
		Use:   "update <memo-id>",
		Short: "Update a memo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if content == "" && visibility == "" {
				return errors.New("update requires --content or --visibility")
			}
			client, err := loadClientFromEnv()
			if err != nil {
				return err
			}
			memo, err := client.UpdateMemo(args[0], memos.UpdateMemoPayload{Content: content, Visibility: visibility})
			if err != nil {
				return err
			}
			if getCommandContext(cmd.Context()).jsonOutput {
				return output.WriteJSON(cmd.OutOrStdout(), memo)
			}
			return output.WriteMemoDetail(cmd.OutOrStdout(), memo)
		},
	}
	cmd.Flags().StringVar(&content, "content", "", "Updated memo content")
	cmd.Flags().StringVar(&visibility, "visibility", "", "Updated memo visibility")
	return cmd
}

func newMemoDeleteCommand() *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "delete <memo-id>",
		Short: "Delete a memo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				return errors.New("delete requires --yes")
			}
			client, err := loadClientFromEnv()
			if err != nil {
				return err
			}
			if err := client.DeleteMemo(args[0]); err != nil {
				return err
			}
			if getCommandContext(cmd.Context()).jsonOutput {
				return output.WriteJSON(cmd.OutOrStdout(), map[string]any{"deleted": true, "memo": args[0]})
			}
			_, err = cmd.OutOrStdout().Write([]byte("deleted\n"))
			return err
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm delete")
	return cmd
}
