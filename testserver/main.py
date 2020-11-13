import sys
from opcua import ua, Server
import time
server = Server()
server.set_endpoint("opc.tcp://127.0.0.1:4840/")
server.set_server_name("Server")
objects = server.get_objects_node()
#создаем объект и присваиваем ему имя 
uri = "http://server"
idx = server.register_namespace(uri)
Object_1 =objects.add_object(idx,"MyFirstObject")
#теперь создаем переменные
a ={} 
f = open("data.txt", "r")
for line in f:
    a[line.split(":")[0]] = (Object_1.add_variable(idx,line.split(":")[0],float(line.split(":")[1])))  
    a[line.split(":")[0]].set_writable() 
f.close()    
server.start() 
while True:
    f = open("data.txt","r")
    for line in f:
        a[line.split(":")[0]].set_value(float(line.split(":")[1]))
    f.close() 
    time.sleep(5)

