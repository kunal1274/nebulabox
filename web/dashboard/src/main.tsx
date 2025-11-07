import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { initTestAware } from './lib/test-aware'
import { initApiTracking } from './lib/api-test-aware'

// Initialize test-aware utilities
initTestAware()
initApiTracking()

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)

