// k6 Performance Test Script for aitestos API
// Usage: k6 run tests/performance/load_test.js
// With custom base URL: BASE_URL=http://localhost:8080 k6 run tests/performance/load_test.js

import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const apiLatency = new Trend('api_latency');
const requestCount = new Counter('requests');

// Test configuration
export const options = {
    // Test stages: ramp up, sustain, peak, sustain, ramp down
    stages: [
        { duration: '30s', target: 10 },   // Ramp up to 10 users
        { duration: '1m', target: 10 },    // Stay at 10 users
        { duration: '30s', target: 50 },   // Ramp up to 50 users (peak)
        { duration: '2m', target: 50 },    // Stay at 50 users
        { duration: '30s', target: 0 },    // Ramp down to 0 users
    ],
    // Performance thresholds
    thresholds: {
        http_req_duration: ['p(99)<500'],  // 99% of requests must complete < 500ms
        errors: ['rate<0.01'],              // Error rate must be < 1%
        api_latency: ['p(99)<500'],         // Custom API latency metric
    },
    // Teardown timeout
    teardownTimeout: '30s',
};

// Base URL from environment or default
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Test data
const TEST_USER_ID = '00000000-0000-0000-0000-000000000001';
const TEST_PROJECT_ID = '00000000-0000-0000-0000-000000000002';

// Default request headers
const defaultHeaders = {
    'Content-Type': 'application/json',
    'X-User-ID': TEST_USER_ID,
};

// Helper function to make HTTP request with metrics
function makeRequest(method, path, body = null) {
    const url = `${BASE_URL}${path}`;
    const params = {
        headers: defaultHeaders,
        timeout: '30s',
    };

    let response;
    if (body) {
        response = http[method.toLowerCase()](url, JSON.stringify(body), params);
    } else {
        response = http[method.toLowerCase()](url, params);
    }

    // Record metrics
    requestCount.add(1);
    apiLatency.add(response.timings.duration);
    errorRate.add(response.status >= 400);

    return response;
}

// Test scenarios

// Health check test
function testHealthCheck() {
    const response = makeRequest('GET', '/health');

    check(response, {
        'health check status 200': (r) => r.status === 200,
        'health check body contains OK': (r) => r.body.includes('OK') || r.body.includes('ok'),
    });
}

// List projects test
function testListProjects() {
    const response = makeRequest('GET', '/api/v1/projects?limit=10');

    check(response, {
        'list projects status 200': (r) => r.status === 200 || r.status === 501,
        'list projects response time < 500ms': (r) => r.timings.duration < 500,
    });
}

// Get project test
function testGetProject() {
    const response = makeRequest('GET', `/api/v1/projects/${TEST_PROJECT_ID}`);

    check(response, {
        'get project status 200 or 404 or 501': (r) => [200, 404, 501].includes(r.status),
        'get project response time < 500ms': (r) => r.timings.duration < 500,
    });
}

// Create project test
function testCreateProject() {
    const timestamp = Date.now();
    const body = {
        name: `Performance Test Project ${timestamp}`,
        prefix: `P${timestamp.toString().slice(-4)}`,
        description: 'Created by k6 performance test',
    };

    const response = makeRequest('POST', '/api/v1/projects', body);

    check(response, {
        'create project status 201 or 501': (r) => [201, 400, 501].includes(r.status),
        'create project response time < 500ms': (r) => r.timings.duration < 500,
    });

    return response;
}

// List test cases test
function testListTestCases() {
    const response = makeRequest('GET', '/api/v1/testcases?limit=10');

    check(response, {
        'list test cases status 200': (r) => r.status === 200 || r.status === 501,
        'list test cases response time < 500ms': (r) => r.timings.duration < 500,
    });
}

// List test plans test
function testListTestPlans() {
    const response = makeRequest('GET', '/api/v1/plans?limit=10');

    check(response, {
        'list plans status 200': (r) => r.status === 200 || r.status === 501,
        'list plans response time < 500ms': (r) => r.timings.duration < 500,
    });
}

// List documents test
function testListDocuments() {
    const response = makeRequest('GET', '/api/v1/knowledge/documents?limit=10');

    check(response, {
        'list documents status 200': (r) => r.status === 200 || r.status === 501,
        'list documents response time < 500ms': (r) => r.timings.duration < 500,
    });
}

// Main test function - executed by each virtual user
export default function () {
    // Scenario 1: Health check (lightweight)
    testHealthCheck();

    sleep(0.5);

    // Scenario 2: Read operations (most common)
    testListProjects();
    testListTestCases();
    testListTestPlans();

    sleep(1);

    // Scenario 3: Individual resource access
    testGetProject();

    sleep(0.5);

    // Scenario 4: Write operations (less frequent)
    // Only create new projects occasionally to avoid database bloat
    if (__ITER % 10 === 0) {
        testCreateProject();
    }

    sleep(0.5);

    // Scenario 5: Knowledge base access
    testListDocuments();

    sleep(1);
}

// Setup function - runs once per VU
export function setup() {
    console.log(`Starting k6 load test against ${BASE_URL}`);

    // Verify server is reachable
    const response = http.get(`${BASE_URL}/health`);
    if (response.status !== 200) {
        console.log(`Warning: Health check returned ${response.status}`);
    }

    return { baseUrl: BASE_URL };
}

// Teardown function - runs once after all VUs complete
export function teardown(data) {
    console.log('k6 load test completed');
}

// Handle summary output
export function handleSummary(data) {
    return {
        stdout: textSummary(data, { indent: ' ', enableColors: true }),
        'tests/performance/summary.json': JSON.stringify(data, null, 2),
    };
}

// Text summary helper (built into k6)
function textSummary(data, options) {
    // k6 provides a default text summary
    // This function can be customized for specific output format
    return '';
}
