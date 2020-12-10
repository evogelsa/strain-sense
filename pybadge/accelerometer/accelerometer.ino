#include <Adafruit_Arcada.h>
#include <elapsedMillis.h>
#include <Filters.h>
#include <Filters/Butterworth.hpp>
#include <AH/Timing/MillisMicrosTimer.hpp>
#include "audio.h"

//#define DO_DISPLAY

// system configuration parameters
const float FS = 100;                 // Sample rate, Hz
const float LP_FC = 40;               // Cutoff frequency for low pass, Hz
const float HP_FC = 240;              // Cutoff frequency for high pass, Hz
const float LP_FN = 2 * LP_FC / FS;   // LP normalized freq
const float HP_FN = 2 * HP_FC / FS;   // HP normalized freq
const float MOVE_TIME = 0.25*60*1000; // time in ms to decide about notification

elapsedMillis t;
const int BEEP_TIME = 5 * 1000; // beep every so and so seconds until movement
elapsedMillis beep_t;           // for beeping every beep

Adafruit_Arcada arcada;

// define pins and constants relating to flex sensor
const int FLEX_PIN = A0;
const float VCC = 4.16;   // configure to vcc of your setup
const float DIV_R = 10e3; // second resistor in voltage divider resistance

// create window for running average, used for movement detection
const int ACCEL_N_INPUTS = FS * 2; // FS * seconds of data desired
float accelInputs[ACCEL_N_INPUTS];
float accelSum = 0;
unsigned int accelIdx = 0;

// define movement detection variables
const float MOVE_THRESHOLD = 0.17; // below threshold means not moving
const int STAND_THRESHOLD = 22000;

// define timer to control sample rate
Timer<micros> samplingTimer = std::round(1e6 / FS);

// define second order low pass butterworth filter
auto lpFilter = butter<2>(LP_FN);

void updateAccel(float *accelMag, float *accelAvg) {
    if (arcada.hasAccel()) {
        // get pybadge events
        sensors_event_t event;
        arcada.accel->getEvent(&event);

        // measure the squared acceleration of each axis
        float xsq = sq(event.acceleration.x);
        float ysq = sq(event.acceleration.y);
        float zsq = sq(event.acceleration.z);

        // combine each direction into single magnitude
        *accelMag = sqrt(xsq + ysq + zsq);

        // filter acceleration magnitude through lowpass butterworth
        *accelMag = lpFilter(*accelMag);

        // add acceleration magnitude to rolling sum
        accelSum += *accelMag;

        // check if we have reached the limit of values we want to average
        // subtract the oldest input if we have
        if (accelIdx >= ACCEL_N_INPUTS) {
            accelSum -= accelInputs[accelIdx % ACCEL_N_INPUTS];
        }

        // track most recent input
        accelInputs[accelIdx % ACCEL_N_INPUTS] = *accelMag;
        accelIdx++;

        // calculate the average
        *accelAvg = accelSum / ACCEL_N_INPUTS;
    }
}

void setup() {
    // initialize serial connection
    Serial.begin(9600);

    // init pybadge
    arcada.arcadaBegin();

#ifdef DO_DISPLAY
    // initialize pybadge display and turn on backlight
    arcada.displayBegin();
    for (int i = 0; i <= 255; i++) {
        arcada.setBacklight(i);
        delay(1);
    }
#endif
}

void loop() {
    if (samplingTimer) {
        float *accelMag, *accelAvg;
        updateAccel(accelMag, accelAvg);

        // measure resistance from flex sensor voltage divider
        int adc = analogRead(FLEX_PIN);
        float voltage = adc * VCC / 1023.0;
        float resistance = DIV_R * (VCC / voltage - 1.0);

        // determine if standing and movement has been detected
        bool standingDetected = (resistance < STAND_THRESHOLD);
        bool movementDetected = (*accelAvg > MOVE_THRESHOLD);

        // output data to serial
        Serial.print(*accelMag);
        Serial.print(",");
        Serial.print(*accelAvg);
        Serial.print(",");
        Serial.print(standingDetected * 3 *1000);
        Serial.print(",");
        Serial.print(movementDetected * 4 *1000);
        Serial.print(",");
        Serial.print((movementDetected && standingDetected)* 5 *1000);
        Serial.print(",");
        Serial.print(resistance);
        Serial.print("\n");

        if (movementDetected && standingDetected) {
            t = 0; // reset timer because user has moved
        } else {
            if (t > MOVE_TIME) //user hasn't moved recently enough
            {
                if (beep_t > BEEP_TIME) // device hasn't beeped in a bit
                {
                    // beep device, but causes sensor funkiness
                    beep_t = 0;
                    arcada.enableSpeaker(true);
                    play_tune(audio, sizeof(audio));
                    arcada.enableSpeaker(false);
                }
            }
        }
    }

#ifdef DO_DISPLAY
    // display the accel average
    arcada.display->fillRect(0, 0, 160, 128, ARCADA_BLACK); // clear a spot on screen
    arcada.display->setCursor(0, 0);
    char a[10];
    sprintf(a, "Z: %8.1f", accelAvg);
    arcada.display->print(a);

    // show resistance on pybadge
    arcada.display->setCursor(0, 16);
    char r[10];
    sprintf(r, "R: %6.0f", resistance);
    arcada.display->print(r);
#endif
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
