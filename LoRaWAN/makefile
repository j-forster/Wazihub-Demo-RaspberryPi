#options are: -DPABOOST -DMAX_DBM=v
#
#CFLAGS=-DPABOOST -DMAX_DBM=14

lora_gateway: lora_gateway.o arduPi.o SX1272.o
	g++ -lrt -lpthread lora_gateway.o arduPi.o SX1272.o -o lora_gateway

lora_gateway.o: lora_gateway.cpp
	g++ $(CFLAGS) -DRASPBERRY -DRASPBERRY2 -DIS_RCV_GATEWAY -c lora_gateway.cpp -o lora_gateway.o

arduPi.o: arduPi_pi2.cpp arduPi_pi2.h
	g++ -c arduPi_pi2.cpp -o arduPi.o	

SX1272.o: SX1272.cpp SX1272.h
	g++ -DRASPBERRY2 -c SX1272.cpp -o SX1272.o

clean:
	rm *.o lora_gateway
