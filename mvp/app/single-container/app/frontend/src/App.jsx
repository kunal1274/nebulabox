import { useState, useEffect } from 'react'
import axios from 'axios'
import './index.css'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:5000'

function App() {
  const [todos, setTodos] = useState([])
  const [newTodo, setNewTodo] = useState('')
  const [editingId, setEditingId] = useState(null)
  const [editText, setEditText] = useState('')
  const [loading, setLoading] = useState(true)
  const [connected, setConnected] = useState(false)

  // Check API connection
  useEffect(() => {
    const checkHealth = async () => {
      try {
        const response = await axios.get(`${API_URL}/health`)
        setConnected(response.status === 200)
      } catch (error) {
        setConnected(false)
        console.error('API health check failed:', error)
      }
    }
    checkHealth()
    const interval = setInterval(checkHealth, 5000)
    return () => clearInterval(interval)
  }, [])

  // Load todos
  useEffect(() => {
    loadTodos()
  }, [])

  const loadTodos = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`${API_URL}/api/todos`)
      setTodos(response.data)
    } catch (error) {
      console.error('Failed to load todos:', error)
    } finally {
      setLoading(false)
    }
  }

  const addTodo = async (e) => {
    e.preventDefault()
    if (!newTodo.trim()) return

    try {
      const response = await axios.post(`${API_URL}/api/todos`, {
        title: newTodo.trim(),
        completed: false
      })
      setTodos([...todos, response.data])
      setNewTodo('')
    } catch (error) {
      console.error('Failed to add todo:', error)
      alert('Failed to add todo')
    }
  }

  const updateTodo = async (id, updates) => {
    try {
      const response = await axios.put(`${API_URL}/api/todos/${id}`, updates)
      setTodos(todos.map(t => t._id === id ? response.data : t))
    } catch (error) {
      console.error('Failed to update todo:', error)
      alert('Failed to update todo')
    }
  }

  const deleteTodo = async (id) => {
    try {
      await axios.delete(`${API_URL}/api/todos/${id}`)
      setTodos(todos.filter(t => t._id !== id))
    } catch (error) {
      console.error('Failed to delete todo:', error)
      alert('Failed to delete todo')
    }
  }

  const toggleComplete = (todo) => {
    updateTodo(todo._id, { completed: !todo.completed })
  }

  const startEdit = (todo) => {
    setEditingId(todo._id)
    setEditText(todo.title)
  }

  const saveEdit = (id) => {
    if (!editText.trim()) return
    updateTodo(id, { title: editText.trim() })
    setEditingId(null)
    setEditText('')
  }

  const cancelEdit = () => {
    setEditingId(null)
    setEditText('')
  }

  return (
    <div className="container">
      <h1>🚀 MERN Todo App - NebulaBox MVP</h1>
      
      <div className={`status ${connected ? 'connected' : 'disconnected'}`}>
        {connected ? '✅ Connected to Backend API' : '❌ Backend API Disconnected'}
      </div>

      <form onSubmit={addTodo} className="todo-form">
        <input
          type="text"
          value={newTodo}
          onChange={(e) => setNewTodo(e.target.value)}
          placeholder="Add a new todo..."
        />
        <button type="submit">Add</button>
      </form>

      {loading ? (
        <div className="loading">Loading todos...</div>
      ) : todos.length === 0 ? (
        <div className="loading">No todos yet. Add one above!</div>
      ) : (
        <ul className="todo-list">
          {todos.map(todo => (
            <li key={todo._id} className={`todo-item ${todo.completed ? 'completed' : ''}`}>
              <input
                type="checkbox"
                checked={todo.completed}
                onChange={() => toggleComplete(todo)}
              />
              {editingId === todo._id ? (
                <>
                  <input
                    type="text"
                    value={editText}
                    onChange={(e) => setEditText(e.target.value)}
                    onKeyPress={(e) => {
                      if (e.key === 'Enter') saveEdit(todo._id)
                      if (e.key === 'Escape') cancelEdit()
                    }}
                    autoFocus
                    style={{ flex: 1, padding: '8px', border: '2px solid #667eea', borderRadius: '3px' }}
                  />
                  <div className="todo-actions">
                    <button className="btn-edit" onClick={() => saveEdit(todo._id)}>Save</button>
                    <button className="btn-delete" onClick={cancelEdit}>Cancel</button>
                  </div>
                </>
              ) : (
                <>
                  <span className="todo-text">{todo.title}</span>
                  <div className="todo-actions">
                    <button className="btn-edit" onClick={() => startEdit(todo)}>Edit</button>
                    <button className="btn-delete" onClick={() => deleteTodo(todo._id)}>Delete</button>
                  </div>
                </>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

export default App

