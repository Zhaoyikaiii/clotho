// Copyright (c) 2026 Clotho contributors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	// "github.com/caarlos0/env/v11"
	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml/v2"

	"github.com/Zhaoyikaiii/clotho/pkg/fileutil"
)

// rrCounter is a global counter for round-robin load balancing across models.
var rrCounter atomic.Uint64

// FlexibleStringSlice is a []string that also accepts JSON numbers,
// so allow_from can contain both "123" and 123.
type FlexibleStringSlice []string

func (f *FlexibleStringSlice) UnmarshalJSON(data []byte) error {
	// Try []string first
	var ss []string
	if err := json.Unmarshal(data, &ss); err == nil {
		*f = ss
		return nil
	}

	// Try []interface{} to handle mixed types
	var raw []any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	result := make([]string, 0, len(raw))
	for _, v := range raw {
		switch val := v.(type) {
		case string:
			result = append(result, val)
		case float64:
			result = append(result, fmt.Sprintf("%.0f", val))
		default:
			result = append(result, fmt.Sprintf("%v", val))
		}
	}
	*f = result
	return nil
}

type Config struct {
	Agents    AgentsConfig    `mapstructure:"agents"`
	Bindings  []AgentBinding  `mapstructure:"bindings"`
	Session   SessionConfig   `mapstructure:"session"`
	Channels  ChannelsConfig  `mapstructure:"channels"`
	Providers ProvidersConfig `mapstructure:"providers"`
	ModelList []ModelConfig   `mapstructure:"model_list"`
	Gateway   GatewayConfig   `mapstructure:"gateway"`
	Tools     ToolsConfig     `mapstructure:"tools"`
	Heartbeat HeartbeatConfig `mapstructure:"heartbeat"`
	Devices   DevicesConfig   `mapstructure:"devices"`
}

