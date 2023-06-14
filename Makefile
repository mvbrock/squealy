CZMQ_BASE=/opt/homebrew/Cellar/czmq/4.2.1
CZMQ_INC=$(CZMQ_BASE)/include
CZMQ_LIB=$(CZMQ_BASE)/lib

ZEROMQ_BASE=/opt/homebrew/Cellar/zeromq/4.3.4
ZEROMQ_INC=$(ZEROMQ_BASE)/include
ZEROMQ_LIB=$(ZEROMQ_BASE)/lib

FLAGS = -std=c++20 -stdlib=libc++
GXX = g++ $(FLAGS)
SRC = $(wildcard src/*.cpp)
OBJ = $(patsubst src/%.cpp, obj/%.o, $(SRC))
INC = -I $(ZEROMQ_INC) -I $(CZMQ_INC)
LIBDIR = -L $(ZEROMQ_LIB) -L $(CZMQ_LIB)
LIBS = -lzmq -lczmq

bin/squealy: $(OBJ)
	mkdir -p bin
	$(GXX) -o $@ $? $(LIBDIR) $(LIBS)

obj/%.o: src/%.cpp
	mkdir -p obj
	$(GXX) -o $@ -c $< $(INC)

.PHONY: clean
clean:
	rm -f obj/*
	rm -f bin/squealy
