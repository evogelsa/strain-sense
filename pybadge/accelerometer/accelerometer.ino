#include <Adafruit_Arcada.h>
#include <Filters.h>
#include <Filters/Butterworth.hpp>
#include <elapsedMillis.h>
#include "audio.h"

#define DO_DISPLAY

Adafruit_Arcada arcada;

// system configuration parameters
const double FS = 100;                 // Sample rate, Hz
const double LP_FC = 20;               // Cutoff frequency for low pass, Hz
const double LP_FN = 2 * LP_FC / FS;   // LP normalized freq
const float MOVE_TIME = .25*60*1000;  // time to decide about notification, ms
const float BEEP_TIME = 5 * 1000;      // beep every 5 seconds until movement
const float SAMPLE_TIME = 10 / FS;

// define pins and constants relating to flex sensor
const int FLEX_PIN = A2;
const float VCC = 4.16;   // configure to vcc of your setup
const float DIV_R = 1; // second resistor in voltage divider resistance in 10 kOhms

// create window for running average, used for movement detection
const int ACCEL_N_INPUTS = FS * .35; // FS * seconds of data desired
float accelInputs[ACCEL_N_INPUTS];
float accelSum = 0;
unsigned int accelIdx = 0;

// define movement detection variables
const float MOVE_THRESHOLD = 0.17; // below threshold means not moving
const float STAND_THRESHOLD = 2.2;

// define timer to control sample rate
elapsedMillis samplingTimer;
elapsedMillis beepTimer;
elapsedMillis moveTimer;

// define second order low pass butterworth filter
auto lpFilter = butter<2>(LP_FN);

void playSound(const uint8_t *audio, uint32_t audio_length) {
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
        *accelMag = abs(sqrt(xsq + ysq + zsq) - 9.3);

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
    Serial.begin(115200);

    arcada.arcadaBegin();

#ifdef DO_DISPLAY
    // initialize pybadge display and turn on backlight
    arcada.displayBegin();
    for (int i = 0; i <= 255; i++) {
        arcada.setBacklight(i);
        delay(1);
    }
    arcada.display->fillRect(0, 0, 160, 128, ARCADA_BLACK); // clear screen
#endif
}

void loop() {
    if (samplingTimer >= SAMPLE_TIME) {
        samplingTimer -= SAMPLE_TIME;

        float *accelMag = (float*) malloc(sizeof(float));
        float *accelAvg = (float*) malloc(sizeof(float));
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
        Serial.print(standingDetected * 3);
        Serial.print(",");
        Serial.print(movementDetected * 4);
        Serial.print(",");
        Serial.print((movementDetected && standingDetected)* 5 );
        Serial.print(",");
        Serial.print(resistance);
        Serial.print("\n");

        if (movementDetected && standingDetected) {
            moveTimer = 0;
        } else if ((moveTimer >= MOVE_TIME) && (beepTimer >= BEEP_TIME)) {
            // if user hasn't moved recently enough and device hasn't beeped
            // recently then beep device, beware causes sensor funkiness
            //beepTimer -= BEEP_TIME;
            beepTimer = 0;
            arcada.enableSpeaker(true);
            playSound(audio, sizeof(audio));
            arcada.enableSpeaker(false);
        }

#ifdef DO_DISPLAY
        // display the accel average
        arcada.display->fillRect(0, 0, 160, 128, ARCADA_BLACK); // clear screen
        arcada.display->setCursor(0, 0);
        char a[10];
        sprintf(a, "Z: %6.3f", *accelAvg);
        arcada.display->print(a);

        // show resistance on pybadge
        arcada.display->setCursor(0, 16);
        char r[10];
        sprintf(r, "R: %6.3f", resistance);
        arcada.display->print(r);
#endif

        free(accelMag);
        free(accelAvg);
    }
}