// MarshalJSON implements custom JSON marshaling for Config
// to omit providers section when empty and session when empty
func (c Config) MarshalJSON() ([]byte, error) {
	type Alias Config
	aux := &struct {
		Providers *ProvidersConfig `mapstructure:"providers,omitempty"`
		Session   *SessionConfig   `mapstructure:"session,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(&c),
	}

	// Only include providers if not empty
	if !c.Providers.IsEmpty() {
		aux.Providers = &c.Providers
	}

	// Only include session if not empty
	if c.Session.DMScope != "" || len(c.Session.IdentityLinks) > 0 {
		aux.Session = &c.Session
	}

	return json.Marshal(aux)
}

type AgentsConfig struct {
	Defaults AgentDefaults `mapstructure:"defaults"`
	List     []AgentConfig `mapstructure:"list,omitempty"`
}

// AgentModelConfig supports both string and structured model config.
// String format: "gpt-4" (just primary, no fallbacks)
// Object format: {"primary": "gpt-4", "fallbacks": ["claude-haiku"]}
type AgentModelConfig struct {
	Primary   string   `mapstructure:"primary,omitempty"`
	Fallbacks []string `mapstructure:"fallbacks,omitempty"`
}

func (m *AgentModelConfig) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		m.Primary = s
		m.Fallbacks = nil
		return nil
	}
	type raw struct {
		Primary   string   `mapstructure:"primary"`
		Fallbacks []string `mapstructure:"fallbacks"`
	}
	var r raw
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	m.Primary = r.Primary
	m.Fallbacks = r.Fallbacks
	return nil
}

func (m AgentModelConfig) MarshalJSON() ([]byte, error) {
	if len(m.Fallbacks) == 0 && m.Primary != "" {
		return json.Marshal(m.Primary)
	}
	type raw struct {
		Primary   string   `mapstructure:"primary,omitempty"`
		Fallbacks []string `mapstructure:"fallbacks,omitempty"`
	}
	return json.Marshal(raw{Primary: m.Primary, Fallbacks: m.Fallbacks})
}

type AgentConfig struct {
	ID        string            `mapstructure:"id"`
	Default   bool              `mapstructure:"default,omitempty"`
	Name      string            `mapstructure:"name,omitempty"`
	Workspace string            `mapstructure:"workspace,omitempty"`
	Model     *AgentModelConfig `mapstructure:"model,omitempty"`
	Skills    []string          `mapstructure:"skills,omitempty"`
	Subagents *SubagentsConfig  `mapstructure:"subagents,omitempty"`
}

type SubagentsConfig struct {
	AllowAgents []string          `mapstructure:"allow_agents,omitempty"`
	Model       *AgentModelConfig `mapstructure:"model,omitempty"`
}

type PeerMatch struct {
	Kind string `mapstructure:"kind"`
	ID   string `mapstructure:"id"`
}

type BindingMatch struct {
	Channel   string     `mapstructure:"channel"`
	AccountID string     `mapstructure:"account_id,omitempty"`
	Peer      *PeerMatch `mapstructure:"peer,omitempty"`
	GuildID   string     `mapstructure:"guild_id,omitempty"`
	TeamID    string     `mapstructure:"team_id,omitempty"`
}

type AgentBinding struct {
	AgentID string       `mapstructure:"agent_id"`
	Match   BindingMatch `mapstructure:"match"`
}

type SessionConfig struct {
	DMScope       string              `mapstructure:"dm_scope,omitempty"`
	IdentityLinks map[string][]string `mapstructure:"identity_links,omitempty"`
}

type AgentDefaults struct {
	Workspace                 string   `mapstructure:"workspace"                       env:"CLOTHO_AGENTS_DEFAULTS_WORKSPACE"`
	RestrictToWorkspace       bool     `mapstructure:"restrict_to_workspace"           env:"CLOTHO_AGENTS_DEFAULTS_RESTRICT_TO_WORKSPACE"`
	AllowReadOutsideWorkspace bool     `mapstructure:"allow_read_outside_workspace"    env:"CLOTHO_AGENTS_DEFAULTS_ALLOW_READ_OUTSIDE_WORKSPACE"`
	Provider                  string   `mapstructure:"provider"                        env:"CLOTHO_AGENTS_DEFAULTS_PROVIDER"`
	ModelName                 string   `mapstructure:"model_name,omitempty"            env:"CLOTHO_AGENTS_DEFAULTS_MODEL_NAME"`
	Model                     string   `mapstructure:"model"                           env:"CLOTHO_AGENTS_DEFAULTS_MODEL"` // Deprecated: use model_name instead
	ModelFallbacks            []string `mapstructure:"model_fallbacks,omitempty"`
	ImageModel                string   `mapstructure:"image_model,omitempty"           env:"CLOTHO_AGENTS_DEFAULTS_IMAGE_MODEL"`
	ImageModelFallbacks       []string `mapstructure:"image_model_fallbacks,omitempty"`
	MaxTokens                 int      `mapstructure:"max_tokens"                      env:"CLOTHO_AGENTS_DEFAULTS_MAX_TOKENS"`
	Temperature               *float64 `mapstructure:"temperature,omitempty"           env:"CLOTHO_AGENTS_DEFAULTS_TEMPERATURE"`
	MaxToolIterations         int      `mapstructure:"max_tool_iterations"             env:"CLOTHO_AGENTS_DEFAULTS_MAX_TOOL_ITERATIONS"`
	SummarizeMessageThreshold int      `mapstructure:"summarize_message_threshold"     env:"CLOTHO_AGENTS_DEFAULTS_SUMMARIZE_MESSAGE_THRESHOLD"`
	SummarizeTokenPercent     int      `mapstructure:"summarize_token_percent"         env:"CLOTHO_AGENTS_DEFAULTS_SUMMARIZE_TOKEN_PERCENT"`
	MaxMediaSize              int      `mapstructure:"max_media_size,omitempty"        env:"CLOTHO_AGENTS_DEFAULTS_MAX_MEDIA_SIZE"`
}

const DefaultMaxMediaSize = 20 * 1024 * 1024 // 20 MB

func (d *AgentDefaults) GetMaxMediaSize() int {
	if d.MaxMediaSize > 0 {
		return d.MaxMediaSize
	}
	return DefaultMaxMediaSize
}

// GetModelName returns the effective model name for the agent defaults.
// It prefers the new "model_name" field but falls back to "model" for backward compatibility.
func (d *AgentDefaults) GetModelName() string {
	if d.ModelName != "" {
		return d.ModelName
	}
	return d.Model
}

type ChannelsConfig struct {
	WhatsApp   WhatsAppConfig   `mapstructure:"whatsapp"`
	Telegram   TelegramConfig   `mapstructure:"telegram"`
	Feishu     FeishuConfig     `mapstructure:"feishu"`
	Discord    DiscordConfig    `mapstructure:"discord"`
	MaixCam    MaixCamConfig    `mapstructure:"maixcam"`
	QQ         QQConfig         `mapstructure:"qq"`
	DingTalk   DingTalkConfig   `mapstructure:"dingtalk"`
	Slack      SlackConfig      `mapstructure:"slack"`
	LINE       LINEConfig       `mapstructure:"line"`
	OneBot     OneBotConfig     `mapstructure:"onebot"`
	WeCom      WeComConfig      `mapstructure:"wecom"`
	WeComApp   WeComAppConfig   `mapstructure:"wecom_app"`
	WeComAIBot WeComAIBotConfig `mapstructure:"wecom_aibot"`
}

// GroupTriggerConfig controls when the bot responds in group chats.
type GroupTriggerConfig struct {
	MentionOnly bool     `mapstructure:"mention_only,omitempty"`
	Prefixes    []string `mapstructure:"prefixes,omitempty"`
}

// TypingConfig controls typing indicator behavior (Phase 10).
type TypingConfig struct {
	Enabled bool `mapstructure:"enabled,omitempty"`
}

// PlaceholderConfig controls placeholder message behavior (Phase 10).
type PlaceholderConfig struct {
	Enabled bool   `mapstructure:"enabled,omitempty"`
	Text    string `mapstructure:"text,omitempty"`
}

type WhatsAppConfig struct {
	Enabled            bool                `mapstructure:"enabled"              env:"CLOTHO_CHANNELS_WHATSAPP_ENABLED"`
	BridgeURL          string              `mapstructure:"bridge_url"           env:"CLOTHO_CHANNELS_WHATSAPP_BRIDGE_URL"`
	UseNative          bool                `mapstructure:"use_native"           env:"CLOTHO_CHANNELS_WHATSAPP_USE_NATIVE"`
	SessionStorePath   string              `mapstructure:"session_store_path"   env:"CLOTHO_CHANNELS_WHATSAPP_SESSION_STORE_PATH"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"           env:"CLOTHO_CHANNELS_WHATSAPP_ALLOW_FROM"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id" env:"CLOTHO_CHANNELS_WHATSAPP_REASONING_CHANNEL_ID"`
}

type TelegramConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_TELEGRAM_ENABLED"`
	Token              string              `mapstructure:"token"                   env:"CLOTHO_CHANNELS_TELEGRAM_TOKEN"`
	BaseURL            string              `mapstructure:"base_url"                env:"CLOTHO_CHANNELS_TELEGRAM_BASE_URL"`
	Proxy              string              `mapstructure:"proxy"                   env:"CLOTHO_CHANNELS_TELEGRAM_PROXY"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_TELEGRAM_ALLOW_FROM"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	Typing             TypingConfig        `mapstructure:"typing,omitempty"`
	Placeholder        PlaceholderConfig   `mapstructure:"placeholder,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_TELEGRAM_REASONING_CHANNEL_ID"`
}

type FeishuConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_FEISHU_ENABLED"`
	AppID              string              `mapstructure:"app_id"                  env:"CLOTHO_CHANNELS_FEISHU_APP_ID"`
	AppSecret          string              `mapstructure:"app_secret"              env:"CLOTHO_CHANNELS_FEISHU_APP_SECRET"`
	EncryptKey         string              `mapstructure:"encrypt_key"             env:"CLOTHO_CHANNELS_FEISHU_ENCRYPT_KEY"`
	VerificationToken  string              `mapstructure:"verification_token"      env:"CLOTHO_CHANNELS_FEISHU_VERIFICATION_TOKEN"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_FEISHU_ALLOW_FROM"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	Placeholder        PlaceholderConfig   `mapstructure:"placeholder,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_FEISHU_REASONING_CHANNEL_ID"`
}

type DiscordConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_DISCORD_ENABLED"`
	Token              string              `mapstructure:"token"                   env:"CLOTHO_CHANNELS_DISCORD_TOKEN"`
	Proxy              string              `mapstructure:"proxy"                   env:"CLOTHO_CHANNELS_DISCORD_PROXY"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_DISCORD_ALLOW_FROM"`
	MentionOnly        bool                `mapstructure:"mention_only"            env:"CLOTHO_CHANNELS_DISCORD_MENTION_ONLY"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	Typing             TypingConfig        `mapstructure:"typing,omitempty"`
	Placeholder        PlaceholderConfig   `mapstructure:"placeholder,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_DISCORD_REASONING_CHANNEL_ID"`
}

type MaixCamConfig struct {
	Enabled            bool                `mapstructure:"enabled"              env:"CLOTHO_CHANNELS_MAIXCAM_ENABLED"`
	Host               string              `mapstructure:"host"                 env:"CLOTHO_CHANNELS_MAIXCAM_HOST"`
	Port               int                 `mapstructure:"port"                 env:"CLOTHO_CHANNELS_MAIXCAM_PORT"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"           env:"CLOTHO_CHANNELS_MAIXCAM_ALLOW_FROM"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id" env:"CLOTHO_CHANNELS_MAIXCAM_REASONING_CHANNEL_ID"`
}

type QQConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_QQ_ENABLED"`
	AppID              string              `mapstructure:"app_id"                  env:"CLOTHO_CHANNELS_QQ_APP_ID"`
	AppSecret          string              `mapstructure:"app_secret"              env:"CLOTHO_CHANNELS_QQ_APP_SECRET"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_QQ_ALLOW_FROM"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_QQ_REASONING_CHANNEL_ID"`
}

type DingTalkConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_DINGTALK_ENABLED"`
	ClientID           string              `mapstructure:"client_id"               env:"CLOTHO_CHANNELS_DINGTALK_CLIENT_ID"`
	ClientSecret       string              `mapstructure:"client_secret"           env:"CLOTHO_CHANNELS_DINGTALK_CLIENT_SECRET"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_DINGTALK_ALLOW_FROM"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_DINGTALK_REASONING_CHANNEL_ID"`
}

type SlackConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_SLACK_ENABLED"`
	BotToken           string              `mapstructure:"bot_token"               env:"CLOTHO_CHANNELS_SLACK_BOT_TOKEN"`
	AppToken           string              `mapstructure:"app_token"               env:"CLOTHO_CHANNELS_SLACK_APP_TOKEN"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_SLACK_ALLOW_FROM"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	Typing             TypingConfig        `mapstructure:"typing,omitempty"`
	Placeholder        PlaceholderConfig   `mapstructure:"placeholder,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_SLACK_REASONING_CHANNEL_ID"`
}

type LINEConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_LINE_ENABLED"`
	ChannelSecret      string              `mapstructure:"channel_secret"          env:"CLOTHO_CHANNELS_LINE_CHANNEL_SECRET"`
	ChannelAccessToken string              `mapstructure:"channel_access_token"    env:"CLOTHO_CHANNELS_LINE_CHANNEL_ACCESS_TOKEN"`
	WebhookHost        string              `mapstructure:"webhook_host"            env:"CLOTHO_CHANNELS_LINE_WEBHOOK_HOST"`
	WebhookPort        int                 `mapstructure:"webhook_port"            env:"CLOTHO_CHANNELS_LINE_WEBHOOK_PORT"`
	WebhookPath        string              `mapstructure:"webhook_path"            env:"CLOTHO_CHANNELS_LINE_WEBHOOK_PATH"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_LINE_ALLOW_FROM"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	Typing             TypingConfig        `mapstructure:"typing,omitempty"`
	Placeholder        PlaceholderConfig   `mapstructure:"placeholder,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_LINE_REASONING_CHANNEL_ID"`
}

type OneBotConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_ONEBOT_ENABLED"`
	WSUrl              string              `mapstructure:"ws_url"                  env:"CLOTHO_CHANNELS_ONEBOT_WS_URL"`
	AccessToken        string              `mapstructure:"access_token"            env:"CLOTHO_CHANNELS_ONEBOT_ACCESS_TOKEN"`
	ReconnectInterval  int                 `mapstructure:"reconnect_interval"      env:"CLOTHO_CHANNELS_ONEBOT_RECONNECT_INTERVAL"`
	GroupTriggerPrefix []string            `mapstructure:"group_trigger_prefix"    env:"CLOTHO_CHANNELS_ONEBOT_GROUP_TRIGGER_PREFIX"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_ONEBOT_ALLOW_FROM"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	Typing             TypingConfig        `mapstructure:"typing,omitempty"`
	Placeholder        PlaceholderConfig   `mapstructure:"placeholder,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_ONEBOT_REASONING_CHANNEL_ID"`
}

type WeComConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_WECOM_ENABLED"`
	Token              string              `mapstructure:"token"                   env:"CLOTHO_CHANNELS_WECOM_TOKEN"`
	EncodingAESKey     string              `mapstructure:"encoding_aes_key"        env:"CLOTHO_CHANNELS_WECOM_ENCODING_AES_KEY"`
	WebhookURL         string              `mapstructure:"webhook_url"             env:"CLOTHO_CHANNELS_WECOM_WEBHOOK_URL"`
	WebhookHost        string              `mapstructure:"webhook_host"            env:"CLOTHO_CHANNELS_WECOM_WEBHOOK_HOST"`
	WebhookPort        int                 `mapstructure:"webhook_port"            env:"CLOTHO_CHANNELS_WECOM_WEBHOOK_PORT"`
	WebhookPath        string              `mapstructure:"webhook_path"            env:"CLOTHO_CHANNELS_WECOM_WEBHOOK_PATH"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_WECOM_ALLOW_FROM"`
	ReplyTimeout       int                 `mapstructure:"reply_timeout"           env:"CLOTHO_CHANNELS_WECOM_REPLY_TIMEOUT"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_WECOM_REASONING_CHANNEL_ID"`
}

