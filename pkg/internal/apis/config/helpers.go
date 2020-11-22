/*
Copyright © 2020 Yashvardhan Kukreja <yash.kukreja.98@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package config

import (
	"fmt"
	"strings"
	"whatsupstream/pkg/internal/apis/config/github"
)

func FetchNotificationData(config Config) (Notification, error) {
	apiUrls := generateAllApiUrlsForConfig(config)
	allIssues := []github.Issue{}
	for i, apiUrl := range apiUrls {
		issues, err := github.FetchTopIssues(apiUrl, config.IssueConfigs[i].MaxIssuesCount)
		if err != nil {
			return Notification{}, fmt.Errorf("error occurred while fetching data for notification: %w", err)
		}
		allIssues = append(allIssues, issues...)
	}
	return Notification{
		Issues: allIssues,
	}, nil
}

func generateAllApiUrlsForConfig(config Config) []string {
	var urls []string
	for _, issueConfig := range config.IssueConfigs {
		urls = append(urls, generateApiUrlForIssueConfig(issueConfig))
	}
	return urls
}

func generateApiUrlForIssueConfig(issueConfig IssueConfig) string {
	var params []string
	params = append(params, "state="+string(issueConfig.State))

	if issueConfig.Assignee != "" && issueConfig.Assignee != "*" {
		params = append(params, "assignee="+issueConfig.Assignee)
	}
	if issueConfig.Creator != "" {
		params = append(params, "creator="+issueConfig.Creator)
	}
	if len(issueConfig.Labels) > 0 {
		params = append(params, "labels="+strings.Join(issueConfig.Labels, ","))
	}
	if issueConfig.Since != "" {
		params = append(params, "since="+issueConfig.Since)
	}

	return fmt.Sprintf("%s/%s/%s/issues?%s",
		API_BASE_URL_REPOS, issueConfig.Owner, issueConfig.RepoName, strings.Join(params, "&"))
}
