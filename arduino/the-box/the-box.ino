#include "DHT.h"
#define DHTPIN 2

DHT dht(DHTPIN, DHT11);

const int LED = 13;

int val; 
float t;   
char buffer[5];    

void setup() {
  Serial.begin(9600);
  pinMode(LED, OUTPUT);
  digitalWrite(LED, HIGH);
  dht.begin();
}

void loop() {
  if (Serial.available()) {
    val = Serial.read();
    if (val == '2') {
      t = dht.readTemperature();
      if (isnan(t)) {
        Serial.write("read temperature error");
      } else {
        Serial.write(dtostrf(t, 4, 1, buffer));
      }
     }
     if (val == '1') {
      digitalWrite(LED, HIGH);
      Serial.write("diode one");
     } else if (val == '0') {
      digitalWrite(LED, LOW);
      Serial.write("diode off");
     }
  }
}
