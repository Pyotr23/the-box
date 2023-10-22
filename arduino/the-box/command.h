typedef enum 
{
  UNKNOWN, 
  TEMPERATURE, 
  RELAY_ON, 
  RELAY_OFF, 
  SET_ID, 
  GET_ID, 
  GET_LOWER_TEMPERATURE_THRESHOLD, 
  GET_HIGHER_TEMPERATURE_THRESHOLD, 
  SET_LOWER_TEMPERATURE_THRESHOLD,
  SET_HIGHER_TEMPERATURE_THRESHOLD,
  SET_MODE,
  GET_MODE,
  BLINK
} Command;

Command commands[] = {
  UNKNOWN, 
  TEMPERATURE, 
  RELAY_ON, 
  RELAY_OFF, 
  SET_ID, 
  GET_ID, 
  GET_LOWER_TEMPERATURE_THRESHOLD, 
  GET_HIGHER_TEMPERATURE_THRESHOLD, 
  SET_LOWER_TEMPERATURE_THRESHOLD,
  SET_HIGHER_TEMPERATURE_THRESHOLD,
  SET_MODE,
  GET_MODE,
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
 
