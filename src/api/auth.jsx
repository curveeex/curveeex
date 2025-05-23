import axios from 'axios'

const API_URL = 'http://localhost:8080'

export const register = async (email, password) => {
    const response = await axios.post(`${API_URL}/signup`, {
        email,
        password
    })
    return response.data
}

export const login = async (email, password) => {
    const response = await axios.post(`${API_URL}/signin`, {
        email,
        password
    })
    if (response.data.access_token) {
        localStorage.setItem('user', JSON.stringify(response.data))
    }
    return response.data
}

export const logout = async () => {
    const user = JSON.parse(localStorage.getItem('user'))
    if (user?.access_token) {
        await axios.post(`${API_URL}/logout`, {}, {
            headers: { Authorization: `Bearer ${user.access_token}` }
        })
    }
    localStorage.removeItem('user')
}

export const getProfile = async () => {
    const user = JSON.parse(localStorage.getItem('user'))
    const response = await axios.get(`${API_URL}/profile`, {
        headers: { Authorization: `Bearer ${user?.access_token}` }
    })
    return response.data
}