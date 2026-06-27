#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>

// Configuration
const char* SSID = "";  // TODO: WiFi SSID
const char* PASSWD = "";  // TODO: WiFi Password
const char* LINE_NOTIFY_TOKEN = "";  // TODO: LINE Notify Token
const char* BACKEND_URL = "";  // TODO: Cloud Run Backend URL

// FSK Decoding
const int ADC_PIN = 34;  // ADC input for FSK signal
const int SAMPLE_RATE = 8000;  // 8kHz sampling
const int BUFFER_SIZE = 8000;

void setup() {
  Serial.begin(115200);
  delay(1000);

  Serial.println("\n\nPhone Guard Firmware Starting...");

  // TODO: Initialize ADC
  // TODO: Connect to WiFi
  // TODO: Load whitelist from EEPROM
}

void loop() {
  // TODO: Wait for incoming call (FSK signal detection)
  // TODO: Decode FSK → extract phone number
  // TODO: Check local whitelist
  // TODO: If not in whitelist, send LINE notification
  // TODO: Log to backend (Phase 2)

  delay(100);
}

// FSK Decoder Stub
void decodeFSK() {
  // Reference: GitHub FSK decoder implementations for ESP32
  // 1. Capture ADC samples during FSK period
  // 2. Apply Goertzel algorithm or FFT
  // 3. Extract mark (1) and space (0) frequencies
  // 4. Decode DTMF-like sequence → phone number

  Serial.println("FSK Decoding stub - implement from GitHub reference");
}

// LINE Notification
void sendLineNotification(const char* phoneNumber) {
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("WiFi not connected");
    return;
  }

  HTTPClient http;
  String url = String(BACKEND_URL) + "/api/notify";

  http.begin(url);
  http.addHeader("Content-Type", "application/json");
  http.addHeader("Authorization", "Bearer " + String(LINE_NOTIFY_TOKEN));

  String payload = "{\"number\":\"" + String(phoneNumber) + "\"}";

  int httpCode = http.POST(payload);
  Serial.printf("LINE Notify Response: %d\n", httpCode);

  http.end();
}
