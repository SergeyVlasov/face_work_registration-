# face_work_registration-


1) тренировка модели 


 создаем виртуальное кружение
- python3 -m virtualenv /path/MyEnv
активируем виртуальное окружение
- source /path/MyEnv/bin/activate
устанавливаем все необходимые библиотеки
- pip3 install opencv-python
- pip3 install dlib
- pip3 install face_recognition
- pip3 install imutils
** если возникают проблемы на этапе установки dlib, то попробуйте установить pip install cmake

моздаем директорию для проекта
- mkdir face_recognition
- cd face_recognition
создаем папку для обучающей выборки а в ней подпапки с соответствующей персоной, и наполняем их соответствующими фото
- mkdir dataset
- cd dataset
- mkdir person1
- mkdir person2 ....
копируем github
- git clone https://github.com/SergeyVlasov/face_recognition

загружаем файл модели
- wget https://raw.githubusercontent.com/opencv/opencv/master/data/haarcascades/haarcascade_frontalface_default.xml
проект должен иметь структуру
```
├── dataset
│   ├── Person1
│   │   ├── 001.jpg
│   │   ├── 002.jpg
│   │   ├── 003.jpg
│   │   ├── 004.jpg
│   │   ├── 005.jpg
│   │   ├── 006.jpg
│   │   ├── 007.jpg
│   │   └── 008.jpeg
│   └── Person2
│      ├── 001.jpg
│      ├── 002.jpg
│      ├── 003.jpg
│      ├── 004.jpg
│      └── 005.jpg
├── train.py
└── haarcascade_frontalface_default.xml
```
подготавливаем датасет и тренируем модель
- python3 train.py --dataset dataset --encodings encodings.pickle --detection-method hog

---------------------------------------------------------------------------------


2) распознавание лиц
на компьютере, котором осуществляется распознавание, тоже нужно установить все библиотеки
создаем виртуальное кружение
- python3 -m virtualenv /path/MyEnv
активируем виртуальное окружение
- source /path/MyEnv/bin/activate
устанавливаем все необходимые библиотеки
- pip3 install opencv-python
- pip3 install dlib
- pip3 install face_recognition
- pip3 install imutils
** если возникают проблемы на этапе установки dlib, то попробуйте установить pip install cmake

копируем 2 файла (с системы где осуществлялось обучение)

- haarcascade_frontalface_default.xml
- encodings.pickle

запускаем файл распознавания

- python3 video_recognition.py --cascade haarcascade_frontalface_default.xml --encodings encodings.pickle


---------------------------------------------------------------------------------

3) создание БД в Postgres

```
sudo -i -u postgres psql

create database worktime;
\connect worktime;

CREATE TABLE users (
id SERIAL,
number int,
name1 varchar,
name2  varchar,
name3  varchar,
note  varchar);

INSERT INTO public.users(number, name1, name2, name3)	VALUES ( 001, 'Иванов', 'Иван', 'Иванович');

CREATE TABLE checktime (
id SERIAL,
date varchar,
time varchar,
iduser int,
inout int);

INSERT INTO public.checktime( date, time, iduser, inout) VALUES ('20-10-20', '8:00', 1, 1);
```







