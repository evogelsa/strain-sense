#include <Adafruit_Arcada.h>

Adafruit_Arcada arcada;

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
        Serial.print("\tX:");
        Serial.print(event.acceleration.x);

        Serial.print("\tY:");
        Serial.print(event.acceleration.y);

        Serial.print("\tZ:");
        Serial.println(event.acceleration.z);

    }
    delay(1/60);
}
