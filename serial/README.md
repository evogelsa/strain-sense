# Serial Communication Protocols

This directory contains the client serial protocols for communicating with the
wearable device.

## Use Instructions

Make sure that the requirements outlined in [requirements.txt][requirements]
are installed. If using pip3, these can be automatically installed using
`pip3 install -r requirements.txt`.

To run the client, make sure the wearable is connected to the client PC with a
serial connection. Then run [main.py][main.py] with python3. The user will need
to supply their login credentials to the pop up before sending data to the
server. Data collection begins when pressing the start buttong, and upon
pressing the stop buttong, the data collection will stop, be parsed, and sent
to the server with the provided credentials.

## Use with a Local Server

If running the server locally, follow the above steps to configure your
environment. When running the main client file, provide the '-local' flag to
make sure the data is sent to the local server.

[requirements]: /serial/requirements.txt
[main.py]: /serial/main.py
