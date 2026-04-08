import React, { useState } from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import SpaceDetail from './pages/SpaceDetail'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(!!localStorage.getItem('token'))

  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Routes>
          <Route 
            path="/login" 
            element={!isAuthenticated ? <Login onLogin={() => setIsAuthenticated(true)} /> : <Navigate to="/" />} 
          />
          <Route 
            path="/" 
            element={isAuthenticated ? <Dashboard initialMode="personal" /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/collaboration" 
            element={isAuthenticated ? <Dashboard initialMode="collaboration" /> : <Navigate to="/login" />} 
          />
          <Route 
            path="/space/:id" 
            element={isAuthenticated ? <SpaceDetail /> : <Navigate to="/login" />} 
          />
        </Routes>
      </div>
    </Router>
  )
}

export default App
