function stop(state) {
    let url = "http://ms-loadrtst038:9999/destroy?state=" + state
    console.log("stop URL: ", url)
    fetch(url)
        .then(response => {
            if (!response.ok) {
                throw new Error('STOP: Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            alert("Тест %s успешно остановлен!" % state)
            console.log('STOP: Request succeeded with JSON response', data);
        })
        .catch(error => {
            console.error('STOP: There has been a problem with your fetch operation:', error);
        });
}