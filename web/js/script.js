async function sendData() {
    var add = document.getElementById('add').value;
    var sub = document.getElementById('sub').value;
    var mult = document.getElementById('mult').value;
    var dev = document.getElementById('dev').value;
    var inputData = document.getElementById('dataInput').value;

    try {
        const response = await fetch('http://localhost:9999/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ add, sub, mult, dev, inputData }),
        });

        if (!response.ok) {
            throw new Error('Ошибка при отправке данных на сервер');
        }

        const responseData = await response.json();

        // Далее обрабатываем ответ от сервера и создаем блок с информацией
        var responseBlocks = document.getElementById('responseBlocks');

        var newBlock = document.createElement('div');
        newBlock.className = 'responseBlock';

        var id = responseData.reqID;
        var example = 'Пример: ' + inputData;
        var answer = 'Ответ: Ожидаем';
        var solveTime = 'Время решения: Ожидаем';
        var status = 'Статус: выполняется';

        newBlock.innerHTML = `<p>${id}</p><p>${example}</p><p>${answer}</p><p>${status}</p><p>${solveTime}</p>`;

        responseBlocks.appendChild(newBlock);
    } catch (error) {
        console.error('Ошибка:', error.message);
    }
}