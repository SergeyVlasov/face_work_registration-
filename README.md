# face_work_registration-



1) создание БД в Postgres

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


```
source /home/serg/Desktop/face_recognition/bin/activate


python3 video_recognition.py --cascade haarcascade_frontalface_default.xml --encodings encodings.pickle
```



