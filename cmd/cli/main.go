package main

import (
	"encoding/json"
	"fmt"
	"github.com/owenrumney/go-github-pr-commenter/commenter"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Starting the github commenter")

	token := os.Getenv("INPUT_GITHUB_TOKEN")

	if len(token) == 0 {
		fail("the INPUT_GITHUB_TOKEN has not been set")
	}

	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	split := strings.Split(githubRepository, "/")
	if len(split) != 2 {
		fail(fmt.Sprintf("unexpected value for GITHUB_REPOSITORY. Expected <organisation/name>, found %v", split))
	}
	owner := split[0]
	repo := split[1]

	fmt.Printf("Working in repository %s\n", repo)

	prNo, err := extractPullRequestNumber()
	if err != nil {
		fmt.Println("Not a PR, nothing to comment on, exiting")
		return
	}
	fmt.Printf("Working in PR %v\n", prNo)

	results, err := loadResultsFile()
	if err != nil {
		fail(fmt.Sprintf("failed to load results. %s", err.Error()))
	}

	if len(results) == 0 {
		fmt.Println("No issues found.")
		os.Exit(0)
	}

	fmt.Printf("Horusec found %v issues\n", len(results))

	c, err := commenter.NewCommenter(token, owner, repo, prNo)
	if err != nil {
		fail(fmt.Sprintf("could not connect to GitHub (%s)", err.Error()))
	}

	workspacePath := fmt.Sprintf("%s/", os.Getenv("GITHUB_WORKSPACE"))
	fmt.Printf("Working in GITHUB_WORKSPACE %s\n", workspacePath)

	workingDir := os.Getenv("INPUT_WORKING_DIRECTORY")
	if workingDir != "" {
		workingDir = strings.TrimPrefix(workingDir, "./")
		workingDir = strings.TrimSuffix(workingDir, "/") + "/"
	}

	var errMessages []string
	var validCommentWritten bool

	for _, result := range results {
		result.File = path.Join(workingDir, strings.ReplaceAll(result.File, workspacePath, ""))
		comment := generateErrorMessage(result)
		fmt.Printf("Preparing comment for violation of rule %v in %v\n", result.RuleId, result.File)
		if result.IsMultiLine() {
			err = c.WriteMultiLineComment(result.File, comment, result.StartLine(), result.EndLine())
		} else {
			err = c.WriteLineComment(result.File, comment, result.StartLine())
		}
		if err != nil {
			// don't error if its simply that the comments aren't valid for the PR
			switch err.(type) {
			case commenter.CommentAlreadyWrittenError:
				fmt.Println("Ignoring - comment already written")
				validCommentWritten = true
			case commenter.CommentNotValidError:
				fmt.Println("Ignoring - change not part of the current PR")
				continue
			default:
				errMessages = append(errMessages, err.Error())
			}
		} else {
			validCommentWritten = true
		}
	}

	if len(errMessages) > 0 {
		fmt.Printf("There were %d errors:\n", len(errMessages))
		for _, err := range errMessages {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	if validCommentWritten || len(errMessages) > 0 {
		if softFail, ok := os.LookupEnv("INPUT_SOFT_FAIL_COMMENTER"); ok && strings.ToLower(softFail) == "true" {
			return
		}
		os.Exit(1)
	}
}

func generateErrorMessage(result result) string {
	str := "## Oops!! 🧐\n"
	str += fmt.Sprintf("Language: **%s**\n\n", result.Language)
	str += fmt.Sprintf("Severity: **%s**\n\n", result.Severity)
	str += fmt.Sprintf("Security Tool: **%s**\n\n", result.SecurityTool)
	str += fmt.Sprintf("Type: **%s**\n\n", result.Type)
	str += "Description:\n"
	str += fmt.Sprintf(">*%s", result.Details)
	return str
}

func extractPullRequestNumber() (int, error) {
	githubEventFile := "/github/workflow/event.json"
	file, err := ioutil.ReadFile(githubEventFile)
	if err != nil {
		fail(fmt.Sprintf("GitHub event payload not found in %s", githubEventFile))
		return -1, err
	}

	var data interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return -1, err
	}
	payload := data.(map[string]interface{})

	prNumber, err := strconv.Atoi(fmt.Sprintf("%v", payload["number"]))
	if err != nil {
		return 0, fmt.Errorf("not a valid PR")
	}
	return prNumber, nil
}

func fail(err string) {
	fmt.Printf("Error: %s\n", err)
	os.Exit(-1)
}
