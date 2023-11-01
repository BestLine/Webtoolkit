document.addEventListener("DOMContentLoaded", function() {
    var form = document.querySelector(".adm_c_project");

    form.addEventListener("submit", function(event) {
        event.preventDefault(); // Предотвращаем отправку формы по умолчанию
        form.action = "/beeload/add/project";

        var formData = new FormData(form);
        var jsonObject = {};
        formData.forEach(function(value, key){
            jsonObject[key] = value;
        });

        var jsonData = JSON.stringify(jsonObject);

        // Далее можно выполнить запрос на сервер, например, с использованием fetch или других методов

        // Пример использования fetch:
        fetch(form.action, {
            method: form.method,
            body: jsonData,
            headers: {
                'Content-Type': 'application/json'
            }
        })
            .then(response => response.text())
            .then(data => {
                console.log(data); // Вы можете обработать ответ сервера здесь
            })
            .catch(error => {
                console.error('Ошибка:', error);
            });
    });
});