type WeComAppConfig struct {
	Enabled            bool                `mapstructure:"enabled"                 env:"CLOTHO_CHANNELS_WECOM_APP_ENABLED"`
	CorpID             string              `mapstructure:"corp_id"                 env:"CLOTHO_CHANNELS_WECOM_APP_CORP_ID"`
	CorpSecret         string              `mapstructure:"corp_secret"             env:"CLOTHO_CHANNELS_WECOM_APP_CORP_SECRET"`
	AgentID            int64               `mapstructure:"agent_id"                env:"CLOTHO_CHANNELS_WECOM_APP_AGENT_ID"`
	Token              string              `mapstructure:"token"                   env:"CLOTHO_CHANNELS_WECOM_APP_TOKEN"`
	EncodingAESKey     string              `mapstructure:"encoding_aes_key"        env:"CLOTHO_CHANNELS_WECOM_APP_ENCODING_AES_KEY"`
	WebhookHost        string              `mapstructure:"webhook_host"            env:"CLOTHO_CHANNELS_WECOM_APP_WEBHOOK_HOST"`
	WebhookPort        int                 `mapstructure:"webhook_port"            env:"CLOTHO_CHANNELS_WECOM_APP_WEBHOOK_PORT"`
	WebhookPath        string              `mapstructure:"webhook_path"            env:"CLOTHO_CHANNELS_WECOM_APP_WEBHOOK_PATH"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"              env:"CLOTHO_CHANNELS_WECOM_APP_ALLOW_FROM"`
	ReplyTimeout       int                 `mapstructure:"reply_timeout"           env:"CLOTHO_CHANNELS_WECOM_APP_REPLY_TIMEOUT"`
	GroupTrigger       GroupTriggerConfig  `mapstructure:"group_trigger,omitempty"`
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id"    env:"CLOTHO_CHANNELS_WECOM_APP_REASONING_CHANNEL_ID"`
}

type WeComAIBotConfig struct {
	Enabled            bool                `mapstructure:"enabled"              env:"CLOTHO_CHANNELS_WECOM_AIBOT_ENABLED"`
	Token              string              `mapstructure:"token"                env:"CLOTHO_CHANNELS_WECOM_AIBOT_TOKEN"`
	EncodingAESKey     string              `mapstructure:"encoding_aes_key"     env:"CLOTHO_CHANNELS_WECOM_AIBOT_ENCODING_AES_KEY"`
	WebhookPath        string              `mapstructure:"webhook_path"         env:"CLOTHO_CHANNELS_WECOM_AIBOT_WEBHOOK_PATH"`
	AllowFrom          FlexibleStringSlice `mapstructure:"allow_from"           env:"CLOTHO_CHANNELS_WECOM_AIBOT_ALLOW_FROM"`
	ReplyTimeout       int                 `mapstructure:"reply_timeout"        env:"CLOTHO_CHANNELS_WECOM_AIBOT_REPLY_TIMEOUT"`
	MaxSteps           int                 `mapstructure:"max_steps"            env:"CLOTHO_CHANNELS_WECOM_AIBOT_MAX_STEPS"`       // Maximum streaming steps
	WelcomeMessage     string              `mapstructure:"welcome_message"      env:"CLOTHO_CHANNELS_WECOM_AIBOT_WELCOME_MESSAGE"` // Sent on enter_chat event; empty = no welcome
	ReasoningChannelID string              `mapstructure:"reasoning_channel_id" env:"CLOTHO_CHANNELS_WECOM_AIBOT_REASONING_CHANNEL_ID"`
}

type HeartbeatConfig struct {
	Enabled  bool `mapstructure:"enabled"  env:"CLOTHO_HEARTBEAT_ENABLED"`
	Interval int  `mapstructure:"interval" env:"CLOTHO_HEARTBEAT_INTERVAL"` // minutes, min 5
}

type DevicesConfig struct {
	Enabled    bool `mapstructure:"enabled"     env:"CLOTHO_DEVICES_ENABLED"`
	MonitorUSB bool `mapstructure:"monitor_usb" env:"CLOTHO_DEVICES_MONITOR_USB"`
}

type ProvidersConfig struct {
	Anthropic     ProviderConfig       `mapstructure:"anthropic"`
	OpenAI        OpenAIProviderConfig `mapstructure:"openai"`
	LiteLLM       ProviderConfig       `mapstructure:"litellm"`
	OpenRouter    ProviderConfig       `mapstructure:"openrouter"`
	Groq          ProviderConfig       `mapstructure:"groq"`
	Zhipu         ProviderConfig       `mapstructure:"zhipu"`
	VLLM          ProviderConfig       `mapstructure:"vllm"`
	Gemini        ProviderConfig       `mapstructure:"gemini"`
	Nvidia        ProviderConfig       `mapstructure:"nvidia"`
	Ollama        ProviderConfig       `mapstructure:"ollama"`
	Moonshot      ProviderConfig       `mapstructure:"moonshot"`
	ShengSuanYun  ProviderConfig       `mapstructure:"shengsuanyun"`
	DeepSeek      ProviderConfig       `mapstructure:"deepseek"`
	Cerebras      ProviderConfig       `mapstructure:"cerebras"`
	VolcEngine    ProviderConfig       `mapstructure:"volcengine"`
	GitHubCopilot ProviderConfig       `mapstructure:"github_copilot"`
	Antigravity   ProviderConfig       `mapstructure:"antigravity"`
	Qwen          ProviderConfig       `mapstructure:"qwen"`
	Mistral       ProviderConfig       `mapstructure:"mistral"`
	Avian         ProviderConfig       `mapstructure:"avian"`
}

