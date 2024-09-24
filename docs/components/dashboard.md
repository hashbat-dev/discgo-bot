# disc-go Dashboard
The Dashboard runs alongside the bot on port :3333. 

This can be accessed here when running locally: [http://localhost:3333/](http://localhost:3333/)

## Overview
A package sits in the Bot which gathers data and compiles it into a struct variable.

The most recent Packet stays in a cached variable in the bot.

A HTML/JS dashboard calls the exposed /getData endpoint in the Go service which returns this cached variable as JSON.

The Dashboard parses this JSON via Javascript, updating local array variables with the new data before refreshing the different Dashboard components.

## Data Gathering
**Dashboard/packet.go** contains the Data Structs used to form the outbound packet, the functions used to Populate the sections of the packet and the **GetPacket()** function used to generate a new packet.

**Scheduler/2seconds.go** contains the **saveNextDashboardInfo()** internal function which it calls every 2 seconds to save a new packet to the cache.

## Data Requesting
**Dashboard/dashboard.go** contains the Http ListenAndServe handler which processes requests for the **/getData** URL, as well as returning the HTML/CSS/JS resources to display the Dashboard. This has been created in a way which allows for additional handlers in future.

## Data Viewing
**Dashboard/Pages/dashboard.html** is the HTML Front-end to the Dashboard.

**Dashboard/Resources/scripts.js** is where the Javascript for the Dashboard resides.

**Dashboard/Resources/styles.css** is where the Styles for the Dashboard are.

## Dashboard HTML/CSS
`<div class="FlexContainer">` is the wrapper which holds the Dashboard components.

Inside this container, multiple `<div class="FlexBase">` elements are used to denote the Dashboard Widgets.

One section of the Dashboard can contain more than one Widget, by using `<div class="FlexContainer FlexBase">` then inserting mulitple FlexBase entries inside of it.

## Dashboard JS
The Javascript assigns page handlers, then runs the **pageLoop()** function every 900ms until closed.

If the **Write Updates** toggle is turned off, the function breaks before starting and will try again on the next 900ms interval. This is to allow for selecting/searching through text, as updates to components will reset user input.

The **pageLoop()** function requests new JSON Data from the **/getData** handler then does the following.
1. Check that static variables are set from the JSON (ie. LogLevels, CommandTypes)
1. Updates the CPU/RAM Arrays and redraws the graphs if there are new entries
1. Updates the ActiveGuilds array and redraws the table
1. Inserts data into the Filter selections if they have not yet been populated.
1. Checks if any new "Recent Logs" exist and if so redraws the table.
1. Checks if any new "Recent Command" entries exist and if so redraws the table.
1. Updates the "Command Info" table.