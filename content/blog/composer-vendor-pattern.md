---
title: "The Vendor Pattern with Composer: When and Why to Commit Dependencies"
date: 2025-01-20
description: "Explore the vendor pattern in Composer—committing dependencies to your repository instead of fetching them dynamically. Learn the tradeoffs and when this approach makes sense."
tags:
  - php
  - composer
  - dependency-management
  - deployment
---

When working with PHP projects using Composer, you typically have two approaches to managing your `vendor/` directory:

1. **Standard approach**: Add `vendor/` to `.gitignore`, commit only `composer.json` and `composer.lock`, and run `composer install` during deployment
2. **Vendor pattern**: Commit the entire `vendor/` directory to your repository

While the standard approach is conventional and recommended for most projects, the vendor pattern offers specific advantages in certain scenarios. Let's explore what it is, when to use it, and the tradeoffs involved.

## What is the Vendor Pattern?

The vendor pattern means committing your project's entire `vendor/` directory—all installed Composer dependencies—directly into your version control system alongside your source code.

This approach flips the typical workflow:

- **Standard**: Source code → Git → Deploy → Run `composer install`
- **Vendor pattern**: Source code + vendor/ → Git → Deploy (no install needed)

## Why Consider the Vendor Pattern?

### 1. **Eliminate Dependency on Package Repositories During Deployment**

With the standard approach, deployment relies on external package repositories (Packagist, private registries, etc.) being available. If a package repository is down, slow, or unreachable during your deployment window, your deploy fails.

The vendor pattern removes this risk—everything needed to run your application is already in version control.

### 2. **Deployment Simplicity**

No need to run `composer install` during deployment. Your application code and all its dependencies are ready to go immediately after pulling from the repository. This reduces deployment complexity and time.

### 3. **Consistency Across Environments**

Everyone working on the project has identical dependencies without the possibility of version mismatches due to timing or environment differences. The `composer.lock` file guarantees reproducibility, but adding the actual files provides absolute consistency.

## The Tradeoffs

Before adopting the vendor pattern, understand the significant drawbacks:

### 1. **Repository Bloat**

This is the biggest downside. Composer packages can be large, and committing all of them inflates your repository:

- **Clone time**: New developers download significantly more data
- **Storage**: Backup and hosting costs increase with repository size
- **Git performance**: Operations slow down with larger repos
- **Network bandwidth**: Cloning becomes expensive on slower connections

A typical project with moderate dependencies might add 50-500 MB (or more) to your repository.

### 2. **Merge Conflicts**

If multiple developers update dependencies or work on different branches, merging `vendor/` becomes painful. You might end up with conflicts in hundreds of package files.

### 3. **Maintenance Burden**

You're now responsible for:
- Monitoring when new dependency versions are available
- Deciding when to update (rather than auto-pulling on install)
- Committing potentially massive changesets when updating dependencies
- Managing the complexity of reviewing what actually changed in dependencies

### 4. **Security Updates**

Security patches require a full update, commit, and redeploy cycle rather than simply running `composer update` locally and redeploying.

## When to Use the Vendor Pattern

Given the tradeoffs, the vendor pattern makes sense in specific scenarios:

### 1. **Unreliable Deployment Environments**

If your deployment infrastructure has:
- No internet access to package repositories
- Unreliable network connectivity during deploys
- Offline deployment requirements

The vendor pattern becomes valuable.

### 2. **Critical Stability Requirements**

For mission-critical applications where:
- You want absolute certainty nothing changes during deployment
- Deployment must be as fast as possible
- Repository availability isn't guaranteed at deploy time

### 3. **Distributed or Offline Teams**

If developers work offline or in environments with poor network access, having dependencies committed can reduce friction.

## Implementation Tips

If you decide to use the vendor pattern:

### 1. **Create Separate Composer Files**

Maintain both configurations:

```json
// composer.prod.json - for deployment
{
    "repositories": [
        {
            "type": "vcs",
            "url": "git@github.com:yourorg/private-package.git"
        }
    ],
    "require": {
        "yourorg/package": "^1.0.5"
    }
}
```

```json
// composer.dev.json - for local development
{
    "repositories": [
        {
            "type": "path",
            "url": "../private-package"
        }
    ],
    "require": {
        "yourorg/package": "*@dev"
    }
}
```

Use symlinks to switch between them:

```bash
ln -s composer.dev.json composer.json  # For development
rm composer.json && cp composer.prod.json composer.json  # For production
```

### 2. **Clean Up `.git` Directories**

After running `composer install`, remove nested `.git` directories in vendor packages to prevent submodule issues:

```bash
find vendor -type d -name ".git" -exec rm -rf {} +
```

### 3. **Use `.gitignore` Selectively**

If partially committing vendor, be explicit about what to include:

```gitignore
# Exclude vendor by default
/vendor/

# But include specific essential packages
!/vendor/yourorg/
```

This is fragile and generally not recommended—commit all or nothing.

### 4. **Document the Approach**

Make it clear to your team that this is your deployment strategy. Include instructions for switching between dev and prod configurations.

## The Recommendation

**For most projects**, stick with the standard approach:
- Add `vendor/` to `.gitignore`
- Commit `composer.json` and `composer.lock`
- Run `composer install` during deployment
- Use a reliable CI/CD pipeline to handle the install step

**Consider the vendor pattern only if you have a specific, well-documented reason** and understand the maintenance burden it introduces.

## See Also

- [Composer Documentation](https://getcomposer.org/doc/)
- [PSR-4 Autoloading Standard](https://www.php-fig.org/psr/psr-4/)
- [Semantic Versioning](https://semver.org/)
