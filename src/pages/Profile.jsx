import { useContext } from 'react'
import { AuthContext } from '../context/AuthContext'

export default function Profile() {
    const { user, logout } = useContext(AuthContext)

    return (
        <div className="profile">
            <h2>Profile</h2>
            {user && (
                <>
                    <p>User ID: {user.user_id}</p>
                    <button onClick={logout}>Logout</button>
                </>
            )}
        </div>
    )
}