// Import the functions you need from the SDKs you need
import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";
// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
  apiKey: "AIzaSyBnN693MGJZ2p1yRICP6IzoOgEKj3ZBsKI",
  authDomain: "horo-d47b1.firebaseapp.com",
  projectId: "horo-d47b1",
  storageBucket: "horo-d47b1.firebasestorage.app",
  messagingSenderId: "796546060161",
  appId: "1:796546060161:web:c244aafd8ebcd96c2ce280",
  measurementId: "G-S14W0PE44V"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);
const auth = getAuth(app);

export {app, auth}