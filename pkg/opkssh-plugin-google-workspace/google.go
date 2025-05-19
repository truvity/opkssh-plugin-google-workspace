package opksshplugingoogleworkspace

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"golang.org/x/oauth2/jwt"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

type (
	ConfigGoogleOAuthApp struct {
		ClientID string `json:"client_id" yaml:"client_id"`
	}

	ConfigGoogleWorkspace struct {
		CustomerID string `json:"customer_id" yaml:"customer_id"`
	}

	ConfigGoogleServiceAccount struct {
		Email   string      `json:"email"    yaml:"email"`
		KeyFile string      `json:"key_file" yaml:"key_file"`
		Key     *jwt.Config `json:"-"        yaml:"-"`
	}

	ConfigGoogle struct {
		OAuth          ConfigGoogleOAuthApp       `json:"oauth"           yaml:"oauth"`
		ServiceAccount ConfigGoogleServiceAccount `json:"service_account" yaml:"service_account"`
		Workspace      ConfigGoogleWorkspace      `json:"workspace"       yaml:"workspace"`
	}

	GoogleFetcher struct {
		Token *jwt.Config
	}
)

var (
	_ GroupMembersFetcher = &GoogleFetcher{}
)

func NewGooglFetcher(config ConfigGoogleServiceAccount) *GoogleFetcher {
	return &GoogleFetcher{
		Token: config.Key,
	}
}

func (gf *GoogleFetcher) GroupMembers(ctx context.Context, logger *slog.Logger, groupEmail string) ([]*Member, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger.DebugContext(ctx, "create Google Workspace Admin Service",
		slog.Any("service_account", gf.Token.Email),
	)

	svc, err := admin.NewService(ctx,
		option.WithTokenSource(gf.Token.TokenSource(ctx)),
		option.WithScopes(admin.AdminDirectoryGroupMemberReadonlyScope),
	)
	if err != nil {
		const message = "failed to create Google Workspace Admin Service"
		logger.ErrorContext(ctx,
			message,
			slog.Any("service_account", gf.Token.Email),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s service account %s %w",
			message,
			gf.Token.Email,
			err,
		)
		return nil, err
	}

	call := svc.Members.List(groupEmail)
	call = call.IncludeDerivedMembership(true)

	logger.DebugContext(ctx, "fetch group's members",
		slog.Any("service_account", gf.Token.Email),
		slog.Any("group", groupEmail),
	)

	var set = make(map[string]*Member)
	err = call.Pages(ctx, func(members *admin.Members) error {
		for _, member := range members.Members {
			set[member.Email] = &Member{
				Id:     member.Id,
				Email:  member.Email,
				Status: member.Status,
				Type:   member.Type,
			}
		}
		return nil
	})
	if err != nil {
		const message = "failed to fetch group's member"
		logger.ErrorContext(ctx, message,
			slog.Any("service_account", gf.Token.Email),
			slog.Any("group", groupEmail),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s service account %s group %s %w",
			message,
			gf.Token.Email,
			groupEmail,
			err,
		)
		return nil, err
	}

	var result = make([]*Member, 0, len(set))
	for _, member := range set {
		result = append(result, member)
	}
	sort.Slice(result, func(i, j int) bool {
		left, right := result[i], result[j]
		return strings.Compare(left.Email, right.Email) < 0
	})

	logger.InfoContext(ctx, "fetch group's members completed",
		slog.String("group", groupEmail),
		slog.Int("members_count", len(result)),
	)

	return result, nil
}
