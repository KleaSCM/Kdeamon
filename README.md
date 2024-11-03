# Kdeamon - Cute App Launcher Deamon 

`Kdeamon` is a minimalist, daemonized application launcher for Linux built with the Fyne GUI toolkit. It provides an intuitive search bar with fuzzy matching, allowing users to quickly find and launch applications installed on their system. `Kdeamon` runs as a background process (daemon) on system startup.

## Features

- **Fuzzy Search**: Type a few letters to find any installed application.
- **Autocomplete Suggestions**: Inline autocomplete feature helps you find the closest match.
- **Daemonized Operation**: Automatically starts in the background on system boot.

## Getting Started

### Prerequisites

- **Go**: Ensure Go is installed (version 1.16 or higher).
- **Fyne GUI Toolkit**: The launcher uses [Fyne](https://fyne.io/) for its graphical interface.
- **Systemd**: Needed to run `Kdeamon` as a daemon on Linux.

### Installing Dependencies

1. Clone this repository:
   ```git clone <repository-url>```
   ```cd Kdeamon```
Install required packages:
```go mod tidy```

install fyne
```go get fyne.io/fyne/v2```


## Building the Project

```go build -o Kdeamon```


## Setup Deamon

Move Kdeamon executable

```sudo mv Kdeamon /usr/local/bin/```


## Create a systemd service file:

``sudo vim /etc/systemd/system/kdeamon.service``


## Add the service file or copy

[Unit]
Description=Kdeamon Application Launcher
After=network.target

[Service]
ExecStart=/usr/local/bin/Kdeamon
Environment=DISPLAY=:0
Environment=XAUTHORITY=/home/ADDUSERNAMEHERE/.Xauthority
Restart=on-failure
Type=oneshot
RemainAfterExit=true

[Install]
WantedBy=default.target


## Enable and start the service:

```sudo systemctl daemon-reload```
```sudo systemctl enable kdeamon```
```sudo systemctl start kdeamon```

then check if its running 

```sudo systemctl status kdeamon```


# Usage

Open Kdeamon:
It will automatically open on system boot.

## Search Applications:

  Type into the search bar to find applications.
  Select an application by typing its name or using the autocomplete suggestion.
        
  ## Launch Application:
  Press Enter to launch the selected application.

## Troubleshooting

## Daemon Not Starting: 
  Ensure that DISPLAY and XAUTHORITY are correctly set in the service file. The XAUTHORITY path may vary based on your setup.
  Autocomplete Not Working: Check that applications are correctly loaded from .desktop files located in /usr/share/applications, /usr/local/share/applications, and ~/.local/share/applications.

