package framework

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGithubRepositoryPullRequestResource(t *testing.T) {
	t.Run("manages the pull request lifecycle", func(t *testing.T) {
		randomID := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

		config := fmt.Sprintf(`
			resource "github_repository" "test" {
				name      = "tf-acc-test-%s"
				auto_init = true
			}

			resource "github_branch" "test" {
				repository    = github_repository.test.name
				branch        = "test"
				source_branch = github_repository.test.default_branch
			}

			resource "github_repository_file" "test" {
				repository     = github_repository.test.name
				branch         = github_branch.test.branch
				file           = "test"
				content        = "bar"
			}

			resource "github_repository_pull_request" "test" {
				base_repository = github_repository_file.test.repository
				base_ref        = github_repository.test.default_branch
				head_ref        = github_branch.test.branch
				title           = "test title"
				body            = "test body"
			}
		`, randomID)

		const resourceName = "github_repository_pull_request.test"

		check := resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(
				resourceName, "base_repository",
				fmt.Sprintf("tf-acc-test-%s", randomID),
			),
			resource.TestCheckResourceAttr(resourceName, "base_ref", "main"),
			resource.TestCheckResourceAttr(resourceName, "head_ref", "test"),
			resource.TestCheckResourceAttr(resourceName, "title", "test title"),
			resource.TestCheckResourceAttr(resourceName, "body", "test body"),
			resource.TestCheckResourceAttr(resourceName, "maintainer_can_modify", "false"),
			resource.TestCheckResourceAttrSet(resourceName, "base_sha"),
			resource.TestCheckResourceAttr(resourceName, "draft", "false"),
			resource.TestCheckResourceAttrSet(resourceName, "head_sha"),
			resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
			resource.TestCheckResourceAttrSet(resourceName, "number"),
			resource.TestCheckResourceAttrSet(resourceName, "opened_at"),
			resource.TestCheckResourceAttrSet(resourceName, "opened_by"),
			resource.TestCheckResourceAttr(resourceName, "state", "open"),
			resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
		)

		testCase := func(t *testing.T, mode string) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t, mode) },
				ProtoV6ProviderFactories: testAccMuxedProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config: config,
						Check:  check,
					},
					{
						ResourceName:      resourceName,
						ImportState:       true,
						ImportStateVerify: true,
					},
				},
			})
		}

		t.Run("with an anonymous account", func(t *testing.T) {
			t.Skip("anonymous account not supported for this operation")
		})

		t.Run("with an individual account", func(t *testing.T) {
			testCase(t, individual)
		})

		t.Run("with an organization account", func(t *testing.T) {
			testCase(t, organization)
		})
	})
}

func TestAccGithubRepositoryPullRequestResource_update(t *testing.T) {
	randomID := acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum)

	configStep1 := fmt.Sprintf(`
		resource "github_repository" "test" {
			name      = "tf-acc-test-%s"
			auto_init = true
		}

		resource "github_branch" "test" {
			repository    = github_repository.test.name
			branch        = "test"
			source_branch = github_repository.test.default_branch
		}

		resource "github_repository_file" "test" {
			repository     = github_repository.test.name
			branch         = github_branch.test.branch
			file           = "test"
			content        = "bar"
		}

		resource "github_repository_pull_request" "test" {
			base_repository = github_repository_file.test.repository
			base_ref        = github_repository.test.default_branch
			head_ref        = github_branch.test.branch
			title           = "test title"
			body            = "test body"
		}
	`, randomID)

	configStep2 := fmt.Sprintf(`
		resource "github_repository" "test" {
			name      = "tf-acc-test-%s"
			auto_init = true
		}

		resource "github_branch" "test" {
			repository    = github_repository.test.name
			branch        = "test"
			source_branch = github_repository.test.default_branch
		}

		resource "github_repository_file" "test" {
			repository     = github_repository.test.name
			branch         = github_branch.test.branch
			file           = "test"
			content        = "bar"
		}

		resource "github_repository_pull_request" "test" {
			base_repository      = github_repository_file.test.repository
			base_ref             = github_repository.test.default_branch
			head_ref             = github_branch.test.branch
			title                = "updated title"
			body                 = "updated body"
			maintainer_can_modify = true
		}
	`, randomID)

	const resourceName = "github_repository_pull_request.test"

	testCase := func(t *testing.T, mode string) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t, mode) },
			ProtoV6ProviderFactories: testAccMuxedProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: configStep1,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "title", "test title"),
						resource.TestCheckResourceAttr(resourceName, "body", "test body"),
						resource.TestCheckResourceAttr(resourceName, "maintainer_can_modify", "false"),
					),
				},
				{
					Config: configStep2,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "title", "updated title"),
						resource.TestCheckResourceAttr(resourceName, "body", "updated body"),
						resource.TestCheckResourceAttr(resourceName, "maintainer_can_modify", "true"),
					),
				},
			},
		})
	}

	t.Run("with an individual account", func(t *testing.T) {
		testCase(t, individual)
	})

	t.Run("with an organization account", func(t *testing.T) {
		testCase(t, organization)
	})
}
