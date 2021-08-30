package users

import (
	"github.com/colinheathman/commit-svc/pkg/dateutil"
	"github.com/colinheathman/commit-svc/pkg/github"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CalculateUsers determines the unique users in a range of Github commits
func CalculateUsers(reader github.CommitReader, dateRange dateutil.DateRange) ([]User, error) {

	// Iterate through all the commits
	// Save unique Github login usernames to a set
	// On saving a new to the set, save the author name and email to a result list

	usernames := make(map[string]bool)
	result := []User{}

	for page := range reader.StreamCommits(dateRange) {
		if page.Err != nil {
			return nil, page.Err
		}
		for _, commit := range page.Commits {

			// Only add authors to the list who haven't yet been seen
			_, ok := usernames[commit.Author.Login]
			if !ok {
				usernames[commit.Author.Login] = true

				result = append(result, User{
					Name:  commit.Commit.Author.Name,
					Email: commit.Commit.Author.Email,
				})
			}

		}
	}

	return result, nil

}
