package git

import (
	"fmt"
	"strings"
	"time"
)

type CommitInfo struct {
	Author         string
	AuthorEmail    string
	Committer      string
	CommitterEmail string
	Subject        string
	Timestamp      time.Time
}

// keep same order as struct
var commitInfoArgs = []string{
	"%an", // Author Name
	"%ae", // Author Email
	"%cn", // Committer Name
	"%ce", // Comitter Email
	"%s",  // Subject
	"%cI", // Commit Time
}

const commitInfoSep = ";-;"

func (h Handler) FetchCommitInfo(hash string) (CommitInfo, error) {
	res, err := h.exec("show", "-s", fmt.Sprintf(`--format=%s`, strings.Join(commitInfoArgs, commitInfoSep)), hash)
	if err != nil {
		return CommitInfo{}, err
	}

	parts := strings.Split(res, commitInfoSep)

	if len(parts) != len(commitInfoArgs) {
		return CommitInfo{}, fmt.Errorf("unexpected number of results received, expected %d got %d", len(commitInfoArgs), len(parts))
	}

	ts, err := time.Parse(time.RFC3339, parts[5])
	if err != nil {
		return CommitInfo{}, fmt.Errorf("could not parse timestamp: %w", err)
	}

	return CommitInfo{
		Author:         parts[0],
		AuthorEmail:    parts[1],
		Committer:      parts[2],
		CommitterEmail: parts[3],
		Subject:        parts[4],
		Timestamp:      ts,
	}, nil
}
