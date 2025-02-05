package notification

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/discord"
	"github.com/nikoksr/notify/service/mail"
	"github.com/nikoksr/notify/service/slack"
	"github.com/nikoksr/notify/service/telegram"
)

type NotificationType string

const (
	ScanFinished    NotificationType = "scan_finished"
	HighFinding     NotificationType = "high_finding"
	CriticalFinding NotificationType = "critical_finding"
	ScanFailed      NotificationType = "scan_failed"
	ScanStarted     NotificationType = "scan_started"
	ScanStopped     NotificationType = "scan_stopped"
)

type NotificationService struct {
	notifier *notify.Notify
	config   *NotificationConfig
	mu       sync.RWMutex
}

type NotificationConfig struct {
	Email    *EmailConfig       `json:"email"`
	Slack    *SlackConfig       `json:"slack"`
	Discord  *DiscordConfig     `json:"discord"`
	Telegram *TelegramConfig    `json:"telegram"`
	Rules    []NotificationRule `json:"rules"`
}

type EmailConfig struct {
	Enabled  bool     `json:"enabled"`
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	From     string   `json:"from"`
	To       []string `json:"to"`
}

type SlackConfig struct {
	Enabled bool   `json:"enabled"`
	Token   string `json:"token"`
	Channel string `json:"channel"`
}

type DiscordConfig struct {
	Enabled   bool   `json:"enabled"`
	WebhookID string `json:"webhook_id"`
	Token     string `json:"token"`
}

type TelegramConfig struct {
	Enabled bool    `json:"enabled"`
	Token   string  `json:"token"`
	ChatIDs []int64 `json:"chat_ids"`
}

type NotificationRule struct {
	ID       string           `json:"id"`
	Type     NotificationType `json:"type"`
	Severity []string         `json:"severity"`
	Channels []string         `json:"channels"`
	Enabled  bool             `json:"enabled"`
}

func NewNotificationService(config *NotificationConfig) (*NotificationService, error) {
	n := &NotificationService{
		notifier: notify.New(),
		config:   config,
	}

	if err := n.configureServices(); err != nil {
		return nil, fmt.Errorf("failed to configure notification services: %v", err)
	}

	return n, nil
}

func (n *NotificationService) configureServices() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Create a new notifier instance
	n.notifier = notify.New()

	// Configure Email
	if n.config.Email != nil && n.config.Email.Enabled {
		// Create mail service
		mailService := mail.New(n.config.Email.Username, n.config.Email.Password)
		// Add receivers
		mailService.AddReceivers(n.config.Email.To...)
		n.notifier.UseServices(mailService)
	}

	// Configure Slack
	if n.config.Slack != nil && n.config.Slack.Enabled {
		slackService := slack.New(n.config.Slack.Token)
		slackService.AddReceivers(n.config.Slack.Channel)
		n.notifier.UseServices(slackService)
	}

	// Configure Discord
	if n.config.Discord != nil && n.config.Discord.Enabled {
		webhookURL := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s",
			n.config.Discord.WebhookID,
			n.config.Discord.Token)
		discordService := discord.New()
		discordService.AddReceivers(webhookURL)
		n.notifier.UseServices(discordService)
	}

	// Configure Telegram
	if n.config.Telegram != nil && n.config.Telegram.Enabled {
		telegramService, err := telegram.New(n.config.Telegram.Token)
		if err != nil {
			return fmt.Errorf("failed to create telegram service: %v", err)
		}
		for _, chatID := range n.config.Telegram.ChatIDs {
			telegramService.AddReceivers(chatID)
		}
		n.notifier.UseServices(telegramService)
	}

	return nil
}

func (n *NotificationService) UpdateConfig(config *NotificationConfig) error {
	n.mu.Lock()
	n.config = config
	n.mu.Unlock()

	return n.configureServices()
}

func (n *NotificationService) ShouldNotify(notificationType NotificationType, severity string) bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	for _, rule := range n.config.Rules {
		if !rule.Enabled {
			continue
		}

		if rule.Type == notificationType {
			// For scan notifications that don't have severity
			if severity == "" {
				return true
			}

			// Check if severity matches
			for _, s := range rule.Severity {
				if s == severity {
					return true
				}
			}
		}
	}

	return false
}

func (n *NotificationService) Notify(ctx context.Context, subject, message string) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.notifier.Send(
		ctx,
		subject,
		message,
	)
}

func (n *NotificationService) NotifyScanFinished(ctx context.Context, scanID, scanName string) error {
	if !n.ShouldNotify(ScanFinished, "") {
		return nil
	}

	subject := fmt.Sprintf("Scan Finished: %s", scanName)
	message := fmt.Sprintf("The scan %s (ID: %s) has finished at %s",
		scanName,
		scanID,
		time.Now().Format(time.RFC3339),
	)

	return n.Notify(ctx, subject, message)
}

func (n *NotificationService) NotifyFinding(ctx context.Context, findingID, scanName, severity, title string) error {
	if !n.ShouldNotify(NotificationType(severity+"_finding"), severity) {
		return nil
	}

	subject := fmt.Sprintf("New %s Finding in %s", severity, scanName)
	message := fmt.Sprintf("A new %s severity finding has been detected:\nTitle: %s\nID: %s\nScan: %s",
		severity,
		title,
		findingID,
		scanName,
	)

	return n.Notify(ctx, subject, message)
}

func (n *NotificationService) NotifyScanFailed(ctx context.Context, scanID, scanName, error string) error {
	if !n.ShouldNotify(ScanFailed, "") {
		return nil
	}

	subject := fmt.Sprintf("Scan Failed: %s", scanName)
	message := fmt.Sprintf("The scan %s (ID: %s) has failed at %s\nError: %s",
		scanName,
		scanID,
		time.Now().Format(time.RFC3339),
		error,
	)

	return n.Notify(ctx, subject, message)
}

func (n *NotificationService) NotifyScanStarted(ctx context.Context, scanID, scanName string) error {
	if !n.ShouldNotify(ScanStarted, "") {
		return nil
	}

	subject := fmt.Sprintf("Scan Started: %s", scanName)
	message := fmt.Sprintf("The scan %s (ID: %s) has started at %s",
		scanName,
		scanID,
		time.Now().Format(time.RFC3339),
	)

	return n.Notify(ctx, subject, message)
}

func (n *NotificationService) NotifyScanStopped(ctx context.Context, scanID, scanName string) error {
	if !n.ShouldNotify(ScanStopped, "") {
		return nil
	}

	subject := fmt.Sprintf("Scan Stopped: %s", scanName)
	message := fmt.Sprintf("The scan %s (ID: %s) was manually stopped at %s",
		scanName,
		scanID,
		time.Now().Format(time.RFC3339),
	)

	return n.Notify(ctx, subject, message)
}
