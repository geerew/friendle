package cmd

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/dao"
	"github.com/geerew/friendle/database"
	"github.com/geerew/friendle/models"
	"github.com/geerew/friendle/utils/auth"
	"github.com/geerew/friendle/utils/filesystem"
	"github.com/geerew/friendle/utils/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// adminResetPasswordCmd resets the password for an admin user
var adminResetPasswordCmd = &cobra.Command{
	Use:   "reset-password <username>",
	Short: "Reset password for an admin user",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := runAdminResetPassword(args[0]); err != nil {
			errorMessage("%s", err)
			os.Exit(1)
		}
	},
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// runAdminResetPassword prompts for a new password and applies it via the recovery API
func runAdminResetPassword(username string) error {
	fmt.Println()
	fmt.Println("Admin Password Reset")
	fmt.Println("====================")
	fmt.Println()

	dataDir := viper.GetString("data-dir")
	httpAddr := viper.GetString("http")
	fs := filesystem.New(afero.NewOsFs())

	if err := verifyAdminUser(username, dataDir); err != nil {
		return err
	}

	var password string
	for {
		password = questionPassword("New Password")
		if password != "" {
			break
		}
		errorMessage("Password cannot be empty")
	}

	for {
		confirmPassword := questionPassword("Confirm Password")
		if confirmPassword == password {
			break
		}
		errorMessage("Passwords do not match")
	}

	fmt.Println()

	if err := resetPasswordViaAPI(fs, username, password, dataDir, httpAddr); err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	successMessage("Password reset successfully for '%s'", username)

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// verifyAdminUser ensures a user exists and is admin
func verifyAdminUser(username, dataDir string) error {
	ctx := context.Background()
	fs := filesystem.New(afero.NewOsFs())

	dbManager, err := database.NewSQLiteManager(&database.DatabaseManagerConfig{
		DataDir: dataDir,
		FS:      fs,
		Testing: false,
	})
	if err != nil {
		return fmt.Errorf("failed to create database manager: %w", err)
	}

	appDao := dao.New(dbManager.DataDb)
	dbOpts := dao.NewOptions().WithWhere(squirrel.Eq{models.USER_TABLE_USERNAME: username})
	user, err := appDao.GetUser(ctx, dbOpts)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user '%s' not found", username)
		}

		return fmt.Errorf("failed to lookup user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("user '%s' not found", username)
	}

	if user.SiteRole != types.SiteRoleAdmin {
		return fmt.Errorf("user '%s' is not an admin user", username)
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// resetPasswordViaAPI generates a recovery token on disk, then makes an HTTP request to
// the running application to reset the password
func resetPasswordViaAPI(fs *filesystem.FS, username, password, dataDir, httpAddr string) error {
	recoveryToken, err := auth.GenerateRecoveryToken(fs, username, password, dataDir)
	if err != nil {
		return fmt.Errorf("failed to generate recovery token: %w", err)
	}

	defer func() {
		auth.DeleteRecoveryToken(fs, dataDir)
	}()

	jsonData, err := json.Marshal(map[string]string{"token": recoveryToken.Token})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("http://%s/api/admin/recovery", httpAddr)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("application is not running or not accessible at %s: %w", httpAddr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("recovery request failed with status %d", resp.StatusCode)
	}

	return nil
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// init adds the admin reset password command to the admin command
func init() {
	adminCmd.AddCommand(adminResetPasswordCmd)

	adminResetPasswordCmd.Flags().String("http", "127.0.0.1:9081", "TCP address of the running HTTP server")
	adminResetPasswordCmd.Flags().String("data-dir", "./friendle_data", "Directory to store data files")

	_ = viper.BindPFlag("http", adminResetPasswordCmd.Flags().Lookup("http"))
	_ = viper.BindPFlag("data-dir", adminResetPasswordCmd.Flags().Lookup("data-dir"))
}
