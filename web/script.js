document.addEventListener('DOMContentLoaded', () => {
    const fetchBtn = document.getElementById('fetchOrderBtn');
    const orderIdInput = document.getElementById('orderIdInput');
    const resultDiv = document.getElementById('result');

    fetchBtn.addEventListener('click', fetchOrder);
    orderIdInput.addEventListener('keypress', (event) => {
        if (event.key === 'Enter') {
            fetchOrder();
        }
    });

    async function fetchOrder() {
        const orderId = orderIdInput.value.trim();
        if (!orderId) {
            resultDiv.innerHTML = '<p class="error">Пожалуйста, введите ID заказа.</p>';
            return;
        }

        resultDiv.innerHTML = '<p>Загрузка данных...</p>';

        try {
            
            const response = await fetch(`http://localhost:8081/order/${orderId}`);

            if (response.status === 404) {
                 resultDiv.innerHTML = `<p class="error">Заказ с ID <strong>${orderId}</strong> не найден.</p>`;
                 return;
            }
            
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Ошибка сервера: ${response.status}. ${errorText}`);
            }

            const data = await response.json();
            displayOrderData(data);

        } catch (error) {
            console.error('Ошибка при получении данных:', error);
            resultDiv.innerHTML = `<p class="error">Не удалось получить данные. Проверьте консоль для деталей.</p>`;
        }
    }

    function displayOrderData(data) {
        const formattedJson = JSON.stringify(data, null, 2);
        resultDiv.innerHTML = `
            <h2>Детали заказа: ${data.order_uid}</h2>
            <pre>${formattedJson}</pre>
        `;
    }
});