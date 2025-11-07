#!/usr/bin/env node

/**
 * Code Generator from NebulaBox Schema
 * 
 * Generates:
 * - TypeScript types for frontend
 * - Go types for backend
 * - API client SDK
 * - Database migrations
 * - UI components
 */

const fs = require('fs');
const path = require('path');

const SCHEMA_PATH = path.join(__dirname, '../schema/nebulabox.schema.json');
const OUTPUT_DIR = path.join(__dirname, '../generated');

// Load schema
const schema = JSON.parse(fs.readFileSync(SCHEMA_PATH, 'utf8'));

// Ensure output directory exists
if (!fs.existsSync(OUTPUT_DIR)) {
  fs.mkdirSync(OUTPUT_DIR, { recursive: true });
}

// Generate TypeScript types
function generateTypeScriptTypes() {
  const output = [];
  output.push('// Auto-generated from schema/nebulabox.schema.json');
  output.push('// DO NOT EDIT MANUALLY');
  output.push('');
  
  // Generate type definitions
  for (const [name, def] of Object.entries(schema.definitions)) {
    output.push(`export interface ${name} {`);
    
    for (const [prop, spec] of Object.entries(def.properties || {})) {
      const required = (def.required || []).includes(prop);
      const optional = required ? '' : '?';
      
      let tsType = 'any';
      if (spec.type === 'string') {
        tsType = 'string';
      } else if (spec.type === 'integer' || spec.type === 'number') {
        tsType = 'number';
      } else if (spec.type === 'boolean') {
        tsType = 'boolean';
      } else if (spec.type === 'array') {
        tsType = `Array<${spec.items?.type || 'any'}>`;
      } else if (spec.type === 'object') {
        tsType = 'Record<string, any>';
      } else if (Array.isArray(spec.type)) {
        tsType = spec.type.map(t => t === 'null' ? 'null' : t).join(' | ');
      }
      
      output.push(`  ${prop}${optional}: ${tsType};`);
    }
    
    output.push('}');
    output.push('');
  }
  
  const filePath = path.join(OUTPUT_DIR, 'types.ts');
  fs.writeFileSync(filePath, output.join('\n'));
  console.log(`✅ Generated TypeScript types: ${filePath}`);
}

// Generate Go types
function generateGoTypes() {
  const output = [];
  output.push('// Auto-generated from schema/nebulabox.schema.json');
  output.push('// DO NOT EDIT MANUALLY');
  output.push('package generated');
  output.push('');
  output.push('import "time"');
  output.push('');
  
  // Generate struct definitions
  for (const [name, def] of Object.entries(schema.definitions)) {
    output.push(`type ${name} struct {`);
    
    for (const [prop, spec] of Object.entries(def.properties || {})) {
      const required = (def.required || []).includes(prop);
      
      let goType = 'interface{}';
      if (spec.type === 'string') {
        goType = 'string';
      } else if (spec.type === 'integer' || spec.type === 'number') {
        goType = 'int';
      } else if (spec.type === 'boolean') {
        goType = 'bool';
      } else if (spec.type === 'array') {
        goType = '[]string'; // Simplified
      } else if (Array.isArray(spec.type) && spec.type.includes('null')) {
        goType = '*string'; // Pointer for nullable
      }
      
      if (spec.format === 'date-time') {
        goType = 'time.Time';
        if (!required) {
          goType = '*time.Time';
        }
      }
      
      const jsonTag = required ? `json:"${prop}"` : `json:"${prop},omitempty"`;
      const propName = prop.charAt(0).toUpperCase() + prop.slice(1);
      output.push(`\t${propName} ${goType} \`${jsonTag}\``);
    }
    
    output.push('}');
    output.push('');
  }
  
  const filePath = path.join(OUTPUT_DIR, 'types.go');
  fs.writeFileSync(filePath, output.join('\n'));
  console.log(`✅ Generated Go types: ${filePath}`);
}

// Generate API client SDK
function generateAPIClient() {
  const output = [];
  output.push('// Auto-generated from schema/nebulabox.schema.json');
  output.push('// DO NOT EDIT MANUALLY');
  output.push('');
  output.push('import { ApiClient } from "../web/dashboard/src/lib/api";');
  output.push('import type { Container, Workspace } from "./types";');
  output.push('');
  output.push('export class NebulaBoxAPI {');
  output.push('  constructor(private client: ApiClient) {}');
  output.push('');
  
  // Generate methods for each endpoint
  for (const [resource, endpoints] of Object.entries(schema.api.endpoints)) {
    for (const [action, endpoint] of Object.entries(endpoints)) {
      const method = endpoint.method.toLowerCase();
      const path = endpoint.path.replace(/:(\w+)/g, '${$1}');
      
      const methodName = `${resource}${action.charAt(0).toUpperCase() + action.slice(1)}`;
      
      output.push(`  async ${methodName}(`);
      
      // Add parameters
      const params = [];
      if (endpoint.params) {
        for (const [param, spec] of Object.entries(endpoint.params)) {
          params.push(`${param}: string`);
        }
      }
      if (endpoint.request) {
        params.push(`data: ${endpoint.request.type === 'object' ? 'any' : 'any'}`);
      }
      if (endpoint.query) {
        params.push(`query?: { ${Object.keys(endpoint.query).map(k => `${k}?: ${endpoint.query[k].type}`).join('; ')} }`);
      }
      
      output.push(`    ${params.join(', ')}`);
      output.push(`  ): Promise<${endpoint.response?.$ref ? endpoint.response.$ref.split('/').pop() : 'any'}> {`);
      
      // Build path with template literal
      let pathVar = endpoint.path.replace(/:(\w+)/g, '${$1}');
      
      // Build request parameters
      const requestParams = [];
      if (endpoint.request) {
        requestParams.push('data');
      }
      if (endpoint.query) {
        requestParams.push('query');
      }
      
      const paramsStr = requestParams.length > 0 ? ', ' + requestParams.join(', ') : '';
      output.push(`    return this.client.request('${endpoint.method}', \`${pathVar}\`${paramsStr});`);
      output.push('  }');
      output.push('');
    }
  }
  
  output.push('}');
  
  const filePath = path.join(OUTPUT_DIR, 'api-client.ts');
  fs.writeFileSync(filePath, output.join('\n'));
  console.log(`✅ Generated API client SDK: ${filePath}`);
}

// Main execution
console.log('🚀 Generating code from schema...\n');

try {
  generateTypeScriptTypes();
  generateGoTypes();
  generateAPIClient();
  
  console.log('\n✅ Code generation complete!');
  console.log(`📁 Output directory: ${OUTPUT_DIR}`);
} catch (error) {
  console.error('❌ Error generating code:', error);
  process.exit(1);
}

