#include "DHT.h"
#include "command.h"
#define DHTPIN 2

DHT dht(DHTPIN, DHT11);

const int LED = 13;

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
      Serial.write("diode off");
      return;
    case RELAY_ON:
      digitalWrite(LED, HIGH);
      Serial.write("diode one");
      return;
    case TEMPERATURE:
      t = dht.readTemperature();
      if (isnan(t)) {
        Serial.write("read temperature error");
      } else {
        Serial.write(dtostrf(t, 4, 1, buffer));
      }
      return;
    case UNKNOWN:
      Serial.write(CommandToInt(UNKNOWN));
      return;      
  }  
}
