let chartData = {
    labels: [],
    datasets: []
};

let ctx = document.getElementById('myChart').getContext('2d');
let myChart = new Chart(ctx, {
    type: 'line',
    data: chartData,
    options: {
        scales: {
            x: {
                type: 'linear',
                position: 'bottom',
                scaleLabel: {
                    display: true,
                    labelString: 'Duration'
                }
            },
            y: {
                beginAtZero: true,
                title: {
                    display: true,
                    text: 'Thread count'
                }
            }
        }
    }
});

function randomColor() {
    var r = Math.floor(Math.random() * 256);
    var g = Math.floor(Math.random() * 256);
    var b = Math.floor(Math.random() * 256);
    return 'rgba(' + r + ',' + g + ',' + b + ', 1)';
}

// Функция для обновления данных графика
function updateChartData() {
    chartData.labels = [];
    chartData.datasets = [];

    var rows = document.getElementById('threadGroupTable').querySelectorAll('tbody tr');

    rows.forEach(function (row) {
        var name = row.cells[0].querySelector('input').value;
        var threadCount = parseInt(row.cells[2].querySelector('input.threads').value);
        var rampUp = parseInt(row.cells[3].querySelector('input.rampup').value);
        var duration = parseInt(row.cells[5].querySelector('input.duration').value);

        var borderColor = randomColor();

        var dataValues = [];
        for (var t = 0; t <= duration; t++) {
            var currentThreadCount = Math.min(threadCount, (t / rampUp) * threadCount);
            dataValues.push({ x: t, y: currentThreadCount });
        }

        dataValues.push({ x: duration, y: 0 });

        chartData.datasets.push({
            label: name,
            data: dataValues,
            fill: false,
            borderColor: borderColor,
            borderWidth: 1
        });
    });
    myChart.update();
}

function add_row() {
    console.log('add_row!')
    const table = document.getElementById('threadGroupTable');
    const row = table.insertRow();
    row.insertCell().innerHTML = '<input type="text" value="New Thread Group">';
    row.insertCell().innerHTML = `<select>
                                        <option value="Module1">Module1</option>
                                        <option value="Module2">Module2</option>
                                        <option value="Module3">Module3</option>
                                        <option value="Module4">Module4</option>
                                        <option value="Module5">Module5</option>
                                      </select>`;
    row.insertCell().innerHTML = `
                <input class="threads" type="number" value="0" min="0">
                <select>
                    <option value="continue">Continue</option>
                    <option value="startNextLoop">Start Next Thread Loop</option>
                    <option value="stopThread">Stop Thread</option>
                    <option value="stopTest">Stop Test</option>
                    <option value="stopTestNow">Stop Test Now</option>
                </select>
            `;
    row.insertCell().innerHTML = `<input class="rampup" type="number" value="10">`;
    row.insertCell().innerHTML = `<input class="delay" type="number" value="10">`;
    row.insertCell().innerHTML = `<input class="duration" type="number" value="10">`;

    updateChartData();
}

function delete_row() {
    const table = document.getElementById('threadGroupTable');
    const row = table.querySelector('tbody tr:last-child');
    if (row) row.remove();

    updateChartData();
}