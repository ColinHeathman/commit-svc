package frequency

import (
	"container/heap"

	"github.com/colinheathman/commit-svc/pkg/dateutil"
	"github.com/colinheathman/commit-svc/pkg/github"
)

type CommitCount struct {
	Name        string `json:"name"`
	CommitCount uint64 `json:"commits"`
}

// CalculateMostFrequentUsers calculates the total number of commits for each username, then returns the top 5
func CalculateMostFrequentUsers(reader github.CommitReader, dateRange dateutil.DateRange) ([]CommitCount, error) {

	// Count up the commits for each login username, and save the names at the same time
	commitCounts := make(map[string]uint64)
	names := make(map[string]string)

	for page := range reader.StreamCommits(dateRange) {
		if page.Err != nil {
			return nil, page.Err
		}
		for _, commit := range page.Commits {

			// Only add authors to the list who haven't yet been seen
			_, ok := commitCounts[commit.Author.Login]
			if !ok {
				commitCounts[commit.Author.Login] = 1
				names[commit.Author.Login] = commit.Commit.Author.Name
			} else {
				commitCounts[commit.Author.Login] = commitCounts[commit.Author.Login] + 1
			}

		}
	}

	// Heap sort the result by most commits
	commitCountHeap := &CommitCountHeap{
		[]CommitCount{},
	}

	// heap push each commits tuple
	for login, commits := range commitCounts {
		heap.Push(commitCountHeap, CommitCount{
			Name:        names[login],
			CommitCount: commits,
		})
	}

	result := []CommitCount{}

	// heap pop the top 5 commit tuples
	for len(commitCountHeap.Heap) > 0 && len(result) < 5 {
		result = append(result, heap.Pop(commitCountHeap).(CommitCount))
	}

	return result, nil
}
