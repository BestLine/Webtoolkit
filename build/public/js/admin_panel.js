document.addEventListener("DOMContentLoaded", function() {
    let form_add_project = document.querySelector(".adm_c_project");
    let form_add_user_to_project = document.querySelector(".adm_add_user_to_project");

    form_add_project.addEventListener("submit", function(event) {
        event.preventDefault(); // Предотвращаем отправку формы по умолчанию
        form_add_project.action = "/beeload/add/project";
        var formData = new FormData(form_add_project);
        var jsonObject = {};
        formData.forEach(function(value, key){
            jsonObject[key] = value;
        });
        var jsonData = JSON.stringify(jsonObject);
        let msg = document.querySelector('.error_message')
        fetch(form_add_project.action, {
            method: form_add_project.method,
            body: jsonData,
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(response => response.text())
            .then(data => {
                console.log(data);
                msg.textContent = data
            })
            .catch(error => {
                console.error('Ошибка:', error);
                msg.textContent = ('Ошибка:' + error);
            });
    });

    document.getElementById('syncProjectsBtn').addEventListener('click', function() {
        console.log('Синхронизация списка проектов');
        let xhr = new XMLHttpRequest();
        let msg = document.querySelector('.error_message')
        let url = "/bucket/sync"
        let message = "Синхронизация прошла успешно"
        xhr.open("GET", url, true);
        xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
        xhr.addEventListener("load", function() {
            if (xhr.status === 200) { // Коды ответов
                msg.textContent = message
            } else if (xhr.status === 404) {
                msg.textContent = "Ресурс не найден: " + xhr.status
            } else if (xhr.status === 500) {
                msg.textContent = "Внутренняя ошибка сервера: " + xhr.status
            } else {
                msg.textContent = "Неизвестный код ответа: " + xhr.status
            }
        });
        xhr.send();
    });

    form_add_user_to_project.addEventListener("submit", function(event) {
        event.preventDefault();
        form_add_user_to_project.action = "/beeload/add/user_to_project";
        var formData = new FormData(form_add_user_to_project);
        var jsonObject = {};
        formData.forEach(function(value, key){
            jsonObject[key] = value;
        });
        var jsonData = JSON.stringify(jsonObject);
        let msg = document.querySelector('.error_message')
        fetch(form_add_user_to_project.action, {
            method: form_add_project.method,
            body: jsonData,
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(response => response.text())
            .then(data => {
                console.log(data);
                msg.textContent = data
            })
            .catch(error => {
                console.error('Ошибка:', error);
                msg.textContent = ('Ошибка:' + error);
            });
    });
});