// IsEmpty checks if all provider configs are empty (no API keys or API bases set)
// Note: WebSearch is an optimization option and doesn't count as "non-empty"
func (p ProvidersConfig) IsEmpty() bool {
	return p.Anthropic.APIKey == "" && p.Anthropic.APIBase == "" &&
		p.OpenAI.APIKey == "" && p.OpenAI.APIBase == "" &&
		p.LiteLLM.APIKey == "" && p.LiteLLM.APIBase == "" &&
		p.OpenRouter.APIKey == "" && p.OpenRouter.APIBase == "" &&
		p.Groq.APIKey == "" && p.Groq.APIBase == "" &&
		p.Zhipu.APIKey == "" && p.Zhipu.APIBase == "" &&
		p.VLLM.APIKey == "" && p.VLLM.APIBase == "" &&
		p.Gemini.APIKey == "" && p.Gemini.APIBase == "" &&
		p.Nvidia.APIKey == "" && p.Nvidia.APIBase == "" &&
		p.Ollama.APIKey == "" && p.Ollama.APIBase == "" &&
		p.Moonshot.APIKey == "" && p.Moonshot.APIBase == "" &&
		p.ShengSuanYun.APIKey == "" && p.ShengSuanYun.APIBase == "" &&
		p.DeepSeek.APIKey == "" && p.DeepSeek.APIBase == "" &&
		p.Cerebras.APIKey == "" && p.Cerebras.APIBase == "" &&
		p.VolcEngine.APIKey == "" && p.VolcEngine.APIBase == "" &&
		p.GitHubCopilot.APIKey == "" && p.GitHubCopilot.APIBase == "" &&
		p.Antigravity.APIKey == "" && p.Antigravity.APIBase == "" &&
		p.Qwen.APIKey == "" && p.Qwen.APIBase == "" &&
		p.Mistral.APIKey == "" && p.Mistral.APIBase == "" &&
		p.Avian.APIKey == "" && p.Avian.APIBase == ""
}

// MarshalJSON implements custom JSON marshaling for ProvidersConfig
// to omit the entire section when empty
func (p ProvidersConfig) MarshalJSON() ([]byte, error) {
	if p.IsEmpty() {
		return []byte("null"), nil
	}
	type Alias ProvidersConfig
	return json.Marshal((*Alias)(&p))
}

type ProviderConfig struct {
	APIKey         string `mapstructure:"api_key"                   env:"CLOTHO_PROVIDERS_{{.Name}}_API_KEY"`
	APIBase        string `mapstructure:"api_base"                  env:"CLOTHO_PROVIDERS_{{.Name}}_API_BASE"`
	Proxy          string `mapstructure:"proxy,omitempty"           env:"CLOTHO_PROVIDERS_{{.Name}}_PROXY"`
	RequestTimeout int    `mapstructure:"request_timeout,omitempty" env:"CLOTHO_PROVIDERS_{{.Name}}_REQUEST_TIMEOUT"`
	AuthMethod     string `mapstructure:"auth_method,omitempty"     env:"CLOTHO_PROVIDERS_{{.Name}}_AUTH_METHOD"`
	ConnectMode    string `mapstructure:"connect_mode,omitempty"    env:"CLOTHO_PROVIDERS_{{.Name}}_CONNECT_MODE"` // only for Github Copilot, `stdio` or `grpc`
}

type OpenAIProviderConfig struct {
	ProviderConfig
	WebSearch bool `mapstructure:"web_search" env:"CLOTHO_PROVIDERS_OPENAI_WEB_SEARCH"`
}

// ModelConfig represents a model-centric provider configuration.
// It allows adding new providers (especially OpenAI-compatible ones) via configuration only.
// The model field uses protocol prefix format: [protocol/]model-identifier
// Supported protocols: openai, anthropic, antigravity, claude-cli, codex-cli, github-copilot
// Default protocol is "openai" if no prefix is specified.
type ModelConfig struct {
	// Required fields
	ModelName string `mapstructure:"model_name"` // User-facing alias for the model
	Model     string `mapstructure:"model"`      // Protocol/model-identifier (e.g., "openai/gpt-4o", "anthropic/claude-sonnet-4.6")

	// HTTP-based providers
	APIBase string `mapstructure:"api_base"` // API endpoint URL
	APIKey  string `mapstructure:"api_key"` // API authentication key
	Proxy   string `mapstructure:"proxy"`  // HTTP proxy URL

	// Special providers (CLI-based, OAuth, etc.)
	AuthMethod  string `mapstructure:"auth_method"`  // Authentication method: oauth, token
	ConnectMode string `mapstructure:"connect_mode"` // Connection mode: stdio, grpc
	Workspace   string `mapstructure:"workspace"`    // Workspace path for CLI-based providers

	// Optional optimizations
	RPM            int    `mapstructure:"rpm"`              // Requests per minute limit
	MaxTokensField string `mapstructure:"max_tokens_field"` // Field name for max tokens (e.g., "max_completion_tokens")
	RequestTimeout int    `mapstructure:"request_timeout"`
	ThinkingLevel  string `mapstructure:"thinking_level"` // Extended thinking: off|low|medium|high|xhigh|adaptive
}

// Validate checks if the ModelConfig has all required fields.
func (c *ModelConfig) Validate() error {
	if c.ModelName == "" {
		return fmt.Errorf("model_name is required")
	}
	if c.Model == "" {
		return fmt.Errorf("model is required")
	}
	return nil
}

type GatewayConfig struct {
	Host string `mapstructure:"host" env:"CLOTHO_GATEWAY_HOST"`
	Port int    `mapstructure:"port" env:"CLOTHO_GATEWAY_PORT"`
}

type ToolConfig struct {
	Enabled bool `mapstructure:"enabled" env:"ENABLED"`
}

type BraveConfig struct {
	Enabled    bool   `mapstructure:"enabled"     env:"CLOTHO_TOOLS_WEB_BRAVE_ENABLED"`
	APIKey     string `mapstructure:"api_key"     env:"CLOTHO_TOOLS_WEB_BRAVE_API_KEY"`
	MaxResults int    `mapstructure:"max_results" env:"CLOTHO_TOOLS_WEB_BRAVE_MAX_RESULTS"`
}

type TavilyConfig struct {
	Enabled    bool   `mapstructure:"enabled"     env:"CLOTHO_TOOLS_WEB_TAVILY_ENABLED"`
	APIKey     string `mapstructure:"api_key"     env:"CLOTHO_TOOLS_WEB_TAVILY_API_KEY"`
	BaseURL    string `mapstructure:"base_url"    env:"CLOTHO_TOOLS_WEB_TAVILY_BASE_URL"`
	MaxResults int    `mapstructure:"max_results" env:"CLOTHO_TOOLS_WEB_TAVILY_MAX_RESULTS"`
}

