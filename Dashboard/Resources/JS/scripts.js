let connectedIconOn = 'Resources/Img/link_on.svg';
let connectedIconOff = 'Resources/Img/link_off.svg';
let connectedIconAmber = 'Resources/Img/link_amber.svg';
let setIntervalIds = [];
let sessionId = '';
let initialWidgetOrdering = false;

// Page Load init --------------------------------------------------
window.onload = function() {
    const pageLoad = document.getElementById('PageLoading');
    pageLoad.style.display = 'block';
    assignEventListeners();
    loadWidgetStructure();
    pageLoad.style.display = 'none';
    setInterval(loadWidgetStructure, 5000); // Check for new Widgets
}

function assignEventListeners() {
    const mainToggle = document.getElementById('WriteUpdates');
    mainToggle.addEventListener('change', function() {
        globalToggle(mainToggle.checked)
    });
}

function loadWidgetOrder() {
    var order = JSON.parse(localStorage.getItem('widgetOrder'));
    if (order) {
        order.forEach(function(id) {
            var widget = $('#' + id);
            $('#WidgetContainer').append(widget); // Append in the correct order
        });
    }
}

function saveWidgetOrder() {
    var order = [];
    $('#WidgetContainer .Widget').each(function(index, element) {
        order.push($(element).attr('id'));
    });
    localStorage.setItem('widgetOrder', JSON.stringify(order));
}

function resetPage() {
    sessionId = '';
    const widgetContainer = document.getElementById('WidgetContainer');

    // Unset the Widget update functions
    for (const id of setIntervalIds) {
        clearInterval(id);
    }
    setIntervalIds.length = 0;

    // Clear Chart Instances
    chartInstances.clear();
    // Delete the Widgets
    while (widgetContainer.firstChild) {
        widgetContainer.removeChild(widgetContainer.firstChild);
    }
}

// Widgets ------------------------------------------------------------------
function loadWidgetStructure() {    
    // Get information on what Widgets we have
    fetchData('').then(widgets => {

        if (!widgets) {
            return
        }

        // Check the SessionID hasn't changed, if it has we need to Wipe the Page as the Bot has restarted
        if (widgets && widgets.length > 0) {
            inboundSessionId = widgets[0].SessionID;
            if (sessionId == '') {
                sessionId = inboundSessionId;
            } else {
                if (sessionId != inboundSessionId) {
                    resetPage();
                    return
                }
            }
        }

        const writeUpdates = document.getElementById('WriteUpdates');
        if (!writeUpdates.checked) {
            return
        }

        const widgetContainer = document.getElementById('WidgetContainer');

        widgets.forEach(widget => {

            const widgetName = widget.Widget;
            const refreshMs = widget.RefreshMs;
            // Check if the holder div already exists
            let existingWidget = widgetContainer.querySelector(`div[widgetName="${widgetName}"]`);
            if (!existingWidget) {
                // Create it
                const newWidgetDiv = document.createElement('div');
                newWidgetDiv.setAttribute('widgetName', widgetName);
                newWidgetDiv.setAttribute('widgetInit', false);
                newWidgetDiv.id = 'Widget' + widget.Widget.replace(/\s+/g, '');
                
                widgetContainer.appendChild(newWidgetDiv);

                updateWidget(widgetName, newWidgetDiv);
                
                // Use the specific refresh interval from RefreshMs
                setIntervalIds.push(setInterval(() => updateWidget(widgetName, newWidgetDiv), refreshMs));
            }
        });

        if (!initialWidgetOrdering) {
            loadWidgetOrder();
        }
    }).catch(error => {
        console.error(`Error fetching widget data: ${error}`);
    });

    loadPageElements();
}

function loadPageElements() {
    $('#WidgetContainer').sortable({
        placeholder: "DragDropPlaceholder",
        helper: 'clone',
        scroll: false,
        handle: '.DragHandle',
        start: function(event, ui) {
            ui.helper.addClass('DragDropCursor');
            ui.item.addClass('DragDropHide');
            
            var item = ui.item;
            var width = item.outerWidth();
            var height = item.outerHeight();
            
            $('.DragDropPlaceholder').css({
                width: (width-8) + 'px',
                height: (height-8) + 'px'
            });
        },
        stop: function(event, ui) {
            ui.item.removeClass('DragDropHide');
            $('.DragDropPlaceholder').css({
                width: '',
                height: ''
            });
            saveWidgetOrder();
        }
    });
}

