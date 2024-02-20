document.addEventListener('DOMContentLoaded', function () {
    let form = document.querySelector(".adm_subs");

    form.addEventListener('submit', function (event) {
        event.preventDefault();

        const phone = document.querySelector(".phone");
        const projectCheckboxes = document.querySelectorAll('input[name="projects"]:checked');

        const userData = {
            user: phone.value,
            projects: Array.from(projectCheckboxes).map(cb => cb.value).join(',')
        };

        // Добавьте код для отправки данных на сервер, например, с использованием Fetch API
        console.log('Sending data:', userData);

        // Пример использования Fetch API
        fetch('/assignProjects', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: new URLSearchParams(userData)

        })
            .then(response => document.querySelector('.error_message').textContent = response.status)
            .then(data => document.querySelector('.error_message').textContent =  data)
            .catch(error => document.querySelector('.error_message').textContent =  error);
    });
});