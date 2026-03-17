package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pingo/core/api"

	"github.com/spf13/cobra"
)

var apiAddr string

var apiCmd = &cobra.Command{
	Use:   "api [config-file]",
	Short: "Start the HTTP API to manage config items",
	Long:  "Start an HTTP API server that can add, update and remove items from the JSON config file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configFile := args[0]

		store := api.NewConfigStore(configFile)
		if _, err := store.Load(); err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		server := &http.Server{
			Addr:    apiAddr,
			Handler: api.NewHandler(store),
		}

		go func() {
			log.Printf("API listening on %s", apiAddr)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("API server error: %v", err)
				os.Exit(1)
			}
		}()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	},
}

func init() {
	apiCmd.Flags().StringVar(&apiAddr, "addr", ":8080", "HTTP listen address")
	rootCmd.AddCommand(apiCmd)
}
