package cmd

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/geerew/friendle/api"
	"github.com/geerew/friendle/app"
	"github.com/geerew/friendle/cron"
	"github.com/geerew/friendle/service"
	"github.com/geerew/friendle/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// serveCmd starts the application
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the application",
	Run: func(cmd *cobra.Command, args []string) {
		runCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		httpAddr := viper.GetString("http")
		dataDir := viper.GetString("data-dir")
		enableSignup := viper.GetBool("enable-signup")
		isDev := viper.GetBool("dev")
		debug := viper.GetBool("debug")

		appMode := app.AppModeProd
		if isDev {
			appMode = app.AppModeDev
		}

		appConfig := &app.Config{
			HttpAddr:     httpAddr,
			DataDir:      dataDir,
			AppMode:      appMode,
			EnableSignup: enableSignup,
			Debug:        debug,
		}

		application, err := app.New(runCtx, appConfig)
		if err != nil {
			os.Stderr.WriteString("Failed to initialize app: " + err.Error() + "\n")
			os.Exit(1)
		}

		appLogger := application.Logger.WithComponent(string(app.ComponentApp))

		appLogger.Info().
			Str("version", version.GetVersion()).
			Str("commit", version.GetCommit()).
			Msg("Starting Friendle")

		// Router
		router := api.New(application, nil)
		appSvc := service.New(application.DbManager.DataDb)

		// Cron
		cron.NewAndStart(runCtx, &cron.Config{
			Rounds:           appSvc.Rounds,
			DailyRoundLogger: application.Logger.WithComponent(string(app.ComponentCron)),
		})

		var wg sync.WaitGroup
		wg.Add(2)

		// Listen for shutdown signals
		go func() {
			defer wg.Done()
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
			<-quit
			cancel()
		}()

		go func() {
			defer wg.Done()
			if err := router.Serve(); err != nil {
				appLogger.Error().Err(err).Msg("Failed to start router")
				os.Exit(1)
			}
		}()

		wg.Wait()

		appLogger.Info().Msg("Shutting down...")

		if err := application.Close(); err != nil {
			appLogger.Error().Err(err).Msg("Failed to close application resources")
		}

	},
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// init adds the serve command to the root command
func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolP("dev", "d", false, "Run in development mode")
	serveCmd.Flags().String("http", "127.0.0.1:9081", "TCP address to listen for the HTTP server")
	serveCmd.Flags().String("data-dir", "./friendle_data", "Directory to store data files")
	serveCmd.Flags().Bool("enable-signup", false, "Allow users to create new accounts")
	serveCmd.Flags().Bool("debug", false, "Enable debug logging")

	_ = viper.BindPFlag("dev", serveCmd.Flags().Lookup("dev"))
	_ = viper.BindPFlag("http", serveCmd.Flags().Lookup("http"))
	_ = viper.BindPFlag("data-dir", serveCmd.Flags().Lookup("data-dir"))
	_ = viper.BindPFlag("enable-signup", serveCmd.Flags().Lookup("enable-signup"))
	_ = viper.BindPFlag("debug", serveCmd.Flags().Lookup("debug"))
}
