import express from 'express'
import cors from 'cors'
import { MongoClient } from 'mongodb'
import dotenv from 'dotenv'

dotenv.config()

const app = express()
const PORT = process.env.PORT || 5000

// MongoDB connection
const MONGO_HOST = process.env.MONGO_HOST || 'localhost'
const MONGO_PORT = process.env.MONGO_PORT || 27017
const MONGO_DB = process.env.MONGO_DB || 'todos'

const mongoUrl = `mongodb://${MONGO_HOST}:${MONGO_PORT}/${MONGO_DB}`
let db = null
let client = null

// Middleware
app.use(cors())
app.use(express.json())

// Connect to MongoDB
async function connectDB() {
  try {
    client = new MongoClient(mongoUrl)
    await client.connect()
    db = client.db(MONGO_DB)
    console.log(`✅ Connected to MongoDB at ${mongoUrl}`)
    
    // Create todos collection if it doesn't exist
    await db.collection('todos').createIndex({ _id: 1 })
  } catch (error) {
    console.error('❌ MongoDB connection error:', error)
    // Retry connection after 5 seconds
    setTimeout(connectDB, 5000)
  }
}

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({
    status: 'healthy',
    service: 'mern-backend',
    database: db ? 'connected' : 'disconnected',
    timestamp: new Date().toISOString()
  })
})

// GET /api/todos - List all todos
app.get('/api/todos', async (req, res) => {
  try {
    if (!db) {
      return res.status(503).json({ error: 'Database not connected' })
    }
    const todos = await db.collection('todos').find({}).toArray()
    res.json(todos)
  } catch (error) {
    console.error('Error fetching todos:', error)
    res.status(500).json({ error: 'Failed to fetch todos' })
  }
})

// POST /api/todos - Create a new todo
app.post('/api/todos', async (req, res) => {
  try {
    if (!db) {
      return res.status(503).json({ error: 'Database not connected' })
    }
    const { title, completed } = req.body
    if (!title) {
      return res.status(400).json({ error: 'Title is required' })
    }
    
    const todo = {
      title: title.trim(),
      completed: completed || false,
      createdAt: new Date(),
      updatedAt: new Date()
    }
    
    const result = await db.collection('todos').insertOne(todo)
    const createdTodo = await db.collection('todos').findOne({ _id: result.insertedId })
    res.status(201).json(createdTodo)
  } catch (error) {
    console.error('Error creating todo:', error)
    res.status(500).json({ error: 'Failed to create todo' })
  }
})

// PUT /api/todos/:id - Update a todo
app.put('/api/todos/:id', async (req, res) => {
  try {
    if (!db) {
      return res.status(503).json({ error: 'Database not connected' })
    }
    const { id } = req.params
    const { title, completed } = req.body
    
    const update = { updatedAt: new Date() }
    if (title !== undefined) update.title = title.trim()
    if (completed !== undefined) update.completed = completed
    
    const result = await db.collection('todos').findOneAndUpdate(
      { _id: id },
      { $set: update },
      { returnDocument: 'after' }
    )
    
    if (!result.value) {
      return res.status(404).json({ error: 'Todo not found' })
    }
    
    res.json(result.value)
  } catch (error) {
    console.error('Error updating todo:', error)
    res.status(500).json({ error: 'Failed to update todo' })
  }
})

// DELETE /api/todos/:id - Delete a todo
app.delete('/api/todos/:id', async (req, res) => {
  try {
    if (!db) {
      return res.status(503).json({ error: 'Database not connected' })
    }
    const { id } = req.params
    
    const result = await db.collection('todos').findOneAndDelete({ _id: id })
    
    if (!result.value) {
      return res.status(404).json({ error: 'Todo not found' })
    }
    
    res.json({ message: 'Todo deleted', id })
  } catch (error) {
    console.error('Error deleting todo:', error)
    res.status(500).json({ error: 'Failed to delete todo' })
  }
})

// Start server
async function startServer() {
  await connectDB()
  
  app.listen(PORT, '0.0.0.0', () => {
    console.log(`🚀 Backend server running on http://0.0.0.0:${PORT}`)
    console.log(`📡 Health check: http://localhost:${PORT}/health`)
    console.log(`📝 API endpoint: http://localhost:${PORT}/api/todos`)
  })
}

// Handle shutdown
process.on('SIGTERM', async () => {
  console.log('SIGTERM received, closing connections...')
  if (client) {
    await client.close()
  }
  process.exit(0)
})

startServer().catch(console.error)

