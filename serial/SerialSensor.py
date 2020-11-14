import serial
import io
import time
import FindSerial
from MyUtils import *
#from matplotlib import pyplot as plt
port = FindSerial.serial_ports()[0]
ser = serial.Serial(
	port=port,\
	baudrate=115200,\
	parity=serial.PARITY_NONE,\
	stopbits=serial.STOPBITS_ONE,\
	bytesize=serial.EIGHTBITS,\
	timeout=0)
plotter = Plotter((500,300), npoints=100, nlines=3, names=["X","Y","Z"], colors=[(255,0,0),(0,255,0),(0,0,255)], p_min=-180, p_max=180)
print("connected to: " + ser.portstr)
count=1
sensor = "A0"
text = ""
timer = Timer()

gyro_components = ['x', 'y', 'z']
gyro_values = [0,0,0]
dist_val = 0

splt = ','
delim='\t'

#plt.axis([0,100,0,360])
#x=np.linspace(0, 20, 1000)

print_str = ""

while True:
	timer.upd()
	count += 1
	plotter.setupPage()
	key = cv2.waitKey(1) & 0xFF
	text = ""
	print_str = ""
	got_serial = False
	for line in ser.readline():
		count = count+1
		text += chr(line)
		got_serial = True
	if got_serial:
		text = text[:-1]
		#print(text)
		#divide input data
		sensor_datas = text.split(delim)
		#print(sensor_datas)
		for sensor_data in sensor_datas:
			#print(sensor_data)#			gyro\tx\ty\tz\t
			if len(sensor_data) > 1:
				datas = sensor_data.split(splt)
				data_type = datas[0]
				datas = datas[1:]
				#datas = datas.split("\t")
				#print(datas)
				if data_type == "READ":
					print_str += "Arduino recieved: {}"
				if data_type == "GYRO":
					for data_idx in range(len(datas)):
						gyro_values[data_idx] = float(datas[data_idx])
	#					print("data value: {}\tindex: {}".format(gyro_values[data_idx], data_idx))
						plotter.addVal(gyro_values[2-data_idx], data_idx)
						#plotter.graphLine(data_idx)
					print_str += "X:{: >7.2f}\tY:{: >7.2f}\tZ:{: >7.2f}".format(gyro_values[0], gyro_values[1], gyro_values[2]) + '\t'
				elif data_type == "DIST":
					dist_val = datas[1]
					print_str += "\tdist:{: >7}".format(dist_val) + '\t'
				elif data_type == 'A0':
					val = datas[0]
					print_str += "A0:{: >4}".format(val) + '\t'
				plotter.graphPoints()################################################## maybe unindent this
			cv2.imshow("plot", plotter.img)
		if len(print_str) > 0:
			print(print_str)
	ser.flushInput()
	if key == ord("q"):
		break

	#ser.flush()
plt.show()
ser.close()
