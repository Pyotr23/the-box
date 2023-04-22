#include "DHT.h"
#include "command.h"
#define DHTPIN 2

DHT dht(DHTPIN, DHT11);

const int LED = 13;

const byte SUCCESS = 1;
const byte ERROR = 0;

Command command;
float t;   
char buffer[5];    

void setup() {
  Serial.begin(9600);
  pinMode(LED, OUTPUT);
  digitalWrite(LED, HIGH);
  dht.begin();
}

void loop() {
  if (!Serial.available()) {
    return;
  }
  
  switch (IntToCommand(Serial.read())) {
    case RELAY_OFF: 
      digitalWrite(LED, LOW);
      sendSuccess();
      return;
    case RELAY_ON:
      digitalWrite(LED, HIGH);
      sendSuccess();
      return;
    case TEMPERATURE:
      t = dht.readTemperature();
      if (isnan(t)) {
        writeErrorMsg("read temperature error");
      } else {
       writeSuccessMsg(dtostrf(t, 4, 1, buffer));
      }
      return;
    case UNKNOWN:
      writeErrorMsg("unknown command");
      return;      
  }  
}

void sendSuccess() {
  Serial.write(SUCCESS);
}

void writeSuccessMsg(char msg[]) {
   sendSuccess();
   Serial.write(msg);
}

void writeErrorMsg(char msg[]) {
   Serial.write(ERROR);
   Serial.write(msg);
}
