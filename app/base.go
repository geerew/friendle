package app

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/cron"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/filesystem"
	"github.com/geerew/friendle/utils/logger"
	"github.com/geerew/friendle/utils/types"
	"github.com/geerew/friendle/utils/wordgame"
	"github.com/spf13/afero"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AppMode selects runtime behaviour for filesystem, logging, and database setup
type AppMode int

const (
	AppModeProd AppMode = iota
	AppModeDev
	AppModeTest
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Config holds application-wide settings passed to New
type Config struct {
	HttpAddr     string
	DataDir      string
	EnableSignup bool
	Debug        bool
	AppMode      AppMode
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// App wires shared runtime dependencies for the API, cron jobs, and CLI
type App struct {
	Logger       *logger.Logger
	FS           *filesystem.FS
	DbManager    *database.DatabaseManager
	Dictionary   *wordgame.Dictionary
	Cron         *cron.Cron
	Config       *Config
	bootstrapped atomic.Int32
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// New creates the application, runs first-start bootstrap checks, and returns a ready App
func New(ctx context.Context, config *Config) (*App, error) {
	logLevel := logger.LevelInfo
	if config.Debug {
		logLevel = logger.LevelDebug
	}

	var appLogger *logger.Logger
	if config.AppMode == AppModeTest {
		appLogger = logger.NilLogger()
	} else {
		loggerConfig := &logger.LoggerConfig{
			Level:         logLevel,
			ConsoleOutput: true,
		}
		appLogger = logger.New(loggerConfig)
	}

	var fs *filesystem.FS
	if config.AppMode == AppModeTest {
		fs = filesystem.New(afero.NewMemMapFs())
	} else {
		fs = filesystem.New(afero.NewOsFs())
	}

	dbConfig := &database.DatabaseManagerConfig{
		DataDir: config.DataDir,
		FS:      fs,
		Testing: config.AppMode == AppModeTest,
	}

	dbManager, err := database.NewSQLiteManager(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database manager: %w", err)
	}

	dict, err := wordgame.LoadDictionary()
	if err != nil {
		return nil, fmt.Errorf("failed to load dictionary: %w", err)
	}

	cronScheduler := cron.New(&cron.Config{
		DataDb: dbManager.DataDb,
		Logger: appLogger.WithComponent(string(ComponentCron)),
	})

	application := &App{
		Logger:     appLogger,
		FS:         fs,
		DbManager:  dbManager,
		Dictionary: dict,
		Config:     config,
		Cron:       cronScheduler,
	}

	if err := application.bootstrap(); err != nil {
		return nil, err
	}

	return application, nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Close releases application resources
func (a *App) Close() error {
	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// IsBootstrapped reports whether the application has a site admin and completed first-run setup
func (a *App) IsBootstrapped() bool {
	return a.bootstrapped.Load() == 1
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SetBootstrapped marks the application as bootstrapped after the first admin is created
func (a *App) SetBootstrapped() {
	a.bootstrapped.Store(1)
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// RefreshBootstrapped syncs the in-memory bootstrapped flag with the site admin count in the database
func (a *App) RefreshBootstrapped() error {
	count, err := a.siteAdminCount(context.Background())
	if err != nil {
		return err
	}

	if count == 0 {
		a.bootstrapped.Store(0)
	} else {
		a.bootstrapped.Store(1)
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// bootstrap checks for a site admin, generates a bootstrap token when missing, and syncs state
func (a *App) bootstrap() error {
	count, err := a.siteAdminCount(context.Background())
	if err != nil {
		return err
	}

	if count == 0 {
		a.bootstrapped.Store(0)

		bootstrapToken, err := auth.GenerateBootstrapToken(a.Config.DataDir, a.FS)
		if err != nil {
			return fmt.Errorf("failed to generate bootstrap token: %w", err)
		}

		bootstrapURL := fmt.Sprintf("http://%s/auth/bootstrap/%s", a.Config.HttpAddr, bootstrapToken.Token)
		a.Logger.WithComponent(string(ComponentApp)).Info().
			Str("bootstrap_url", bootstrapURL).
			Str("expires_in", "5 minutes").
			Msg("Bootstrap required")
	} else {
		a.bootstrapped.Store(1)
		_ = auth.DeleteBootstrapToken(a.Config.DataDir, a.FS)
		a.Logger.WithComponent(string(ComponentApp)).Info().Msg("Application bootstrapped")
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// siteAdminCount returns how many site admins exist in the database
func (a *App) siteAdminCount(ctx context.Context) (int, error) {
	appDao := dao.New(a.DbManager.DataDb)
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_SITE_ROLE: types.SiteRoleAdmin})

	count, err := appDao.CountUsers(ctx, dbOpts)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin users: %w", err)
	}

	return count, nil
}
