// обработка перехода на регистрацию
const registerLink = document.querySelector('.register_href');
registerLink.addEventListener('click', (event) => {
    event.preventDefault(); // Отменяем переход по ссылке (если она имеет атрибут href="#")
    $( "#wrapper" ).load( "/register");
    console.log('Ссылка "Здесь вы можете его создать" была нажата!');
});

const loginLink = document.querySelector('.login_href');
if (loginLink) {
    loginLink.addEventListener('click', (event) => {
        event.preventDefault(); // Отменяем переход по ссылке (если она имеет атрибут href="#")
        $( "#wrapper" ).load( "/login");
        console.log('Ссылка "Войти" была нажата!');
    });
}