import serial
import io
import time
import FindSerial
from serutils import *
#from matplotlib import pyplot as plt
port = FindSerial.serial_ports()[0]
ser = serial.Serial(
	port=port,\
	baudrate=115200,\
	parity=serial.PARITY_NONE,\
	stopbits=serial.STOPBITS_ONE,\
	bytesize=serial.EIGHTBITS,\
	timeout=0)
plotter = Plotter((800,300), npoints=800, nlines=3, names=["X","Y","Z"], colors=[(255,0,0),(0,255,0),(0,0,255)], p_min=-180, p_max=180)
print("connected to: " + ser.portstr)
count=1
sensor = "A0"
text = ""

gyro_components = ['x', 'y', 'z']
gyro_values = [0,0,0]
dist_val = 0

splt = ','
delim='\t'

#plt.axis([0,100,0,360])
#x=np.linspace(0, 20, 1000)

print_str = ""

while True:
	key = cv2.waitKey(1) & 0xFF

	count += 1
	plotter.setupPage()
	text = ""
	print_str = ""

	got_serial = False
	# for line in ser.readline():
	# 	count = count+1
	# 	text += chr(line)
	# 	got_serial = True
	line = ser.readline().decode()[1:][:-2]

	# print(line.decode())
	if line != "":
		# print("\"", end="")
		# print(line, end = "\"\n")
		try:
			x_bit, y_bit, z_bit = line.split("\t")
			dat = [x_bit, y_bit, z_bit]
			for i in range(3):
				plotter.addVal(float(dat[i].split(":")[1]), i)
		except:
			# print("not enough")
			pass

	ser.flushInput()
	plotter.graphPoints()
	cv2.imshow("plot", plotter.img)
	if key == ord("q"):
		break
	if key == ord("s"):
		vals = plotter.get_data()
		# np.savez("data/valz", vals)
		np.savetxt("data/accel_data.csv", vals, delimiter=",")
		# print(vals.shape)
