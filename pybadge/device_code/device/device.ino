#include <Adafruit_Arcada.h>
#include <math.h>
#include <elapsedMillis.h>
#define win_len 15
#define move_time_min 1 // 30
Adafruit_Arcada arcada;
elapsedMillis t;


float run_x[win_len];
float run_y[win_len];
float run_z[win_len];
float run_mag[win_len];

float avg_list(float* list)
{
  float run_sum = 0;
  for (int i = 0; i < win_len; i++)
  {
    run_sum += list[i];
  }
  return run_sum / win_len;
}

void add_2_list(float value, float* list)
{
  for (int i = 0; i < win_len-1; i++)
  {
    list[i] = list[i+1];
  }
  list[win_len-1] = value;
}

float accel_mag(float x, float y, float z)
{
  return sqrt(sq(x) + sq(y) + sq(z));
}



volatile uint16_t milliseconds = 0;
void timercallback() 
{
  analogWrite(13, milliseconds);  // pulse the LED
  if (milliseconds == 0) {
    milliseconds = 255;
  } else {
    milliseconds--;
  }
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

    arcada.timerCallback(1000, timercallback);
}

int txt_y = 60;
float x_ac = 0;
float y_ac = 0;
float z_ac = 0;
float ac_mag = 0;

float av_x = 0;
float av_y = 0;
float av_z = 0;
float av_mag = 0;

void loop() {
    if (arcada.hasAccel()) {
        // get pybadge events
        sensors_event_t event;
        arcada.accel->getEvent(&event);
        
        x_ac = event.acceleration.x;
        y_ac = event.acceleration.y;
        z_ac = event.acceleration.z;

        ac_mag = accel_mag(x_ac, y_ac, z_ac);
        
        // clear a spot on screen
        arcada.display->fillScreen(ARCADA_BLACK);
        //arcada.display->fillRect(0, 0, 160, 8, ARCADA_BLACK);

        // display the x accel data
        arcada.display->setCursor(0, txt_y);
        char x[9];
        sprintf(x, "X: %4.1f", event.acceleration.x);
        arcada.display->print(x);
        //arcada.display->print(event.acceleration.x, 1);

        // display the y accel data
        arcada.display->setCursor(50, txt_y);
        char y[9];
        sprintf(y, "Y: %4.1f", event.acceleration.y);
        arcada.display->print(y);
        //arcada.display->print(event.acceleration.y, 1);

        // display the z accel data
        arcada.display->setCursor(100, txt_y);
        char z[9];
        sprintf(z, "Z: %4.1f", event.acceleration.z);
        arcada.display->print(z);
        //arcada.display->print(event.acceleration.z, 1);

        // print the accelerometer data to serial plotter
        /* Serial.print("\tX:"); */
        /* Serial.print(event.acceleration.x); */
        /*  */
        /* Serial.print("\tY:"); */
        /* Serial.print(event.acceleration.y); */
        /*  */
        /* Serial.print("\tZ:"); */
        /* Serial.println(event.acceleration.z); */
        add_2_list(x_ac, run_x);
        add_2_list(y_ac, run_y);
        add_2_list(z_ac, run_z);
        add_2_list(ac_mag, run_mag);
        

        av_x = avg_list(run_x);
        av_y = avg_list(run_y);
        av_z = avg_list(run_z);
        av_mag = avg_list(run_mag);
        
        if (t > move_time_min * 60*1000)
        {
          t = 0;
          //    OH YEAH NOTIFICATION TIME
          // TIME TO SING THE SONG OF MY PEOPLE
          //   (if they haven't moved enough)
        }
        else
        {
          //check if they're moving enough
          //maybe count total movement?
          //accumulate changes in acceleration?
        }
        
        Serial.print(av_x);//event.acceleration.x);
        Serial.print(",");
        Serial.print(av_y);//event.acceleration.y);
        Serial.print(",");
        Serial.print(av_z);//event.acceleration.z);
        Serial.print(",");
        Serial.print(av_mag);
        
//        Serial.print(",");
//        Serial.print(event.acceleration.x);
//        Serial.print(",");
//        Serial.print(event.acceleration.y);
//        Serial.print(",");
//        Serial.print(event.acceleration.z);
//        Serial.print(",");
//        Serial.print(ac_mag);
        
        Serial.print("\n");
    }
    delay(25);
}
