package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"spec-viewer/internal/server"
	"spec-viewer/internal/socket"
	"spec-viewer/internal/watcher"
	"spec-viewer/pkg/logger"
	"spec-viewer/pkg/ui"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the spec viewer server",
	Long:  `Starts the web server and file watcher to serve your markdown specs.`,
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := os.Stat(folder); os.IsNotExist(err) {
			logger.Fatal("Folder does not exist", "folder", folder, "error", err)
		}

		PrintBanner(port, folder)

		// Create context that listens for the interrupt signal from the OS.
		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		hub := socket.NewHub()

		go watcher.Watch(ctx, folder, hub)

		srv := server.New(hub, server.Config{
			Port:   port,
			Folder: folder,
		})

		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatal("Server failed", "error", err)
			}
		}()

		// Listen for the interrupt signal.
		<-ctx.Done()

		// Restore default behavior on the interrupt signal and notify user of shutdown.
		stop()
		logger.Info("press Ctrl+C again to force")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatal("Server forced to shutdown", "error", err)
		}

		ui.PrintSuccess("Server stopped")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVarP(&port, "port", "p", "9091", "Port to run the server on")
	serveCmd.Flags().StringVarP(&folder, "folder", "f", "./specs", "Folder to watch for specs")
}
