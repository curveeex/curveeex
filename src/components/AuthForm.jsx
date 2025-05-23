import { useState, useContext } from 'react'
import { AuthContext } from '../context/AuthContext'
import { useNavigate } from 'react-router-dom'

export default function AuthForm({ isLogin }) {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')
    const { login, register } = useContext(AuthContext)
    const navigate = useNavigate()

    const handleSubmit = async (e) => {
        e.preventDefault()
        try {
            if (isLogin) {
                await login(email, password)
            } else {
                await register(email, password)
                await login(email, password)
            }
            navigate('/profile')
        } catch (err) {
            setError('Invalid credentials or server error')
            console.error('Auth error:', err)
        }
    }

    return (
        <div className="auth-form">
            <h2>{isLogin ? 'Login' : 'Register'}</h2>
            {error && <div className="error">{error}</div>}
            <form onSubmit={handleSubmit}>
                <input
                    type="email"
                    placeholder="Email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    minLength={8}
                />
                <button type="submit">{isLogin ? 'Login' : 'Register'}</button>
            </form>
        </div>
    )
}