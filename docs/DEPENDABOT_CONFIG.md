# Dependabot Configuration Guide

## Overview

This document describes the Dependabot configuration for the Pokedex repository, set up to automatically manage dependency updates for GitHub Actions and Go modules.

## Configuration Files

### `.github/dependabot.yml`

The main Dependabot configuration file that defines how dependencies are managed.

#### Key Settings

**GitHub Actions Updates:**
- **Package Ecosystem:** `github-actions`
- **Schedule:** Weekly on Mondays at 09:00 UTC
- **PR Limit:** 5 open pull requests maximum
- **Labels:** `dependencies`, `ci-cd`
- **Commit Message Prefix:** `chore(deps)`

**Go Modules Updates:**
- **Package Ecosystem:** `gomod`
- **Schedule:** Weekly on Mondays at 10:00 UTC
- **PR Limit:** 10 open pull requests maximum
- **Labels:** `dependencies`, `go`
- **Commit Message Prefix:** `chore(deps)`

#### Dependency Grouping

Go module updates are grouped for better organization:

1. **bubbletea Group** - Groups all charmbracelet dependencies (`github.com/charmbracelet/*`)
2. **test-deps Group** - Groups test and mock dependencies (patterns: `*test*`, `*mock*`)

#### Ignore Rules

- Major version updates (`version-update:semver-major`) are ignored by default to maintain stability
- Only patch and minor updates are included in automatic updates

### `.github/workflows/dependabot-auto-merge.yml`

Workflow for automatically merging Dependabot pull requests that meet criteria.

#### Features

- **Trigger:** Runs on all pull requests
- **Auto-Merge Criteria:** 
  - Only runs for `dependabot[bot]` actor
  - Auto-merges patch and minor version updates
  - Skips major version updates (requires manual review)
- **Merge Strategy:** Rebase merge
- **Permissions:** Includes write access to contents and pull requests

## Repository Settings

### ✅ Enabled Features

Based on API verification:
- ✅ **Dependabot Security Updates** - Enabled
- ✅ **Secret Scanning** - Enabled
- ✅ **Push Protection** - Enabled

### Manual Configuration Steps (if needed)

If Dependabot is not automatically enabled:

1. **Enable Dependabot via GitHub UI:**
   - Navigate to Repository Settings → Code security and analysis
   - Enable "Dependabot alerts"
   - Enable "Dependabot security updates"

2. **Configure Branch Protection (Recommended):**
   - Go to Settings → Branches → Branch protection rules
   - Create a rule for `main` and `develop` branches
   - Require status checks to pass before merging
   - Allow auto-merge for specific workflows

3. **Grant Workflow Permissions:**
   - Settings → Actions → General → Workflow permissions
   - Select "Read and write permissions"
   - Enable "Allow GitHub Actions to create and approve pull requests"

## How It Works

### Weekly Dependency Checks

1. **Monday 09:00 UTC** - GitHub Actions dependencies are checked
2. **Monday 10:00 UTC** - Go modules are checked
3. Dependabot creates PRs for available updates

### PR Creation Process

1. Dependabot creates a pull request with:
   - Title: `chore(deps): bump <dependency>`
   - Labels: `dependencies` and ecosystem-specific labels
   - Automatic detection of update type (patch/minor/major)

### Auto-Merge Logic

1. **Patch & Minor Updates:** Automatically merged via rebase
2. **Major Updates:** Created as PRs for manual review
3. **Grouped Updates:** Each group may generate separate PRs

## Managing Dependabot

### View Dependabot Activity

```bash
# List recent Dependabot PRs
gh pr list --author dependabot --state all

# Get details about Dependabot configuration status
gh api repos/pfjndev/pokedex --jq '.security_and_analysis'
```

### Customization Options

To modify the configuration:

1. Edit `.github/dependabot.yml`
2. Commit changes to main branch
3. Changes take effect within 24 hours

Common customizations:
- **Change Schedule:** Modify `interval`, `day`, `time`
- **Adjust PR Limits:** Update `open-pull-requests-limit` values
- **Add Package Ecosystems:** Add new entries for npm, docker, etc.
- **Modify Ignore Rules:** Update patterns to include/exclude dependencies

### Disabling Auto-Merge for Specific PRs

If a Dependabot PR is created but shouldn't auto-merge:
1. Add label `do-not-merge` or similar
2. Manually dismiss auto-merge with `gh pr merge --disable-auto`

## Best Practices

1. **Review Major Updates:** Always manually review major version updates for breaking changes
2. **Test Before Merging:** Ensure CI/CD passes before merging updates
3. **Group Related Updates:** Use dependency groups to keep related updates together
4. **Monitor Security Alerts:** Check Dependabot alerts tab for vulnerabilities
5. **Regular Maintenance:** Review ignored rules quarterly to ensure they're still needed

## Troubleshooting

### Dependabot Not Creating PRs

1. Verify `.github/dependabot.yml` syntax (YAML validation)
2. Check repository has public visibility or GitHub organization plan with Dependabot
3. Ensure branch is set to receive updates
4. Wait up to 24 hours for initial check

### Auto-Merge Failures

1. Check workflow logs: `gh run list --workflow dependabot-auto-merge.yml`
2. Verify workflow has necessary permissions
3. Check branch protection rules don't block auto-merge
4. Ensure workflow file is in default branch

### Duplicate or Unwanted PRs

1. Close unwanted PRs manually
2. Add `ignore` rule to `dependabot.yml` if pattern needs exclusion
3. Increase `open-pull-requests-limit` if PRs are being delayed

## References

- [Dependabot Documentation](https://docs.github.com/en/code-security/dependabot)
- [Dependabot Configuration Options](https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/configuration-options-for-dependency-updates)
- [Dependabot Auto-Merge Setup](https://docs.github.com/en/code-security/dependabot/working-with-dependabot/automating-dependabot-with-github-actions)

## Configuration History

- **Created:** Feature branch `feature/pokemon-models-rebased`
- **Ecosystems:** GitHub Actions, Go modules
- **Default Schedule:** Weekly on Mondays
- **Auto-Merge:** Enabled for patch and minor updates only
