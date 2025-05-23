let accessToken = null;

document.getElementById('signup-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('signup-email').value;
    const password = document.getElementById('signup-password').value;

    try {
        const response = await fetch('/signup', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        const data = await response.json();
        showResult(`Регистрация успешна: ${JSON.stringify(data)}`);
    } catch (error) {
        showResult(`Ошибка регистрации: ${error.message}`);
    }
});

document.getElementById('signin-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const email = document.getElementById('signin-email').value;
    const password = document.getElementById('signin-password').value;

    try {
        const response = await fetch('/signin', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });
        const data = await response.json();
        accessToken = data.access_token;
        showResult(`Вход успешен. Токен сохранён: ${accessToken}`);
    } catch (error) {
        showResult(`Ошибка входа: ${error.message}`);
    }
});

document.getElementById('profile-btn').addEventListener('click', async () => {
    if (!accessToken) {
        showResult('Сначала войдите в систему');
        return;
    }

    try {
        const response = await fetch('/profile', {
            headers: { 'Authorization': `Bearer ${accessToken}` }
        });
        const data = await response.json();
        showResult(`Данные профиля: ${JSON.stringify(data)}`);
    } catch (error) {
        showResult(`Ошибка получения профиля: ${error.message}`);
    }
});

document.getElementById('logout-btn').addEventListener('click', () => {
    accessToken = null;
    showResult('Вы вышли из системы');
});

function showResult(message) {
    document.getElementById('result').textContent = message;
}