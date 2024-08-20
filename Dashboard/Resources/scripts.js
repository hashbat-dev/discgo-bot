
// Run then setup the Loop for data checks every 2 seconds
window.onload = function() {
    setupHandlers();
    pageLoop(); 
    setInterval(pageLoop, 900); 
}

// Global Variables
let serverNameSet = false;
let filtersSet = false;
let chartInit = false;
let dashboardMaxLogs = -1;
let dashboardMaxCommands = -1;
let hardwareInterval = -1;
let hardwareMax = -1;
let logLevels = [];
let commandTypes = [];
let cpuValuesArray = [];
let ramValuesArray = [];
let activeGuilds = [];
let logs = [];
let commands = [];
let commandInfo = [];

function pageLoop() {
    fetchData().then(newData => {
        
        const writeUpdates = document.getElementById('WriteUpdates');
        if (!writeUpdates.checked) {
            return
        }

        // Do any one-time variable updates need to happen?
        if (!serverNameSet) {
            document.getElementById('ServerName').innerText = newData['PacketInfo']['ServerName'];
            serverNameSet = true;
        }  
        if (logLevels.length === 0) {
            logLevels = newData['Logging']['LogLevels'];
        } 
        if (commandTypes.length === 0) {
            commandTypes = newData['Commands']['CommandTypes'];
        }
        if (dashboardMaxLogs < 1) {
            dashboardMaxLogs = newData['PacketInfo']['MaxLogs'];
        }
        if (dashboardMaxCommands < 1) {
            dashboardMaxCommands = newData['PacketInfo']['MaxCommands'];
        }
        if (hardwareInterval < 1) {
            hardwareInterval = newData['PacketInfo']['HardwareInterval'];
        }
        if (hardwareMax < 1) {
            hardwareMax = newData['PacketInfo']['HardwareMax'];
        }

        let hardwareUpdate = false;
        if (cpuValuesArray.length < hardwareMax && cpuValuesArray.length != newData['HardwareInfo']['CpuValues'].length) {
            cpuValuesArray = newData['HardwareInfo']['CpuValues'];
            hardwareUpdate = true;
        }
        if (ramValuesArray.length < hardwareMax && ramValuesArray.length != newData['HardwareInfo']['RamValues'].length) {
            ramValuesArray = newData['HardwareInfo']['RamValues'];
            hardwareUpdate = true;
        }
        if (hardwareUpdate) {
            updateCpuRamChart();
        }

        activeGuilds = newData['ActiveGuilds'];
        updateGuildTable();

        if (!filtersSet) {
            populateFilters();
            filtersSet = true;
        }

        if (logs.length < dashboardMaxLogs && logs.length != newData['Logging']['LogEntries'].length) {
            logs = newData['Logging']['LogEntries'];
            updateLogTable();
        }

        if (newData['Commands']['Commands'] != null) {
            if (commands.length < dashboardMaxCommands && commands.length != newData['Commands']['Commands'].length) {
                commands = newData['Commands']['Commands'];
                updateCommandsTable();
            }
        }

        if (newData['Commands']['CommandInfo'] != null) {
            commandInfo = newData['Commands']['CommandInfo'];
            updateCommandInfoTable();
        }
       
    }).catch(error => {
        console.error('Error processing data: ', error);
    });

}

// Returns JSON 
async function fetchData() {
    try {
        // 1. Get the getData URL
        const currentUrl = window.location.href;
        const dataUrl = `${currentUrl}/getData`;

        // 2. Fetch the JSON data 
        const response = await fetch(dataUrl);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        // 3. Parse the JSON data and save it as an Object
        const dataObject = await response.json();
        return dataObject;

    } catch (error) {
        console.error('Error fetching data: ', error);
    }
}

function setupHandlers() {

    const logLevelFilter = document.getElementById('LogFilterLevel');
    logLevelFilter.addEventListener('change', updateLogTable);
    
    const logLevelGuild = document.getElementById('LogFilterGuild');
    logLevelGuild.addEventListener('change', updateLogTable);

}

let cpuChart;
let ramChart;

