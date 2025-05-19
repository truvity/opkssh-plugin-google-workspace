package opksshplugingoogleworkspace

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type (
	Member struct {
		Id     string `json:"id"`
		Email  string `json:"email"` // user's email
		Status string `json:"status"`
		Type   string `json:"type"`
	}

	GroupMembersFetcher interface {
		GroupMembers(ctx context.Context, logger *slog.Logger, groupEmail string) ([]*Member, error)
	}
)

func Verify(ctx context.Context, logger *slog.Logger, fetcher GroupMembersFetcher, config *Config, request *Request) (bool, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	startTime := time.Now()

	result, err := func() (bool, error) {
		inform := logger.With(
			slog.String("principal", request.Principal),
			slog.String("email", request.Email),
			slog.Bool("email_verified", request.EmailVerified),
			slog.String("aud", request.ClientID),
		)

		if !request.EmailVerified {
			const decision = "deny"
			const reason = "email not verified"
			inform.WarnContext(ctx, decision,
				slog.String("decision", decision),
				slog.String("reason", reason),
			)
			return false, nil
		}

		if request.ClientID != config.Google.OAuth.ClientID {
			const decision = "deny"
			const reason = "client_id and aud mismatch"
			inform.WarnContext(ctx, decision,
				slog.String("decision", decision),
				slog.String("reason", reason),
				slog.String("client_id", config.Google.OAuth.ClientID),
			)
			return false, nil
		}

		policy := config.Policy[request.Principal]

		if policy == nil {
			const decision = "deny"
			const reason = "principal does not have any policy"
			inform.WarnContext(ctx, decision,
				slog.String("decision", decision),
				slog.String("reason", reason),
			)
			return false, nil
		}

		for index := range policy.User {
			userEmail := policy.User[index]
			if request.Email == userEmail {
				const decision = "allow"
				const reason = "user's policy of principal"
				inform.InfoContext(ctx, decision,
					slog.String("decision", decision),
					slog.String("reason", reason),
				)
				return true, nil
			}
		}

		for index := range policy.Group {
			groupEmail := policy.Group[index]
			memberList, err := fetcher.GroupMembers(ctx, logger, groupEmail)
			if err != nil {
				const message = "failed to fetch group's members"
				inform.ErrorContext(ctx, message,
					slog.String("group", groupEmail),
					slog.Any("error", err),
				)
				err = fmt.Errorf("%s group %s %w",
					message,
					groupEmail,
					err,
				)
				return false, err
			}
			for _, member := range memberList {
				if request.Email == member.Email {
					const decision = "allow"
					const reason = "group's policy of principal"
					inform.InfoContext(ctx, decision,
						slog.String("decision", decision),
						slog.String("reason", reason),
						slog.String("group", groupEmail),
					)
					return true, nil
				}
			}
		}

		const message = "deny - no policy to allow"
		inform.WarnContext(ctx, message)
		return false, nil
	}()
	if !result || err != nil {
		// we need this delay to avoid timing attack based on negative resulsts
		round := time.Second * 5
		delay := (time.Since(startTime) + round).Truncate(round)
		time.Sleep(delay)
	}
	return result, err

}
