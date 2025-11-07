#!/usr/bin/env node

/**
 * Collaboration Sync Layer
 * 
 * Implements stateful collaboration like Cursor/Replit
 * - Real-time file synchronization
 * - Presence indicators
 * - Operational Transform for conflict resolution
 */

const WebSocket = require('ws');
const fs = require('fs');
const path = require('path');
const chokidar = require('chokidar');

class CollaborationSync {
  constructor(config = {}) {
    this.wsUrl = config.wsUrl || 'ws://localhost:8081/ws';
    this.workspaceId = config.workspaceId;
    this.userId = config.userId || `user-${Date.now()}`;
    this.ws = null;
    this.fileWatchers = new Map();
    this.pendingChanges = new Map();
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 10;
  }

  connect() {
    return new Promise((resolve, reject) => {
      this.ws = new WebSocket(`${this.wsUrl}?workspace=${this.workspaceId}&user=${this.userId}`);
      
      this.ws.on('open', () => {
        console.log('✅ Connected to collaboration server');
        this.reconnectAttempts = 0;
        this.sendPresence('online');
        resolve();
      });
      
      this.ws.on('message', (data) => {
        this.handleMessage(JSON.parse(data));
      });
      
      this.ws.on('error', (error) => {
        console.error('❌ WebSocket error:', error);
        reject(error);
      });
      
      this.ws.on('close', () => {
        console.log('⚠️  Disconnected from collaboration server');
        this.handleReconnect();
      });
    });
  }

  handleMessage(message) {
    switch (message.type) {
      case 'file-change':
        this.handleRemoteFileChange(message);
        break;
      case 'presence-update':
        this.handlePresenceUpdate(message);
        break;
      case 'cursor-update':
        this.handleCursorUpdate(message);
        break;
      case 'conflict':
        this.handleConflict(message);
        break;
      default:
        console.warn('Unknown message type:', message.type);
    }
  }

  handleRemoteFileChange(message) {
    const { file, content, userId, timestamp } = message;
    
    // Skip if this is our own change
    if (userId === this.userId) {
      return;
    }
    
    // Apply operational transform to merge changes
    const mergedContent = this.operationalTransform(
      this.getFileContent(file),
      content,
      timestamp
    );
    
    // Update file
    this.writeFile(file, mergedContent);
    console.log(`📥 Synced file: ${file} (from ${userId})`);
  }

  handlePresenceUpdate(message) {
    const { users } = message;
    console.log(`👥 Active users: ${users.length}`);
    // Emit presence event for UI
    this.emit('presence', users);
  }

  handleCursorUpdate(message) {
    const { userId, file, position } = message;
    // Emit cursor event for UI
    this.emit('cursor', { userId, file, position });
  }

  handleConflict(message) {
    const { file, localVersion, remoteVersion } = message;
    console.warn(`⚠️  Conflict detected in ${file}`);
    // Emit conflict event for manual resolution
    this.emit('conflict', { file, localVersion, remoteVersion });
  }

  operationalTransform(localContent, remoteContent, remoteTimestamp) {
    // Simple last-write-wins for now
    // In production, use proper OT library like ShareJS
    return remoteContent;
  }

  watchFile(filePath) {
    if (this.fileWatchers.has(filePath)) {
      return; // Already watching
    }
    
    const watcher = chokidar.watch(filePath, {
      persistent: true,
      ignoreInitial: true
    });
    
    watcher.on('change', () => {
      this.syncFile(filePath);
    });
    
    this.fileWatchers.set(filePath, watcher);
  }

  syncFile(filePath) {
    // Debounce file syncs
    if (this.pendingChanges.has(filePath)) {
      clearTimeout(this.pendingChanges.get(filePath));
    }
    
    const timeout = setTimeout(() => {
      const content = this.getFileContent(filePath);
      this.sendFileChange(filePath, content);
      this.pendingChanges.delete(filePath);
    }, 500); // 500ms debounce
    
    this.pendingChanges.set(filePath, timeout);
  }

  sendFileChange(filePath, content) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      return;
    }
    
    this.ws.send(JSON.stringify({
      type: 'file-change',
      file: filePath,
      content: content,
      userId: this.userId,
      timestamp: Date.now()
    }));
  }

  sendPresence(status) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      return;
    }
    
    this.ws.send(JSON.stringify({
      type: 'presence',
      userId: this.userId,
      status: status,
      timestamp: Date.now()
    }));
  }

  getFileContent(filePath) {
    try {
      return fs.readFileSync(filePath, 'utf8');
    } catch (error) {
      return '';
    }
  }

  writeFile(filePath, content) {
    try {
      fs.writeFileSync(filePath, content, 'utf8');
    } catch (error) {
      console.error(`Error writing file ${filePath}:`, error);
    }
  }

  handleReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('❌ Max reconnection attempts reached');
      return;
    }
    
    this.reconnectAttempts++;
    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
    
    console.log(`🔄 Reconnecting in ${delay}ms... (attempt ${this.reconnectAttempts})`);
    
    setTimeout(() => {
      this.connect().catch(() => {
        // Retry handled by handleReconnect
      });
    }, delay);
  }

  emit(event, data) {
    // Event emitter for UI integration
    if (this.onEvent) {
      this.onEvent(event, data);
    }
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
    }
    this.fileWatchers.forEach(watcher => watcher.close());
    this.fileWatchers.clear();
  }
}

// CLI usage
if (require.main === module) {
  const sync = new CollaborationSync({
    workspaceId: process.env.NEBULABOX_WORKSPACE_ID || 'default',
    userId: process.env.NEBULABOX_USER_ID || `user-${Date.now()}`
  });
  
  sync.connect()
    .then(() => {
      // Watch project files
      sync.watchFile('internal/api/server.go');
      sync.watchFile('web/dashboard/src/App.tsx');
      // Add more files as needed
      
      console.log('👀 Watching files for collaboration...');
    })
    .catch((error) => {
      console.error('Failed to connect:', error);
      process.exit(1);
    });
  
  // Graceful shutdown
  process.on('SIGINT', () => {
    console.log('\n👋 Disconnecting...');
    sync.disconnect();
    process.exit(0);
  });
}

module.exports = CollaborationSync;