// Shared Functions ---------------------------------------------------------
// Gets Data from the GoLang service via a HTTP Call
async function fetchData(widget) {
    try {

        // 1. Get the getData URL
        const currentUrl = window.location.href;
        let dataUrl = `${currentUrl}getData`;
        if (widget != '') {
            dataUrl += "?widget=" + encodeURIComponent(widget);
        }

        // 2. Fetch the JSON data 
        const response = await fetch(dataUrl);
        if (!response.ok) {
            console.error(new Error(`HTTP error! status: ${response.status}`));
            return false;
        }

        // 3. Parse the JSON data and save it as an Object
        const dataObject = await response.json();
        return dataObject;

    } catch (error) {
        console.error(`Error fetching data: ${error}`);
        return false;
    }
}

async function updateWidget(widget, widgetDiv) {
    try {
        const data = await fetchData(widget);

        if (data == false) {
            const connectedIconOFF = widgetDiv.querySelector('.ConnectedOFF');
            const connectedIconAMBER = widgetDiv.querySelector('.ConnectedAMBER');
            const connectedIconON = widgetDiv.querySelector('.ConnectedON');

            connectedIconOFF.style.display = 'block';
            connectedIconAMBER.style.display = 'none';
            connectedIconON.style.display = 'none';
            connectedIconOFF.title = 'Update failed at: ' + getCurrentTime();
            return
        }

        if (data.SessionID != sessionId) {
            resetPage();
            return
        }

        // Has the Widget had its first time setup?
        if (widgetDiv.getAttribute('widgetInit') == 'false') {
            if (data.Type != 'info') {
                // => Set Widget class and Width
                widgetDiv.classList.add('Widget')
                if (data.Options.Width != null && data.Options.Width != '') {
                    switch(data.Options.Width) {
                        case '100%':
                            widgetDiv.classList.add('WidgetWidth100');
                            break;
                        case '75%':
                            widgetDiv.classList.add('WidgetWidth75');
                            break;
                        case '50%':
                            widgetDiv.classList.add('WidgetWidth50');
                            break;
                        case '25%':
                            widgetDiv.classList.add('WidgetWidth25');
                            break;
                        default:
                            widgetDiv.style.width = data.Options.Width;
                    }
                    
                }
            }
  
            switch(data.Type) {
                case 'table':
                    initTableWidget(data, widgetDiv);
                    break;
                case 'graph':
                    initGraphWidget(data, widgetDiv);
                    break;
                case 'info':
                    initInfoWidget(data, widgetDiv);
                    break;
                default:
                    console.error(`Unknown data type: ${data.Type}`);
                    console.error(data);
                }
        }

        switch(data.Type) {
            case 'table':
                updateTable(data, widgetDiv);
                break;
            case 'graph':
                updateGraph(data, widgetDiv);
                break;
            default:
            }

    } catch (error) {
        console.error(`updateWidget Error:\nError: ${error}\nData:`, data, `\nWidget Div:`, widgetDiv);
    }
}

