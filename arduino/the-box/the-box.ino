#include <EEPROM.h>
#include <DHT.h>

#include "command.h"

#define DHTPIN 2


DHT dht(DHTPIN, DHT11);

const bool IS_DEBUG = false;  

const int LED = 13;

const byte SUCCESS = 1;
const byte ERROR = 0;

const int ID_ADDRESS = 0;
const int LOWER_TEMPERATURE_THRESHOLD_ADDRESS = 1;
const int HIGHER_TEMPERATURE_THRESHOLD_ADDRESS = 2;

const int WAITING_TIMEOUT_MS = 5000;
const int WAITING_SLEEP_TIMEOUT_MS = 100;

int waitingCount;

byte id, lowerTemperatureThreshold, higherTemperatureThreshold;  

void setup() {
  Serial.begin(9600);
  
  pinMode(LED, OUTPUT);
  
  digitalWrite(LED, HIGH);
  
  dht.begin();

  waitingCount = WAITING_TIMEOUT_MS / WAITING_SLEEP_TIMEOUT_MS;

  lowerTemperatureThreshold = EEPROM.read(LOWER_TEMPERATURE_THRESHOLD_ADDRESS);
  higherTemperatureThreshold = EEPROM.read(HIGHER_TEMPERATURE_THRESHOLD_ADDRESS);
  id = EEPROM.read(ID_ADDRESS);
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
  } else if (command == SET_ID) {    
    id = waitNumber();  
    if (id == 0) {
      writeErrorMsg("id not waited");
      return;
    }   
    EEPROM.write(ID_ADDRESS, id);    
    sendSuccess();
  } else if (command == GET_ID) {    
    writeSuccessMsg(id);  
  } else if (command == GET_LOWER_TEMPERATURE_THRESHOLD) {
    writeSuccessMsg(lowerTemperatureThreshold);
  } else if (command == GET_HIGHER_TEMPERATURE_THRESHOLD) {
    writeSuccessMsg(higherTemperatureThreshold);
  } else if (command == UNKNOWN) {
    writeErrorMsg("unknown command");
  } else if (command == SET_LOWER_TEMPERATURE_THRESHOLD) {    
    lowerTemperatureThreshold = waitNumber();  
    if (lowerTemperatureThreshold == 0) {
      writeErrorMsg("lower temperature threshold not waited");
      return;
    }   
    EEPROM.write(LOWER_TEMPERATURE_THRESHOLD_ADDRESS, lowerTemperatureThreshold);    
    sendSuccess();
  } else if (command == SET_HIGHER_TEMPERATURE_THRESHOLD) {    
    higherTemperatureThreshold = waitNumber();  
    if (higherTemperatureThreshold == 0) {
      writeErrorMsg("higher temperature threshold not waited");
      return;
    }   
    EEPROM.write(HIGHER_TEMPERATURE_THRESHOLD_ADDRESS, higherTemperatureThreshold);    
    sendSuccess();
  }
}

byte waitNumber() {  
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
  return 0;
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
