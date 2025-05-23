const express = require('express');
const axios = require('axios');
const cors = require('cors');

const app = express();
const PORT = 3000;

// Настройка CORS
app.use(cors());
app.use(express.json());

// Прокси для бэкенда
const API_URL = 'http://localhost:8080';

// Роут для регистрации
app.post('/signup', async (req, res) => {
    try {
        const response = await axios.post(`${API_URL}/signup`, req.body);
        res.json(response.data);
    } catch (error) {
        res.status(error.response?.status || 500).json({ error: error.message });
    }
});

// Роут для входа
app.post('/signin', async (req, res) => {
    try {
        const response = await axios.post(`${API_URL}/signin`, req.body);
        res.json(response.data);
    } catch (error) {
        res.status(error.response?.status || 500).json({ error: error.message });
    }
});

// Роут для профиля (защищённый)
app.get('/profile', async (req, res) => {
    try {
        const token = req.headers.authorization?.split(' ')[1];
        const response = await axios.get(`${API_URL}/profile`, {
            headers: { Authorization: `Bearer ${token}` }
        });
        res.json(response.data);
    } catch (error) {
        res.status(error.response?.status || 500).json({ error: error.message });
    }
});

// Статические файлы (HTML, JS)
app.use(express.static('public'));

app.listen(PORT, () => {
    console.log(`Фронтенд запущен на http://localhost:${PORT}`);
});