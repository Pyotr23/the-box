typedef enum 
{  
  NONE, 
  CONTINUOUS_HEATING, 
  KEEP_RANGE,
} Mode;

Mode commands[] = {
  NONE, 
  CONTINUOUS_HEATING, 
  KEEP_RANGE,
};


int ModeToInt (Mode mode) {
  int i;
  for (i = 0; i < 5; i = i + 1) {
    if (mode == modes[i]) {
      return i;
    }
  }
  return i;
}

Mode IntToMode (int num) {
  return modes[num];
}