import cv2
import numpy as np
from collections import deque

class Plotter:
	def __init__(
			self,
			dim,
			npoints=30,
			nlines=1,
			plotpadding=30,
			textsps=5,
			sps=10,
			edgsps=30,
			p_min=None,
			p_max=None,
			names=["default"],
			fontsize=1,
			fontthickness=2,
			font=cv2.FONT_HERSHEY_DUPLEX,
			colors=[(255, 0, 0)],
			tickHeight=10,
			tickSps=20,
			tickCol=(0, 0, 0),
			bgColor=(255, 255, 255)):
		self.plotWidth = int(dim[0])
		self.plotHeight = int(dim[1])
		self.plot_minimum = p_min
		self.plot_maximum = p_max
		self.npoints = npoints
		self.plotPad = plotpadding
		self.names = names
		self.sps = sps
		self.edgsps  = edgsps
		self.textsps = textsps
		self.fontsize = fontsize
		self.fontthickness = fontthickness
		self.font = font
		self.textCol = (0,0,0)
		self.nlines = nlines
		self.bgCol = bgColor
		self.tickCol = tickCol
		# self.vals = [deque(maxlen=npoints) for i in range(nlines)]
		self.vals = [np.zeros((3, npoints)) for i in range(nlines)]
		for i in self.vals:
			for j in range(len(i)):
				i[j] = 0
		self.cols = colors
		self.img = np.zeros((dim[1], dim[0], 3), np.uint8)
		self.centX = int(self.plotWidth / 2)
		self.centY = int(self.plotHeight / 2)
		self.tickHeight = tickHeight
		self.tickSps = tickSps
	def setupPage(self, axis=False):
		"""plotmax = max(self.vals)[0]
		plotmin = min(self.vals)[0]
		#print("{}, {}".format(plotmin, plotmax))
		graphRange = 2 * (self.plotPad + max(plotmax, abs(plotmin)))
		#print("{}, {:.2f}, {:.2f}".format(graphRange, plotmin, plotmax))
		#print(plotmax)
		self.plotHeight = int(max(graphRange, self.plotHeight))
		self.centY = int(self.plotHeight / 2)"""
		self.img = np.zeros((self.plotHeight, self.plotWidth, 3), np.uint8)
		self.img[:,:] = self.bgCol

		#draw labels
		textLocX = self.plotWidth
		textLocY = 0#self.plotHeight
		for nameIndex in range(len(self.names) - 1, -1, -1):
			namesize = cv2.getTextSize(self.names[nameIndex], self.font, self.fontsize, self.fontthickness)
			nameW = int(namesize[0][0])
			nameH = int(namesize[1])
			if nameIndex == len(self.names) - 1:
				textLocX -= nameW + self.edgsps + self.textsps
				textLocY += nameH + self.edgsps + self.textsps
			else:
				textLocX -= nameW + self.sps + self.textsps * 2

			DrawRectangle(
				self.img,
				(textLocX - self.textsps, textLocY + self.textsps),
				(textLocX + nameW + self.textsps, self.sps),
				color=self.cols[nameIndex],
				thickness=-1
			)
			DrawRectangle(
				self.img,
				(textLocX - self.textsps, textLocY + self.textsps),
				(textLocX + nameW + self.textsps, self.sps),
				color=self.textCol,
				thickness=2
			)
			avgCol = 0
			for _ in range(3):
				avgCol += self.cols[nameIndex][_]
			avgCol /= 3
			cv2.putText(
				self.img,
				self.names[nameIndex],
				(textLocX, textLocY),
				self.font,
				self.fontsize,
				(0,0,0) if avgCol > 127 else (255, 255, 255)
			)

		#setup axis
		if axis:
			DrawLine(self.img, (0, self.centY), (self.plotWidth, self.centY),color=self.tickCol)
			DrawLine(self.img, (int(self.plotWidth/10), 0), (int(self.plotWidth/10), self.plotHeight), color=self.tickCol)
			vertTicks(self.img, 0, self.plotHeight, int(self.plotWidth/10), self.tickHeight, int(self.tickSps), col=self.tickCol)
			horizTicks(self.img, 0, self.plotWidth, self.centY, self.tickHeight, int(self.tickSps), col=self.tickCol)
	def printRange(self):
		plotmax = int(abs(np.amax(np.asarray(self.vals))))
		plotmin = int(abs(np.amin(np.asarray(self.vals))))
		graphRange = 0 * self.plotPad + plotmax + plotmin
		print("{:.2f}, {:.2f}, {:.2f}".format(plotmin, plotmax, graphRange))
	def graphLine(self, lineIndex=0):
		line = self.vals[lineIndex]
		xSpace = int(self.plotWidth / (self.npoints))
		for i in range(1, len(line)):
			if line[i-1] is None or line[i] is None:
				break
			xVal = (i) * xSpace
			if self.plot_minimum != None and self.plot_maximum != None:
				yVal0 = int(mapVal(line[i-1], self.plot_minimum, self.plot_maximum, 0, self.plotHeight))
				yVal1 = int(mapVal(line[i], self.plot_minimum, self.plot_maximum, 0, self.plotHeight))
			else:
				yVal0 = int(line[i-1] + self.centY)
				yVal1 = int(line[i] + self.centY)

			cv2.line(self.img, (xVal, yVal0), (xVal + xSpace, yVal1), self.cols[lineIndex], 2)
	def graphPoints(self):
		self.setupPage()
		for line in range(self.nlines):
			self.graphLine(line)
	def addVal(self, value, lineIndex=0, left = True):
		#value += int(self.img.shape[0]/2)
		self.vals[lineIndex].append(float(value)) if left else self.vals[lineIndex].append((value))
	def get_data(self):
		arr = np.zeros((self.vals.shape))
		try:
		 	arr = np.asarray(self.vals, dtype=np.float64)
		except:
			print(self.vals[-1])
			arr = []
		return arr

