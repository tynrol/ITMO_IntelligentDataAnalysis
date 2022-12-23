import * as React from 'react';
import { Routes, Route } from 'react-router-dom';

import Access from "./Access";
import Detect from "./Detect";

export default function App() {
    return (
        <div className="App">
            <Routes>
                <Route path="/" element={<Detect />} />
                <Route path="access" element={<Access />} />
            </Routes>
        </div>
    );
}

