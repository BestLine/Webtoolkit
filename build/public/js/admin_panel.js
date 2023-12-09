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
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
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
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    });
});