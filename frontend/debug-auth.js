// Auth Debug Script - Run in Browser Console

console.log('=== Auth Debug ===');

// 1. Check localStorage
const accessToken = localStorage.getItem('access_token');
const refreshToken = localStorage.getItem('refresh_token');
console.log('1. localStorage state:');
console.log('  - access_token:', accessToken ? '✓ exists' : '✗ missing');
console.log('  - refresh_token:', refreshToken ? '✓ exists' : '✗ missing');

// 2. Parse token
if (accessToken) {
  try {
    const parts = accessToken.split('.');
    if (parts.length === 3) {
      const payload = JSON.parse(atob(parts[1]));
      console.log('2. Token payload:', payload);
      console.log('  - User:', payload.user_id || payload.sub);
      console.log('  - Expires:', new Date(payload.exp * 1000).toLocaleString());
      console.log('  - Expired:', payload.exp < Math.floor(Date.now() / 1000));
    } else {
      console.log('2. ✗ Invalid token format');
    }
  } catch (e) {
    console.log('2. ✗ Failed to parse token:', e);
  }
} else {
  console.log('2. ✗ No token to parse');
}

// 3. Check Zustand store
const authStore = window.__ZUSTAND_STORES__?.find(s => s?.state?.user !== undefined);
if (authStore) {
  console.log('3. Zustand auth store:', authStore.getState());
} else {
  console.log('3. ✗ No auth store found');
}

// 4. Test API call
console.log('4. Testing API call...');
fetch('/api/v1/projects', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': accessToken ? `Bearer ${accessToken}` : ''
  },
  body: JSON.stringify({ name: 'DebugTest', prefix: 'DBG', description: 'Debug' })
})
.then(r => {
  console.log('  - Status:', r.status);
  return r.text();
})
.then(text => {
  console.log('  - Response:', text);
})
.catch(e => {
  console.log('  - Error:', e);
});

console.log('=== End Debug ===');