// Update Widgets -----------------------------------------------------------
async function updateTable(data, widgetDiv) {
    const loadingBlock = widgetDiv.querySelector('.Loading');
    const connectedIconOFF = widgetDiv.querySelector('.ConnectedOFF');
    const connectedIconAMBER = widgetDiv.querySelector('.ConnectedAMBER');
    const connectedIconON = widgetDiv.querySelector('.ConnectedON');

    // Are updates disabled?
    const updateToggle = widgetDiv.querySelector('.updateToggle');
    if (!updateToggle.checked) {
        connectedIconOFF.style.display = 'none';
        connectedIconAMBER.style.display = 'block';
        connectedIconON.style.display = 'none';

        connectedIconAMBER.title = 'Updates manually disabled at: ' + getCurrentTime();
        return
    }

    try {
        // Set as loading
        loadingBlock.style.display = 'block';
        const tbody = widgetDiv.querySelector("tbody");

        // Clear current Table Entries
        tbody.innerHTML = '';

        if (data.Rows == null) {
            // Work out how many Columns our table has
            const theadRow = widgetDiv.querySelector("thead tr");
            const columnCount = theadRow ? theadRow.children.length : 0;
            
            // Create an "empty" row
            const tr = document.createElement("tr");
            const td = document.createElement("td");
            td.colSpan = columnCount;
            td.textContent = 'No rows';
            tr.appendChild(td);
            tbody.appendChild(tr);            
        } else {

            const colFilters = []; // [filterValue (str), Columns(int[]), FullMatchOnly(bool)]

            // Assign Indexes to the Columns (used for filtering)
            data.Columns.forEach((column, index) => {
                column.Index = index;
            });

            // Check for set Filters
            const filters = widgetDiv.querySelectorAll(".WidgetHeader .WidgetFilter");
            filters.forEach(filter => {
                if (filter.querySelector("select")) {
                    const selectFilter = filter.querySelector("select");
                    if (selectFilter) {
                        // Value
                        const selectValue = selectFilter.value;
                        if (!selectValue || selectValue == '') {
                            return
                        }

                        // Column Indexes
                        const selectColumnIndexes = [];
                        const selectColumns = selectFilter.getAttribute("Columns");
                        if (selectColumns) {
                            const columnsArray = selectColumns.split(",");
                            columnsArray.forEach(column => {
                                const trimmedColumn = column.trim();
                                const index = data.Columns.findIndex(column => column.Name === trimmedColumn);
                                if (index !== -1) {
                                    selectColumnIndexes.push(index);
                                }
                            });
                        } else {
                            console.error("Columns attribute not found");
                            return;
                        }

                        // Full Match Only
                        const selectFullMatch = selectFilter.getAttribute("FullMatchOnly");
                        let selectInsertFullMatch = false;
                        if (selectFullMatch) {
                            selectInsertFullMatch = selectFullMatch;
                        }

                        colFilters.push([selectValue, selectColumnIndexes, selectInsertFullMatch]);
                    }
                } else if (filter.querySelector('input[type="text"]')) {
                    const textFilter = filter.querySelector('input[type="text"]');
                    if (textFilter) {
                        // Value
                        const textValue = textFilter.value.trim();
                        if (textValue.length < 3) {
                            return
                        }

                        // Column Indexes
                        const textColumnIndexes = [];
                        const textColumns = textFilter.getAttribute("Columns");
                        if (textColumns) {
                            const columnsArray = textColumns.split(",");
                            columnsArray.forEach(column => {
                                const trimmedColumn = column.trim();
                                const index = data.Columns.findIndex(column => column.Name === trimmedColumn);
                                if (index !== -1) {
                                    textColumnIndexes.push(index);
                                }
                            });
                        } else {
                            console.error("Columns attribute not found");
                            return;
                        }

                        // Full Match Only
                        const textFullMatch = textFilter.getAttribute("FullMatchOnly");
                        let textInsertFullMatch = false;
                        if (textFullMatch) {
                            textInsertFullMatch = textFullMatch;
                        }

                        colFilters.push([textValue, textColumnIndexes, textInsertFullMatch]);
                    }
                } else if (filter.querySelector(".SelectCheckboxList")) {
                    const checkListFilter = filter.querySelector('.SelectCheckboxList');
                    if (checkListFilter) {
                        // Values
                        const valueArray = [];
                        const checkboxes = checkListFilter.querySelectorAll("input[type='checkbox']");
                        checkboxes.forEach(checkbox => {
                            if (checkbox.checked) {
                                valueArray.push(checkbox.value);
                            }
                        });
                        const checkListValue = valueArray.join('|');

                        // Column Indexes
                        const checkListColumnIndexes = [];
                        const checkListColumns = checkListFilter.getAttribute("Columns");
                        if (checkListColumns) {
                            const columnsArray = checkListColumns.split(",");
                            columnsArray.forEach(column => {
                                const trimmedColumn = column.trim();
                                const index = data.Columns.findIndex(column => column.Name === trimmedColumn);
                                if (index !== -1) {
                                    checkListColumnIndexes.push(index);
                                }
                            });
                        } else {
                            console.error("Columns attribute not found");
                            return;
                        }

                        // Full Match Only
                        const checkListFullMatch = checkListFilter.getAttribute("FullMatchOnly");
                        let checkListInsertFullMatch = false;
                        if (checkListFullMatch) {
                            checkListInsertFullMatch = checkListFullMatch;
                        }

                        colFilters.push([checkListValue, checkListColumnIndexes, checkListFullMatch]);
                    }
                    
                } else {
                    console.error("Unhandled filter found");
                }
            });

            // Iterate through our Rows
            data.Rows.forEach(row => {
                // Before processing the row, does it pass filtering?
                let skipRow = false;
                if (colFilters && colFilters.length > 0) {
                    let matchedFilters = 0;
                    colFilters.forEach(([matchValueRaw, colIndexes, fullMatch]) => {
                        let matchValue = matchValueRaw.trim().toLowerCase();
                        for (let c = 0; c < colIndexes.length; c++) {
                            let rowValue = row.Values[colIndexes[c]].Value.trim().toLowerCase();
                            if (row.Values[colIndexes[c]].HoverText) {
                                rowValue = row.Values[colIndexes[c]].HoverText.trim().toLowerCase();
                            }
                            if (matchValue.includes('|')) {
                                // Match any of the split values
                                const multiMatch = matchValue.split('|');
                                for (let p = 0; p < multiMatch.length; p++) {
                                    if (fullMatch) {
                                        if (rowValue == multiMatch[p].toLowerCase()) {
                                            matchedFilters++;
                                            return;
                                        }
                                    } else {
                                        if (rowValue.includes(multiMatch[p].toLowerCase())) {
                                            matchedFilters++;
                                            return;
                                        }
                                    }
                                }
                            } else {
                                // Match the Value directly
                                if (parseBoolean(fullMatch)) {
                                    if (rowValue == matchValue.toLowerCase()) {
                                        matchedFilters++;
                                        return;
                                    }
                                } else {
                                    if (rowValue.includes(matchValue.toLowerCase())) {
                                        matchedFilters++;
                                        return;
                                    }
                                }
                            }
                        }
                        
                    });
                    if (colFilters.length != matchedFilters) {
                        skipRow = true;
                    }
                }
                if (skipRow) {
                    return;
                }
                //---------------------------------------------------
                // Create the row

                // Create the <tr>
                const tr = document.createElement("tr");

                // Colour the Row?
                if (row.TextColour && row.TextColour.Html) {
                    tr.style.color = row.TextColour.Html;
                }

                // Iterate through each value
                row.Values.forEach(value => {

                    // Create the <td>
                    const td = document.createElement("td");
                    td.textContent = value.Value;

                     // Colour the Cell?
                    if (value.TextColour && value.TextColour.Html) {
                        td.style.color = value.TextColour.Html;
                    }

                    // Hover text
                    if (value.HoverText) {
                        td.classList.add('ToolTip');
                        td.title = value.HoverText;
                    }

                    // Apply coluring
                    if (value.TextColour && value.TextColour.Color) {
                        td.style.color = value.TextColour.Color;
                    }

                    tr.appendChild(td);
                });

                tbody.appendChild(tr);
            });
        }

        connectedIconOFF.style.display = 'none';
        connectedIconAMBER.style.display = 'none';
        connectedIconON.style.display = 'block';

        connectedIconON.title = 'Last update: ' + getCurrentTime();
        loadingBlock.style.display = 'none';
    } catch (error) {
        loadingBlock.style.display = 'none';
        
        connectedIconOFF.style.display = 'block';
        connectedIconAMBER.style.display = 'none';
        connectedIconON.style.display = 'none';
        connectedIconOFF.title = 'Error obtaining data at: ' + getCurrentTime() + ", check console for information";
        console.error(`updateTable Error: `, error);
        console.error(data);
        console.error(widgetDiv);
    }
}

