document.addEventListener('DOMContentLoaded', function () {
    let form = document.querySelector(".adm_subs");

    form.addEventListener('submit', function (event) {
        event.preventDefault();

        const userSelect = document.querySelector(".user");
        const projectCheckboxes = document.querySelectorAll('input[name="projects"]:checked');

        const userData = {
            user: userSelect.value,
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
            .then(response => response.json())
            .then(data => console.log('Server response:', data))
            .catch(error => console.error('Error:', error));
    });
});