def mapVal(val, inMin, inMax, outMin = 0, outMax = 1):
	if inMin < inMax:
		if val >=  inMax:
			return outMax
		if val <=  inMin:
			return outMin

	inRange = inMax - inMin
	outRange = outMax - outMin
	inPerc = (val - inMin) / inRange
	outVal = outMin + outRange * inPerc
	if outVal > outMax:
		return outMax
	elif outVal < outMin:
		return outMin
	return outVal
def DrawLine(img, pnt1, pnt2, color=(0,0,0),thickness=2):
	cv2.line(img, (int(pnt1[0]), int(pnt1[1])), (int(pnt2[0]), int(pnt2[1])), color=color, thickness=thickness)
def DrawLineAtAngle(img, pnt, length, theta, color=(0, 0, 255), thickness=2):
	cv2.line(img, (int(pnt[0]), int(pnt[1])), (int(pnt[0] + length * math.cos(theta)), int(pnt[1] + length * math.sin(theta))), color=color, thickness=thickness)
def DrawCircle(img,x,y,rad = 2,col = (0,0,255), thickness = 2, lineType = 8, shift = 0):
	if rad < 0:
		raise NegativeMagnitude("Radius cannot be negative")
	return cv2.circle(img, (int(x),int(y)), rad, col, thickness = thickness, lineType = lineType, shift = shift)
def DrawRectangle(image, p0, p1, color = (0,255,0),thickness = 2):
	cv2.rectangle(image, p0, p1, color, thickness)
def DrawGrid(image, x, y, width, height, rows, columns):
	col_spacing = (width / columns)
	row_spacing = (height / rows)
	#DrawCircle(image, x, y, rad=5)
	#DrawCircle(image, x+width/2, y+height/2, rad=5)
	DrawRectangle(image, (x, y), (x + width, y + height))
	for i in range(rows):
		DrawLine(image,
			(int(x), int(i*row_spacing + y)),
			(int(x + width), int(i*row_spacing + y)))
	for i in range(columns):
		DrawLine(image,
			(int(x + i*col_spacing), int(y)),
			(int(x + i*col_spacing), int(y + height)))
