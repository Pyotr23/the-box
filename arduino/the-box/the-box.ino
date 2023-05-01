#include <EEPROM.h>
#include <DHT.h>

#include "command.h"

#define DHTPIN 2

DHT dht(DHTPIN, DHT11);

const int LED = 13;

const byte SUCCESS = 1;
const byte ERROR = 0;

const int ID_ADDRESS = 0;

const bool IS_DEBUG = false;   

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
  }
  
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
    int id = 0;
    for (int i = 0; i < 30; i++) {
      delay(100);
      if (Serial.available() == 0) {
        continue;
      }
      
      id = Serial.read();
      if (IS_DEBUG) {        
        id -= 48;
      }
      break;
    }
  
    if (id == 0) {
      writeErrorMsg("id not waited");
      return;
    }
   
    EEPROM.write(ID_ADDRESS, id);
    sendSuccess();
  } else if (command == GET_ID) {
    int currentId = EEPROM.read(ID_ADDRESS);
    writeSuccessMsg(currentId);
  } else if (command == UNKNOWN) {
    writeErrorMsg("unknown command");
  }  
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
