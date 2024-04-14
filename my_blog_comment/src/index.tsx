import React from 'react';
import App from './App';
import { createRoot } from 'react-dom/client';
import './Reset.css'


const container = document.getElementById('comment');
const root = createRoot(container!); // createRoot(container!) if you use TypeScript
root.render(<App />);
