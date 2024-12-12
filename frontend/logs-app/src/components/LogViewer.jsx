import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { API_URL } from '../utils/keys';

const LogViewer = () => {
  const [logs, setLogs] = useState([]);
  const [selectedLog, setSelectedLog] = useState('');
  const [logContent, setLogContent] = useState('');

  useEffect(() => {
    axios.get(`${API_URL}/api/logs`)
      .then((res) => setLogs(res.data))
      .catch((err) => console.error('Error fetching logs', err));
  }, []);

  const handleViewLog = (log) => {
    axios.get(`${API_URL}/api/logs/${log}`)
      .then((res) => setLogContent(res.data))
      .catch((err) => console.error('Error fetching log content', err));
  };

  return (
    <div className="container mt-4">
        <h1 className="mb-3">Log Viewer</h1>
        <ul className="list-group mb-3">
        {logs.map((log) => (
            <li className="list-group-item d-flex justify-content-between align-items-center" key={log}>
            {log}
            <button className="btn btn-secondary btn-sm" onClick={() => handleViewLog(log)}>
                View Log
            </button>
            </li>
        ))}
        </ul>
        <pre className="bg-light p-3 border rounded">{logContent}</pre>
    </div>
  );
};

export default LogViewer;