type DuckDuckGoConfig struct {
	Enabled    bool `mapstructure:"enabled"     env:"CLOTHO_TOOLS_WEB_DUCKDUCKGO_ENABLED"`
	MaxResults int  `mapstructure:"max_results" env:"CLOTHO_TOOLS_WEB_DUCKDUCKGO_MAX_RESULTS"`
}

type PerplexityConfig struct {
	Enabled    bool   `mapstructure:"enabled"     env:"CLOTHO_TOOLS_WEB_PERPLEXITY_ENABLED"`
	APIKey     string `mapstructure:"api_key"     env:"CLOTHO_TOOLS_WEB_PERPLEXITY_API_KEY"`
	MaxResults int    `mapstructure:"max_results" env:"CLOTHO_TOOLS_WEB_PERPLEXITY_MAX_RESULTS"`
}

type SearXNGConfig struct {
	Enabled    bool   `mapstructure:"enabled"     env:"CLOTHO_TOOLS_WEB_SEARXNG_ENABLED"`
	BaseURL    string `mapstructure:"base_url"    env:"CLOTHO_TOOLS_WEB_SEARXNG_BASE_URL"`
	MaxResults int    `mapstructure:"max_results" env:"CLOTHO_TOOLS_WEB_SEARXNG_MAX_RESULTS"`
}

type GLMSearchConfig struct {
	Enabled bool   `mapstructure:"enabled"  env:"CLOTHO_TOOLS_WEB_GLM_ENABLED"`
	APIKey  string `mapstructure:"api_key"  env:"CLOTHO_TOOLS_WEB_GLM_API_KEY"`
	BaseURL string `mapstructure:"base_url" env:"CLOTHO_TOOLS_WEB_GLM_BASE_URL"`
	// SearchEngine specifies the search backend: "search_std" (default),
	// "search_pro", "search_pro_sogou", or "search_pro_quark".
	SearchEngine string `mapstructure:"search_engine" env:"CLOTHO_TOOLS_WEB_GLM_SEARCH_ENGINE"`
	MaxResults   int    `mapstructure:"max_results"   env:"CLOTHO_TOOLS_WEB_GLM_MAX_RESULTS"`
}

type WebToolsConfig struct {
	ToolConfig `                 envPrefix:"CLOTHO_TOOLS_WEB_"`
	Brave      BraveConfig      `                                json:"brave"`
	Tavily     TavilyConfig     `                                json:"tavily"`
	DuckDuckGo DuckDuckGoConfig `                                json:"duckduckgo"`
	Perplexity PerplexityConfig `                                json:"perplexity"`
	SearXNG    SearXNGConfig    `                                json:"searxng"`
	GLMSearch  GLMSearchConfig  `                                json:"glm_search"`
	// Proxy is an optional proxy URL for web tools (http/https/socks5/socks5h).
	// For authenticated proxies, prefer HTTP_PROXY/HTTPS_PROXY env vars instead of embedding credentials in config.
	Proxy           string `mapstructure:"proxy,omitempty"             env:"CLOTHO_TOOLS_WEB_PROXY"`
	FetchLimitBytes int64  `mapstructure:"fetch_limit_bytes,omitempty" env:"CLOTHO_TOOLS_WEB_FETCH_LIMIT_BYTES"`
}

type CronToolsConfig struct {
	ToolConfig         `    envPrefix:"CLOTHO_TOOLS_CRON_"`
	ExecTimeoutMinutes int `                                 env:"CLOTHO_TOOLS_CRON_EXEC_TIMEOUT_MINUTES" json:"exec_timeout_minutes"` // 0 means no timeout
}

type ExecConfig struct {
	ToolConfig          `         envPrefix:"CLOTHO_TOOLS_EXEC_"`
	EnableDenyPatterns  bool     `                                 env:"CLOTHO_TOOLS_EXEC_ENABLE_DENY_PATTERNS"  json:"enable_deny_patterns"`
	CustomDenyPatterns  []string `                                 env:"CLOTHO_TOOLS_EXEC_CUSTOM_DENY_PATTERNS"  json:"custom_deny_patterns"`
	CustomAllowPatterns []string `                                 env:"CLOTHO_TOOLS_EXEC_CUSTOM_ALLOW_PATTERNS" json:"custom_allow_patterns"`
}

type SkillsToolsConfig struct {
	ToolConfig            `                       envPrefix:"CLOTHO_TOOLS_SKILLS_"`
	Registries            SkillsRegistriesConfig `                                   json:"registries"`
	MaxConcurrentSearches int                    `                                   json:"max_concurrent_searches" env:"CLOTHO_TOOLS_SKILLS_MAX_CONCURRENT_SEARCHES"`
	SearchCache           SearchCacheConfig      `                                   json:"search_cache"`
}

type MediaCleanupConfig struct {
	ToolConfig `    envPrefix:"CLOTHO_MEDIA_CLEANUP_"`
	MaxAge     int `                                    env:"CLOTHO_MEDIA_CLEANUP_MAX_AGE"  json:"max_age_minutes"`
	Interval   int `                                    env:"CLOTHO_MEDIA_CLEANUP_INTERVAL" json:"interval_minutes"`
}

