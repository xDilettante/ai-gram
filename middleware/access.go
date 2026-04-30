package middleware

import (
	"context"
	stderrors "errors"

	"ai-gram/dispatch"
	"ai-gram/telegram"
)

// AccessMode controls how Access middleware authorizes incoming updates.
type AccessMode string

const (
	// AccessModeAdmin allows only configured admins, users, or chats.
	AccessModeAdmin AccessMode = "admin"
	// AccessModePublic allows every update.
	AccessModePublic AccessMode = "public"
	// AccessModeOff disables access checks and allows every update.
	AccessModeOff AccessMode = "off"
)

// AccessConfig configures static access control for update handlers.
type AccessConfig struct {
	// Mode selects the access strategy. Empty mode is treated as admin mode.
	Mode AccessMode

	// AdminUserIDs are user IDs that always pass admin-mode checks.
	AdminUserIDs []int64
	// AllowedUserIDs are additional user IDs allowed in admin mode.
	AllowedUserIDs []int64
	// AllowedChatIDs are chat IDs allowed in admin mode.
	AllowedChatIDs []int64

	// OnDeny is called for denied updates. Nil means denied updates are ignored.
	OnDeny func(context.Context, telegram.Update) error
}

// AccessPolicy decides whether an update is allowed to reach the next handler.
type AccessPolicy interface {
	IsAllowed(context.Context, telegram.Update) bool
}

// AccessFunc adapts a function to AccessPolicy.
type AccessFunc func(context.Context, telegram.Update) bool

// IsAllowed calls f when f is non-nil. A nil AccessFunc denies the update.
func (f AccessFunc) IsAllowed(ctx context.Context, update telegram.Update) bool {
	if f == nil {
		return false
	}
	return f(ctx, update)
}

// Access returns middleware that authorizes updates using a static config.
func Access(config AccessConfig) dispatch.Middleware {
	return AccessWithPolicy(accessConfigPolicy{config: config}, config.OnDeny)
}

// AccessWithPolicy returns middleware that authorizes updates using policy.
// A nil policy denies every update and calls onDeny when it is provided.
func AccessWithPolicy(policy AccessPolicy, onDeny func(context.Context, telegram.Update) error) dispatch.Middleware {
	return func(next dispatch.Handler) dispatch.Handler {
		if next == nil {
			return dispatch.HandlerFunc(func(context.Context, telegram.Update) error {
				return stderrors.New("handler is required")
			})
		}

		return dispatch.HandlerFunc(func(ctx context.Context, update telegram.Update) error {
			if policy == nil || !policy.IsAllowed(ctx, update) {
				if onDeny != nil {
					return onDeny(ctx, update)
				}
				return nil
			}

			return next.HandleUpdate(ctx, update)
		})
	}
}

type accessConfigPolicy struct {
	config AccessConfig
}

func (p accessConfigPolicy) IsAllowed(ctx context.Context, update telegram.Update) bool {
	return accessConfigAllows(p.config, update)
}

func accessConfigAllows(config AccessConfig, update telegram.Update) bool {
	switch config.Mode {
	case "", AccessModeAdmin:
		return isAllowedByAdminConfig(config, update)
	case AccessModePublic, AccessModeOff:
		return true
	default:
		return false
	}
}

func isAllowedByAdminConfig(config AccessConfig, update telegram.Update) bool {
	if user := update.EffectiveUser(); user != nil {
		if containsInt64(config.AdminUserIDs, user.ID) || containsInt64(config.AllowedUserIDs, user.ID) {
			return true
		}
	}
	if chat := update.EffectiveChat(); chat != nil && containsInt64(config.AllowedChatIDs, chat.ID) {
		return true
	}

	return false
}

func containsInt64(values []int64, target int64) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
