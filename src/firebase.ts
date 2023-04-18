import { initializeApp } from 'firebase/app';

initializeApp({
	apiKey: import.meta.env['VITE_FIREBASE_API_KEY'],
	projectId: import.meta.env['VITE_FIREBASE_PROJECT_ID'],
});
