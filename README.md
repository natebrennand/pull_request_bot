
#Pull Request Bot

A tool to enforce code review.

Using a GitHub webhook, the bot listens for comments on open Pull Requests.
Using a set of "approval messages", the bot waits for N approved users to give the pull request their blessing.
After N approvals, the bot merges the pull request.


##Settings

You will need a GitHub API token in the settings that has permissions to alter repos.
Add all repos that you wish the bot to listen to in the [config file](config.json).

You only need to alter the [settings file](settings.dev) if you wish to run the project on a different port or rename your configuration json.

##Webhook

For every repo you wish you use the utility, you need to setup a webhook to receive all `Issue Comment` events.

