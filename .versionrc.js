module.exports = {
    // Types for changelog sections
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

    // URL formats
    commitUrlFormat: "https://github.com/adcondev/poster/commit/{{hash}}",
    compareUrlFormat: "https://github.com/adcondev/poster/compare/{{previousTag}}...{{currentTag}}",
    issueUrlFormat: "https://github.com/adcondev/poster/issues/{{id}}",
    userUrlFormat: "https://github.com/{{user}}",

    // Release commit message
    releaseCommitMessageFormat: "chore(release): {{currentTag}} [skip ci]",

    // Custom header for CHANGELOG
    header: "# Changelog\n\nAll notable changes to the POS Printer library will be documented in this file.\n\nThe format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),\nand this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).\n",

    // ✅ CORRECCIÓN: Solo una configuración de archivos
    bumpFiles: [
        {
            filename: "package.json",
            type: "json"
        }
    ],

    // ✅ CORRECCIÓN: Lifecycle scripts
    scripts: {
        prebump: "echo 'Preparing release...'",
        postbump: "echo 'Version bumped successfully'",
        precommit: "go mod tidy && git add go.mod go.sum",
        postcommit: "echo 'Release commit created'",
        pretag: "echo 'Creating tag...'",
        posttag: "echo 'Release v{{currentTag}} tagged successfully'"
    }
};