from FindSerial import *

sp = serial_ports()[0]
console = serial.Serial(
	port = sp,
	baudrate = 115200
)
plotter = Plotter((500,300), npoints=100, nlines=3, names=["X","Y","Z"], colors=[(255,0,0),(0,255,0),(0,0,255)], p_min=-180, p_max=180)
