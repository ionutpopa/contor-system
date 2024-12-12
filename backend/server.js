const express = require("express");
const fs = require("fs");
const path = require("path");
const cors = require("cors");

const app = express();
app.use(express.json());
app.use(
  cors({
    origin: "http://localhost:3000", // Replace with your app's URL
  })
);

const LOGS_DIR = path.join(__dirname, "../logs");
const CONFIG_FILE = path.join(__dirname, "../config.json");

app.get("/api/logs", (req, res) => {
  fs.readdir(LOGS_DIR, (err, files) => {
    if (err) return res.status(500).send(err.message);
    res.json(files.filter((file) => file.endsWith(".txt")));
  });
});

app.get("/api/logs/:log", (req, res) => {
  const logPath = path.join(LOGS_DIR, req.params.log);
  fs.readFile(logPath, "utf-8", (err, content) => {
    if (err) return res.status(500).send(err.message);
    res.send(content);
  });
});

app.get("/api/config", (req, res) => {
  const configPath = path.join(CONFIG_FILE);
  fs.readFile(configPath, "utf-8", (err, content) => {
    if (err) return res.status(500).send(err.message);
    res.send(JSON.stringify(content));
  });
});

app.post("/api/config", (req, res) => {
  fs.writeFile(CONFIG_FILE, JSON.stringify(req.body, null, 2), (err) => {
    if (err) return res.status(500).send(err.message);
    res.send("Config saved");
  });
});

app.listen(4000, () => console.log("API running on port 4000"));