async function updateGraph(data, widgetDiv) {
    const loadingBlock = widgetDiv.querySelector('.Loading');
    const connectedIconOFF = widgetDiv.querySelector('.ConnectedOFF');
    const connectedIconAMBER = widgetDiv.querySelector('.ConnectedAMBER');
    const connectedIconON = widgetDiv.querySelector('.ConnectedON');

    // Are updates disabled?
    const updateToggle = widgetDiv.querySelector('.updateToggle');
    if (!updateToggle.checked) {
        connectedIconOFF.style.display = 'none';
        connectedIconAMBER.style.display = 'block';
        connectedIconON.style.display = 'none';
        connectedIconAMBER.title = 'Updates manually disabled at: ' + getCurrentTime();
        return;
    }

    try {
        loadingBlock.style.display = 'block';

        const canvas = widgetDiv.querySelector('canvas');
        if (!canvas) {
            throw new Error('Could not obtain <canvas> tag');
        }

        const chartId = canvas.id;
        const ctx = canvas.getContext('2d');

        const chartDatasets = data.Options.Datasets.map(dataset => ({
            backgroundColor: dataset.BackgroundColour[0],
            borderColor: dataset.BorderColour[0],
            borderWidth: dataset.BorderWidth,
            data: dataset.Data,
            fill: dataset.Fill,
            label: dataset.Label,
            pointRadius: dataset.PointRadius
        }));

        // Check if a chart instance exists for this canvas
        if (chartInstances.has(chartId)) {

            let existingChart = Chart.getChart(canvas);

            // Only update the datasets
            existingChart.data.labels = data.Options.ChartLabels;
            existingChart.data.datasets = chartDatasets;
            existingChart.options.animation = false;

            // Trigger an update on the chart
            existingChart.update();
        } else {
            // Create a new chart (if it doesn't already exist)
            let chartConfig;
            if (data.Options.MinValue || data.Options.MaxValue) {
                chartConfig = {
                    type: data.Options.GraphWidgetChartType,
                    data: {
                        labels: data.Options.ChartLabels,
                        datasets: chartDatasets
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        animation: false,
                        scales: {
                            x: {
                                ticks: {
                                    display: false
                                },
                                title: {
                                    display: true,
                                    text: data.Options.XLabel
                                }
                            },
                            y: {
                                min: data.Options.MinValue,
                                max: data.Options.MaxValue,
                                ticks: {
                                    beginAtZero: true,
                                },
                                title: {
                                    display: true,
                                    text: data.Options.YLabel
                                }
                            }
                        }
                    }
                };
            } else {
                chartConfig = {
                    type: data.Options.GraphWidgetChartType,
                    data: {
                        labels: data.Options.ChartLabels,
                        datasets: chartDatasets
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
                                    text: data.Options.XLabel
                                }
                            },
                            y: {
                                ticks: {
                                    beginAtZero: true,
                                },
                                title: {
                                    display: true,
                                    text: data.Options.YLabel
                                }
                            }
                        }
                    }
                };
            }
            
            const newChart = new Chart(ctx, chartConfig);
            chartInstances.set(chartId, newChart);
        }

        connectedIconOFF.style.display = 'none';
        connectedIconAMBER.style.display = 'none';
        connectedIconON.style.display = 'block';
        connectedIconON.title = 'Last update: ' + getCurrentTime();
        loadingBlock.style.display = 'none';
    }
    catch (error) {
        loadingBlock.style.display = 'none';
        connectedIconOFF.style.display = 'block';
        connectedIconAMBER.style.display = 'none';
        connectedIconON.style.display = 'none';
        connectedIconOFF.title = 'Error obtaining data at: ' + getCurrentTime() + ", check console for information";
        console.error(`updateGraph Error:\nError: ${error}\nData:`, data, `\nWidget Div:`, widgetDiv);
    }
}

