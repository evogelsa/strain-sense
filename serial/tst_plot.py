from serutils import *
import time, math
plotter = Plotter((1000,700), npoints=500, nlines=3, names=["X","Y","Z"], colors=[(255,0,0),(0,255,0),(0,0,255)])#, p_min=-180, p_max=180)

f1 =1
f2 = 2
f3 = .5
A1 = 200
A2 = 200
A3 = 200
while True:
	k = cv2.waitKey(1) & 0xFF
	t = time.time()
	if k == ord("q"):
		break
	# plotter.setupPage()
	plotter.addVal(A1 * math.sin(2*math.pi*f1 * t), 0)
	plotter.addVal(A2 * math.sin(2*math.pi*f2 * t), 1)
	plotter.addVal(A3 * math.sin(2*math.pi*f3 * t), 2)
	plotter.graphPoints()
	cv2.imshow("plot", plotter.img)
