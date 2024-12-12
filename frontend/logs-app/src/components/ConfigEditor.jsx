import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { API_URL } from '../utils/keys';

const ConfigEditor = () => {
  const [config, setConfig] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    axios.get(`${API_URL}/api/config`)
        .then((res) => setConfig(res.data))
        .catch((err) => console.error('Error fetching logs', err));
  }, [])

  const handleSave = () => {
    try {
      const parsedConfig = JSON.parse(config);
      setError('');
      // Save config via API
      axios.post(`${API_URL}/api/config`, parsedConfig)
        .then(() => alert('Config saved successfully!'))
        .catch(err => console.error('Error saving config', err));
    } catch (e) {
      setError('Invalid JSON format');
    }
  };

  return (
    <div className="container mt-4">
        <h1 className="mb-3">Edit Configuration</h1>
        <textarea
            className="form-control mb-3"
            rows="10"
            value={config}
            onChange={(e) => setConfig(e.target.value)}
        />
        {error && <p className="text-danger">{error}</p>}
        <button className="btn btn-primary" onClick={handleSave}>
            Save Config
        </button>
    </div>
  );
};

export default ConfigEditor;