// Initialise Widgets --------------------------------------------------------
async function initTableWidget(data, widgetDiv) {
    try {
        // Update the container <div>
        widgetDiv.classList.add('WidgetTable')  

        // Add the Header
        // => <div>
        const widgetHeader = document.createElement('div');
        widgetHeader.classList.add('WidgetHeader');

        // => <h2>
        const widgetHeading = Object.assign(document.createElement('h2'), {
            innerText: data.Options.Name
        });
        widgetHeader.appendChild(widgetHeading);

        // => Connected Icons
        const widgetConnIconOff = Object.assign(document.createElement('img'), {
            src: connectedIconOff,
            title: 'Not yet connected',
            className: 'ConnectedIcon ConnectedOFF'
        });
        widgetHeader.appendChild(widgetConnIconOff);
        const widgetConnIconAmber = Object.assign(document.createElement('img'), {
            src: connectedIconAmber,
            title: 'Updated stopped',
            className: 'ConnectedIcon ConnectedAMBER'
        });
        widgetHeader.appendChild(widgetConnIconAmber);
        const widgetConnIconOn = Object.assign(document.createElement('img'), {
            src: connectedIconOn,
            title: 'Connected',
            className: 'ConnectedIcon ConnectedON'
        });
        widgetHeader.appendChild(widgetConnIconOn);

        // => Drag Handle
        const widgetDragHandle = Object.assign(document.createElement('div'), {
            className: 'DragHandle'
        });
        widgetHeader.appendChild(widgetDragHandle);

        // => Update Switch 
        // =>   => (Container)
        const switchContainer = document.createElement('div');
        switchContainer.classList.add('SwitchContainer')
        
        // =>   => (Item Wrapper)
        const switchSwitchWrap = document.createElement('label');
        switchSwitchWrap.classList.add('switch');

        // =>   => (Switch)
        const switchSwitchWrapInput = Object.assign(document.createElement('input'), {
            className: 'updateToggle',
            type: 'checkbox',
            checked: true
        });
        switchSwitchWrap.appendChild(switchSwitchWrapInput)
        
        const switchSwitchWrapSpan = document.createElement('span');
        switchSwitchWrapSpan.classList.add('slider', 'round');
        switchSwitchWrap.appendChild(switchSwitchWrapSpan);
        switchContainer.appendChild(switchSwitchWrap);
        widgetHeader.appendChild(switchContainer);
        widgetDiv.appendChild(widgetHeader);
        
                
        // => Any Filters?
        if (data.Filters != null && data.Filters.length > 0) {
            data.Filters.forEach(filter => {
                // Create the container
                const filterContainer = Object.assign(document.createElement('div'), {
                    className: 'WidgetFilter'
                });
                let fullMatch = false;
                if (filter.FullMatchOnly != 'undefined') {
                    fullMatch = filter.FullMatchOnly;
                }
                switch(filter.FilterType) {
                    case 1: // SelectRegular
                        const selectRegular = document.createElement('select');
                        const optionElement = document.createElement('option');
                        optionElement.value = '';
                        optionElement.textContent = filter.Name + '...';  
                        selectRegular.appendChild(optionElement); 
                        filter.Values.forEach(function(optionValue) {
                            const optionElement = document.createElement('option');
                            optionElement.value = optionValue;
                            optionElement.textContent = optionValue;  
                            selectRegular.appendChild(optionElement); 
                        });
                        selectRegular.setAttribute('FullMatchOnly', fullMatch);
                        selectRegular.setAttribute('Columns', filter.ColumnNames.join(','));

                        selectRegular.addEventListener('change', function() {
                            updateWidget(data.Options.Name, widgetDiv);
                        });

                        filterContainer.appendChild(selectRegular);
                        break;
                    case 2: // SelectCheckbox
                        const selectCheckbox = document.createElement('div');
                        selectCheckbox.classList.add('SelectCheckbox');

                        const selectCheckboxLabel = document.createElement('div');
                        selectCheckboxLabel.classList.add('SelectCheckboxLabel');
                        selectCheckboxLabel.innerText = filter.Name + '...';
                        selectCheckbox.appendChild(selectCheckboxLabel); 

                        const selectCheckboxList = document.createElement('div');
                        selectCheckboxList.setAttribute('FullMatchOnly', fullMatch);
                        selectCheckboxList.setAttribute('Columns', filter.ColumnNames.join(','));
                        selectCheckboxList.classList.add('SelectCheckboxList');
                        filter.Values.forEach(function(optionValue) {
                            let selectCheckboxListItem = document.createElement('div');
                            selectCheckboxListItem.classList.add('SelectCheckboxListItem');
                            let checkName = optionValue.replace(/\s+/g, '');

                            let sclCheckbox = document.createElement('input');
                            sclCheckbox.setAttribute('type', 'checkbox');
                            sclCheckbox.setAttribute('name', checkName);
                            sclCheckbox.setAttribute('checked', 'checked');
                            sclCheckbox.setAttribute('value', optionValue);
                            sclCheckbox.addEventListener('change', function() {
                                updateWidget(data.Options.Name, widgetDiv);
                            });
                            selectCheckboxListItem.appendChild(sclCheckbox); 

                            let sclLabel = document.createElement('label');
                            sclLabel.setAttribute('for', checkName);
                            sclLabel.innerText = optionValue;
                            selectCheckboxListItem.appendChild(sclLabel); 


                            selectCheckboxList.appendChild(selectCheckboxListItem); 
                        });
                        selectCheckbox.appendChild(selectCheckboxList);
                        filterContainer.appendChild(selectCheckbox);
                        break;
                    case 3: // FreeText
                        const freeText = document.createElement('input');
                        freeText.setAttribute('type', 'text');
                        freeText.setAttribute('FullMatchOnly', fullMatch);
                        freeText.setAttribute('placeholder', filter.Name + '...');
                        freeText.setAttribute('Columns', filter.ColumnNames.join(','));
                        freeText.addEventListener('input', function() {
                            updateWidget(data.Options.Name, widgetDiv);
                        });

                        filterContainer.appendChild(freeText);
                        break;
                    default:
                        console.error('Unknown FilterType iota of ' + filter.FilterType);
                }

                widgetHeader.appendChild(filterContainer);
            });
        }

        // Create the content
        // => <div>
        const widgetContentContainer = document.createElement('div');
        widgetContentContainer.classList.add('WidgetContent');

        // Add the Loading <div>
        const widgetLoader = Object.assign(document.createElement('div'), {
            className: 'Loading',
            style: { display: 'none' }
        });
        widgetContentContainer.appendChild(widgetLoader);

        // => <table>
        const widgetTable = document.createElement('table');

        // => <thead>
        const widgetTableHead = document.createElement('thead');

        // => <tr>
        const widgetTableHeadRow = document.createElement('tr');
        data.Columns.forEach(column => {
            // => <th>
            const widgetColHeader = Object.assign(document.createElement('th'), {
                innerText: column.Name
            });
            widgetTableHeadRow.appendChild(widgetColHeader);
        });

        widgetTableHead.appendChild(widgetTableHeadRow);
        widgetTable.appendChild(widgetTableHead);

        // => <tbody>
        const widgetTableBody = document.createElement('tbody');
        widgetTable.appendChild(widgetTableBody);

        widgetContentContainer.appendChild(widgetTable);
        widgetDiv.appendChild(widgetContentContainer);

        // Set Init as completed
        widgetDiv.setAttribute('widgetInit', true);
    } catch (error) {
        console.error(`initTableWidget Error: `, error);
        console.error(data);
        console.error(widgetDiv);
    }
}

