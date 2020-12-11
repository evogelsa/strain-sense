# Strain Sense Wearable Project

## How it works

The wearable device is attached to the user with an adhesive, but in future
version it would be integrated into clothing. The accelerometer and main device
is attached to the lower back region like in the below example.

![Lower back][PyBadge]

The flex sensor gets mounted below the knee of the user in a similar fashion.

![Knee][flex]

Data from the two sensors is processed onboard the PyBadge device. Two states
are considered during operation of the device. Using the accelerometer, we are
able to detect if the user is moving by comparing the average changes in
acceleration magnitude with a threshold value that was tuned to be indicative of
movement. By reading the output voltage across the flex sensor and through the
low pass filters, we can determine the flex resistance, which is directly
proportional to the bend angle. This allows us to detect whether or not the user
is extending their leg, and therefore, standing.

The device monitors these two states and determines if the user has been idle
and sitting for an adjustable amount of time. If the device determines the user
to be inactive, it will beep periodically to signal the user to stand up and
move. However, if the user stands and moves around at any point, the device will
reset the inactivity timer.

On top of this, we have a webportal where the user can visualize their activity
and log their lower back pain (LBP) on a day to day basis. While we were unable
to finish collecting data, we expect that wearing the device will correlate
with an increased activity while working and sitting, and in turn this will
lead to a decreased lower back pain.

## Design

The device utilizes two main sensors: a variable resistance flex sensor, and an
accelerometer. The flex sensor is mounted beneath the knee of the user via an
adhesive, and the accelerometer is integrated into the main body of the device
which rests on the lower back of the user.

The flex sensor has a second order hardware lowpass filter at the end of its
output. The schematic as is shown below.

![Flex Sensor Filter][flex-schem]

The accelerometer data is combined into a single magnitude and filtered with a
second order butterworth filter in software.

## Operation

To create this wearable device you must first design the flex sensor hardware
to attach to the PyBadge device. We recommend soldering the required components
to a protoboard and using header pins to connect directly to the PyBadge. An
example is provided below.

![lowpass][lowpass]

The device will then need to be mounted to the body and plugged into a
computer. The code inside the [pybadge](/pybadge) directory can be uploaded
to the device via the Arduino IDE. Some configuration will be required to get
the PyBadge working with Arduino. [Adafruit has provided documentation on how to
do this][adafruit-pybadge].

Then once the device has been configured and the code is running, the user
should run the Python client code inside the [serial](/serial) directory. This
will create a connection to the device which will log the data. Upon stopping
the program it will send the data to the [server][server] to be logged.

For best results we recommend the user log their LBP at least once daily. This
can be done using the local server inside the [server](/server) source code, or
on the public version [here][server]. Doing so will provide useful feedback on
how the device has helped reduce LBP. The webpotal provides a dashboard which
gives the user a calendar heatmap to log and display LBP, and it will also
show the most recent sensor readings from the device that have been logged.

![example server dash][server-dash]

The data visualization is normalized to simplify its presentation, and for this
reason it doesn't represent real numeric values. However, the user is able to
see if they have been moving regularly by viewing the peaks in flex data and in
accelerometer data.

## Demonstration

Example videos of the device and of webportal can be viewed in the
[examples](/examples) directory.

[flex-schem]: https://i.imgur.com/NrKf4Wi.png
[PyBadge]: https://cdn.discordapp.com/attachments/539516208161226754/787091452118695976/image0.jpg
[flex]: https://cdn.discordapp.com/attachments/539516208161226754/787091452706553866/image1.jpg
[lowpass]: https://i.imgur.com/hZTa0sJ.jpeg
[adafruit-pybadge]: https://learn.adafruit.com/adafruit-pybadge/using-with-arduino-ide
[server]: https://ethanvogelsang.com/wearables/login
[server-dash]: https://i.imgur.com/bA6Lpi2.png
