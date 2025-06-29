package gotify

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gotify/pkg/config"
	"gotify/pkg/notifier/services"
)

var (
	configPath string
	service    string
	receivers  []int64
)

var rootCmd = &cobra.Command{
	Use:   "gotify",
	Short: "Gotify - A simple CLI tool to send notifications",
	Long:  `Gotify is a CLI tool that allows you to send notifications to various messaging services.`,
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a notification",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendNotification(args[0])
	},
}

func init() {
	sendCmd.Flags().StringVarP(&service, "service", "s", "", "Service to use (e.g., telegram)")
	sendCmd.Flags().Int64SliceVarP(&receivers, "receivers", "r", []int64{}, "Receiver IDs (e.g., 123456789)")
	sendCmd.MarkFlagRequired("service")
	rootCmd.AddCommand(sendCmd)

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file (default is $HOME/.config/gotify/config.yaml)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func sendNotification(message string) error {
	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize notifier based on the service
	switch service {
	case "telegram":
		// Check environment variable first
		apiKey := os.Getenv("TELEGRAM_API_KEY")
		
		// If not in environment, try to get from config
		if apiKey == "" {
			if svcCfg, exists := cfg.Services["telegram"]; exists && svcCfg.APIKey != "" {
				apiKey = svcCfg.APIKey
				// Use receivers from config if not provided via CLI
				if len(receivers) == 0 {
					receivers = svcCfg.Receivers
				}
			} else {
				return fmt.Errorf("telegram API key not found. Set TELEGRAM_API_KEY environment variable or add it to the config file")
			}
		}

		if len(receivers) == 0 {
			return fmt.Errorf("no receivers specified. Use --receivers flag or configure them in the config file")
		}

		tg, err := services.NewTelegram(apiKey, receivers)
		if err != nil {
			return fmt.Errorf("failed to create telegram client: %w", err)
		}

		ctx := context.Background()
		return tg.Send(ctx, "Notification", message)

	default:
		return fmt.Errorf("unsupported service: %s", service)
	}
}
