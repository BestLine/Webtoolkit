function getMonitoringUrls () {
    const selectField = document.getElementById('settings_monitoring_project');
    console.log(selectField.value)
    const url='/getMonitoringUrl'
    let data = {
        project: selectField.value,
    };
    let xhr = new XMLHttpRequest();
    xhr.open('POST', url, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function() {
        if (xhr.status === 200) {
            console.log('Запрос успешно отправлен');
            let response = JSON.parse(xhr.responseText);
            updateTextarea(response);
        } else {
            console.log('Ошибка: ' + xhr.statusText);
        }
    }
    xhr.send(JSON.stringify(data));
}

function addMonitoringUrl() {
    console.log("addMonitoringUrl")
    const selectField = document.getElementById('settings_monitoring_project');
    const new_url = document.querySelector(".monitoring-text-field")
    console.log((selectField.value).length)
    console.log((selectField.value))
    if (!selectField.value || selectField.value === "Выберите проект") {
        // Останавливаем отправку формы
        console.log("Проверка формы")
        selectField.setCustomValidity('Выбери проект.');
        selectField.reportValidity();
        return
    } else {
        console.log("не проверка формы")
        selectField.setCustomValidity('');
    }
    if (!new_url.value) {
        console.log("Проверка формы2")
        new_url.setCustomValidity('Вставь ссылку.');
        new_url.reportValidity();
        return
    } else {
        if (isValidURL(new_url.value) === false) {
            console.log("Проверка формы3")
            new_url.setCustomValidity('Не корректный формат.');
            new_url.reportValidity();
            return
        }
        console.log(isValidURL(new_url.value))
        console.log("не проверка формы")
        new_url.setCustomValidity('');
    }
    console.log("New url: ", new_url.value)
    let url = "/addMonitoringUrl"
    let data = {
        project: selectField.value,
        new_url: new_url.value
    };
    let xhr = new XMLHttpRequest();
    xhr.open('POST', url, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function() {
        if (xhr.status === 200) {
            console.log('Запрос успешно отправлен');
            getMonitoringUrls()
        } else {
            console.log('Ошибка: ' + xhr.statusText);
        }
    }
    xhr.send(JSON.stringify(data));
    //ouiehfn
}

function updateTextarea(urls) {
    const textarea = document.querySelector('.monitoring-text-area');
    if (urls !== null) {
        textarea.value = urls.join('\n');
    } else {
        textarea.value = 'Нет доступных ссылок';
    }
}

function isValidURL(url) {
    const pattern = new RegExp('^(https?:\\/\\/)?' + // protocol
        '((([a-zA-Z\\d]([a-zA-Z\\d-]*[a-zA-Z\\d])*)\\.?)+[a-zA-Z]{2,}|'+ // domain name
        '((\\d{1,3}\\.){3}\\d{1,3}))' + // OR ip (v4) address
        '(\\:\\d+)?(\\/[-a-zA-Z\\d%_.~+]*)*' + // port and path
        '(\\?[;&a-zA-Z\\d%_.~+=-]*)?' + // query string
        '(\\#[-a-zA-Z\\d_]*)?$', 'i'); // fragment locator
    return !!pattern.test(url);
}

$(document).ready(function (){
    console.log('monitoring-project.onchange READY!')
    $(".monitoring-project").on('change',function (){
        console.log('Change!')
        getMonitoringUrls()
    })
})