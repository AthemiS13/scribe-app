# Scribe App

## Overview

The Scribe App is the official companion application for the [Scribe Hardware](../scribe/README.md). It allows users to seamlessly send text, notes, and data to their Scribe device via Bluetooth Low Energy (BLE).

**Hardware:** [Scribe Device](../scribe/README.md)

## Features

*   **Cross-Platform:** Built with [Wails](https://wails.io/), running on macOS, Windows, and Linux.
*   **Bluetooth Connectivity:** Automatically detects and connects to the Scribe device.
*   **Text Management:** Type, format, and send text strings directly to the Scribe OLED display.
*   **User Interface:** Clean and intuitive React-based frontend.

## Application

The primary function of the Scribe App is to assist with memory and organization. Users can send text strings via the app, which the Scribe firmware processes and organizes into pages. This helps users keep track of important notes, dates, and names, enhancing their ability to recall information easily.

## Development

This project is built using **Go** (backend) and **React** (frontend) via the **Wails** framework.

### Prerequisites
*   Go 1.18+
*   Node.js & npm

### Setup

1.  Clone the repository.
2.  Navigate to the `scribe-app` directory.
3.  Install frontend dependencies:
    ```bash
    cd frontend
    npm install
    ```

### Running in Development Mode

To run the app in live development mode:

```bash
wails dev
```

This will start a Vite development server for the frontend and a Go backend server.

### Building

To build the application for production:

```bash
wails build
```

## Credits

*   **App Backend & Architecture:** [George](https://github.com/freddycz)
*   **Hardware & Concept:** [AthemiS13](https://github.com/AthemiS13)
