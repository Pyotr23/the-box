typedef enum 
{
  UNKNOWN, 
  TEMPERATURE, 
  PIN_ON, 
  PIN_OFF, 
  CHECK_PIN
} Command;

Command commands[] = {
  UNKNOWN, 
  TEMPERATURE, 
  PIN_ON, 
  PIN_OFF,
  CHECK_PIN  
};


int CommandToInt (Command command) {
  int i;
  for (i = 0; i < 5; i = i + 1) {
    if (command == commands[i]) {
      return i;
    }
  }
  return i;
}

Command IntToCommand (int num) {
  return commands[num];
}
 
