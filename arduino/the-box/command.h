typedef enum 
{
  UNKNOWN, 
  TEMPERATURE, 
  RELAY_ON, 
  RELAY_OFF, 
  BLINK
} Command;

Command commands[] = {
  UNKNOWN, 
  TEMPERATURE, 
  RELAY_ON, 
  RELAY_OFF,
  BLINK,
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
 