function setupCharts() {

    try {

        const ctxCpu = document.getElementById('chartCpu').getContext('2d');
        cpuChart = new Chart(ctxCpu, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                        label: 'CPU',
                        data: [],
                        borderColor: 'rgba(23, 213, 235, 1)',
                        backgroundColor: 'rgba(23, 213, 235, 0.2)',
                        fill: false,
                        pointRadius: 0
                    }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        ticks: {
                            display: false
                        },
                        title: {
                            display: true,
                            text: 'Time'
                        }
                    },
                    y: {
                        min: 0,
                        max: 100,
                        ticks: {
                            beginAtZero: true,
                        },
                        title: {
                            display: true,
                            text: 'Usage (%)'
                        }
                    }
                }
            }
        });

        const ctxRam = document.getElementById('chartRam').getContext('2d');
        ramChart = new Chart(ctxRam, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                        label: 'RAM',
                        data: [],
                        borderColor: 'rgba(38, 199, 54, 1)',
                        backgroundColor: 'rgba(38, 199, 54, 0.2)',
                        fill: false,
                        pointRadius: 0
                    }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: {
                        ticks: {
                            display: false
                        },
                        title: {
                            display: true,
                            text: 'Time'
                        }
                    },
                    y: {
                        ticks: {
                            beginAtZero: true,
                        },
                        title: {
                            display: true,
                            text: 'Usage (MB)'
                        }
                    }
                }
            }
        });

        chartInit = true;
    }
    catch (error) {
        console.error('Error initalising Charts: ', error);
    }
    
}

function updateCpuRamChart() {

    if (!chartInit) {
        setupCharts();
    }

    if (!chartInit) {
        return;
    }

    if (cpuChart) {
        cpuChart.data.datasets[0].data = cpuValuesArray;
        cpuChart.data.labels = getLabels(cpuValuesArray);
        cpuChart.update();
    } else {
        console.error('CPU Chart instance is not initialized.');
    }

    if (ramChart) {
        ramChart.data.datasets[0].data = ramValuesArray;
        ramChart.data.labels = getLabels(ramValuesArray);
        ramChart.update();
    } else {
        console.error('RAM Chart instance is not initialized.');
    }

}

function getLabels(dataArray) {
    const length = dataArray.length;
    const incrementingArray = Array.from({ length }, (_, index) => index);
    return incrementingArray;
}

function updateGuildTable() {
    const tableBody = document.querySelector('#GuildTable tbody');
    tableBody.innerHTML = ''; // Clear existing rows

    activeGuilds.forEach(guild => {
        const row = document.createElement('tr');
        const dbIdCell = document.createElement('td');
        const discordIdCell = document.createElement('td');
        const nameCell = document.createElement('td');
        const callsCell = document.createElement('td');
        const lastCmdCell = document.createElement('td');

        nameCell.textContent = guild.Name;
        discordIdCell.innerHTML = getConcatIDElement(guild.DiscordID);
        dbIdCell.textContent = guild.DbID;
        callsCell.textContent = guild.CommandCount;

        const dateTimeString = formatDateTime(guild.LastCommand);
        if (dateTimeString === '') {
            lastCmdCell.textContent = 'N/A';
        } else {
            lastCmdCell.textContent = dateTimeString;
        }

        row.appendChild(nameCell);
        row.appendChild(discordIdCell);
        row.appendChild(dbIdCell);
        row.appendChild(callsCell);
        row.appendChild(lastCmdCell);

        tableBody.appendChild(row);
    });
}

function updateLogTable() {
    const tableBody = document.querySelector('#LogTable tbody');
    tableBody.innerHTML = ''; // Clear existing rows

    // Get any Filters which are set
    const filterLevel = document.getElementById('LogFilterLevel').value;
    const guildLevel = document.getElementById('LogFilterGuild').value;

    logs.forEach(log => {

        // Level filter?
        if (filterLevel > -1) {
            if (filterLevel != log['LogLevel']) {
                return;
            }
        }

        // Guild filter?
        if (guildLevel != '') {
            if (guildLevel != log['LogInfo']['GuildID']) {
                return;
            }
        }

        const row = document.createElement('tr');
        const datetimeCell = document.createElement('td');
        const levelCell = document.createElement('td');
        const guildIdCell = document.createElement('td');
        const sourceCell = document.createElement('td');
        const logTextCell = document.createElement('td');

        datetimeCell.textContent = formatDateTime(log['LogInfo']['DateTime']);
        levelCell.textContent = logLevels[log['LogLevel']];
        guildIdCell.innerHTML =  getConcatIDElement(log['LogInfo']['GuildID']);
        sourceCell.textContent = log['LogInfo']['CodeSource'];
        logTextCell.textContent = log['LogText'];

        row.classList.add('Log'+logLevels[log['LogLevel']]);
        guildIdCell.classList.add('FontSize10');

        row.appendChild(datetimeCell);
        row.appendChild(levelCell);
        row.appendChild(guildIdCell);
        row.appendChild(sourceCell);
        row.appendChild(logTextCell);

        tableBody.appendChild(row);
    });
}

