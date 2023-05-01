typedef enum {UNKNOWN, TEMPERATURE, RELAY_ON, RELAY_OFF, SET_ID, GET_ID} Command;

const int UNKNOWN_NUM = 0;
const int TEMPERATURE_NUM = 1;
const int RELAY_ON_NUM = 2;
const int RELAY_OFF_NUM = 3;
const int SET_ID_NUM = 4;
const int GET_ID_NUM = 5;

int CommandToInt (Command command) {
  switch (command) {
    case TEMPERATURE: return TEMPERATURE_NUM;
    case RELAY_ON: return RELAY_ON_NUM;
    case RELAY_OFF: return RELAY_OFF_NUM;
    case SET_ID: return SET_ID_NUM;
    case GET_ID: return GET_ID_NUM;
    case UNKNOWN: return UNKNOWN_NUM;
  }
}

Command IntToCommand (int num) {
  switch (num) {
    case TEMPERATURE_NUM: return TEMPERATURE;
    case RELAY_ON_NUM: return RELAY_ON;
    case RELAY_OFF_NUM: return RELAY_OFF;
    case SET_ID_NUM: return SET_ID;
    case GET_ID_NUM: return GET_ID;
    default: return UNKNOWN;
  }
}
