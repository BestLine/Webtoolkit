function stop(state) {
    let url = "/destroy?state=" + state
    console.log("stop URL: ", url)
    showLoading()
    fetch(url)
        .then(response => {
            if (!response.ok) {
                alert("STOP: Network response was not ok")
                throw new Error('STOP: Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            alert("Тест %s успешно остановлен!" % state)
            console.log('STOP: Request succeeded with JSON response', data);
        })
        .catch(error => {
            alert('STOP: There has been a problem with your fetch operation: %s' % error)
            console.error('STOP: There has been a problem with your fetch operation:', error);
        });
    hideLoading()
}