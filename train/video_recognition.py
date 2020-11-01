

# python3 video_recognition.py --cascade haarcascade_frontalface_default.xml --encodings encodings.pickle

# import library yang di perlukan
from imutils.video import VideoStream
from imutils.video import FPS
import face_recognition
import argparse
import imutils
import pickle
import time
import cv2

import datetime
import requests


# Parsing Argumen
ap = argparse.ArgumentParser()
ap.add_argument("-c", "--cascade", required=True,
	help = "path to where the face cascade resides")
ap.add_argument("-e", "--encodings", required=True,
	help="path to serialized db of facial encodings")
args = vars(ap.parse_args())

# load file cascade OpenCV
print("[INFO] loading encodings + face detector...")
data = pickle.loads(open(args["encodings"], "rb").read())
detector = cv2.CascadeClassifier(args["cascade"])

# Camera
print("[INFO] Stream Camera...")
vs = VideoStream(src=0).start()
time.sleep(2.0)

# Penghitung FPS (Frame per Second)
fps = FPS().start()

# loop 
while True:
	
	frame = vs.read()
	frame = imutils.resize(frame, width=500)
	
	# RGB
	gray = cv2.cvtColor(frame, cv2.COLOR_BGR2GRAY)
	rgb = cv2.cvtColor(frame, cv2.COLOR_BGR2RGB)

	# detect
	rects = detector.detectMultiScale(gray, scaleFactor=1.1, 
		minNeighbors=5, minSize=(30, 30),
		flags=cv2.CASCADE_SCALE_IMAGE)

   # boxes
	boxes = [(y, x + w, y + h, x) for (x, y, w, h) in rects]

	encodings = face_recognition.face_encodings(rgb, boxes)
	names = []

	# loop 
	for encoding in encodings:
		matches = face_recognition.compare_faces(data["encodings"],
			encoding)
		name = "Unknown"

		# check i
		if True in matches:
			matchedIdxs = [i for (i, b) in enumerate(matches) if b]
			counts = {}
			for i in matchedIdxs:
				name = data["names"][i]
                                
				counts[name] = counts.get(name, 0) + 1
			name = max(counts, key=counts.get)
		print(name + " " + datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S'))
		if name == 'yakubovich':
			requests.get('http://localhost:8080/'+ datetime.datetime.now().strftime('%Y-%m-%d') +'/'+ datetime.datetime.now().strftime('%H:%M:%S') +'/1/1')
		names.append(name)
                

	# loop 
	for ((top, right, bottom, left), name) in zip(boxes, names):
		
		cv2.rectangle(frame, (left, top), (right, bottom),
			(0, 255, 0), 2)
		y = top - 15 if top - 15 > 15 else top + 15
		cv2.putText(frame, name, (left, y), cv2.FONT_HERSHEY_SIMPLEX,
			0.75, (0, 255, 0), 2)

	
	cv2.imshow("Frame", frame)
	key = cv2.waitKey(1) & 0xFF

	
	if key == ord("q"):
		break

	# update FPS
	fps.update()

#  info FPS
fps.stop()
print("[INFO] elasped time: {:.2f}".format(fps.elapsed()))
print("[INFO] approx. FPS: {:.2f}".format(fps.fps()))

# cleanup
cv2.destroyAllWindows()
vs.stop()
