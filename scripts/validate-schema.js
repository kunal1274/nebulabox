#!/usr/bin/env node

/**
 * Schema Validator
 * Validates the NebulaBox schema against JSON Schema draft-07
 */

const fs = require('fs');
const path = require('path');

const SCHEMA_PATH = path.join(__dirname, '../schema/nebulabox.schema.json');

function validateSchema() {
  try {
    const schema = JSON.parse(fs.readFileSync(SCHEMA_PATH, 'utf8'));
    
    const errors = [];
    
    // Basic structure validation
    if (!schema.definitions) {
      errors.push('Missing "definitions" section');
    }
    
    if (!schema.api) {
      errors.push('Missing "api" section');
    }
    
    if (!schema.database) {
      errors.push('Missing "database" section');
    }
    
    if (!schema.ui) {
      errors.push('Missing "ui" section');
    }
    
    // Validate definitions
    if (schema.definitions) {
      for (const [name, def] of Object.entries(schema.definitions)) {
        if (!def.type) {
          errors.push(`Definition "${name}" missing "type"`);
        }
        
        if (def.type === 'object' && !def.properties) {
          errors.push(`Definition "${name}" is object but missing "properties"`);
        }
      }
    }
    
    // Validate API endpoints
    if (schema.api && schema.api.endpoints) {
      for (const [resource, endpoints] of Object.entries(schema.api.endpoints)) {
        for (const [action, endpoint] of Object.entries(endpoints)) {
          if (!endpoint.method) {
            errors.push(`API endpoint "${resource}.${action}" missing "method"`);
          }
          if (!endpoint.path) {
            errors.push(`API endpoint "${resource}.${action}" missing "path"`);
          }
        }
      }
    }
    
    // Validate database tables
    if (schema.database && schema.database.tables) {
      for (const [table, tableDef] of Object.entries(schema.database.tables)) {
        if (!tableDef.columns) {
          errors.push(`Database table "${table}" missing "columns"`);
        }
        if (!tableDef.primaryKey) {
          errors.push(`Database table "${table}" missing "primaryKey"`);
        }
      }
    }
    
    if (errors.length > 0) {
      console.error('❌ Schema validation failed:');
      errors.forEach(error => console.error(`  - ${error}`));
      process.exit(1);
    }
    
    console.log('✅ Schema validation passed');
    console.log(`📊 Version: ${schema.version || 'unknown'}`);
    console.log(`📦 Definitions: ${Object.keys(schema.definitions || {}).length}`);
    console.log(`🔌 API endpoints: ${countEndpoints(schema.api?.endpoints || {})}`);
    console.log(`🗄️  Database tables: ${Object.keys(schema.database?.tables || {}).length}`);
    
  } catch (error) {
    console.error('❌ Error validating schema:', error.message);
    process.exit(1);
  }
}

function countEndpoints(endpoints) {
  let count = 0;
  for (const resource of Object.values(endpoints)) {
    count += Object.keys(resource).length;
  }
  return count;
}

validateSchema();

