# Triggers
Triggers are the concept of tracking phrases/words from a user who posts in a guild.

## Detecting Triggers
This is done in the NewMessage handler (**Bot/Handlers/newMessage.go**), each Guild has its own list of possible Triggers which are definable by Guild management. 

The Trigger collection is stored in the Cache. This is populated on Bot startup in the NewGuild handler (**Bot/handlers/newGuild.go**), and will be updated whenever a new entry is inserted into the Database by the Guild.

## Trigger Options
Triggers have some options which are defined by Guild management. These are as follows:
* **NotifyOnDetection**: When True, sends a "_PHRASE_ MENTIONED" message to the user when detected.
* **WordOnlyMatch**: When True, will only match if used as a standalone word. eg. For 'Phrase', "My phrase" would be accepted but "phraseing" would not.

## Adding Triggers (Under Construction)
This will be done via a **/slash** command. Guild management will have the following options:
* Phrase (string, max 50 characters)
* Notification message when Triggered (boolean)
* Match only when used as a word (boolean)