package app

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

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
	"github.com/stretchr/testify/require"
)

type App struct {
	Logger     *logger.Logger
	FS         *filesystem.FS
	DbManager  *database.DatabaseManager
	Dictionary *wordgame.Dictionary
	Cron       *cron.Cron
	Config     *Config
	bootstrapped atomic.Int32
}

type AppMode int

const (
	AppModeProd AppMode = iota
	AppModeDev
	AppModeTest
)

type Config struct {
	HttpAddr     string
	DataDir      string
	EnableSignup bool
	Debug        bool
	AppMode      AppMode
}

func NewApp(ctx context.Context, config *Config) (*App, error) {
	logLevel := logger.LevelInfo
	if config.Debug {
		logLevel = logger.LevelDebug
	}

	var appLogger *logger.Logger
	if config.AppMode == AppModeTest {
		appLogger = logger.NilLogger()
	} else {
		appLogger = logger.New(&logger.LoggerConfig{
			Level:         logLevel,
			ConsoleOutput: true,
		})
	}

	var fs *filesystem.FS
	if config.AppMode == AppModeTest {
		fs = filesystem.New(afero.NewMemMapFs())
	} else {
		fs = filesystem.New(afero.NewOsFs())
	}

	dbManager, err := database.NewSQLiteManager(&database.DatabaseManagerConfig{
		DataDir: config.DataDir,
		FS:      fs,
		Testing: config.AppMode == AppModeTest,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create database manager: %w", err)
	}

	dict, err := wordgame.LoadDictionary()
	if err != nil {
		return nil, fmt.Errorf("failed to load dictionary: %w", err)
	}

	app := &App{
		Logger:     appLogger,
		FS:         fs,
		DbManager:  dbManager,
		Dictionary: dict,
		Config:     config,
		Cron: cron.NewCronScheduler(&cron.CronConfig{
			DataDb: dbManager.DataDb,
			FS:     fs,
			DataDir: config.DataDir,
			Logger: appLogger.WithComponent(string(ComponentCron)),
		}),
	}

	if err := app.bootstrap(); err != nil {
		return nil, err
	}

	return app, nil
}

func NewTestApp(t *testing.T) *App {
	t.Helper()
	app, err := NewApp(context.Background(), &Config{
		HttpAddr:     "127.0.0.1:9081",
		DataDir:      "./friendle_data",
		AppMode:      AppModeTest,
		EnableSignup: true,
	})
	require.NoError(t, err)
	return app
}

func (a *App) Close() error { return nil }

func (a *App) IsBootstrapped() bool { return a.bootstrapped.Load() == 1 }
func (a *App) SetBootstrapped()     { a.bootstrapped.Store(1) }
func (a *App) UnsetBootstrapped()   { a.bootstrapped.Store(0) }

func (a *App) bootstrap() error {
	appDao := dao.New(a.DbManager.DataDb)
	count, err := appDao.CountUsers(
		context.Background(),
		dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_SITE_ROLE: types.SiteRoleAdmin}),
	)
	if err != nil {
		return fmt.Errorf("failed to count admin users: %w", err)
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
