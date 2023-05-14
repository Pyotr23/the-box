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
} Command;

const int UNKNOWN_NUM = 0;
const int TEMPERATURE_NUM = 1;
const int RELAY_ON_NUM = 2;
const int RELAY_OFF_NUM = 3;
const int SET_ID_NUM = 4;
const int GET_ID_NUM = 5;
const int GET_LOWER_TEMPERATURE_THRESHOLD_NUM = 6;
const int GET_HIGHER_TEMPERATURE_THRESHOLD_NUM = 7;
const int SET_LOWER_TEMPERATURE_THRESHOLD_NUM = 8;
const int SET_HIGHER_TEMPERATURE_THRESHOLD_NUM = 9;

int CommandToInt (Command command) {
  switch (command) {
    case TEMPERATURE: return TEMPERATURE_NUM;
    case RELAY_ON: return RELAY_ON_NUM;
    case RELAY_OFF: return RELAY_OFF_NUM;
    case SET_ID: return SET_ID_NUM;
    case GET_ID: return GET_ID_NUM;
    case UNKNOWN: return UNKNOWN_NUM;
    case GET_LOWER_TEMPERATURE_THRESHOLD: return GET_LOWER_TEMPERATURE_THRESHOLD_NUM;
    case SET_LOWER_TEMPERATURE_THRESHOLD: return SET_LOWER_TEMPERATURE_THRESHOLD_NUM;
    case GET_HIGHER_TEMPERATURE_THRESHOLD: return GET_HIGHER_TEMPERATURE_THRESHOLD_NUM;
    case SET_HIGHER_TEMPERATURE_THRESHOLD: return SET_HIGHER_TEMPERATURE_THRESHOLD_NUM;
  }
}

Command IntToCommand (int num) {
  switch (num) {
    case TEMPERATURE_NUM: return TEMPERATURE;
    case RELAY_ON_NUM: return RELAY_ON;
    case RELAY_OFF_NUM: return RELAY_OFF;
    case SET_ID_NUM: return SET_ID;
    case GET_ID_NUM: return GET_ID;
    case GET_LOWER_TEMPERATURE_THRESHOLD_NUM: return GET_LOWER_TEMPERATURE_THRESHOLD;
    case SET_LOWER_TEMPERATURE_THRESHOLD_NUM: return SET_LOWER_TEMPERATURE_THRESHOLD;
    case GET_HIGHER_TEMPERATURE_THRESHOLD_NUM: return GET_HIGHER_TEMPERATURE_THRESHOLD;
    case SET_HIGHER_TEMPERATURE_THRESHOLD_NUM: return SET_HIGHER_TEMPERATURE_THRESHOLD;
    default: return UNKNOWN;
  }
}
