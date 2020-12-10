#include <Adafruit_Arcada.h>
#include <elapsedMillis.h>
#include <Filters.h>
#include "audio.h"
#define mvmnt_win_len 20 //length of array for movement detection
#define move_time 0.25 *60*1000 // time in ms to decide about notification
//                30
//#define DO_DISPLAY

elapsedMillis t;
const int beep_time = 5 * 1000; // beep every so and so seconds until movement
elapsedMillis beep_t; //for beeping every beep


Adafruit_Arcada arcada;

const int FLEX_PIN = A0;
const float VCC = 4.16; // measure vcc for best accuracy
const float DIV_R = 10000; // measure divider resistance for best accuracy

float mvmnt_detect_wndw[mvmnt_win_len];

float mvmnt_roll_sum = 0;
float mvmnt_avg = 0;
float accel_out;


bool movement_detected = false;
const float mvmnt_thresh = 0.17; //value found below which the rolling average for acceleration indicates stillness
bool standing_detected = false;
const int standing_thresh = 22000;


const float hp_freq = 0.01;// hp frequency to filter
FilterOnePole hp_mag(HIGHPASS, hp_freq);

const float lp_freq = 10;
FilterOnePole lp_mag(LOWPASS, lp_freq);



float avg_list(float* list, int len)
{
  float run_sum = 0;
  for (int i = 0; i < len; i++)
  {
    run_sum += list[i];
  }
  return run_sum / len;
}
void add_2_list(float value, float* list, int len)
{
  for (int i = 0; i < len - 1; i++)
  {
    list[i] = list[i+1];
  }
  list[len-1] = value;
}


void update_accel()
{
  if (arcada.hasAccel()) {
    // get pybadge events
    sensors_event_t event;
    arcada.accel->getEvent(&event);

    float xsq = sq(event.acceleration.x);
    float ysq = sq(event.acceleration.y);
    float zsq = sq(event.acceleration.z);
    float accelMag = sqrt(xsq + ysq + zsq);

    hp_mag.input(accelMag);
    lp_mag.input(hp_mag.output());//accelMag);
    //*
    accel_out = hp_mag.output();
    /*/
    accel_out = lp_mag.output();
    
    //*/
    mvmnt_roll_sum += abs(accel_out) - mvmnt_detect_wndw[0];
    add_2_list(abs(accel_out), mvmnt_detect_wndw, mvmnt_win_len); //add the filtered accel magnitude to the detection list
    mvmnt_avg = mvmnt_roll_sum / mvmnt_win_len;
  }
}


void setup() {
    // initialize serial connection
    Serial.begin(9600);
    for (int i = 0; i < mvmnt_win_len; i++)
    {
      mvmnt_detect_wndw[i] = 0;
    }

    // init pybadge
    arcada.arcadaBegin();

    
    // initialize pybadge display and turn on backlight
    #ifdef DO_DISPLAY
      arcada.displayBegin();
    for (int i = 0; i <= 255; i++) {
        arcada.setBacklight(i);
        delay(1);
    }
    #endif
    arcada.timerCallback(10, update_accel);
}

void loop() {
//    if (arcada.hasAccel()) {
//        // get pybadge events
//        sensors_event_t event;
//        arcada.accel->getEvent(&event);
//
//        float xsq = sq(event.acceleration.x);
//        float ysq = sq(event.acceleration.y);
//        float zsq = sq(event.acceleration.z);
//        float accelMag = sqrt(xsq + ysq + zsq);
//
//        hp_mag.input(accelMag);
//        lp_mag.input(hp_mag.output());//accelMag);
//        /*
//        accel_out = hp_mag.output();
//        /*/
//        accel_out = lp_mag.output();
//        
//        //*/
//        mvmnt_roll_sum += abs(accel_out) - mvmnt_detect_wndw[0];
//        add_2_list(abs(accel_out), mvmnt_detect_wndw, mvmnt_win_len); //add the filtered accel magnitude to the detection list
//        mvmnt_avg = mvmnt_roll_sum / mvmnt_win_len;
//    }
    // display the accel data or dont to save resources
    #ifdef DO_DISPLAY
      arcada.display->fillRect(0, 0, 160, 128, ARCADA_BLACK); // clear a spot on screen
      arcada.display->setCursor(0, 0);
      char a[10];
      sprintf(a, "Z: %8.1f", accel_out);
      arcada.display->print(a);
    #endif
    
    // measure resistance from flex sensor voltage divider
    int adc = analogRead(FLEX_PIN);
    float voltage = adc * VCC / 1023.0;
    float resistance = DIV_R * (VCC / voltage - 1.0);
<<<<<<< HEAD
    //hp_flex.input(resistance);
    //resistance = hp_flex.output();
=======
>>>>>>> 09d97136969c4891c618fa68a3efb4289181cf76

    // show resistance on pybadge or dont to save resources
    #ifdef DO_DISPLAY
      arcada.display->setCursor(0, 16);
      char r[10];
      sprintf(r, "R: %6.0f", resistance);
      arcada.display->print(r);
    #endif

    if (resistance < standing_thresh)
    {
      standing_detected = true;
    }
    else
    {
      standing_detected = false;
    }
    if (mvmnt_avg > mvmnt_thresh)
    {
      movement_detected = true;
    }
    else
    {
      movement_detected = false;
    }

    // output data to serial
    Serial.print(accel_out);
    Serial.print(",");
    Serial.print(mvmnt_avg*10000.0);
    Serial.print(",");
    Serial.print(standing_detected * 3 *1000);
    Serial.print(",");
    Serial.print(movement_detected * 4 *1000);
    Serial.print(",");
    Serial.print((movement_detected && standing_detected)* 5 *1000);
    Serial.print(",");
    Serial.print(resistance);
    Serial.print("\n");
//    delay(25);

    if (movement_detected && standing_detected)
    {
      t = 0; // reset timer because user has moved
    }
    else
    {
      if (t > move_time) //user hasn't moved in enough time
      {
        if (beep_t > beep_time) // device hasn't beeped in a bit
        {
          beep_t = 0;
          arcada.enableSpeaker(true);
          play_tune(audio, sizeof(audio)); // beep device. discombobulates sensors, beware
          arcada.enableSpeaker(false);
        }
      }
    }

}


void play_tune(const uint8_t *audio, uint32_t audio_length) {
  uint32_t t;
  uint32_t prior, usec = 1000000L / SAMPLE_RATE;
  analogWriteResolution(8);
  for (uint32_t i=0; i<audio_length; i++) {
    while((t = micros()) - prior < usec);
    analogWrite(A0, (uint16_t)audio[i] / 8);
    analogWrite(A1, (uint16_t)audio[i] / 8);
    prior = t;
  }
}