type ToolsConfig struct {
	AllowReadPaths  []string           `mapstructure:"allow_read_paths"  env:"CLOTHO_TOOLS_ALLOW_READ_PATHS"`
	AllowWritePaths []string           `mapstructure:"allow_write_paths" env:"CLOTHO_TOOLS_ALLOW_WRITE_PATHS"`
	Web             WebToolsConfig     `mapstructure:"web"`
	Cron            CronToolsConfig    `mapstructure:"cron"`
	Exec            ExecConfig         `mapstructure:"exec"`
	Skills          SkillsToolsConfig  `mapstructure:"skills"`
	MediaCleanup    MediaCleanupConfig `mapstructure:"media_cleanup"`
	MCP             MCPConfig          `mapstructure:"mcp"`
	AppendFile      ToolConfig         `mapstructure:"append_file"                                              envPrefix:"CLOTHO_TOOLS_APPEND_FILE_"`
	EditFile        ToolConfig         `mapstructure:"edit_file"                                                envPrefix:"CLOTHO_TOOLS_EDIT_FILE_"`
	FindSkills      ToolConfig         `mapstructure:"find_skills"                                              envPrefix:"CLOTHO_TOOLS_FIND_SKILLS_"`
	I2C             ToolConfig         `mapstructure:"i2c"                                                      envPrefix:"CLOTHO_TOOLS_I2C_"`
	InstallSkill    ToolConfig         `mapstructure:"install_skill"                                            envPrefix:"CLOTHO_TOOLS_INSTALL_SKILL_"`
	ListDir         ToolConfig         `mapstructure:"list_dir"                                                 envPrefix:"CLOTHO_TOOLS_LIST_DIR_"`
	Message         ToolConfig         `mapstructure:"message"                                                  envPrefix:"CLOTHO_TOOLS_MESSAGE_"`
	ReadFile        ToolConfig         `mapstructure:"read_file"                                                envPrefix:"CLOTHO_TOOLS_READ_FILE_"`
	Spawn           ToolConfig         `mapstructure:"spawn"                                                    envPrefix:"CLOTHO_TOOLS_SPAWN_"`
	SPI             ToolConfig         `mapstructure:"spi"                                                      envPrefix:"CLOTHO_TOOLS_SPI_"`
	Subagent        ToolConfig         `mapstructure:"subagent"                                                 envPrefix:"CLOTHO_TOOLS_SUBAGENT_"`
	WebFetch        ToolConfig         `mapstructure:"web_fetch"                                                envPrefix:"CLOTHO_TOOLS_WEB_FETCH_"`
	WriteFile       ToolConfig         `mapstructure:"write_file"                                               envPrefix:"CLOTHO_TOOLS_WRITE_FILE_"`
}

type SearchCacheConfig struct {
	MaxSize    int `mapstructure:"max_size"    env:"CLOTHO_SKILLS_SEARCH_CACHE_MAX_SIZE"`
	TTLSeconds int `mapstructure:"ttl_seconds" env:"CLOTHO_SKILLS_SEARCH_CACHE_TTL_SECONDS"`
}

type SkillsRegistriesConfig struct {
	ClawHub ClawHubRegistryConfig `mapstructure:"clawhub"`
}

type ClawHubRegistryConfig struct {
	Enabled         bool   `mapstructure:"enabled"           env:"CLOTHO_SKILLS_REGISTRIES_CLAWHUB_ENABLED"`
	BaseURL         string `mapstructure:"base_url"          env:"CLOTHO_SKILLS_REGISTRIES_CLAWHUB_BASE_URL"`
	AuthToken       string `mapstructure:"auth_token"        env:"CLOTHO_SKILLS_REGISTRIES_CLAWHUB_AUTH_TOKEN"`
	SearchPath      string `mapstructure:"search_path"       env:"CLOTHO_SKILLS_REGISTRIES_CLAWHUB_SEARCH_PATH"`
	SkillsPath      string `mapstructure:"skills_path"       env:"CLOTHO_SKILLS_REGISTRIES_CLAWHUB_SKILLS_PATH"`
	DownloadPath    string `mapstructure:"download_path"     env:"CLOTHO_SKILLS_REGISTRIES_CLAWHUB_DOWNLOAD_PATH"`
	Timeout         int    `mapstructure:"timeout"           env:"CLOTHO_SKILLS_REGISTRIES_CLAWHUB_TIMEOUT"`
	MaxZipSize      int    `mapstructure:"max_zip_size"      env:"CLOTHO_SKILLS_REGISTRIES_CLAWHUB_MAX_ZIP_SIZE"`
	MaxResponseSize int    `mapstructure:"max_response_size" env:"CLOTHO_SKILLS_REGISTRIES_CLAWHUB_MAX_RESPONSE_SIZE"`
}

// MCPServerConfig defines configuration for a single MCP server
type MCPServerConfig struct {
	// Enabled indicates whether this MCP server is active
	Enabled bool `mapstructure:"enabled"`
	// Command is the executable to run (e.g., "npx", "python", "/path/to/server")
	Command string `mapstructure:"command"`
	// Args are the arguments to pass to the command
	Args []string `mapstructure:"args,omitempty"`
	// Env are environment variables to set for the server process (stdio only)
	Env map[string]string `mapstructure:"env,omitempty"`
	// EnvFile is the path to a file containing environment variables (stdio only)
	EnvFile string `mapstructure:"env_file,omitempty"`
	// Type is "stdio", "sse", or "http" (default: stdio if command is set, sse if url is set)
	Type string `mapstructure:"type,omitempty"`
	// URL is used for SSE/HTTP transport
	URL string `mapstructure:"url,omitempty"`
	// Headers are HTTP headers to send with requests (sse/http only)
	Headers map[string]string `mapstructure:"headers,omitempty"`
}

