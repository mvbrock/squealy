CZMQ_BASE=/opt/homebrew/Cellar/czmq/4.2.1
CZMQ_INC=$(CZMQ_BASE)/include
CZMQ_LIB=$(CZMQ_BASE)/lib

ZEROMQ_BASE=/opt/homebrew/Cellar/zeromq/4.3.4
ZEROMQ_INC=$(ZEROMQ_BASE)/include
ZEROMQ_LIB=$(ZEROMQ_BASE)/lib

SRC = $(wildcard src/*.cpp)
OBJ = $(patsubst src/%.cpp, obj/%.o, $(SRC))
INC = -I $(ZEROMQ_INC) -I $(CZMQ_INC)
LIBDIR = -L $(ZEROMQ_LIB) -L $(CZMQ_LIB)
LIBS = -lzmq -lczmq

bin/squealy: $(OBJ)
	mkdir -p bin
	g++ -o $@ $? $(LIBDIR) $(LIBS)

obj/%.o: src/%.cpp
	mkdir -p obj
	g++ -o $@ -c $< $(INC)

.PHONY: clean
clean:
	rm obj/*
	rm bin/squealy
