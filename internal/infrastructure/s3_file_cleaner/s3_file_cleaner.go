package s3filecleaner

import (
	"context"
	"fmt"
	"log"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/storages"
	"github.com/jackc/pgx/v5"
)

type s3Cleaner struct {
	conn    *pgx.Conn
	storage storages.FileStorage
}

func NewS3Cleaner(conn *pgx.Conn, storage storages.FileStorage) *s3Cleaner {
	return &s3Cleaner{
		conn:    conn,
		storage: storage,
	}
}

func (s3C *s3Cleaner) ListenForFileMetaEvents(ctx context.Context) error {
	_, err := s3C.conn.Exec(ctx, "LISTEN filemeta_events")
	if err != nil {
		return fmt.Errorf("failed to listen to channel: %v", err)
	}
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping listener...")
			return nil
		default:
			notification, err := s3C.conn.WaitForNotification(ctx)
			if err != nil {
				return fmt.Errorf("failed to wait for notification: %w", err)
			}

			uri := notification.Payload

			if uri == "" {
				log.Println("Received empty payload")
				continue
			}

			log.Printf("Received S3 URI: %s", uri)

			if err := s3C.storage.DeleteFile(uri); err != nil {
				log.Printf("Failed to delete from S3 (%s): %v", uri, err)
			} else {
				log.Printf("Successfully deleted from S3: %s", uri)
			}
		}
	}
}