function updateCommandsTable() {
    const tableBody = document.querySelector('#CommandTable tbody');
    tableBody.innerHTML = ''; // Clear existing rows

    commands.forEach(command => {

        const row = document.createElement('tr');
        const typeCell = document.createElement('td');
        const commandCell = document.createElement('td');
        const guildIdCell = document.createElement('td');
        const userIdCell = document.createElement('td');
        const userNameCell = document.createElement('td');
        const callTime = document.createElement('td');
        const callDuration = document.createElement('td');

        typeCell.textContent = commandTypes[command['TypeID']];
        commandCell.textContent = command['Command'];
        guildIdCell.innerHTML = getConcatIDElement(command['GuildID']);
        userIdCell.innerHTML = getConcatIDElement(command['UserID']);
        userNameCell.textContent = command['UserName'];
        callTime.textContent = formatDateTime(command['CallTime']);
        callDuration.textContent = formatDuration(command['CallDuration']);

        guildIdCell.classList.add('FontSize10');
        userIdCell.classList.add('FontSize10');

        row.appendChild(typeCell);
        row.appendChild(commandCell);
        row.appendChild(guildIdCell);
        row.appendChild(userIdCell);
        row.appendChild(userNameCell);
        row.appendChild(callTime);
        row.appendChild(callDuration);

        tableBody.appendChild(row);
    });
}

function updateCommandInfoTable() {
    const tableBody = document.querySelector('#CommandInfoTable tbody');
    tableBody.innerHTML = ''; // Clear existing rows

    commandInfo.forEach(cmdInfo => {

        const row = document.createElement('tr');
        const typeCell = document.createElement('td');
        const commandCell = document.createElement('td');
        const countCell = document.createElement('td');
        const avgTimeCell = document.createElement('td');
        const lastCallCell = document.createElement('td');

        typeCell.textContent = commandTypes[cmdInfo['TypeID']];
        commandCell.textContent = cmdInfo['Command'];
        countCell.textContent = cmdInfo['Count'];
        avgTimeCell.textContent = formatDuration(cmdInfo['AvgDuration']);
        lastCallCell.textContent = formatDateTime(cmdInfo['LastCall']);

        row.appendChild(typeCell);
        row.appendChild(commandCell);
        row.appendChild(countCell);
        row.appendChild(avgTimeCell);
        row.appendChild(lastCallCell);

        tableBody.appendChild(row);
    });
}
function getConcatIDElement(guildId) {

    if (guildId.length <= 5) {
        return guildId;
    }

    const lastFiveChars = guildId.slice(-5);
    return "<span class='GuildID' title='"+guildId+"'>.." + lastFiveChars + "</span>";
}

function formatDateTime(dateTimeString) {
    
    // Create a Date object from the DateTime string
    const date = new Date(dateTimeString);
    
    // Is this a "null" Date/Time? (Writing as 01/01/2000 from Go)
    if (date.getFullYear() === 2000) {
        return '';
    }

    // Get today's date
    const today = new Date();
    
    // Check if the date is today
    const isToday = date.toDateString() === today.toDateString();
    
    // Define options for formatting
    const options = {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        fractionalSecondDigits: 3 // For milliseconds
    };
    
    // Format time part
    const timeFormatter = new Intl.DateTimeFormat('en-GB', options);
    const formattedTime = timeFormatter.format(date);
    
    // Format date part if not today
    if (isToday) {
        return formattedTime; // Only time
    } else {
        const dayOptions = {
            ...options,
            weekday: 'short'
        };
        const dayFormatter = new Intl.DateTimeFormat('en-GB', dayOptions);
        const formattedDay = dayFormatter.format(date);
        
        // Return day and time
        return `${formattedDay} ${formattedTime}`;
    }
}

function formatDuration(ns) {
    const nsPerMs = 1e6;
    const nsPerSec = 1e9;

    // Calculate seconds and milliseconds
    const seconds = Math.floor(ns / nsPerSec);
    ns %= nsPerSec;
    
    const milliseconds = Math.floor(ns / nsPerMs);

    // Return the formatted string
    return `${seconds}s ${milliseconds}ms`;
}

function populateFilters() {

    // Log Level selector
    const selectBox = document.getElementById('LogFilterLevel');
            
    for (const key in logLevels) {
        if (logLevels.hasOwnProperty(key)) {
            const option = document.createElement('option');
            option.value = key;
            option.textContent = logLevels[key];
            
            selectBox.appendChild(option);
        }
    }

    // Guild selector
    const guildBox = document.getElementById('LogFilterGuild');
           
    activeGuilds.forEach(guild => {
        const option = document.createElement('option');
        option.value = guild.DiscordID;
        option.textContent = guild.Name; 
        guildBox.appendChild(option);
    });

}