from kafka import KafkaProducer
from json import dumps
import time

producer = KafkaProducer(acks=0, compression_type='gzip', bootstrap_servers=['localhost:9092'], value_serializer=lambda x: dumps(x).encode('utf-8'))

start = time.time()
for i in range(1):
    data = {
        "node_id" : "DRO2",
        "values" : [1, 2, 3, 4, 5, 1], # [sensor1, sensor2, sensor3, sensor4, sensor5, done]
        "timestamp" : "2021-08-26 16:57:05"
    }
    print("publish: ", type(data), data, "\n")
    producer.send('sensor-data', value=data) # topic name: sensor-data, value: data
    producer.flush()
print("elapsed :", time.time() - start)
