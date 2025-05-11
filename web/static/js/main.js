document.addEventListener('DOMContentLoaded', function() {
    // Проверка, что страница загрузилась
    console.log('Web UI Loaded');

    // Пример обработчика для формы калькуляции
    const calcForm = document.querySelector('form[action="/calculate"]');
    if (calcForm) {
        calcForm.addEventListener('submit', function(event) {
            event.preventDefault();

            const expression = document.querySelector('input[name="expression"]').value;
            if (expression) {
                // Здесь можно добавить логику для отправки выражения на сервер или расчета
                console.log('Sending expression for calculation: ' + expression);

                // Простая отправка с помощью fetch (можно расширить для реальных запросов)
                fetch('/calculate', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: `expression=${encodeURIComponent(expression)}`
                })
                .then(response => response.json())
                .then(data => {
                    // Вывести результат
                    console.log('Result:', data);
                })
                .catch(error => console.error('Error:', error));
            }
        });
    }
});
