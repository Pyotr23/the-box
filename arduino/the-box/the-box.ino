#include <EEPROM.h>
#include <DHT.h>

#include "command.h"

#define DHTPIN 2


DHT dht(DHTPIN, DHT11);

const bool IS_DEBUG = false;  

const int LED = 13;

const byte SUCCESS = 1;
const byte ERROR = 0;

const int BLINK_COUNT = 3;
const int ONE_BLINK_TIMEOUT_MS = 500;

const int TICK_RATE_MS = 500;

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

  // switch/case doesn't working... i don't know why...
  Command command = IntToCommand(intCommand);
  if (command == RELAY_OFF) {
    digitalWrite(LED, LOW);
    sendSuccess();
  } else if (command == RELAY_ON) {
    digitalWrite(LED, HIGH);
    sendSuccess();
  } else if (command == TEMPERATURE) {    
    float t = dht.readTemperature();
    if (isnan(t)) {
      writeErrorMsg("read temperature error");
    } else {  
      char buffer[5]; 
      dtostrf(t, 4, 1, buffer);    
      writeSuccessMsg(buffer);
    } 
  } else if (command == BLINK) {
    for (int i = 0; i < BLINK_COUNT; i++) {
      digitalWrite(LED, HIGH);
      delay(ONE_BLINK_TIMEOUT_MS);
      digitalWrite(LED, LOW);
      delay(ONE_BLINK_TIMEOUT_MS);
    }
    sendSuccess();
  }

  delay(TICK_RATE_MS);
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
