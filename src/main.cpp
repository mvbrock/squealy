#include <unistd.h>
#include <string>
#include <iostream>

#include <zmq.h>
#include <czmq.h>
#include <zgossip.h>
#include <zactor.h>


using namespace std;

int main(int argc, char **argv) {
    string bindaddr = "tcp://*:5555";
    string connectaddr = "tcp://127.0.0.1:5555";
    string tuplekey = "brock-key";
    for(;;) {
        switch(getopt(argc, argv, "b:c:k:")) {
            case 'b':
                bindaddr = string(optarg);
                continue;
            case 'c':
                connectaddr = string(optarg);
                continue;
            case 'k':
                tuplekey = string(optarg);
                continue;
            case -1:
                break;
            default:
                break;
        }
        break;
    }
    // Create the server
    zactor_t *server = zactor_new(zgossip, (void*)"server");
    zstr_send(server, "VERBOSE");
    
    // Bind to an address
    zstr_sendx(server, "BIND", bindaddr.c_str(), NULL);

    // Connect to another address
    zstr_sendx(server, "CONNECT", connectaddr.c_str(), NULL);

    // Set a value upon connecting
    zstr_sendx(server, "PUBLISH", tuplekey.c_str(), "brock-value", NULL);

    // Receive any published values
    while(true) {
        char *cmd, *key, *value;
        zstr_recvx(server, &cmd, &key, &value, NULL);
        cout << "cmd: " << cmd << ", key: " << key << ", value: " << value << endl;
    }
    zactor_destroy(&server);
}
