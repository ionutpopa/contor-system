# Application that runs power flow and monitors powers (apparent, reactive, active) in nodes (TC - current transformer and VT - voltage transformer)

### Requirements: go installed on user machine

WINDOWS NOTE: to both build and start, if you're using wsl you need to run the following command before running any script:
- wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz
- tar -C /usr/local -xzf go1.23.3.linux-amd64.tar.gz
- export PATH=$PATH:/usr/local/go/bin
- delete the file go1.23.3.linux-amd64.tar.gz that is not in the first level of the application folder

## Start the app
For Windows: Open PowerShell -> wsl -> ./start/start.sh
For MacBook: Open Terminal -> ./start/start.sh

## Build the app
For Windows: Open PowerShell -> wsl -> ./build/build.sh
For MacBook: Open Terminal -> ./build/build.sh
