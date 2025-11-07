import { test, expect } from '@playwright/test';
import { HealthChecker } from '../utils/test-helpers';

/**
 * Setup test - ensures services are ready before running tests
 */
test.describe('Setup and Health Checks', () => {
  const healthChecker = new HealthChecker();

  test('API server should be healthy', async ({ request }) => {
    const response = await request.get('http://localhost:8081/api/health');
    expect(response.ok()).toBeTruthy();
    
    const data = await response.json();
    expect(data.status).toBe('healthy');
    expect(data.service).toBe('nebulabox-api');
  });

  test('Frontend should be accessible', async ({ page }) => {
    await page.goto('http://localhost:3001');
    await expect(page).toHaveTitle(/NebulaBox|Dashboard/i);
  });

  test('API endpoints should be accessible', async ({ request }) => {
    // Test buildspec endpoint
    const buildspecResponse = await request.post('http://localhost:8081/api/buildspec/validate', {
      data: {
        spec: {
          version: '1.0',
          name: 'test',
          tag: 'test:latest',
          base: { image: 'alpine', tag: '3.19' },
          steps: [],
        },
      },
    });
    expect(buildspecResponse.ok()).toBeTruthy();
  });
});

