package opksshplugingoogleworkspace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

type (
	ConfigCache struct {
		Path     *string        `json:"path,omitempty"     yaml:"path,omitempty"`
		Duration *time.Duration `json:"duration,omitempty" yaml:"duration,omitempty"`
	}

	CacheFetcher struct {
		customerId string
		now        time.Time
		deadline   time.Time
		fetcher    GroupMembersFetcher
		path       string

		mutex sync.Mutex
		info  *Info
		lock  *flock.Flock
	}
)

var (
	_ GroupMembersFetcher = &CacheFetcher{}
)

func NewCacheFetcher(config ConfigCache, customerId string, fetcher GroupMembersFetcher) *CacheFetcher {
	now := time.Now()
	deadline := now.Add(-1 * *config.Duration)
	path := *config.Path
	return &CacheFetcher{
		// immutable
		customerId: customerId,
		now:        now,
		deadline:   deadline,
		fetcher:    fetcher,
		path:       path,
		// volatile
		lock: flock.New(path + ".filelock"),
	}
}

func (c *CacheFetcher) GroupMembers(
	ctx context.Context,
	logger *slog.Logger,
	groupEmail string,
) ([]*Member, error) {
	result := c.get(ctx, logger, groupEmail)
	if result != nil {
		return result, nil
	}
	members, err := c.fetcher.GroupMembers(ctx, logger, groupEmail)
	if err != nil {
		return nil, err
	}
	err = c.add(ctx, logger, groupEmail, members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (c *CacheFetcher) get(ctx context.Context, logger *slog.Logger, groupEmail string) []*Member {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// check if cache loaded
	if c.info == nil {
		// cache not loaded, load cache
		c.lock.Lock()
		defer c.lock.Unlock()
		c.unsafeLoad(ctx, logger)
	}
	// search in cache
	group := c.info.GetCustomer(c.customerId).GetGroup(c.deadline, groupEmail)
	if group == nil {
		// cache miss
		return nil
	}

	// return result
	result := make([]*Member, 0)
	for _, m := range group.Members {
		result = append(result, m)
	}
	sort.Slice(result, func(i, j int) bool {
		left, right := result[i], result[j]
		return strings.Compare(left.Email, right.Email) < 0
	})

	return result
}

func (c *CacheFetcher) add(ctx context.Context, logger *slog.Logger, groupEmail string, members []*Member) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lock.Lock()
	defer c.lock.Unlock()
	c.unsafeLoad(ctx, logger)
	for index := range members {
		member := members[index]
		c.info.AddCustomer(c.customerId).AddGroup(c.now, groupEmail).AddMember(member)
	}
	return c.unsafeSave(ctx, logger)

}

func (c *CacheFetcher) unsafeSave(ctx context.Context, logger *slog.Logger) error {
	raw, err := json.MarshalIndent(c.info, "", "  ")
	if err != nil {
		const message = "failed to serialize cache file"
		logger.ErrorContext(ctx, message,
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s %w",
			message,
			err,
		)
		return err
	}

	// create cache dir
	parentPath := filepath.Dir(c.path)
	err = os.MkdirAll(parentPath, 0600)
	if err != nil {
		const message = "failed to create cache directory"
		logger.ErrorContext(ctx, message,
			slog.String("path", parentPath),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s directory %s problem %w",
			message,
			parentPath,
			err,
		)
		return err
	}

	fileName := filepath.Base(c.path)

	// create temporary file
	tempPattern := fileName + ".*"
	tempFile, err := os.CreateTemp(parentPath, tempPattern)
	if err != nil {
		const message = "failed to create temporary cache file"
		logger.ErrorContext(ctx, message,
			slog.String("dir", parentPath),
			slog.String("pattern", tempPattern),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s dir %s pattern %s problem %w",
			message,
			parentPath,
			tempPattern,
			err,
		)
		return err
	}

	// get path to temporary file
	pathTemp := tempFile.Name()

	// remove file at the end
	defer func() {
		if tempFile == nil {
			return
		}
		err := os.Remove(pathTemp)
		if err != nil {
			const message = "failed to remove temporary file"
			logger.ErrorContext(ctx, message,
				slog.String("path", pathTemp),
				slog.Any("error", err),
			)
		}
	}()

	// write temporary file
	_, err = io.Copy(tempFile, bytes.NewReader(raw))
	if err != nil {
		const message = "failed to write temporary cache file"
		logger.ErrorContext(ctx, message,
			slog.String("path", pathTemp),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s path %s problem %w",
			message,
			pathTemp,
			err,
		)
		return err
	}

	// close temporary file
	err = tempFile.Close()
	if err != nil {
		const message = "failed to close temporary cache file"
		logger.ErrorContext(ctx, message,
			slog.String("path", pathTemp),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s path %s problem %w",
			message,
			pathTemp,
			err,
		)
		return err
	}

	// rename temporary file
	err = os.Rename(pathTemp, c.path)
	if err != nil {
		const message = "failed to rename temporary cache file"
		logger.ErrorContext(ctx, message,
			slog.String("from", pathTemp),
			slog.String("to", c.path),
			slog.Any("error", err),
		)
		err = fmt.Errorf("%s from %s to %s problem %w",
			message, pathTemp, c.path, err)
		return err
	}
	tempFile = nil

	logger.InfoContext(ctx, "cache saved",
		slog.String("path", c.path),
	)

	return nil
}

func (c *CacheFetcher) unsafeLoad(ctx context.Context, logger *slog.Logger) {
	defer func() {
		if c.info == nil {
			c.info = &Info{}
		}
	}()
	raw, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			const message = "cache file does not exist"
			logger.InfoContext(ctx, message,
				slog.String("path", c.path),
			)
		} else {
			const message = "failed to read cache file"
			logger.ErrorContext(ctx, message,
				slog.String("path", c.path),
			)
		}
		return
	}

	var external Info
	err = json.Unmarshal(raw, &external)
	if err != nil {
		const message = "failed to parse cache file"
		logger.ErrorContext(ctx, message,
			slog.String("path", c.path),
			slog.Any("error", err),
		)
		return
	}

	external.Merge(c.info)
	c.info = &external

	logger.DebugContext(ctx, "cache loaded",
		slog.String("path", c.path),
	)
}
