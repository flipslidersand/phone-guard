#include <Arduino.h>
#include <WiFi.h>
#include <HTTPClient.h>
#include <EEPROM.h>
#include <math.h>

// ===== Configuration =====
const char* SSID = "";              // WiFi SSID (edit here)
const char* PASSWD = "";            // WiFi Password (edit here)
const char* LINE_NOTIFY_TOKEN = ""; // LINE Notify Token (edit or use env)
const char* DEVICE_ID = "device_1"; // Device identifier (for backend)

// ===== Hardware Config =====
const int ADC_PIN = 34;             // ADC input for FSK signal
const int LED_PIN = 2;              // Status LED (optional)

// ===== FSK Parameters =====
const int SAMPLE_RATE = 8000;       // 8kHz sampling
const int SAMPLES_PER_BIT = 8000 / 1200; // 1200 baud
const float FREQ_MARK = 1300.0;    // Mark frequency
const float FREQ_SPACE = 2100.0;   // Space frequency
const int FSK_DETECTION_SAMPLES = SAMPLE_RATE; // 1 second buffer

// ===== EEPROM Config =====
const int EEPROM_SIZE = 4096;
const int WHITELIST_START_ADDR = 0;
const int MAX_WHITELIST_ENTRIES = 50;
struct WhitelistEntry {
  char number[15];
  char label[30];
  uint32_t timestamp;
};

// ===== Global State =====
volatile bool callDetected = false;
unsigned long lastCallTime = 0;
const unsigned long DEBOUNCE_DELAY = 60000; // 60s between notifications

void setup() {
  Serial.begin(115200);
  delay(1000);

  Serial.println("\n\n=== Phone Guard Firmware v1.0 ===");

  // Initialize LED
  pinMode(LED_PIN, OUTPUT);
  digitalWrite(LED_PIN, LOW);

  // Initialize EEPROM
  EEPROM.begin(EEPROM_SIZE);
  Serial.println("[EEPROM] Initialized");

  // Initialize ADC
  pinMode(ADC_PIN, INPUT);
  analogReadResolution(12); // 12-bit ADC (0-4095)
  Serial.println("[ADC] Initialized on pin " + String(ADC_PIN));

  // Load whitelist from EEPROM
  loadWhitelist();

  // Connect to WiFi
  connectWiFi();

  Serial.println("[SETUP] Complete");
  digitalWrite(LED_PIN, HIGH);
}

void loop() {
  // Check WiFi connection
  if (WiFi.status() != WL_CONNECTED) {
    connectWiFi();
  }

  // Listen for FSK signal (simplified detection)
  // In production: use interrupt-based detection
  int adcValue = analogRead(ADC_PIN);

  // Simple threshold detection (FSK signal when ADC > 2048)
  if (adcValue > 2048 && !callDetected) {
    delay(500); // Debounce
    if (analogRead(ADC_PIN) > 2048) {
      callDetected = true;
      digitalWrite(LED_PIN, HIGH);
      Serial.println("[CALL] Incoming call detected!");
      handleIncomingCall();
    }
  } else if (adcValue < 2048 && callDetected) {
    callDetected = false;
    digitalWrite(LED_PIN, LOW);
  }

  delay(100);
}

// ===== WiFi Management =====
void connectWiFi() {
  if (WiFi.status() == WL_CONNECTED) {
    return;
  }

  Serial.printf("[WiFi] Connecting to %s...\n", SSID);
  WiFi.mode(WIFI_STA);
  WiFi.begin(SSID, PASSWD);

  int attempts = 0;
  while (WiFi.status() != WL_CONNECTED && attempts < 20) {
    delay(500);
    Serial.print(".");
    attempts++;
  }

  if (WiFi.status() == WL_CONNECTED) {
    Serial.printf("\n[WiFi] Connected! IP: %s\n", WiFi.localIP().toString().c_str());
  } else {
    Serial.println("\n[WiFi] Failed to connect");
  }
}

