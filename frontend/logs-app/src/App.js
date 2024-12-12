import React from 'react';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import ConfigEditor from './components/ConfigEditor';
import LogViewer from './components/LogViewer';

const App = () => {
  return (
    <Router>
      <nav>
        <Link to="/">Edit Config</Link> | <Link to="/logs">View Logs</Link>
      </nav>
      <Routes>
        <Route path="/" element={<ConfigEditor />} />
        <Route path="/logs" element={<LogViewer />} />
      </Routes>
    </Router>
  );
};

export default App;
