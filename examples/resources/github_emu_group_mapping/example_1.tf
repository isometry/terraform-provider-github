resource "github_emu_group_mapping" "example_emu_group_mapping" {
  team_slug = "emu-test-team" # The GitHub team name to modify
  group_id  = 28836           # The group ID of the external group to link
}

# Note that here GITHUB_OWNER and GITHUB_TOKEN have been set in the environment.
