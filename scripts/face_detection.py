import cv2
import sys

def detect_faces(image_path, output_path):
    # Load the image
    image = cv2.imread(image_path)

    # Load the pre-trained classifier for face detection
    face_cascade = cv2.CascadeClassifier(cv2.data.haarcascades + 'haarcascade_frontalface_default.xml')

    # Convert the image to grayscale
    gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)

    # Detect faces
    faces = face_cascade.detectMultiScale(gray, scaleFactor=1.1, minNeighbors=5, minSize=(30, 30))

    # Draw rectangles around the detected faces
    for (x, y, w, h) in faces:
        cv2.rectangle(image, (x, y), (x+w, y+h), (255, 0, 0), 2)

    # Save the modified image
    cv2.imwrite(output_path, image)

    # Get the number of detected faces
    face_count = len(faces)
    
    # Print the number of detected faces
    print(f'Number of detected faces: {face_count}')

    # Return the number of detected faces
    return face_count

if __name__ == '__main__':
    image_path = sys.argv[1]
    output_path = sys.argv[2]
    detected_faces = detect_faces(image_path, output_path)
    sys.stdout.write(str(detected_faces))  # Send the number of faces as output
