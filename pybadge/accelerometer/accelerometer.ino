#include <Adafruit_Arcada.h>

Adafruit_Arcada arcada;

const int FLEX_PIN = A0;
const float VCC = 4.16; // measure vcc for best accuracy
const float DIV_R = 47000; // measure divider resistance for best accuracy

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

        // display the accel data
        arcada.display->setCursor(0, 0);
        char a[10];
        sprintf(a, "Z: %8.1f", accelMag);
        arcada.display->print(a);

        // print the accelerometer data to serial plotter
        Serial.print(a);
        Serial.print(",");
    }

    // measure resistance from flex sensor voltage divider
    int adc = analogRead(FLEX_PIN);
    float voltage = adc * VCC / 1023.0;
    float resistance = DIV_R * (VCC / voltage - 1.0);

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
