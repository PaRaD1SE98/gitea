// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrations

import (
	"net/http"
	"testing"

	"code.gitea.io/gitea/models/unittest"
	user_model "code.gitea.io/gitea/models/user"

	"github.com/stretchr/testify/assert"
)

func TestAPIReposRaw(t *testing.T) {
	defer prepareTestEnv(t)()
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	// Login as User2.
	session := loginUser(t, user.Name)
	token := getTokenForLoggedInUser(t, session)

	for _, ref := range [...]string{
		"master", // Branch
		"v1.1",   // Tag
		"65f1bf27bc3bf70f64657658635e66094edbcb4d", // Commit
	} {
		req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/raw/%s/README.md?token="+token, user.Name, ref)
		resp := session.MakeRequest(t, req, http.StatusOK)
		assert.EqualValues(t, "file", resp.Header().Get("x-gitea-object-type"))
	}
	// Test default branch
	req := NewRequestf(t, "GET", "/api/v1/repos/%s/repo1/raw/README.md?token="+token, user.Name)
	resp := session.MakeRequest(t, req, http.StatusOK)
	assert.EqualValues(t, "file", resp.Header().Get("x-gitea-object-type"))
}
