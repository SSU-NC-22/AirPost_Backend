from kafka import KafkaProducer
from json import dumps
import time

producer = KafkaProducer(acks=0, compression_type='gzip', bootstrap_servers=['localhost:9092'], value_serializer=lambda x: dumps(x).encode('utf-8'))

start = time.time()
for i in range(5):
    data = {
        "sensor_id" : "sensor-" + str(i),
        "node_id" : "node-" + str(i),
        "values" : [1, 2, 3],
        "timestamp" : "2021-08-16"
    }
    print(data)
    producer.send('sensor-data', value=data) # topic name: sensor-data, value: data
    producer.flush()
print("elapsed :", time.time() - start)
