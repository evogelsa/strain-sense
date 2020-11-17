import serial
import io
import time
import FindSerial
from serutils import *
#from matplotlib import pyplot as plt

N_DATAPOINTS = 800

port = FindSerial.serial_ports()[0]
ser = serial.Serial(
	port=port,
	baudrate=115200)
filename ="data/accel_data_running.csv"
with open(filename, "w") as f:
	f.write("X,Y,Z\n")

while True:
	try:
		line = ser.readline().decode()[:-1]#[1:][:-2]
		if line != "":
			# print("\"", end="")
			print(line)
			try:
				x_bit, y_bit, z_bit = line.split(",")
				x = float(x_bit)
				y = float(y_bit)
				z = float(z_bit)
				data = [x,y,z]
				# print(x_bit)
				# dat = [x_bit, y_bit, z_bit]
				# data = []
				# for i in range(3):
				# 	data.append(float(dat[i].split(":")[1]))
				with open(filename, "a") as f:
					f.write("{},{},{}\n".format(data[0], data[1], data[2]))
			except:
				# print("not enough")
				pass
	except KeyboardInterrupt:
		print("user quit")
		break
	ser.flushInput()
