import socket
import time

HOST = '127.0.0.1'  # The server's hostname or IP address
PORT = 2575        # The port used by the server

START_BLOCK = chr(int("0x0B", 16))
END_BLOCK = chr(int("0x1C", 16))
CR = chr(int("0x0D", 16))

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.connect((HOST, PORT))
    s.sendall(f'{START_BLOCK}Hello{CR}morning{END_BLOCK}{CR}'.encode('latin-1'))
    s.sendall(f'{START_BLOCK}Hello{CR}world{END_BLOCK}{CR}'.encode('latin-1'))
    data = s.recv(1024)
    data2 = s.recv(1024)

print('Received', repr(data))
print('Received', repr(data2))

