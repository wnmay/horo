import http from 'k6/http';

// Firebase configuration
const FIREBASE_API_KEY = __ENV.FIREBASE_API_KEY;
const EMAIL = __ENV.FIREBASE_EMAIL;
const PASSWORD = __ENV.FIREBASE_PASSWORD;

export function getFirebaseToken() {
  console.log('🔐 Authenticating with Firebase...');
  
  const url = `https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=${FIREBASE_API_KEY}`;
  
  const payload = JSON.stringify({
    email: EMAIL,
    password: PASSWORD,
    returnSecureToken: true,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const response = http.post(url, payload, params);

  if (response.status !== 200) {
    console.error('❌ Firebase authentication failed:', response.body);
    throw new Error(`Firebase auth failed with status ${response.status}`);
  }

  const data = JSON.parse(response.body);
  
  console.log('✅ Successfully authenticated!');
  console.log(`   User ID: ${data.localId}`);
  console.log(`   Email: ${data.email}`);
  console.log(`   Token expires in: ${data.expiresIn} seconds (~${Math.floor(data.expiresIn / 60)} minutes)`);
  
  return {
    idToken: data.idToken,
    refreshToken: data.refreshToken,
    expiresIn: data.expiresIn,
    localId: data.localId,
    email: data.email,
  };
}

export function refreshFirebaseToken(refreshToken) {
  console.log('🔄 Refreshing Firebase token...');
  
  const url = `https://securetoken.googleapis.com/v1/token?key=${FIREBASE_API_KEY}`;
  
  const payload = JSON.stringify({
    grant_type: 'refresh_token',
    refresh_token: refreshToken,
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const response = http.post(url, payload, params);

  if (response.status !== 200) {
    console.error('❌ Token refresh failed:', response.body);
    throw new Error(`Token refresh failed with status ${response.status}`);
  }

  const data = JSON.parse(response.body);
  
  console.log('✅ Token refreshed successfully!');
  console.log(`   New token expires in: ${data.expires_in} seconds (~${Math.floor(data.expires_in / 60)} minutes)`);
  
  return {
    idToken: data.id_token,
    refreshToken: data.refresh_token,
    expiresIn: data.expires_in,
  };
}
