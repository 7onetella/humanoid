[![CircleCI](https://circleci.com/gh/7onetella/humanoid/tree/master.svg?style=svg)](https://circleci.com/gh/7onetella/humanoid/tree/master)

# humanoid
humanoid is a generic slack app which utilizes bot account so it can be notified of incoming messages irregardless of intended recipient. humanoid alone does not do any automation. humanoid relies on external tools/commands to do anything useful. 

chatops is the ultimate goal for humanoid. however, humanoid can be repurposed to do something else. I have yet to have needs for anything else but I will come back and update this readme page if new usecases are applicable.

humanoid is simply empty shell that passes messages to morgan at the time of writing this readme.

## taking a different approach to designing of typical chatops bot
hubot is a great platform for all sorts of automation. I wanted hubot version of chatops app in golang. that gave birth to humanoid. I am in much favor of go's single binary software artifact. hubot uses regular expression to match messages to commands. I am taking a different approach. humanoid execution of commands are tightly controlled by humanoid-config.ini file. that allows delegation of parsing and matching to the executing command. there is no ad-hoc "eval" of commands or execution of Linux os commands. I would advise anyone from adding such commands to humanoid-config.ini file.

one drawback is every time new sub commands are added, humanoid-config.ini needs to be updated.

humanoid simply calls exec to replay message to CLI command. first string and it's subsequent strings are checked against the [allowed] portion of humanoid-config.ini file. 

for example,

@morgan aws ec2 start db

@morgan is the bot user handle. this tells humanoid to pick up the message instead of disregarding it.

aws ec2 start needs to be part of humanoid-config.ini
```
[allowed]
aws
aws ec2
aws ec2 start
aws ec2 stop
aws ec2 describe-instances
aws ec2 terminate

[approval required]
aws ec2 terminate

[peer approver]
7onetella
```

### how to build binary
go build -o humanoid

### how to run
```
export SLACK_BOT_USER_OAUTH_ACCESS_TOKEN=xxxxx; \
export SLACK_BOT_MEMBER_ID=XXXXXXX; \
export SLACK_BOT_DEBUG=true; \
./humanoid
```

### Screenshots
![aws ec2 describe-instances](/asset/example.png)