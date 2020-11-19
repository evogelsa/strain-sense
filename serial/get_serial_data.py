import sys
import serial
import io
import time
import FindSerial
import PySimpleGUI as sg
import threading
import server


N_DATAPOINTS = 800

class DataReader():
    def __init__(self, filename='file.csv', serial=None):
        self._running = True
        self.filename = filename
        self.serial = serial
        with open(filename, "w") as f:
            f.write("X,Y,Z\n")

    def stop(self):
        self._running = False

    def run(self):
        while self._running:
            line = self.serial.readline().decode()[:-1]
            if line != "":
                print(line)
                x_bit, y_bit, z_bit = line.split(",")
                x = float(x_bit)
                y = float(y_bit)
                z = float(z_bit)
                data = [x,y,z]
                with open(self.filename, "a") as f:
                    f.write("{},{},{}\n".format(data[0], data[1], data[2]))
            self.serial.flushInput()

def main():
    port = FindSerial.serial_ports()[0]

    ser = serial.Serial(
        port=port,
        baudrate=115200)

    filename ="data/accel_data_running.csv"

    reader = DataReader(filename=filename, serial=ser)
    reader_thread = threading.Thread(target=reader.run)

    sg.theme('DarkGrey4')
    layout = [[sg.Text('Press start to begin data recording and press stop')],
              [sg.Text('to stop recording and send data to server.')],
              [sg.Text('Username'), sg.InputText()],
              [sg.Text('Password'), sg.InputText(password_char='*')],
              [sg.Button('Start'), sg.Button('Stop')]]
    window = sg.Window('Strain Sense', layout)

    if '-local' in sys.argv:
        req_url = r'https://localhost:32321/wearables/dashboard'
    else:
        req_url = r'https://ethanvogelsang.com/wearables/dashboard'

    while True:
        event, values = window.read()
        if event == sg.WIN_CLOSED:
            break
        if event == 'Start':
            reader_thread.start()
        if event == 'Stop':
            reader.stop()
            reader_thread.join()
            usr = str(values[0])
            pwd = str(values[1])
            server.send(usr, pwd, filename, req_url)
            break

if __name__ == '__main__':
    sys.exit(main())
