from kafka import KafkaProducer
from json import dumps
import time

producer = KafkaProducer(acks=0, compression_type='gzip', bootstrap_servers=['localhost:9092'], value_serializer=lambda x: dumps(x).encode('utf-8'))

# start = time.time()
######################################
for i in range(1):
    data = {
        "node_id" : "DRO3",
        "values" : [37.51780779452564, 126.87847049489307, 0, 0, 80, 0], # drone: [lat, long, alt, velocity, batteryPer, done]
        "timestamp" : "2021-09-24 18:09:39"
    }
    print("publish: ", type(data), data, "\n")
    producer.send('sensor-data', value=data) # topic name: sensor-data, value: data
    producer.flush()
    time.sleep(5)

for i in range(1):
    data = {
        "node_id" : "DRO3",
        "values" : [37.51788893119509, 126.87851560992833, 0, 0, 80, 0], # drone: [lat, long, alt, velocity, batteryPer, done]
        "timestamp" : "2021-09-24 18:09:44"
    }
    print("publish: ", type(data), data, "\n")
    producer.send('sensor-data', value=data) # topic name: sensor-data, value: data
    producer.flush()
    time.sleep(5)

for i in range(1):
    data = {
        "node_id" : "DRO3",
        "values" : [37.51789129998933, 126.87862872263788, 0, 0, 80, 0], # drone: [lat, long, alt, velocity, batteryPer, done]
        "timestamp" : "2021-09-24 18:09:49"
    }
    print("publish: ", type(data), data, "\n")
    producer.send('sensor-data', value=data) # topic name: sensor-data, value: data
    producer.flush()
######################################
# print("elapsed :", time.time() - start)

# station: [temperature, humidity, light, lat, long, alt]
# "values" : [37.5176341, 126.8785336, 0, 0, 80, 0], # drone: [lat, long, alt, velocity, batteryPer, done] 1
# "values" : [37.517684, 126.8786976, 0, 0, 75, 1], # drone: [lat, long, alt, velocity, batteryPer, done]