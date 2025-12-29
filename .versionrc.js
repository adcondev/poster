module.exports = {
    // Updated types with library-specific sections
    types: [
        {type: "feat", section: "✨ Features"},
        {type: "fix", section: "🐛 Bug Fixes"},
        {type: "perf", section: "⚡ Performance Improvements"},
        {type: "deps", section: "📦 Dependencies"},
        {type: "revert", section: "⏪ Reverts"},
        {type: "test", section: "✅ Tests"},
        {type: "ci", section: "🤖 Continuous Integration"},
        {type: "build", section: "🏗️ Build System"},
        {type: "refactor", section: "♻️ Code Refactoring"},
        {type: "docs", section: "📝 Documentation"},
        {type: "style", section: "🎨 Code Style"},
        {type: "chore", hidden: true}
    ],

    // Commit URLs
    commitUrlFormat: "https://github.com/adcondev/poster/commit/{{hash}}",
    compareUrlFormat: "https://github.com/adcondev/poster/compare/{{previousTag}}...{{currentTag}}",
    issueUrlFormat: "https://github.com/adcondev/poster/issues/{{id}}",
    userUrlFormat: "https://github.com/{{user}}",

    // Skip CI on release commits
    releaseCommitMessageFormat: "chore(release): {{currentTag}} [skip ci]",

    // Custom header for CHANGELOG
    header: "# Changelog\n\nAll notable changes to the POS Printer library will be documented in this file.\n\nThe format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),\nand this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).\n",

    // Bump files (automatically update version in these files)
    bumpFiles: [
        {
            filename: "package.json",
            type: "json"
        }
    ],

    // Package files to read current version from
    packageFiles: [
        {
            filename: "package.json",
            type: "json"
        }
    ],

    // Scripts to run
    scripts: {
        // Prebump runs before version bump
        prebump: "echo 'Preparing release...'",

        // Postbump runs after version bump but before git operations
        postbump: "echo 'Version bumped successfully'",

        // Precommit runs before git commit
        precommit: "go mod tidy && git add go.mod go.sum",

        // Postcommit runs after git commit but before tag
        postcommit: "echo 'Release commit created'",

        // Pretag runs before git tag
        pretag: "echo 'Creating tag...'",

        // Posttag runs after git tag
        posttag: "echo 'Release v{{currentTag}} tagged successfully'"
    },

    // Skip certain lifecycle steps if needed
    skip: {
        // skip:  {
        //   changelog: false,
        //   commit: false,
        //   tag: false
        // }
    },

    // Conventional commits preset configuration
    preset: {
        name: "conventionalcommits",
        types: [
            {type: "feat", section: "✨ Features"},
            {type: "fix", section: "🐛 Bug Fixes"},
            {type: "perf", section: "⚡ Performance Improvements"},
            {type: "deps", section: "📦 Dependencies"},
            {type: "revert", section: "⏪ Reverts"},
            {type: "test", section: "✅ Tests"},
            {type: "ci", section: "🤖 Continuous Integration"},
            {type: "build", section: "🏗️ Build System"},
            {type: "refactor", section: "♻️ Code Refactoring"},
            {type: "docs", section: "📝 Documentation"},
            {type: "style", section: "🎨 Code Style"}
        ]
    }
};