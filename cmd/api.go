package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"pingo/core/api"
	"pingo/core/app"
	"pingo/core/config"

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

		var (
			application *app.Application
			appMutex    sync.RWMutex
			appOnce     sync.Once
			appErr      error
			appReady    = make(chan struct{})
		)

		// Inicializa a aplicação em background
		go func() {
			appOnce.Do(func() {
				cfg, err := config.LoadConfig(configFile)
				if err != nil {
					appErr = fmt.Errorf("error loading config: %v", err)
					close(appReady)
					return
				}

				newApp, err := app.NewApplication(cfg)
				if err != nil {
					appErr = fmt.Errorf("error creating application: %v", err)
					close(appReady)
					return
				}

				appMutex.Lock()
				application = newApp
				appMutex.Unlock()
				close(appReady)

				// Inicia a aplicação (roda workers)
				if err := newApp.Run(); err != nil {
					log.Printf("Application error: %v", err)
					os.Exit(1)
				}
			})
		}()

		// Aguarda a aplicação estar pronta
		<-appReady
		if appErr != nil {
			log.Printf("%v", appErr)
			os.Exit(1)
		}

		config.SetOnUpdateHandler(func() {
			appMutex.RLock()
			app := application
			appMutex.RUnlock()

			if app == nil {
				log.Println("Warning: Application not ready for config update")
				return
			}

			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				log.Printf("Error reloading config: %v", err)
				return
			}
			app.UpdateConfig(cfg)
		})

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down...")

		// Shutdown do servidor HTTP
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)

		// Notifica a aplicação para encerrar
		appMutex.RLock()
		app := application
		appMutex.RUnlock()

		if app != nil {
			log.Println("Stopping application workers...")
		}
	},
}

func init() {
	apiCmd.Flags().StringVar(&apiAddr, "addr", ":8080", "HTTP listen address")
	rootCmd.AddCommand(apiCmd)
}