// ===== FSK Decoding =====
// Reference: Goertzel algorithm for frequency detection
// For production: use existing Arduino FSK library or implement full Goertzel
String decodeFSK() {
  // Capture ADC samples
  int samples[FSK_DETECTION_SAMPLES];
  for (int i = 0; i < FSK_DETECTION_SAMPLES; i++) {
    samples[i] = analogRead(ADC_PIN);
    delayMicroseconds(125); // 8kHz sampling
  }

  // TODO: Apply Goertzel algorithm to detect 1300Hz (mark) and 2100Hz (space)
  // Goertzel pseudocode:
  // 1. For each frequency (MARK, SPACE):
  //    - Calculate coefficient: 2 * cos(2π * freq / SAMPLE_RATE)
  //    - Process samples through Goertzel filter
  //    - Extract magnitude
  // 2. Determine bit value: MARK > SPACE ? '1' : '0'
  // 3. Accumulate bits → extract phone number

  // Simplified stub: return dummy number for testing
  Serial.println("[FSK] Decoding stub - implement Goertzel algorithm");
  return "09012345678"; // Placeholder
}

// ===== Whitelist Management =====
void loadWhitelist() {
  Serial.println("[WHITELIST] Loading from EEPROM...");
  // Read whitelist entries from EEPROM
  // Format: [entry0][entry1]...[entryN][0xFF=end]
  // For now: stub
  Serial.println("[WHITELIST] Loaded 0 entries");
}

bool isWhitelisted(const char* phoneNumber) {
  // Check if number exists in EEPROM whitelist
  // For now: always return false (all numbers trigger notification)
  return false;
}

void addToWhitelist(const char* phoneNumber, const char* label) {
  Serial.printf("[WHITELIST] Adding: %s (%s)\n", phoneNumber, label);
  // Write to EEPROM
  // TODO: Implement EEPROM write logic
}

// ===== Incoming Call Handler =====
void handleIncomingCall() {
  // Wait for FSK period (~500ms after first ring)
  delay(1500);

  // Decode FSK signal
  String phoneNumber = decodeFSK();
  Serial.printf("[CALL] Extracted number: %s\n", phoneNumber.c_str());

  // Check whitelist
  if (isWhitelisted(phoneNumber.c_str())) {
    Serial.println("[CALL] Number is whitelisted - no notification sent");
    return;
  }

  // Debounce: only notify once per minute
  unsigned long now = millis();
  if (now - lastCallTime < DEBOUNCE_DELAY) {
    Serial.println("[CALL] Debounce: notification already sent recently");
    return;
  }
  lastCallTime = now;

  // Send LINE notification
  sendLineNotification(phoneNumber.c_str());
}

// ===== LINE Notification =====
void sendLineNotification(const char* phoneNumber) {
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("[LINE] WiFi not connected - skipping notification");
    return;
  }

  HTTPClient http;

  // Option 1: Send to backend API (Phase 2)
  // String url = String(BACKEND_URL) + "/api/notify";

  // Option 2: Send directly to LINE Notify (Phase 1)
  String url = "https://notify-api.line.me/api/notify";

  http.begin(url);
  http.addHeader("Content-Type", "application/x-www-form-urlencoded");
  http.addHeader("Authorization", "Bearer " + String(LINE_NOTIFY_TOKEN));

  String payload = "message=%E7%9F%A5%E3%82%89%E3%81%AA%E3%81%84%E7%95%AA%E5%8F%B7%E3%81%8B%E3%82%89%E7%9D%80%E4%BF%A1%0A";
  payload += "%E7%95%AA%E5%8F%B7%3A%20" + String(phoneNumber);

  int httpCode = http.POST(payload);
  Serial.printf("[LINE] Response code: %d\n", httpCode);

  if (httpCode == HTTP_CODE_OK || httpCode == 200) {
    Serial.println("[LINE] Notification sent successfully");
  } else {
    Serial.printf("[LINE] Error: %d\n", httpCode);
  }

  http.end();
}
