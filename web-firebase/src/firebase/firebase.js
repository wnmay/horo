import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
  apiKey: "AIzaSyCy5xdjvetLg8b8oHVfUgjRdpfMM9StkJM",
  authDomain: "food-ordering-48dd9.firebaseapp.com",
  projectId: "food-ordering-48dd9",
  storageBucket: "food-ordering-48dd9.firebasestorage.app",
  messagingSenderId: "839817786451",
  appId: "1:839817786451:web:8bb9401a21a7dc3b73c2a6",
  measurementId: "G-36EMSH6DZQ"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const auth = getAuth(app)

export { app, auth };
