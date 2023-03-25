typedef enum {UNKNOWN, TEMPERATURE, RELAY_ON, RELAY_OFF} Command;

const int UNKNOWN_NUM = '0';
const int TEMPERATURE_NUM = '1';
const int RELAY_ON_NUM = '2';
const int RELAY_OFF_NUM = '3';

int CommandToInt (Command command) {
  switch (command) {
    case TEMPERATURE: return TEMPERATURE_NUM;
    case RELAY_ON: return RELAY_ON_NUM;
    case RELAY_OFF: return RELAY_OFF_NUM;
    case UNKNOWN: return UNKNOWN_NUM;
  }
}

Command IntToCommand (int num) {
  switch (num) {
    case TEMPERATURE_NUM: return TEMPERATURE;
    case RELAY_ON_NUM: return RELAY_ON;
    case RELAY_OFF_NUM: return RELAY_OFF;
    default: return UNKNOWN;
  }
}