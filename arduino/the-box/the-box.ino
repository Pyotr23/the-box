const int LED = 13;

int val;

void setup() {
  Serial.begin(9600);
  pinMode(LED, OUTPUT);
  digitalWrite(LED, HIGH);
}

void loop() {
  if (Serial.available()) {
     val = Serial.read();
     if (val == '1') {
      digitalWrite(LED, HIGH);
      Serial.write("diode one");
     } else if (val == '0') {
      digitalWrite(LED, LOW);
      Serial.write("diode off");
     }
  }
}
