// Phone Guard ESP32 Configuration
// Edit these values or use environment variables from .env

#ifndef CONFIG_H
#define CONFIG_H

// WiFi Configuration
const char* WIFI_SSID = "";        // TODO: Set your WiFi SSID
const char* WIFI_PASSWORD = "";    // TODO: Set your WiFi password

// LINE Notify Configuration
const char* LINE_NOTIFY_TOKEN = "";  // TODO: Set your LINE Notify token

// Backend Configuration (for Phase 2)
const char* BACKEND_URL = "http://localhost:8080";
const char* DEVICE_ID = "device_001";

// Hardware Configuration
const int ADC_PIN = 34;             // ADC pin for FSK signal input
const int LED_PIN = 2;              // Status LED pin

// FSK Parameters
const int SAMPLE_RATE = 8000;       // 8kHz sampling rate
const float FREQ_MARK = 1300.0;     // Mark frequency (1)
const float FREQ_SPACE = 2100.0;    // Space frequency (0)

#endif