let chartInstances = new Map()

async function initGraphWidget(data, widgetDiv) {
    try {
        // Update the container <div>
        widgetDiv.classList.add('WidgetGraph')  

         // Add the Loading <div>
        const widgetLoader = Object.assign(document.createElement('div'), {
            className: 'Loading',
            style: { display: 'none' }
        });
        widgetDiv.appendChild(widgetLoader);

        // Add the Header
        // => <div>
        const widgetHeader = document.createElement('div');
        widgetHeader.classList.add('WidgetHeader');

        // => <h2>
        const widgetHeading = Object.assign(document.createElement('h2'), {
            innerText: data.Options.Name
        });
        widgetHeader.appendChild(widgetHeading);

        // => Connected Icon
        const widgetConnIconOff = Object.assign(document.createElement('img'), {
            src: connectedIconOff,
            title: 'Not yet connected',
            className: 'ConnectedIcon ConnectedOFF'
        });
        widgetHeader.appendChild(widgetConnIconOff);
        const widgetConnIconAmber = Object.assign(document.createElement('img'), {
            src: connectedIconAmber,
            title: 'Updated stopped',
            className: 'ConnectedIcon ConnectedAMBER'
        });
        widgetHeader.appendChild(widgetConnIconAmber);
        const widgetConnIconOn = Object.assign(document.createElement('img'), {
            src: connectedIconOn,
            title: 'Connected',
            className: 'ConnectedIcon ConnectedON'
        });
        widgetHeader.appendChild(widgetConnIconOn);

        // => Update Switch 
        // =>   => (Container)
        const switchContainer = document.createElement('div');
        switchContainer.classList.add('SwitchContainer')
        
        // =>   => (Item Wrapper)
        const switchSwitchWrap = document.createElement('label');
        switchSwitchWrap.classList.add('switch');

        // =>   => (Switch)
        const switchSwitchWrapInput = Object.assign(document.createElement('input'), {
            className: 'updateToggle',
            type: 'checkbox',
            checked: true
        });
        switchSwitchWrap.appendChild(switchSwitchWrapInput)
        
        // => Drag Handle
        const widgetDragHandle = Object.assign(document.createElement('div'), {
            className: 'DragHandle'
        });
        widgetHeader.appendChild(widgetDragHandle);

        const switchSwitchWrapSpan = document.createElement('span');
        switchSwitchWrapSpan.classList.add('slider', 'round');
        switchSwitchWrap.appendChild(switchSwitchWrapSpan);
        switchContainer.appendChild(switchSwitchWrap);
        widgetHeader.appendChild(switchContainer);
        widgetDiv.appendChild(widgetHeader);

        // Add the Content Container <div>
        const widgetContentContainer = document.createElement('div');
        widgetContentContainer.classList.add('WidgetContent');

        // => Add ChartContainer <div>
        const chartContainer = document.createElement('div');
        chartContainer.classList.add('ChartContainer');

        // => Add Canvas <canvas>
        const chartCanvas = document.createElement('canvas');
        chartCanvas.id = data.Options.Name.replace(/\s+/g, '');
        chartCanvas.classList.add('chartCanvas');
        widgetContentContainer.appendChild(chartCanvas);

        widgetDiv.appendChild(widgetContentContainer);

        // Set Init as completed
        widgetDiv.setAttribute('widgetInit', true);
    } catch (error) {
        console.error(`initGraphWidget Error:\nError: ${error}\nData:`, data, `\nWidget Div:`, widgetDiv);
    }
}

