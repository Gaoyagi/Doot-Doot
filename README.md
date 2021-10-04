game plan:
-create discord bot [x]
-figure out how to get it to react to an @'ing in a channel
    - to use the bot you have to: @bot bind
        - this binds the bot to that channel and it will only read commands from that channel from now on
-figure out how to get it to only react to messages from that channel
    -implement fake commands to prove it
-figure out how to add it to a voice call
-figure out how ti to stream music from links
-final step, link the bot to an EC2 instance so it always stays running
    - https://dev.to/rishabk7/host-your-discord-bot-on-ec2-instance-aws-5c07
    - https://towardsaws.com/building-hosting-a-discord-bot-on-aws-e157bd7faf78



-how to kep track of what channel the bot is bound to
-how to keep track of the availible commands
-should i make a struct
    -it would have to inherit from the default session


-each session onject created is basically a user?


