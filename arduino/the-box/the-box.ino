#include <DHT.h>

#include "command.h"

#define DHTPIN 2

const int LED = 13;
const int MAX_PIN = LED;
const int BLUETOOTH_RX = 0;
const int BLUETOOTH_TX = 1;
const int DHT_PIN = 2;

int busyPins[] = {BLUETOOTH_RX, BLUETOOTH_TX, DHT_PIN};

DHT dht(DHT_PIN, DHT11);

const bool IS_DEBUG = false;  

const byte SUCCESS = 1;
const byte ERROR = 0;

const int BLINK_COUNT = 3;
const int ONE_BLINK_TIMEOUT_MS = 500;

const int TICK_RATE_MS = 500;

const int WAITING_TIMEOUT_MS = 5000;
const int WAITING_SLEEP_TIMEOUT_MS = 100;
int waitingCount = WAITING_TIMEOUT_MS / WAITING_SLEEP_TIMEOUT_MS;

int usedPin;
void setup() {
  Serial.begin(9600);
  
  pinMode(LED, OUTPUT);
  
  digitalWrite(LED, HIGH);
  
  dht.begin();
}

void loop() {
  if (Serial.available() == 0) {
    return;
  }

  int intCommand = Serial.read();
  if (IS_DEBUG) {
    intCommand -= 48;
    Serial.println(intCommand);
  }

  switch (IntToCommand(intCommand)) {
    case PIN_OFF:
      usedPin = waitNumber();  
      if (usedPin == -1) {
        writeErrorMsg("pin not waited");
        return;
      }  
      digitalWrite(usedPin, LOW);
      sendSuccess();
      break;
    case PIN_ON:
      usedPin = waitNumber();  
      if (usedPin == -1) {
        writeErrorMsg("pin not waited");
        return;
      }  
      digitalWrite(usedPin, HIGH);
      sendSuccess();
      break;
    case CHECK_PIN:
      usedPin = waitNumber();    
      if (usedPin == -1) {
        writeErrorMsg("pin not waited");
        return;
      }  
      if (isAvailablePin(usedPin)) {
        sendSuccess();     
      } else {
        writeErrorMsg("pin is busy");
      }
      break;
    case TEMPERATURE:
      float t = dht.readTemperature();
      if (isnan(t)) {
        writeErrorMsg("read temperature error");
      } else {  
        char buffer[5]; 
        dtostrf(t, 4, 1, buffer);    
        writeSuccessMsg(buffer);
      } 
      break;
    case BLINK:
      for (int i = 0; i < BLINK_COUNT; i++) {
        digitalWrite(LED, HIGH);
        delay(ONE_BLINK_TIMEOUT_MS);
        digitalWrite(LED, LOW);
        delay(ONE_BLINK_TIMEOUT_MS);
      }
      sendSuccess();
      break;
    default:
      writeErrorMsg("unknown command");
  }
  
  delay(TICK_RATE_MS);
}

bool isAvailablePin(int pin) {  
  if (sizeof(busyPins) == 0) {
    return true;
  } 
  for (int i = 0; i < (sizeof(busyPins) / sizeof(busyPins[0])); i++) {   
    if (busyPins[i] == pin) {      
      return false;
    }
  }  
  return true;
}

int waitNumber() {  
  for (int i = 0; i < waitingCount; i++) {
      delay(WAITING_SLEEP_TIMEOUT_MS);
      if (Serial.available() == 0) {
        continue;
      }
      
      byte n = Serial.read();
      if (IS_DEBUG) {        
        n -= 48;
      }
      return n;
    }
  return -1;
}

void sendSuccess() {
  if (IS_DEBUG) {
    Serial.println(SUCCESS);
    return;
  }
  
  Serial.write(SUCCESS);
}

void writeSuccessMsg(int num) {
  if (IS_DEBUG) {
    sendSuccess();
    Serial.println(num);
    return;
  }
  
  sendSuccess();
  Serial.print(num);
}

void writeSuccessMsg(char msg[]) {
  if (IS_DEBUG) {
    sendSuccess();
    Serial.println(msg);
    return;
  }
  
  sendSuccess();
  Serial.print(msg);
}

void writeErrorMsg(char msg[]) {
  if (IS_DEBUG) {
    Serial.println(ERROR);
    Serial.println(msg);
    return;
  }
  
  Serial.write(ERROR);
  Serial.write(msg);
}

char* toChars(float num) {
  char buffer[5]; 
  dtostrf(num, 4, 1, buffer);
  return buffer;
}

char* toChars(int num) {
  char cstr[16];
  itoa(num, cstr, 10);
  return cstr;
}
