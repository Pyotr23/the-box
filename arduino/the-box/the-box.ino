#include <DHT.h>

#include "command.h"

#define DHTPIN 2

const int LED = 13;
const int MAX_PIN = LED;
const int MIN_PIN = 0;
const int BLUETOOTH_RX = 0;
const int BLUETOOTH_TX = 1;
const int DHT_PIN = 2;

int busyPins[] = {BLUETOOTH_RX, BLUETOOTH_TX, DHT_PIN};

DHT dht(DHT_PIN, DHT11);

const bool IS_DEBUG = false;  

const byte SUCCESS = 1;
const byte ERROR = 0;

const int TICK_RATE_MS = 500;

const int WAITING_TIMEOUT_MS = 4500;

int usedPin;
void setup() {
  Serial.begin(9600);

  for (int i = 0; i <= MAX_PIN; i++) {   
    bool isUsed = false;
    for (int j = 0; j < (sizeof(busyPins) / sizeof(busyPins[0])); j++) {   
      if (busyPins[j] == i) {  
        isUsed = true;    
        break;
      }
    }  
    if (!isUsed) {
      pinMode(i, OUTPUT);
    }    
  }  
    
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
      waitPinAndSetLevel(LOW);
      break;
    case PIN_ON:
      waitPinAndSetLevel(HIGH);
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
    default:
      writeErrorMsg("unknown command");
  }
  
  delay(TICK_RATE_MS);
}

void waitPinAndSetLevel(int value) {
  int p = waitNumber();
  if (p == -1) {
    writeErrorMsg("pin not waited");
    return;
  }  
  if (isAvailablePin(p)) {
    digitalWrite(p, value);  
    sendSuccess();   
  } else {
    writeErrorMsg("pin is busy");
  }
}

bool isAvailablePin(int pin) {  
  if (sizeof(busyPins) == 0) {
    return true;
  } 
  if (pin > MAX_PIN || pin < MIN_PIN) {
    return false; 
  }
  for (int i = 0; i < (sizeof(busyPins) / sizeof(busyPins[0])); i++) {   
    if (busyPins[i] == pin) {      
      return false;
    }
  }  
  return true;
}

int waitNumber() {  
  Serial.setTimeout(WAITING_TIMEOUT_MS);
  unsigned long startTime = millis();
  int res = Serial.parseInt();
  if (millis() - startTime > WAITING_TIMEOUT_MS) {
    return -1;
  }
  Serial.read();
  return res;
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
  Serial.print(msg);
}

void writeErrorMsg(int num) {
  if (IS_DEBUG) {
    Serial.println(ERROR);
    Serial.println(num);
    return;
  }
  
  Serial.write(ERROR);
  Serial.print(num);
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
