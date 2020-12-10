#include <Adafruit_Arcada.h>
#include <elapsedMillis.h>
#include <Filters.h>
#define mvmnt_win_len 15

elapsedMillis t;
Adafruit_Arcada arcada;

const int FLEX_PIN = A0;
const float VCC = 4.16; // measure vcc for best accuracy
const float DIV_R = 10000; // measure divider resistance for best accuracy

float mvmnt_detect_wndw[mvmnt_win_len];
float accel_out;

const float hp_freq = 10;// Hz
FilterOnePole hp_mag(HIGHPASS, hp_freq);

const float flex_hp_freq = 10;// Hz
FilterOnePole hp_flex(HIGHPASS, flex_hp_freq);

float avg_list(float* list, int len)
{
  float run_sum = 0;
  for (int i = 0; i < len; i++)
  {
    run_sum += list[i];
  }
  return run_sum / len;
}

void add_2_list(uint8_t value, uint8_t* list, int len)
{
  for (int i = 0; i < len - 1; i++)
  {
    list[i] = list[i+1];
  }
  list[len-1] = value;
}



void setup() {
    // initialize serial connection
    Serial.begin(9600);
    Serial.println("Serial initialized.");

    // init pybadge
    arcada.arcadaBegin();

    // initialize pybadge display and turn on backlight
    arcada.displayBegin();
    for (int i = 0; i <= 255; i++) {
        arcada.setBacklight(i);
        delay(1);
    }
    Serial.print("flex");
    Serial.print(",");
    Serial.println("accel");
}

void loop() {
    if (arcada.hasAccel()) {
        // get pybadge events
        sensors_event_t event;
        arcada.accel->getEvent(&event);

        // clear a spot on screen
        arcada.display->fillRect(0, 0, 160, 128, ARCADA_BLACK);

        float xsq = sq(event.acceleration.x);
        float ysq = sq(event.acceleration.y);
        float zsq = sq(event.acceleration.z);
        float accelMag = sqrt(xsq + ysq + zsq);

        hp_mag.input(accelMag);
        accel_out = hp_mag.output();
        // display the accel data
        arcada.display->setCursor(0, 0);
        char a[10];
        sprintf(a, "Z: %8.1f", accel_out);
        arcada.display->print(a);

        // print the accelerometer data to serial plotter
        Serial.print(accel_out);
        Serial.print(",");
    }

    // measure resistance from flex sensor voltage divider
    int adc = analogRead(FLEX_PIN);
    float voltage = adc * VCC / 1023.0;
    float resistance = DIV_R * (VCC / voltage - 1.0);
    //hp_flex.input(resistance);
    //resistance = hp_flex.output();

    // show resistance on pybadge
    arcada.display->setCursor(0, 16);
    char r[10];
    sprintf(r, "R: %6.0f", resistance);
    arcada.display->print(r);

    // output resistance to serial
    Serial.print(resistance);
    Serial.print("\n");
    delay(25);
}
