# GitHub Repository Configuration

This document describes the GitHub repository settings, templates, and protections configured for the Pokedex project.

## Configuration Summary

**Repository:** pfjndev/pokedex
**Branch:** main
**Date Configured:** 2024

## Templates

### Pull Request Template (`.github/pull_request_template.md`)

The PR template includes:
- Description section for change summary
- Type of change checklist:
  - 🐛 Bug fix
  - ✨ Feature
  - 💥 Breaking change
  - 📚 Documentation update
  - ♻️ Code refactor
  - 🔧 Configuration or build system change
  - ⚡ Performance improvement
- Testing checklist for validation
- Code quality checklist
- Additional checks before merging

### Bug Report Template (`.github/ISSUE_TEMPLATE/bug_report.md`)

Includes sections for:
- Clear description of the bug
- Steps to reproduce
- Expected vs. actual behavior
- Environment details (Go version, OS, architecture)
- Error messages and logs
- Screenshots and additional context
- Possible solutions

### Feature Request Template (`.github/ISSUE_TEMPLATE/feature_request.md`)

Includes sections for:
- Feature description
- Problem statement
- Proposed solution
- Alternative solutions
- Benefits and use cases
- Acceptance criteria

### Issue Template Config (`.github/ISSUE_TEMPLATE/config.yml`)

Configures:
- Disables blank issues (users must use a template)
- Provides links to documentation and discussions

## Branch Protection Rules

**Main Branch Protection** is configured with the following rules:

### Pull Request Requirements
- ✅ **Require pull request reviews before merging**
  - Minimum approving reviews required: 1
  - Dismiss stale pull request approvals: Yes
  - Require code owner reviews: No

### Status Checks
- ✅ **Require status checks to pass before merging**
  - Strict mode enabled (requires branch to be up to date)
  - Current contexts configured: None (add as CI/CD workflows are added)

### Other Protections
- ✅ **Enforce all the above rules for administrators:** Yes
- ✅ **Restrict who can push to matching branches:** No restrictions
- ✅ **Allow force pushes:** No
- ✅ **Allow deletions:** No
- ✅ **Require conversation resolution before merging:** No
- ✅ **Require linear history:** No
- ✅ **Require signed commits:** No
- ✅ **Lock branch:** No

## README Updates

The README.md has been updated to include:
- Status badges (Build, Go version, Go Report Card, License)
- Project description
- Features overview
- Prerequisites
- Installation and usage instructions
- Project structure
- Development guidelines
- Contributing information
- License information

## License

MIT License added (`LICENSE` file) with:
- Standard MIT license text
- Copyright year: 2024
- Copyright holder: pfjndev

## Next Steps

When you have CI/CD workflows configured, update the branch protection to require status checks by running:

```bash
gh api repos/pfjndev/pokedex/branches/main/protection/required_status_checks/contexts \
  -X POST \
  -f context="<workflow-name>"
```

Replace `<workflow-name>` with actual workflow status check names from GitHub Actions.

## Verification

To verify the current branch protection settings:

```bash
gh api repos/pfjndev/pokedex/branches/main/protection
```

To view the repository in GitHub UI:

```bash
gh repo view pfjndev/pokedex --web
```

## Files Created/Modified

- ✅ `.github/pull_request_template.md` - Created
- ✅ `.github/ISSUE_TEMPLATE/bug_report.md` - Created
- ✅ `.github/ISSUE_TEMPLATE/feature_request.md` - Created
- ✅ `.github/ISSUE_TEMPLATE/config.yml` - Created
- ✅ `LICENSE` - Created (MIT)
- ✅ `README.md` - Updated with badges and documentation
- ✅ Branch protection - Configured via GitHub API

All changes have been committed and pushed to the main branch.
