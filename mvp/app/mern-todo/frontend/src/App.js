import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './App.css';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:5000';

function App() {
  const [todos, setTodos] = useState([]);
  const [text, setText] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Fetch todos on mount
  useEffect(() => {
    fetchTodos();
  }, []);

  const fetchTodos = async () => {
    try {
      setLoading(true);
      const response = await axios.get(`${API_URL}/api/todos`);
      setTodos(response.data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch todos');
      console.error('Error fetching todos:', err);
    } finally {
      setLoading(false);
    }
  };

  const addTodo = async (e) => {
    e.preventDefault();
    if (!text.trim()) return;

    try {
      setLoading(true);
      const response = await axios.post(`${API_URL}/api/todos`, { text });
      setTodos([response.data, ...todos]);
      setText('');
      setError(null);
    } catch (err) {
      setError('Failed to add todo');
      console.error('Error adding todo:', err);
    } finally {
      setLoading(false);
    }
  };

  const toggleTodo = async (id, completed) => {
    try {
      setLoading(true);
      const todo = todos.find(t => t._id === id);
      const response = await axios.put(`${API_URL}/api/todos/${id}`, {
        text: todo.text,
        completed: !completed,
      });
      setTodos(todos.map(t => (t._id === id ? response.data : t)));
      setError(null);
    } catch (err) {
      setError('Failed to update todo');
      console.error('Error updating todo:', err);
    } finally {
      setLoading(false);
    }
  };

  const deleteTodo = async (id) => {
    try {
      setLoading(true);
      await axios.delete(`${API_URL}/api/todos/${id}`);
      setTodos(todos.filter(t => t._id !== id));
      setError(null);
    } catch (err) {
      setError('Failed to delete todo');
      console.error('Error deleting todo:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App">
      <div className="container">
        <h1>📝 MERN Todo App</h1>
        <p className="subtitle">Built with MongoDB, Express, React, Node.js</p>

        {error && <div className="error">{error}</div>}

        <form onSubmit={addTodo} className="todo-form">
          <input
            type="text"
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder="Add a new todo..."
            className="todo-input"
            disabled={loading}
          />
          <button type="submit" disabled={loading || !text.trim()} className="add-button">
            {loading ? 'Adding...' : 'Add Todo'}
          </button>
        </form>

        <div className="todos-container">
          {loading && todos.length === 0 ? (
            <div className="loading">Loading todos...</div>
          ) : todos.length === 0 ? (
            <div className="empty">No todos yet. Add one above! 🎉</div>
          ) : (
            <ul className="todo-list">
              {todos.map((todo) => (
                <li key={todo._id} className={`todo-item ${todo.completed ? 'completed' : ''}`}>
                  <input
                    type="checkbox"
                    checked={todo.completed}
                    onChange={() => toggleTodo(todo._id, todo.completed)}
                    className="todo-checkbox"
                  />
                  <span className="todo-text">{todo.text}</span>
                  <button
                    onClick={() => deleteTodo(todo._id)}
                    className="delete-button"
                    disabled={loading}
                  >
                    🗑️
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="footer">
          <p>Backend: {API_URL}</p>
          <button onClick={fetchTodos} className="refresh-button" disabled={loading}>
            🔄 Refresh
          </button>
        </div>
      </div>
    </div>
  );
}

export default App;

