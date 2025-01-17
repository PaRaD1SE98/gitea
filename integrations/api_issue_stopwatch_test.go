// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrations

import (
	"net/http"
	"testing"

	"code.gitea.io/gitea/models/db"
	issues_model "code.gitea.io/gitea/models/issues"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/unittest"
	user_model "code.gitea.io/gitea/models/user"
	api "code.gitea.io/gitea/modules/structs"

	"github.com/stretchr/testify/assert"
)

func TestAPIListStopWatches(t *testing.T) {
	defer prepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	session := loginUser(t, owner.Name)
	token := getTokenForLoggedInUser(t, session)
	req := NewRequestf(t, "GET", "/api/v1/user/stopwatches?token=%s", token)
	resp := session.MakeRequest(t, req, http.StatusOK)
	var apiWatches []*api.StopWatch
	DecodeJSON(t, resp, &apiWatches)
	stopwatch := unittest.AssertExistsAndLoadBean(t, &issues_model.Stopwatch{UserID: owner.ID})
	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: stopwatch.IssueID})
	if assert.Len(t, apiWatches, 1) {
		assert.EqualValues(t, stopwatch.CreatedUnix.AsTime().Unix(), apiWatches[0].Created.Unix())
		assert.EqualValues(t, issue.Index, apiWatches[0].IssueIndex)
		assert.EqualValues(t, issue.Title, apiWatches[0].IssueTitle)
		assert.EqualValues(t, repo.Name, apiWatches[0].RepoName)
		assert.EqualValues(t, repo.OwnerName, apiWatches[0].RepoOwnerName)
		assert.Greater(t, apiWatches[0].Seconds, int64(0))
	}
}

func TestAPIStopStopWatches(t *testing.T) {
	defer prepareTestEnv(t)()

	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 2})
	_ = issue.LoadRepo(db.DefaultContext)
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: issue.Repo.OwnerID})
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	req := NewRequestf(t, "POST", "/api/v1/repos/%s/%s/issues/%d/stopwatch/stop?token=%s", owner.Name, issue.Repo.Name, issue.Index, token)
	session.MakeRequest(t, req, http.StatusCreated)
	session.MakeRequest(t, req, http.StatusConflict)
}

func TestAPICancelStopWatches(t *testing.T) {
	defer prepareTestEnv(t)()

	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})
	_ = issue.LoadRepo(db.DefaultContext)
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: issue.Repo.OwnerID})
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 1})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	req := NewRequestf(t, "DELETE", "/api/v1/repos/%s/%s/issues/%d/stopwatch/delete?token=%s", owner.Name, issue.Repo.Name, issue.Index, token)
	session.MakeRequest(t, req, http.StatusNoContent)
	session.MakeRequest(t, req, http.StatusConflict)
}

func TestAPIStartStopWatches(t *testing.T) {
	defer prepareTestEnv(t)()

	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 3})
	_ = issue.LoadRepo(db.DefaultContext)
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: issue.Repo.OwnerID})
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	req := NewRequestf(t, "POST", "/api/v1/repos/%s/%s/issues/%d/stopwatch/start?token=%s", owner.Name, issue.Repo.Name, issue.Index, token)
	session.MakeRequest(t, req, http.StatusCreated)
	session.MakeRequest(t, req, http.StatusConflict)
}
