module.exports = {
    tagPrefix: "v",

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

    commitUrlFormat: "https://github.com/adcondev/poster/commit/{{hash}}",
    compareUrlFormat: "https://github.com/adcondev/poster/compare/{{previousTag}}...{{currentTag}}",
    issueUrlFormat: "https://github.com/adcondev/poster/issues/{{id}}",
    userUrlFormat: "https://github.com/{{user}}",

    releaseCommitMessageFormat: "chore(release): v{{currentTag}} [skip ci]",

    header: "# Changelog\n\nAll notable changes to Poster will be documented in this file.\n",

    packageFiles: [
        {
            filename: "package.json",
            type: "json"
        }
    ],

    bumpFiles: [
        {
            filename: "package.json",
            type: "json"
        },
        {
            filename: "package-lock.json",
            type: "json"
        }
    ],

    scripts: {
        prebump: "echo 'Preparing release...'",
        postbump: "echo 'Version bumped to {{currentTag}}'",
        postcommit: "echo 'Release commit created'",
        pretag: "echo 'Creating tag v{{currentTag}}...'",
        posttag: "echo 'Tag v{{currentTag}} created successfully'"
    }
};