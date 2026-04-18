// Package notify provides notification integrations for vaultwatch.
//
// Chat integrations:
//
//   - GoogleChatNotifier: sends alerts to a Google Chat space via
//     an incoming webhook URL.
//
//   - TelegramNotifier: sends alerts via the Telegram Bot API using
//     HTML-formatted messages to a specified chat ID.
//
// Both notifiers implement the Notifier interface and can be composed
// with MultiNotifier to fan out alerts across multiple channels.
package notify
