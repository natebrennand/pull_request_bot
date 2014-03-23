
#Pull Request Bot

A tool to enforce code review.

Using a GitHub webhook, the bot listens for comments on open Pull Requests.
Using a set of "approval messages", the bot waits for N approved users to give the pull request their blessing.
After N approvals, the bot merges the pull request.


#Settings

You will need a GitHub API token in the settings that has permissions to alter repos.

Add all repos that you wish the bot to listen to in the [config file][config.json].