// MCPConfig defines configuration for all MCP servers
type MCPConfig struct {
	ToolConfig `envPrefix:"CLOTHO_TOOLS_MCP_"`
	// Servers is a map of server name to server configuration
	Servers map[string]MCPServerConfig `mapstructure:"servers,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	// Detect file format by extension
	ext := strings.ToLower(filepath.Ext(path))

	if ext == ".toml" {
		// TOML format - use mapstructure tags
		var tomlMap any
		if err := toml.Unmarshal(data, &tomlMap); err != nil {
			return nil, err
		}

		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName:          "mapstructure",
			Result:           cfg,
			WeaklyTypedInput: true,
		})
		if err != nil {
			return nil, err
		}
		if err := decoder.Decode(tomlMap); err != nil {
			return nil, err
		}
	} else {
		// JSON format - convert to map then use mapstructure decoder
		var jsonMap any
		if err := json.Unmarshal(data, &jsonMap); err != nil {
			return nil, err
		}

		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
			TagName:          "mapstructure",
			Result:           cfg,
			WeaklyTypedInput: true,
		})
		if err != nil {
			return nil, err
		}
		if err := decoder.Decode(jsonMap); err != nil {
			return nil, err
		}
	}

	// TODO: re-enable env parsing after fixing
	// if err := env.Parse(cfg); err != nil {
	// 	return nil, err
	// }

	// Migrate legacy channel config fields to new unified structures
	cfg.migrateChannelConfigs()

	// Auto-migrate: if only legacy providers config exists, convert to model_list
	if len(cfg.ModelList) == 0 && cfg.HasProvidersConfig() {
		cfg.ModelList = ConvertProvidersToModelList(cfg)
	}

	// Validate model_list for uniqueness and required fields
	if err := cfg.ValidateModelList(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) migrateChannelConfigs() {
	// Discord: mention_only -> group_trigger.mention_only
	if c.Channels.Discord.MentionOnly && !c.Channels.Discord.GroupTrigger.MentionOnly {
		c.Channels.Discord.GroupTrigger.MentionOnly = true
	}

	// OneBot: group_trigger_prefix -> group_trigger.prefixes
	if len(c.Channels.OneBot.GroupTriggerPrefix) > 0 &&
		len(c.Channels.OneBot.GroupTrigger.Prefixes) == 0 {
		c.Channels.OneBot.GroupTrigger.Prefixes = c.Channels.OneBot.GroupTriggerPrefix
	}
}

func SaveConfig(path string, cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Use unified atomic write utility with explicit sync for flash storage reliability.
	return fileutil.WriteFileAtomic(path, data, 0o600)
}

func (c *Config) WorkspacePath() string {
	return expandHome(c.Agents.Defaults.Workspace)
}

func (c *Config) GetAPIKey() string {
	if c.Providers.OpenRouter.APIKey != "" {
		return c.Providers.OpenRouter.APIKey
	}
	if c.Providers.Anthropic.APIKey != "" {
		return c.Providers.Anthropic.APIKey
	}
	if c.Providers.OpenAI.APIKey != "" {
		return c.Providers.OpenAI.APIKey
	}
	if c.Providers.Gemini.APIKey != "" {
		return c.Providers.Gemini.APIKey
	}
	if c.Providers.Zhipu.APIKey != "" {
		return c.Providers.Zhipu.APIKey
	}
	if c.Providers.Groq.APIKey != "" {
		return c.Providers.Groq.APIKey
	}
	if c.Providers.VLLM.APIKey != "" {
		return c.Providers.VLLM.APIKey
	}
	if c.Providers.ShengSuanYun.APIKey != "" {
		return c.Providers.ShengSuanYun.APIKey
	}
	if c.Providers.Cerebras.APIKey != "" {
		return c.Providers.Cerebras.APIKey
	}
	return ""
}

func (c *Config) GetAPIBase() string {
	if c.Providers.OpenRouter.APIKey != "" {
		if c.Providers.OpenRouter.APIBase != "" {
			return c.Providers.OpenRouter.APIBase
		}
		return "https://openrouter.ai/api/v1"
	}
	if c.Providers.Zhipu.APIKey != "" {
		return c.Providers.Zhipu.APIBase
	}
	if c.Providers.VLLM.APIKey != "" && c.Providers.VLLM.APIBase != "" {
		return c.Providers.VLLM.APIBase
	}
	return ""
}

func expandHome(path string) string {
	if path == "" {
		return path
	}
	if path[0] == '~' {
		home, _ := os.UserHomeDir()
		if len(path) > 1 && path[1] == '/' {
			return home + path[1:]
		}
		return home
	}
	return path
}

// GetModelConfig returns the ModelConfig for the given model name.
// If multiple configs exist with the same model_name, it uses round-robin
// selection for load balancing. Returns an error if the model is not found.
func (c *Config) GetModelConfig(modelName string) (*ModelConfig, error) {
	matches := c.findMatches(modelName)
	if len(matches) == 0 {
		return nil, fmt.Errorf("model %q not found in model_list or providers", modelName)
	}
	if len(matches) == 1 {
		return &matches[0], nil
	}

	// Multiple configs - use round-robin for load balancing
	idx := rrCounter.Add(1) % uint64(len(matches))
	return &matches[idx], nil
}

// findMatches finds all ModelConfig entries with the given model_name.
func (c *Config) findMatches(modelName string) []ModelConfig {
	var matches []ModelConfig
	for i := range c.ModelList {
		if c.ModelList[i].ModelName == modelName {
			matches = append(matches, c.ModelList[i])
		}
	}
	return matches
}

// HasProvidersConfig checks if any provider in the old providers config has configuration.
func (c *Config) HasProvidersConfig() bool {
	return !c.Providers.IsEmpty()
}

// ValidateModelList validates all ModelConfig entries in the model_list.
// It checks that each model config is valid.
// Note: Multiple entries with the same model_name are allowed for load balancing.
func (c *Config) ValidateModelList() error {
	for i := range c.ModelList {
		if err := c.ModelList[i].Validate(); err != nil {
			return fmt.Errorf("model_list[%d]: %w", i, err)
		}
	}
	return nil
}

func (t *ToolsConfig) IsToolEnabled(name string) bool {
	switch name {
	case "web":
		return t.Web.Enabled
	case "cron":
		return t.Cron.Enabled
	case "exec":
		return t.Exec.Enabled
	case "skills":
		return t.Skills.Enabled
	case "media_cleanup":
		return t.MediaCleanup.Enabled
	case "append_file":
		return t.AppendFile.Enabled
	case "edit_file":
		return t.EditFile.Enabled
	case "find_skills":
		return t.FindSkills.Enabled
	case "i2c":
		return t.I2C.Enabled
	case "install_skill":
		return t.InstallSkill.Enabled
	case "list_dir":
		return t.ListDir.Enabled
	case "message":
		return t.Message.Enabled
	case "read_file":
		return t.ReadFile.Enabled
	case "spawn":
		return t.Spawn.Enabled
	case "spi":
		return t.SPI.Enabled
	case "subagent":
		return t.Subagent.Enabled
	case "web_fetch":
		return t.WebFetch.Enabled
	case "write_file":
		return t.WriteFile.Enabled
	case "mcp":
		return t.MCP.Enabled
	default:
		return true
	}
}
