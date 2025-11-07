# NebulaBox CI/CD Pipeline

This document describes the Continuous Integration and Continuous Deployment (CI/CD) setup for NebulaBox.

## Overview

NebulaBox uses automated CI/CD pipelines to:
- Run tests on every commit and pull request
- Build and package the application
- Deploy to staging and production environments
- Generate test coverage reports
- Run performance benchmarks

## Supported Platforms

### GitHub Actions
GitHub Actions workflows are configured in `.github/workflows/`:
- **CI Workflow** (`.github/workflows/ci.yml`): Runs on every push and PR
- **Release Workflow** (`.github/workflows/release.yml`): Creates releases on version tags

### GitLab CI
GitLab CI configuration is in `.gitlab-ci.yml`:
- Runs tests, builds, and deployments
- Supports staging and production environments

## CI Pipeline Stages

### 1. Test Stage
- **Unit Tests**: Fast, isolated tests for individual components
- **Integration Tests**: Tests for component interactions
- **Benchmarks**: Performance benchmarks (non-blocking)

### 2. Build Stage
- Builds the API server binary
- Builds the test runner binary
- Builds the frontend dashboard

### 3. Lint Stage
- Runs `golangci-lint` for code quality checks
- Checks code formatting
- Validates code style

### 4. Frontend Stage
- Installs npm dependencies
- Builds the React dashboard
- Validates frontend builds

### 5. Deploy Stage
- **Staging**: Automatic deployment on `develop` branch (optional)
- **Production**: Manual deployment on `main`/`master` branch

## Local Testing

### Run CI Tests Locally

```bash
# Run the full CI test suite
make ci-test

# Or run the script directly
bash scripts/ci-test.sh
```

### Pre-commit Hooks

Install Git hooks to run tests before commits:

```bash
bash scripts/setup-git-hooks.sh
```

This installs a pre-commit hook that:
- Runs linters
- Runs unit tests
- Checks code formatting
- Validates code quality

To skip hooks: `git commit --no-verify`

## GitHub Actions

### Workflow Triggers
- **Push**: Triggers on pushes to `main`, `develop`, or `master`
- **Pull Request**: Triggers on PRs to `main`, `develop`, or `master`
- **Release**: Triggers on version tags (`v*`)

### Jobs

#### Test Job
- Runs unit tests
- Runs integration tests
- Runs benchmarks (non-blocking)
- Generates coverage reports
- Uploads coverage as artifact

#### Build Job
- Builds API server
- Builds test runner
- Uploads binaries as artifacts

#### Lint Job
- Runs `golangci-lint`
- Validates code quality

#### Frontend Job
- Builds React dashboard
- Uploads dist as artifact

#### Release Job
- Creates release archives
- Uploads to GitHub Releases
- Includes binaries and frontend dist

## GitLab CI

### Pipeline Stages
1. **test**: Unit tests, integration tests, benchmarks
2. **build**: Builds binaries
3. **lint**: Code quality checks
4. **frontend**: Builds dashboard
5. **deploy**: Deployment to environments

### Environments
- **Staging**: Auto-deploys from `develop` branch (manual trigger)
- **Production**: Deploys from `main`/`master` branch (manual trigger)

### Artifacts
- Test coverage reports
- Built binaries
- Frontend dist
- Benchmark results

## Environment Variables

### Required
- `GITHUB_TOKEN`: For GitHub Actions releases (auto-provided)
- `GO_VERSION`: Go version to use (default: 1.21)
- `NODE_VERSION`: Node.js version (default: 20)

### Optional
- `NEBULABOX_REGISTRY_URL`: Registry URL for tests
- `NEBULABOX_ADMIN_USER`: Admin username for tests
- `NEBULABOX_ADMIN_PASS`: Admin password for tests

## Test Coverage

Coverage reports are generated and uploaded as artifacts:
- HTML report: `coverage.html`
- Cobertura XML: `coverage.xml` (GitLab)

View coverage in:
- GitHub Actions artifacts
- GitLab CI/CD pipelines
- Local: `make test-unit-coverage`

## Benchmarking

Benchmarks run in CI but are non-blocking:
- API endpoint benchmarks
- Containerd operation benchmarks
- Results available in artifacts

## Deployment

### Manual Deployment
Deployments are manual by default for safety:
- GitLab: Click "Deploy" button in pipeline
- GitHub: Create release tag

### Automatic Deployment (Optional)
To enable automatic staging deployments:
1. Set up deployment credentials
2. Configure deployment scripts
3. Update CI configuration

## Troubleshooting

### Tests Failing in CI
1. Check test logs in CI output
2. Run tests locally: `make test-unit`
3. Verify environment variables
4. Check for dependency issues

### Build Failures
1. Verify Go version matches CI
2. Check for dependency conflicts
3. Review build logs

### Frontend Build Issues
1. Check Node.js version
2. Clear npm cache: `npm cache clean --force`
3. Verify package-lock.json is up to date

## Best Practices

1. **Run tests locally** before pushing
2. **Keep CI fast** - optimize slow tests
3. **Use artifacts** for debugging failed builds
4. **Monitor coverage** - maintain >80% coverage
5. **Review lint errors** before merging
6. **Test in CI** - don't skip CI checks

## Extending CI/CD

### Add New Test Stage
1. Add job to `.github/workflows/ci.yml`
2. Add stage to `.gitlab-ci.yml`
3. Update `scripts/ci-test.sh` if needed

### Add Deployment Target
1. Create deployment job
2. Configure credentials/secrets
3. Add environment configuration
4. Update documentation

### Custom Scripts
Place custom scripts in `scripts/` directory:
- `ci-test.sh` - Main CI test runner
- `pre-commit.sh` - Pre-commit checks
- `setup-git-hooks.sh` - Hook installer

## Security

- Never commit secrets or credentials
- Use CI/CD secrets management
- Rotate tokens regularly
- Review dependency updates
- Scan for vulnerabilities

## Resources

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [GitLab CI Docs](https://docs.gitlab.com/ee/ci/)
- [Go Testing Guide](https://golang.org/pkg/testing/)
- [Makefile Documentation](./Makefile)

