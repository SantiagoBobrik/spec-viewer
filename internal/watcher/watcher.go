package watcher

import (
	"context"
	"os"
	"path/filepath"
	"spec-viewer/internal/socket"
	"spec-viewer/pkg/logger"
	"spec-viewer/pkg/ui"

	"github.com/fsnotify/fsnotify"
)

func Watch(ctx context.Context, root string, hub *socket.Hub) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal("Failed to create watcher", "error", err)
	}
	defer watcher.Close()

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		logger.Fatal("Error walking directory", "error", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					logger.Info("Watching new directory", "path", event.Name)
					watcher.Add(event.Name)
				}
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
				ui.PrintFileChange(event.Name)
				hub.Broadcast(socket.Events.Reload)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Error("Watcher error", "error", err)
		}
	}
}
