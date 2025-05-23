import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import Navbar from './components/Navbar'
import Login from './pages/Login'
import Register from './pages/Register'
import Profile from './pages/Profile'
import './styles/main.css'

export default function App() {
    return (
        <AuthProvider>
            <BrowserRouter>
                <Navbar />
                <div className="container">
                    <Routes>
                        <Route path="/login" element={<Login />} />
                        <Route path="/register" element={<Register />} />
                        <Route path="/profile" element={<Profile />} />
                        <Route path="*" element={<Login />} />
                    </Routes>
                </div>
            </BrowserRouter>
        </AuthProvider>
    )
}