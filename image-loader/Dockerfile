# syntax=docker/dockerfile:1
FROM python:3.8-buster

WORKDIR /app
COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt
RUN pip3 install python-multipart

COPY . .

CMD [ "python3", "py_inference.py"]

EXPOSE 33334