async function initInfoWidget(data, widgetDiv) {
    try {
        widgetDiv.classList = '';
        data.Items.forEach(item => {
            const element = document.createElement('input');
            element.setAttribute('type', 'hidden');
            element.setAttribute('id', item.Name);
            element.setAttribute('name', item.Name);
            element.setAttribute('value', item.Value);
            element.setAttribute('description', item.Description);
            widgetDiv.appendChild(element);
        });

        widgetDiv.setAttribute('widgetInit', true);
    } catch (error) {
        console.error(`initInfoWidget Error:\nError: ${error}\nData:`, data, `\nWidget Div:`, widgetDiv);
    }
}

// Other Functions ----------------------------------------------------------
function globalToggle(inChecked) {
    const widgetContainer = document.querySelector('div#WidgetContainer');
    const inputs = widgetContainer.querySelectorAll('div.Widget input.updateToggle');
    inputs.forEach(input => {
        input.checked = inChecked;
    });
}
function parseBoolean(str) {
    return str.trim().toLowerCase() === 'true';
}
// Helpers ------------------------------------------------------------------
function getCurrentTime() {
    const now = new Date();
    
    // Get hours, minutes, and seconds
    const hours = String(now.getHours()).padStart(2, '0');
    const minutes = String(now.getMinutes()).padStart(2, '0');
    const seconds = String(now.getSeconds()).padStart(2, '0');
    
    // Format as HH:mm:ss
    return `${hours}:${minutes}:${seconds}`;
}