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
        arcada.display->fillRect(0, 0, 160, 8, ARCADA_BLACK);

        // display the x accel data
        arcada.display->setCursor(0, 0);
        char x[9];
        sprintf(x, "X: %4.1f", event.acceleration.x);
        arcada.display->print(x);
        //arcada.display->print(event.acceleration.x, 1);

        // display the y accel data
        arcada.display->setCursor(50, 0);
        char y[9];
        sprintf(y, "Y: %4.1f", event.acceleration.y);
        arcada.display->print(y);
        //arcada.display->print(event.acceleration.y, 1);

        // display the z accel data
        arcada.display->setCursor(100, 0);
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
        Serial.print(event.acceleration.x);
        Serial.print(",");
        Serial.print(event.acceleration.y);
        Serial.print(",");
        Serial.print(event.acceleration.z);
        Serial.print(",");
    }

    // measure resistance from flex sensor voltage divider
    int adc = analogRead(FLEX_PIN);
    float voltage = adc * VCC / 1023.0;
    float resistance = DIV_R * (VCC / v - 1.0);

    // show resistance on pybadge
    arcada.display->setCursor(0, 16);
    char r[9];
    sprintf(r, "R: %4.1f", resistance);
    arcada.display->print(r);

    // output resistance to serial
    Serial.print(resistance);
    Serial.print("\n");
    delay(25);
}
